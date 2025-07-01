package middleware

import (
	"net/http"
	"strings"
)

type CORSConfig struct {
	AllowOrigin      string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials string
}

// Option defines a function that modifies CORSConfig.
type Option func(*CORSConfig)

// WithAllowOrigin sets the allowed origin.
func WithAllowOrigin(origin string) Option {
	return func(c *CORSConfig) {
		c.AllowOrigin = origin
	}
}

// WithAllowMethods sets the allowed methods.
func WithAllowMethods(methods ...string) Option {
	return func(c *CORSConfig) {
		c.AllowMethods = methods
	}
}

// WithAllowHeaders sets the allowed headers.
func WithAllowHeaders(headers ...string) Option {
	return func(c *CORSConfig) {
		c.AllowHeaders = headers
	}
}

// WithAllowCredentials sets the allow credentials flag.
func WithAllowCredentials(credentials string) Option {
	return func(c *CORSConfig) {
		c.AllowCredentials = credentials
	}
}

// CORS returns a CORS middleware with the given options.
func CORS(opts ...Option) func(http.Handler) http.Handler {
	// Set default values
	cfg := &CORSConfig{
		AllowOrigin:      "*",
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: "true",
	}
	for _, opt := range opts {
		opt(cfg)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", cfg.AllowOrigin)
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.AllowMethods, ","))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.AllowHeaders, ","))
			w.Header().Set("Access-Control-Allow-Credentials", cfg.AllowCredentials)

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
