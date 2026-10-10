package logfmt

import (
	"context"
	"io"
	"log/slog"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/relexec/pkg/log/logutil"
)

var builderPool = sync.Pool{
	New: func() any {
		return new(strings.Builder)
	},
}

// Handler implements a [slog.Handler] that outputs log records in proper
// logfmt format.
//
// Ref: https://github.com/golang/example/blob/master/slog-handler-guide/guide.md
// Ref: https://brandur.org/logfmt
type Handler struct {
	mu   *sync.Mutex
	out  io.Writer
	opts *slog.HandlerOptions
	// preformatted contains the formatted output of Attrs associated with
	// already opened Groups. This prevents the need to format Attrs on each
	// call to Log/LogAttrs for groups and attrs that were associated with the
	// Handler from a call to WithGroup or WithAttrs.
	preformatted string
}

// NewHandler returns a new slog.Handler that abides by the official logfmt
// standard.
func NewHandler(out io.Writer, opts *slog.HandlerOptions) *Handler {
	return &Handler{
		mu:   &sync.Mutex{},
		out:  out,
		opts: opts,
	}
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *Handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	h2 := *h
	return &h2
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	h2 := *h

	b := builderPool.Get().(*strings.Builder)
	defer builderPool.Put(b)
	b.Reset()

	for _, a := range attrs {
		h2.appendAttr(b, a)
	}
	h2.preformatted = b.String()
	return &h2
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	b := builderPool.Get().(*strings.Builder)
	defer builderPool.Put(b)
	b.Reset()
	if !r.Time.IsZero() {
		h.appendAttr(b, slog.Time(slog.TimeKey, r.Time))
	}
	h.appendAttr(b, slog.Any(slog.LevelKey, r.Level))
	if h.Enabled(ctx, slog.LevelDebug) && r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		b.WriteString(slog.SourceKey)
		b.WriteRune('=')
		sl := logutil.FormatSourceFileAndLine(f)
		b.WriteString(strconv.Quote(sl))
		b.WriteRune(' ')
	}

	h.appendAttr(b, slog.String(slog.MessageKey, r.Message))

	b.WriteString(h.preformatted)

	r.Attrs(
		func(a slog.Attr) bool {
			h.appendAttr(b, a)
			return true
		},
	)

	b.WriteRune('\n')

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.out.Write([]byte(b.String()))
	return err
}

func (h *Handler) appendAttr(b *strings.Builder, a slog.Attr) {
	if a.Equal(slog.Attr{}) {
		return
	}
	a.Value = a.Value.Resolve()
	switch a.Value.Kind() {
	case slog.KindDuration:
		b.WriteString(a.Key)
		b.WriteRune('=')
		b.WriteString(strconv.FormatInt(a.Value.Duration().Nanoseconds(), 10))
	case slog.KindTime:
		b.WriteString(a.Key)
		b.WriteRune('=')
		b.WriteString(a.Value.Time().Format(time.RFC3339Nano))
	case slog.KindGroup:
		attrs := a.Value.Group()
		if len(attrs) == 0 {
			return
		}
		for _, ga := range attrs {
			h.appendAttr(b, ga)
		}
	default:
		b.WriteString(a.Key)
		b.WriteRune('=')
		b.WriteString(strconv.Quote(a.Value.String()))
	}
	b.WriteRune(' ')
}
