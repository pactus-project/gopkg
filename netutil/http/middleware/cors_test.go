package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCORSConfigBasicCheck(t *testing.T) {
	origin := "https://example.com"

	tests := []struct {
		name    string
		cfg     *CORSConfig
		wantErr bool
	}{
		{
			name: "default config",
			cfg:  DefaultCORSConfig(),
		},
		{
			name: "explicit origins with credentials",
			cfg:  corsConfig([]string{origin}, []string{"GET"}, true),
		},
		{
			name: "wildcard methods",
			cfg:  corsConfig([]string{origin}, []string{"*"}, false),
		},
		{
			name:    "wildcard methods with credentials",
			cfg:     corsConfig([]string{origin}, []string{"*"}, true),
			wantErr: true,
		},
		{
			name: "arbitrary method allowed",
			cfg:  corsConfig([]string{origin}, []string{"get"}, false),
		},
		{
			name:    "wildcard origin with credentials",
			cfg:     corsConfig([]string{"*"}, nil, true),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.BasicCheck()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func corsConfig(origins, methods []string, credentials bool) *CORSConfig {
	return &CORSConfig{
		AllowedOrigins:   origins,
		AllowedMethods:   methods,
		AllowCredentials: credentials,
	}
}

func TestCORSMiddleware(t *testing.T) {
	config := DefaultCORSConfig()
	config.AllowedOrigins = []string{"https://example.com"}
	config.AllowCredentials = true

	middleware := CORS(config)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://test.com", http.NoBody)
	req.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "https://example.com", res.Header.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "*", res.Header.Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "*", res.Header.Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "true", res.Header.Get("Access-Control-Allow-Credentials"))
	assert.Contains(t, res.Header.Values("Vary"), "Origin")
}

func TestCORSMiddlewareOptionsRequest(t *testing.T) {
	config := DefaultCORSConfig()
	middleware := CORS(config)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodOptions, "http://test.com", http.NoBody)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusNoContent, res.StatusCode)
	assert.Equal(t, "*", res.Header.Get("Access-Control-Allow-Origin"))
}

func TestCORSMiddlewareWildcardOrigin(t *testing.T) {
	config := DefaultCORSConfig()
	middleware := CORS(config)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://test.com", http.NoBody)
	req.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "*", res.Header.Get("Access-Control-Allow-Origin"))
	assert.Empty(t, res.Header.Get("Access-Control-Allow-Credentials"))
}

func TestCORSMiddlewareDisallowedOrigin(t *testing.T) {
	config := DefaultCORSConfig()
	config.AllowedOrigins = []string{"https://example.com"}

	middleware := CORS(config)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://test.com", http.NoBody)
	req.Header.Set("Origin", "https://evil.com")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	res := w.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Empty(t, res.Header.Get("Access-Control-Allow-Origin"))
}
