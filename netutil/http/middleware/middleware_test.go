package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChain(t *testing.T) {
	var order []string

	makeMiddleware := func(name string) Middleware {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, name+":before")
				next.ServeHTTP(w, r)
				order = append(order, name+":after")
			})
		}
	}

	handler := Chain(
		makeMiddleware("first"),
		makeMiddleware("second"),
	)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		order = append(order, "handler")
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://test.com", http.NoBody)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, []string{"first:before", "second:before", "handler", "second:after", "first:after"}, order)
}
