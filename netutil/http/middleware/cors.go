// Package middleware provides HTTP middleware components including logging,
// CORS handling, panic recovery, and middleware chaining.
package middleware

import (
	"errors"
	"net/http"
	"slices"
	"strings"
)

// CORSConfig holds the configuration for Cross-Origin Resource Sharing.
type CORSConfig struct {
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
}

// DefaultCORSConfig returns a default CORS configuration.
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	}
}

// BasicCheck performs basic validation on the configuration.
func (c *CORSConfig) BasicCheck() error {
	if slices.Contains(c.AllowedOrigins, "*") && c.AllowCredentials {
		return errors.New("wildcard origin '*' cannot be used with allow_credentials=true")
	}

	if slices.Contains(c.AllowedMethods, "*") && c.AllowCredentials {
		return errors.New("wildcard method '*' cannot be used with allow_credentials=true")
	}

	return nil
}

// CORS creates middleware to handle CORS requests.
func CORS(config *CORSConfig) Middleware {
	allowAllOrigins := slices.Contains(config.AllowedOrigins, "*")
	allowedOrigins := make(map[string]struct{}, len(config.AllowedOrigins))
	for _, origin := range config.AllowedOrigins {
		allowedOrigins[origin] = struct{}{}
	}

	methodsHeader := strings.Join(config.AllowedMethods, ", ")
	headersHeader := strings.Join(config.AllowedHeaders, ", ")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Vary", "Origin")

			origin := r.Header.Get("Origin")
			switch {
			case allowAllOrigins:
				w.Header().Set("Access-Control-Allow-Origin", "*")
			case origin != "":
				if _, ok := allowedOrigins[origin]; ok {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				}
			}

			w.Header().Set("Access-Control-Allow-Methods", methodsHeader)
			w.Header().Set("Access-Control-Allow-Headers", headersHeader)

			if config.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
