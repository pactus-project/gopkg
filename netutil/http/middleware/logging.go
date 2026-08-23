package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// Logging logs incoming HTTP requests and their duration.
func Logging() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(start)

			slog.Info(
				"http request",
				"method", r.Method,
				"path", r.URL.Path,
				"remote", r.RemoteAddr,
				"duration_ms", duration.Milliseconds(),
			)
		})
	}
}
