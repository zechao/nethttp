package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/zechao158/nethhtp/middleware"
)

func TestCORS_DefaultConfig(t *testing.T) {
	handler := middleware.CORS()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	if got := resp.Header.Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, "*")
	}
	if got := resp.Header.Get("Access-Control-Allow-Methods"); got != "GET,POST,PUT,DELETE,OPTIONS" {
		t.Errorf("Access-Control-Allow-Methods = %q, want %q", got, "GET,POST,PUT,DELETE,OPTIONS")
	}
	if got := resp.Header.Get("Access-Control-Allow-Headers"); got != "Content-Type,Authorization" {
		t.Errorf("Access-Control-Allow-Headers = %q, want %q", got, "Content-Type,Authorization")
	}
	if got := resp.Header.Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Errorf("Access-Control-Allow-Credentials = %q, want %q", got, "true")
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestCORS_OptionsRequest(t *testing.T) {
	handler := middleware.CORS()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called for OPTIONS requests")
	}))

	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusNoContent)
	}
}

func TestCORS_CustomConfig(t *testing.T) {

	handler := middleware.CORS(
		middleware.WithAllowOrigin("https://example.com"),
		middleware.WithAllowMethods(http.MethodPatch, http.MethodHead),
		middleware.WithAllowHeaders("X-Custom-Header"),
		middleware.WithAllowCredentials("false"),
	)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))

	req := httptest.NewRequest(http.MethodPatch, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	if got := resp.Header.Get("Access-Control-Allow-Origin"); got != "https://example.com" {
		t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, "https://example.com")
	}
	if got := resp.Header.Get("Access-Control-Allow-Methods"); got != "PATCH,HEAD" {
		t.Errorf("Access-Control-Allow-Methods = %q, want %q", got, "PATCH,HEAD")
	}
	if got := resp.Header.Get("Access-Control-Allow-Headers"); got != "X-Custom-Header" {
		t.Errorf("Access-Control-Allow-Headers = %q, want %q", got, "X-Custom-Header")
	}
	if got := resp.Header.Get("Access-Control-Allow-Credentials"); got != "false" {
		t.Errorf("Access-Control-Allow-Credentials = %q, want %q", got, "false")
	}
	if resp.StatusCode != http.StatusAccepted {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusAccepted)
	}
}

func TestCORS_AllowMethodsJoin(t *testing.T) {
	methods := []string{"GET", "POST", "PATCH"}
	handler := middleware.CORS(
		middleware.WithAllowMethods(methods...))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	got := resp.Header.Get("Access-Control-Allow-Methods")
	want := strings.Join(methods, ",")
	if got != want {
		t.Errorf("Access-Control-Allow-Methods = %q, want %q", got, want)
	}
}
