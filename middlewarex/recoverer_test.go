package middlewarex

import (
	"encoding/json"
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

func TestRecoverer_WithJSONError(t *testing.T) {
	type ErrorResponse struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}

	jsonErrorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	h := Recoverer(jsonErrorHandler)(next)
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	res := w.Result()
	var responseBody ErrorResponse
	_ = json.NewDecoder(res.Body).Decode(&responseBody)

	assert.Equal(t, res.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, http.StatusInternalServerError, responseBody.Status)
	assert.Equal(t, "test panic", responseBody.Message)
}
