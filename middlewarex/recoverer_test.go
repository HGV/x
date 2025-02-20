package middlewarex

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecoverer_WithPlainTextError(t *testing.T) {
	plainTextErrorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal server error: " + err.Error()))
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	h := Recoverer(plainTextErrorHandler)(next)
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	res := w.Result()
	assert.Equal(t, res.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, w.Body.String(), "internal server error: test panic")
}
