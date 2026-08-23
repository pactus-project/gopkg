package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggingMiddleware(t *testing.T) {
	var logBuffer bytes.Buffer
	originalLogger := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&logBuffer, nil)))
	defer slog.SetDefault(originalLogger)

	middleware := Logging()

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://test.com/foo", http.NoBody)
	req.RemoteAddr = "127.0.0.1"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	logged := logBuffer.String()
	assert.Contains(t, logged, "method=GET")
	assert.Contains(t, logged, "path=/foo")
	assert.Contains(t, logged, "remote=127.0.0.1")
	assert.Contains(t, logged, "duration_ms=")
}
