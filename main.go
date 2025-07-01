package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/zechao158/nethhtp/middleware"
)

type Mux struct {
	http.ServeMux
	mws []middleware.Middleware
}

func NewMux() *Mux {
	return &Mux{
		ServeMux: *http.NewServeMux(),
	}
}

func (m *Mux) Use(mws ...middleware.Middleware) {
	m.mws = append(m.mws, mws...)
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler = &m.ServeMux
	mwStack := middleware.CreateStack(m.mws...)
	handler = mwStack(handler)
	handler.ServeHTTP(w, r)
}

func (m *Mux) Mount(prefix string, sub *Mux) {
	if prefix == "" || prefix[0] != '/' {
		panic("prefix must start with '/'")
	}
	prefix = strings.TrimSuffix(prefix, "/")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, prefix) {
			http.NotFound(w, r)
			return
		}

		// clone request
		r2 := r.Clone(r.Context())

		// strip prefix from Path
		r2.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
		if r2.URL.Path == "" {
			r2.URL.Path = "/" // avoid empty path
		}

		sub.ServeHTTP(w, r2)
	})

	m.Handle(prefix+"/", handler) // match prefix and all sub-paths
}

type Profile struct {
	Name string
	Age  int
	Addr string
}

var profiles = map[string]Profile{
	"alice": {Name: "Alice", Age: 30, Addr: "Wonderland"},
	"bob":   {Name: "Bob", Age: 25, Addr: "Builderland"},
}

func encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

// decode decodes JSON from the request body into the provided type T.
// It returns the decoded value and an error if decoding fails.
func decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

func main() {

	m := NewMux()

	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "all path")
	})

	m.HandleFunc("GET /profile/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		if profile, ok := profiles[name]; ok {
			err := encode(w, r, http.StatusOK, profile)
			if err != nil {
				http.Error(w, fmt.Sprintf("encode profile: %v", err), http.StatusInternalServerError)
				return
			}
			return
		}
		http.Error(w, "Profile not found", http.StatusNotFound)

	})

	m.HandleFunc("GET /hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		fmt.Fprintln(w, "Hello world")
		fmt.Fprintf(w, "Hello %s", name)
	})
	m.HandleFunc("GET /redirect/{domain}", func(w http.ResponseWriter, r *http.Request) {
		domain := r.PathValue("domain")
		http.Redirect(w, r, fmt.Sprintf("https://%s.com", domain), http.StatusMovedPermanently)
	})

	m.Use(middleware.Logging, middleware.CORS(), middleware.TraceMiddleware)

	m.Mount("/v1", m)
	server := http.Server{
		Addr:    ":8080",
		Handler: m,
	}
	fmt.Println("Server serving at port", server.Addr)
	// Graceful shutdown
	stop := make(chan struct{})
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("ListenAndServe error: %v\n", err)
		}
		close(stop)
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	}
	<-stop
	fmt.Println("Server exiting")
}
