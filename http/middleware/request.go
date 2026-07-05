package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"
	chimid "github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/relexec/pkg/http/metrics"
	"github.com/relexec/pkg/http/request"
)

const (
	headerKeyRequestID = "X-Request-ID"
)

// AdaptLoggerRoutePattern is middleware that customizes the request context's
// logger to ensure the RoutePattern is included in the request slog.Group.
// This is needed because of the way Chi lazily evaluates the route pattern.
// This middleware must be added to a chi.Router-specific handler as opposed to
// a global middleware.
func AdaptLoggerRoutePattern(ctx huma.Context, next func(huma.Context)) {
	rctx := ctx.Context()
	pattern := chi.RouteContext(rctx).RoutePattern()
	if pattern == "" {
		reqLogger := request.LoggerFromContext(rctx)
		if reqLogger != nil {
			reqID := request.IDFromContext(rctx)

			reqLogger = reqLogger.With(
				slog.Group("request",
					slog.String("id", reqID),
					slog.String("route", pattern),
				),
			)
			ctx = huma.WithValue(
				ctx,
				request.ContextKeyRequestLogger,
				reqLogger,
			)
		}
	}

	next(ctx)
}

// RequestContext ensures that there is a request ID in the request's context
// and that the request logger and metrics handler is appropriately
// established.
func RequestContext(
	rootLogger *slog.Logger,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			metrics.InstrumentRequestsCurrent.Add(ctx, 1)
			defer metrics.InstrumentRequestsCurrent.Add(ctx, -1)

			requestID := r.Header.Get(headerKeyRequestID)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			wrapped := chimid.NewWrapResponseWriter(w, r.ProtoMajor)

			ctx = context.WithValue(
				ctx,
				request.ContextKeyRequestID,
				requestID,
			)
			reqLogger := rootLogger.With(
				slog.Group(
					"request",
					slog.String("id", requestID),
				),
			)
			ctx = context.WithValue(
				ctx,
				request.ContextKeyRequestLogger,
				reqLogger,
			)

			start := time.Now()
			next.ServeHTTP(wrapped, r.WithContext(ctx))
			elapsed := time.Since(start).Seconds()

			// NOTE(jaypipes): Because of the way Chi lazily evaluates the
			// route pattern, we need to access the route AFTER the call to
			// next.ServeHTTP().
			//
			// See: AdaptLoggerRoutePattern
			route := chi.RouteContext(ctx).RoutePattern()
			statusCode := wrapped.Status()

			logger := request.LoggerFromContext(ctx)

			logger.Debug(
				fmt.Sprintf(
					"%s %s -> %d in %0.5fs",
					r.Method, route, statusCode, elapsed,
				),
			)

			attrs := []attribute.KeyValue{
				metrics.AttributeMethod(r.Method),
				metrics.AttributeRoute(route),
				metrics.AttributeStatusCode(statusCode),
			}
			metrics.InstrumentRequests.Add(
				ctx, 1,
				metric.WithAttributes(attrs...),
			)
			metrics.InstrumentRequestsDuration.Record(
				ctx, elapsed,
				metric.WithAttributes(attrs...),
			)
		})
	}
}
