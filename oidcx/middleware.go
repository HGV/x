package oidcx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
)

type Middleware struct {
	o *middlewareOptions
	v *oidc.IDTokenVerifier
}

type middlewareOptions struct {
	ClientID                   string
	SkipClientIDCheck          bool
	Email                      string
	SkipEmailCheck             bool
	InsecureSkipSignatureCheck bool
	AuthFailedHandler          func(error) http.HandlerFunc
}

type MiddlewareOption func(*middlewareOptions)

type idTokenContextKey struct{}

func NewMiddleware(ctx context.Context, issuer string, opts ...MiddlewareOption) *Middleware {
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		panic(err)
	}

	o := &middlewareOptions{
		AuthFailedHandler: defaultAuthFailedHandler,
	}

	for _, opt := range opts {
		opt(o)
	}

	mw := Middleware{
		o: o,
		v: provider.VerifierContext(ctx, &oidc.Config{
			ClientID:                   o.ClientID,
			SkipClientIDCheck:          o.SkipClientIDCheck,
			InsecureSkipSignatureCheck: o.InsecureSkipSignatureCheck,
		}),
	}

	return &mw
}

func (mw *Middleware) Handler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		bearerToken, ok := validateAuthHeader(r.Header.Get("Authorization"), "Bearer ")
		if !ok {
			mw.o.AuthFailedHandler(errors.New("bearer token is missing or invalid")).ServeHTTP(w, r)
			return
		}

		idToken, err := mw.v.Verify(ctx, bearerToken)
		if err != nil {
			mw.o.AuthFailedHandler(err).ServeHTTP(w, r)
			return
		}

		if !mw.o.SkipEmailCheck {
			if mw.o.Email == "" {
				mw.o.AuthFailedHandler(errors.New("invalid configuration, Email must be provided or SkipEmailCheck must be set")).ServeHTTP(w, r)
				return
			}

			var claims struct {
				Email string `json:"email"`
			}
			if err = idToken.Claims(&claims); err != nil {
				mw.o.AuthFailedHandler(err).ServeHTTP(w, r)
				return
			}
			if !strings.EqualFold(mw.o.Email, claims.Email) {
				mw.o.AuthFailedHandler(fmt.Errorf("expected email %q got %q", mw.o.Email, claims.Email)).ServeHTTP(w, r)
				return
			}
		}

		ctx = context.WithValue(ctx, idTokenContextKey{}, idToken)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func IDTokenFromContext(ctx context.Context) (*oidc.IDToken, bool) {
	octx, ok := ctx.Value(idTokenContextKey{}).(*oidc.IDToken)
	return octx, ok
}

func WithAuthFailedHandler(h func(error) http.HandlerFunc) MiddlewareOption {
	return func(o *middlewareOptions) {
		if h != nil {
			o.AuthFailedHandler = h
		}
	}
}

func WithClientID(clientID string) MiddlewareOption {
	return func(o *middlewareOptions) {
		o.ClientID = clientID
	}
}

func WithSkipClientIDCheck() MiddlewareOption {
	return func(o *middlewareOptions) {
		o.SkipClientIDCheck = true
	}
}

func WithEmail(email string) MiddlewareOption {
	return func(o *middlewareOptions) {
		o.Email = email
	}
}

func WithSkipEmailCheck() MiddlewareOption {
	return func(o *middlewareOptions) {
		o.SkipEmailCheck = true
	}
}

func withInsecureSkipSignatureCheck() MiddlewareOption {
	return func(o *middlewareOptions) {
		o.InsecureSkipSignatureCheck = true
	}
}

func validateAuthHeader(s, scheme string) (string, bool) {
	if len(s) >= len(scheme) && strings.EqualFold(s[0:len(scheme)], scheme) {
		return s[len(scheme):], true
	}
	return s, false
}

func defaultAuthFailedHandler(err error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}
}
