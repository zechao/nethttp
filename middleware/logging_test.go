package middleware_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/zechao158/nethhtp/middleware"
)

func TestLoggingMiddleware_StatusOK(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello"))
	})

	req := httptest.NewRequest("GET", "/testpath", nil)
	rr := httptest.NewRecorder()

	loggedHandler := middleware.Logging(handler)
	loggedHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	logOutput := buf.String()
	if !strings.Contains(logOutput, "200") {
		t.Errorf("log output does not contain status code: %s", logOutput)
	}
	if !strings.Contains(logOutput, "GET") || !strings.Contains(logOutput, "/testpath") {
		t.Errorf("log output missing method or path: %s", logOutput)
	}
	if !strings.Contains(logOutput, "hello") && rr.Body.String() != "hello" {
		t.Errorf("response body mismatch")
	}
}

func TestLoggingMiddleware_CustomStatus(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	req := httptest.NewRequest("POST", "/teapot", nil)
	rr := httptest.NewRecorder()

	loggedHandler := middleware.Logging(handler)
	loggedHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTeapot {
		t.Errorf("expected status %d, got %d", http.StatusTeapot, rr.Code)
	}

	logOutput := buf.String()
	if !strings.Contains(logOutput, "418") {
		t.Errorf("log output does not contain status code 418: %s", logOutput)
	}
	if !strings.Contains(logOutput, "POST") || !strings.Contains(logOutput, "/teapot") {
		t.Errorf("log output missing method or path: %s", logOutput)
	}
}

func TestLoggingMiddleware_DefaultStatus(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// No WriteHeader called, should default to 200
	})

	req := httptest.NewRequest("GET", "/default", nil)
	rr := httptest.NewRecorder()

	loggedHandler := middleware.Logging(handler)
	loggedHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	logOutput := buf.String()
	if !strings.Contains(logOutput, "200") {
		t.Errorf("log output does not contain status code 200: %s", logOutput)
	}
	if !strings.Contains(logOutput, "/default") {
		t.Errorf("log output missing path: %s", logOutput)
	}
}

func TestLoggingMiddleware_LogsDuration(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
	})

	req := httptest.NewRequest("GET", "/duration", nil)
	rr := httptest.NewRecorder()

	loggedHandler := middleware.Logging(handler)
	loggedHandler.ServeHTTP(rr, req)

	logOutput := buf.String()
	if !strings.Contains(logOutput, "/duration") {
		t.Errorf("log output missing path: %s", logOutput)
	}
	if !strings.Contains(logOutput, "ms") {
		t.Errorf("log output missing duration: %s", logOutput)
	}
}
