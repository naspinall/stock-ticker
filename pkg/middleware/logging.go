package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type ContextKey string

const TraceIDKey = ContextKey("trace-id")

type ContextHandler struct {
	slog.Handler
}

func NewContextHandler(handler slog.Handler) *ContextHandler {
	return &ContextHandler{
		handler,
	}
}

func (ch *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	fmt.Println()
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
		r.AddAttrs(slog.String("trace-id", string(traceID)))
	}

	return ch.Handler.Handle(ctx, r)
}

type Logging struct {
	logger *slog.Logger
}

func NewLogging(logger *slog.Logger) *Logging {
	return &Logging{
		logger: logger,
	}
}

func (l *Logging) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Track request time
		requestStart := time.Now()

		requestID := uuid.NewString()

		r = r.WithContext(
			context.WithValue(r.Context(), TraceIDKey, requestID),
		)

		next.ServeHTTP(w, r)

		// Calculate the time length of the request
		requestEnd := time.Now()
		requestDuration := requestEnd.Sub(requestStart)

		l.logger.InfoContext(r.Context(), "access", "url", r.URL.String(), "duration_seconds", requestDuration.Seconds())

	})
}

func (l *Logging) AttachMiddleware(handlerFunc http.HandlerFunc) http.Handler {
	return l.Middleware(http.HandlerFunc(handlerFunc))
}
