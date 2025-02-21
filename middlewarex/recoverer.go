package middlewarex

import (
	"errors"
	"net/http"
)

type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

func Recoverer(errorHandler ErrorHandler) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					if rvr == http.ErrAbortHandler {
						// Don’t recover http.ErrAbortHandler so the response to
						// the client is aborted.
						panic(rvr)
					}

					var err error
					switch x := rvr.(type) {
					case string:
						err = errors.New(x)
					case error:
						err = x
					default:
						err = errors.New("unknown panic")
					}

					errorHandler(w, r, err)
				}
			}()

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
