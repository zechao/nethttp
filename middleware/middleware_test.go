package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zechao158/nethhtp/middleware"
)

func TestCreateStack_Order(t *testing.T) {
	var callOrder []string

	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callOrder = append(callOrder, "mw1")
			next.ServeHTTP(w, r)
		})
	}
	mw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callOrder = append(callOrder, "mw2")
			next.ServeHTTP(w, r)
		})
	}
	mw3 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callOrder = append(callOrder, "mw3")
			next.ServeHTTP(w, r)
		})
	}

	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callOrder = append(callOrder, "final")
	})

	stack := middleware.CreateStack(mw1, mw2, mw3)
	handler := stack(finalHandler)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	expectedOrder := []string{"mw1", "mw2", "mw3", "final"}
	for i, v := range expectedOrder {
		if callOrder[i] != v {
			t.Errorf("expected callOrder[%d] = %q, got %q", i, v, callOrder[i])
		}
	}
}

func TestCreateStack_Empty(t *testing.T) {
	finalCalled := false
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		finalCalled = true
	})

	stack := middleware.CreateStack()
	handler := stack(finalHandler)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if !finalCalled {
		t.Error("final handler was not called when middleware stack is empty")
	}
}
