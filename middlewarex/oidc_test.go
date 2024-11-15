package middlewarex

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/stretchr/testify/assert"
)

func TestValidateAuthHeader(t *testing.T) {
	tests := []struct {
		authHeader    string
		scheme        string
		expectedToken string
		expectedOk    bool
	}{
		{
			authHeader: "", scheme: "bearer",
			expectedToken: "", expectedOk: false,
		},
		{
			authHeader: "bearer token", scheme: "bearer ",
			expectedToken: "token", expectedOk: true,
		},
		{
			authHeader: "BEARER token", scheme: "bearer ",
			expectedToken: "token", expectedOk: true,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			token, ok := validateAuthHeader(tt.authHeader, tt.scheme)
			assert.Equal(t, tt.expectedOk, ok)
			assert.Equal(t, tt.expectedToken, token)
		})
	}
}

func TestHandler(t *testing.T) {
	issuer := "https://api.accounts.hgv.it"
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.NotNil(t, IDTokenFromContext(r.Context()))
		w.WriteHeader(http.StatusTeapot)
	})

	t.Run("unauthorized", func(t *testing.T) {
		h := OIDC(context.Background(), issuer)(next)
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("overwrite default error handler", func(t *testing.T) {
		h := OIDC(context.Background(), issuer, WithAuthFailedHandler(func(err error) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
			}
		}))(next)
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("valid expired token", func(t *testing.T) {
		h := OIDC(context.Background(), issuer, WithOIDCConfig(oidc.Config{
			SkipClientIDCheck: true,
			SkipExpiryCheck:   true,
		}))(next)
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Add("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6IjJlNTc0NjE3LTJlYzYtNGNhNy1hYTE2LThiYTYyMWRlMGI3YSIsInR5cCI6IkpXVCJ9.eyJhdWQiOltdLCJjbGllbnRfaWQiOiJjMDI2NTZiZC00NzZkLTQ1MGYtOWMwZC0zN2ZiMDhiYTI3MjEiLCJleHAiOjE3MzE2Njk3NDMsImV4dCI6e30sImlhdCI6MTczMTY2NjE0MywiaXNzIjoiaHR0cHM6Ly9hcGkuYWNjb3VudHMuaGd2Lml0IiwianRpIjoiZjk3YWE1ODAtZjZmNC00ZGQ3LTlkMDgtMjM1YTM5ZGU4ZWZlIiwibmJmIjoxNzMxNjY2MTQzLCJzY3AiOltdLCJzdWIiOiJjMDI2NTZiZC00NzZkLTQ1MGYtOWMwZC0zN2ZiMDhiYTI3MjEifQ.IeIc2EWCYjH8EaYClYpaTpYz-DDRbpu4vRuzirmBXZy28r7OazSrJdRSEa2a_G9Yq0UzmJXeBtPAouvsQdwmHX1PdBFzwwqLPT4kXcxMmlX6RvnTy-95wVfXnJJP-cGU5U4sMKKFGnsecAQotesEsYk19Dxylr5RMA-DsgwwpN8GQuf4KdLJk4IDJx8Z-FlfAG4XMODGM2S3sqGCwc6b5nQUXa_cUTIMqJCyUdb3Kd3OcQHKEK0o0esG1CBgqj3RrRE98BejeEjR5LOYiQpY1aAklmxa_3UOtEi9Bej1PRyybRxV7QbNE8_K0WVdj3CCedbtpK7DB0mNGCtas2bjiFxsr9MBHUtDcU3taXEoEkSqye7vIbLgd66SFm5gq78-PeJEvbwYqpt4LB7b7F-ZpyhCU-3T3SNkMPHY-q7hIBPauRbJbtWdK3w_xjjjCJdgjspk-CEyOUfhogjKmavxcuuXOGBphOeJ7WCRMTlmv9ira0DZqwBCQTGitkGGT98l4guaIYoB27Zsl-wdgxK2F0AwjvHFTYNUsG3Nf9NJ4ULjPMusBBA9hHBoO1UrlNWgXEpJWvr5YV_vt0Omlqvv-ci7M3Rx1-MjRyBYTQRxVRLhtDtGK4TbW4jCEIE38_k5IDqH6WxaUsgxTxFu8rx5xWhpRlKuIQRrDyWA1ylMo_U")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		assert.Equal(t, http.StatusTeapot, w.Code)
	})
}
