package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const TraceIDHeader = "X-Trace-Id"

type contextKey string

const traceIDKey contextKey = "traceID"

// TraceIDFromContext retrieves the trace ID from the context.
func TraceIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(traceIDKey).(string); ok {
		return v
	}
	return ""
}

// TraceMiddleware is an HTTP middleware that ensures a trace ID is present.
func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get(TraceIDHeader)
		if traceID == "" {
			traceID = uuid.New().String()
		}
		// Set the trace ID in the response header for visibility
		w.Header().Set(TraceIDHeader, traceID)
		// Store the trace ID in the request context
		ctx := context.WithValue(r.Context(), traceIDKey, traceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
