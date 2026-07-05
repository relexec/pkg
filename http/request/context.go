package request

import (
	"context"
	"log/slog"
)

type contextKey string

const (
	ContextKeyRequestID     contextKey = "request.id.t2"
	ContextKeyRequestLogger contextKey = "request.logger.t2"
	ContextKeyRequestAuth   contextKey = "request.auth.t2"
)

// IDFromContext returns the Request ID contained in the supplied context, or empty
// string if not found in the context.
func IDFromContext(ctx context.Context) string {
	v := ctx.Value(ContextKeyRequestID)
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
}

// LoggerFromContext returns the Request Logger contained in the supplied
// context, or a discard logger if not found in the context.
func LoggerFromContext(ctx context.Context) *slog.Logger {
	v := ctx.Value(ContextKeyRequestLogger)
	logger, ok := v.(*slog.Logger)
	if !ok {
		return slog.New(slog.DiscardHandler)
	}
	return logger
}

// AuthFromContext returns the Auth struct containing authenticated identity
// information contained in the supplied context, or empty string if not found
// in the context.
func AuthFromContext(ctx context.Context) *Auth {
	v := ctx.Value(ContextKeyRequestAuth)
	a, ok := v.(*Auth)
	if !ok {
		return nil
	}
	return a
}
