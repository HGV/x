package middlewarex

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecurityHeaders(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		keys := []string{
			"Content-Security-Policy",
			"Strict-Transport-Security",
			"X-Content-Type-Options",
			"Referrer-Policy",
		}
		for _, key := range keys {
			if w.Header().Get(key) == "" {
				t.Errorf("%s not present", key)
			}
		}
	})
	h := SecurityHeaders(next)
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
}
