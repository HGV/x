package oidcx

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateAuthHeader(t *testing.T) {
	tests := []struct {
		authHeader    string
		scheme        string
		expectedToken string
		expectedOK    bool
	}{
		{
			authHeader: "", scheme: "bearer",
			expectedToken: "", expectedOK: false,
		},
		{
			authHeader: "bearer token", scheme: "bearer ",
			expectedToken: "token", expectedOK: true,
		},
		{
			authHeader: "BEARER token", scheme: "bearer ",
			expectedToken: "token", expectedOK: true,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			token, ok := validateAuthHeader(tt.authHeader, tt.scheme)
			assert.Equal(t, tt.expectedOK, ok)
			assert.Equal(t, tt.expectedToken, token)
		})
	}
}

func TestHandler(t *testing.T) {
	issuer := "https://api.accounts.hgv.it"

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idToken, ok := IDTokenFromContext(r.Context())
		assert.True(t, ok)
		assert.NotNil(t, idToken)
		w.WriteHeader(http.StatusTeapot)
	})

	makeRequest := func(h http.Handler, token string) *httptest.ResponseRecorder {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0LXVzZXIiLCJhdWQiOlsidGVzdC1jbGllbnQiXSwiaXNzIjoiaHR0cHM6Ly9hcGkuYWNjb3VudHMuaGd2Lml0IiwiaWF0IjoxNjAwMDAwMDAwLCJleHAiOjIwMDAwMDAwMDB9.hJREizNgcJpnEEyZ5lE5VC9tPY45JIFJoxm9ZlIPgTI")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		return w
	}

	t.Run("unauthorized with no token", func(t *testing.T) {
		h := NewMiddleware(context.Background(), issuer).Handler(next)
		w := makeRequest(h, "")
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("custom error handler returns 403", func(t *testing.T) {
		h := NewMiddleware(context.Background(), issuer,
			WithAuthFailedHandler(func(err error) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusForbidden)
				}
			}),
		).Handler(next)
		w := makeRequest(h, "")
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("email config required but missing", func(t *testing.T) {
		h := NewMiddleware(context.Background(), issuer,
			WithAuthFailedHandler(func(err error) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte(err.Error()))
				}
			}),
			WithSkipClientIDCheck(),
			withInsecureSkipSignatureCheck(),
		).Handler(next)
		w := makeRequest(h, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0LXVzZXIiLCJhdWQiOlsidGVzdC1jbGllbnQiXSwiaXNzIjoiaHR0cHM6Ly9hcGkuYWNjb3VudHMuaGd2Lml0IiwiaWF0IjoxNjAwMDAwMDAwLCJleHAiOjIwMDAwMDAwMDB9.hJREizNgcJpnEEyZ5lE5VC9tPY45JIFJoxm9ZlIPgTI")
		b, _ := io.ReadAll(w.Body)
		assert.Equal(t, "invalid configuration, Email must be provided or SkipEmailCheck must be set", string(b))
	})

	t.Run("email mismatch", func(t *testing.T) {
		h := NewMiddleware(context.Background(), issuer,
			WithAuthFailedHandler(func(err error) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte(err.Error()))
				}
			}),
			WithSkipClientIDCheck(),
			WithEmail("test@hgv.it"),
			withInsecureSkipSignatureCheck(),
		).Handler(next)
		w := makeRequest(h, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0LXVzZXIiLCJhdWQiOlsidGVzdC1jbGllbnQiXSwiaXNzIjoiaHR0cHM6Ly9hcGkuYWNjb3VudHMuaGd2Lml0IiwiaWF0IjoxNjAwMDAwMDAwLCJleHAiOjIwMDAwMDAwMDB9.hJREizNgcJpnEEyZ5lE5VC9tPY45JIFJoxm9ZlIPgTI")
		b, _ := io.ReadAll(w.Body)
		assert.Equal(t, "expected email \"test@hgv.it\" got \"\"", string(b))
	})

	t.Run("valid expired token without email check", func(t *testing.T) {
		h := NewMiddleware(context.Background(), issuer,
			WithSkipClientIDCheck(),
			WithSkipEmailCheck(),
			withInsecureSkipSignatureCheck(),
		).Handler(next)
		w := makeRequest(h, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0LXVzZXIiLCJhdWQiOlsidGVzdC1jbGllbnQiXSwiaXNzIjoiaHR0cHM6Ly9hcGkuYWNjb3VudHMuaGd2Lml0IiwiaWF0IjoxNjAwMDAwMDAwLCJleHAiOjIwMDAwMDAwMDB9.hJREizNgcJpnEEyZ5lE5VC9tPY45JIFJoxm9ZlIPgTI")
		assert.Equal(t, http.StatusTeapot, w.Code)
	})
}
