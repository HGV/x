package oryx

import (
	"context"
	"errors"
	"net/http"

	ory "github.com/ory/client-go"
)

type Middleware struct {
	o *middlewareOptions
	c *ory.APIClient
}

type middlewareOptions struct {
	AuthFailedHandler func(error) http.HandlerFunc
}

type MiddlewareOption func(*middlewareOptions)

type sessionContextKey struct{}

func NewMiddleware(url string, opts ...MiddlewareOption) *Middleware {
	o := &middlewareOptions{
		AuthFailedHandler: defaultAuthFailedHandler,
	}

	for _, opt := range opts {
		opt(o)
	}

	config := ory.NewConfiguration()
	config.Servers = ory.ServerConfigurations{
		{
			URL: url,
		},
	}

	return &Middleware{
		o: o,
		c: ory.NewAPIClient(config),
	}
}

func (mw *Middleware) Session(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := mw.validateSession(r)
		if err != nil {
			mw.o.AuthFailedHandler(err).ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), sessionContextKey{}, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (mw *Middleware) validateSession(r *http.Request) (*ory.Session, error) {
	cookies := r.Header.Get("Cookie")
	session, _, err := mw.c.FrontendAPI.ToSession(r.Context()).
		Cookie(cookies).
		Execute()
	if err != nil {
		return nil, err
	}
	if session == nil || session.Active == nil || !*session.Active {
		return nil, errors.New("no active session found")
	}
	return session, nil
}

func SessionFromContext(ctx context.Context) (*ory.Session, bool) {
	sctx, ok := ctx.Value(sessionContextKey{}).(*ory.Session)
	return sctx, ok
}

func WithAuthFailedHandler(h func(error) http.HandlerFunc) MiddlewareOption {
	return func(o *middlewareOptions) {
		if h != nil {
			o.AuthFailedHandler = h
		}
	}
}

func defaultAuthFailedHandler(err error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}
}
