package compact

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

const (
	compactUTCNanoFormat = "20060102150405.000000000"
)

var builderPool = sync.Pool{
	New: func() any {
		return new(strings.Builder)
	},
}

// Handler implements a [slog.Handler] that outputs log records in compact
// format mostly meant for testing.
//
// Log records are formatted in the following way:
//
// <short_nano_ts> <log_level> [<source_line>] <msg> <k=v attrs>
//
// <short_name_ts> is a shortened UTC-standardized timestamp with the following
// format:
//
// # YYYYMMDDhhmmss.nnnnnnnnn
//
// <source_line> is a shortened filepath and line number that only appears in
// log lines when the log level is DEBUG.
//
// For example, given a log record with this makeup:
//   - frame:
//     file: /home/jaypipes/src/github.com/relexec/pkg/log/config.go
//     line: 20
//   - time: 2026-10-10T11:08:25.744464185-04:00
//   - level: DEBUG
//   - msg: "received pause signal"
//   - attrs:
//     run_uuid="9ea59727-665d-4be7-8a18-77c2cf0e057c"
//     target_uuid="50c2a5e6-4266-411e-b947-bfab27d10170"
//     target_generation="1"
//
// The compact logger when configured for DEBUG log level would output the
// following log line:
//
// 20261010150825.744464185 DEBUG pkg/log/config.go:20 "received pause signal" run_uuid="9ea59727-665d-4be7-8a18-77c2cf0e057c target_uuid="50c2a5e6-4266-411e-b947-bfab27d10170" target_generation= "1"
//
// Ref: https://github.com/golang/example/blob/master/slog-handler-guide/guide.md
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

// NewHandler returns a new slog.Handler that outputs a compact format mostly
// meant for testing.
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
		b.WriteString(r.Time.UTC().Format(compactUTCNanoFormat))
	}
	b.WriteRune(' ')
	b.WriteString(r.Level.String())
	b.WriteRune(' ')
	if h.Enabled(ctx, slog.LevelDebug) && r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		sl := logutil.FormatSourceFileAndLine(f)
		b.WriteString(strconv.Quote(sl))
		b.WriteRune(' ')
	}

	b.WriteString(strconv.Quote(r.Message))
	b.WriteRune(' ')

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
		b.WriteString(a.Value.Time().UTC().Format(time.RFC3339Nano))
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
