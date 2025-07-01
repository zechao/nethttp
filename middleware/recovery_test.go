package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/zechao158/nethhtp/middleware"
)

func TestRecovery_Panic(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("something went wrong")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	middleware.Recovery(handler).ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Internal Server Error") {
		t.Errorf("expected body to contain 'Internal Server Error', got %q", rec.Body.String())
	}
}

func TestRecovery_NoPanic(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	middleware.Recovery(handler).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if rec.Body.String() != "OK" {
		t.Errorf("expected body 'OK', got %q", rec.Body.String())
	}
}
