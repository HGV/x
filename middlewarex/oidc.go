package middlewarex

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
)

type (
	OIDCMiddlewareOption func(*oidcMiddleware)

	oidcMiddleware struct {
		authFailedHandler func(error) http.HandlerFunc
		config            *oidc.Config
		verifier          *oidc.IDTokenVerifier
	}
	oidcContextKey int
)

const (
	idTokenContextKey oidcContextKey = iota
)

func OIDC(ctx context.Context, issuer string, opts ...OIDCMiddlewareOption) func(next http.Handler) http.Handler {
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		panic(err)
	}

	mw := oidcMiddleware{
		authFailedHandler: oidcAuthFailed,
	}

	for _, opt := range opts {
		opt(&mw)
	}

	if mw.config == nil {
		mw.verifier = provider.Verifier(&oidc.Config{})
	} else {
		mw.verifier = provider.Verifier(mw.config)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			bearerToken, ok := validateAuthHeader(authHeader, "Bearer ")
			if !ok {
				mw.authFailedHandler(errors.New("bearer token is missing or invalid")).ServeHTTP(w, r)
				return
			}

			idToken, err := mw.verifier.Verify(r.Context(), bearerToken)
			if err != nil {
				mw.authFailedHandler(err).ServeHTTP(w, r)
				return
			}

			ctx := contextWithIDToken(r.Context(), idToken)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func contextWithIDToken(ctx context.Context, idToken *oidc.IDToken) context.Context {
	return context.WithValue(ctx, idTokenContextKey, idToken)
}

func IDTokenFromContext(ctx context.Context) *oidc.IDToken {
	if idToken, ok := ctx.Value(idTokenContextKey).(*oidc.IDToken); ok {
		return idToken
	}
	return nil
}

func validateAuthHeader(s, scheme string) (string, bool) {
	if len(s) >= len(scheme) && strings.EqualFold(s[0:len(scheme)], scheme) {
		return s[len(scheme):], true
	}
	return s, false
}

func oidcAuthFailed(err error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func WithAuthFailedHandler(h func(error) http.HandlerFunc) OIDCMiddlewareOption {
	return func(opt *oidcMiddleware) {
		if h != nil {
			opt.authFailedHandler = h
		}
	}
}

func WithOIDCConfig(c oidc.Config) OIDCMiddlewareOption {
	return func(opt *oidcMiddleware) {
		opt.config = &c
	}
}
