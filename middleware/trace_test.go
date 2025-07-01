package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/zechao158/nethhtp/middleware"
)

func TestTraceMiddleware_GeneratesTraceIDIfMissing(t *testing.T) {
	var gotTraceID string

	handler := middleware.TraceMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotTraceID = middleware.TraceIDFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	resp := rr.Result()
	traceIDHeader := resp.Header.Get(middleware.TraceIDHeader)
	if traceIDHeader == "" {
		t.Fatal("expected trace ID header to be set")
	}
	if gotTraceID == "" {
		t.Fatal("expected trace ID in context")
	}
	if traceIDHeader != gotTraceID {
		t.Errorf("trace ID in header and context do not match: header=%q, context=%q", traceIDHeader, gotTraceID)
	}
}

func TestTraceMiddleware_UsesExistingTraceID(t *testing.T) {
	const existingTraceID = "test-trace-id-123"
	var gotTraceID string

	handler := middleware.TraceMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotTraceID = middleware.TraceIDFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(middleware.TraceIDHeader, existingTraceID)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	resp := rr.Result()
	traceIDHeader := resp.Header.Get(middleware.TraceIDHeader)
	if traceIDHeader != existingTraceID {
		t.Errorf("expected trace ID header to be %q, got %q", existingTraceID, traceIDHeader)
	}
	if gotTraceID != existingTraceID {
		t.Errorf("expected trace ID in context to be %q, got %q", existingTraceID, gotTraceID)
	}
}

func TestTraceMiddleware_HeaderIsSetOnResponse(t *testing.T) {
	handler := middleware.TraceMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	resp := rr.Result()
	traceIDHeader := resp.Header.Get(middleware.TraceIDHeader)
	if strings.TrimSpace(traceIDHeader) == "" {
		t.Error("expected trace ID header to be set on response")
	}
}
