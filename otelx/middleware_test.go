package otelx

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
)

func TestRouteTagger(t *testing.T) {
	spans := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spans))

	var labeler *otelhttp.Labeler
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		labeler, _ = otelhttp.LabelerFromContext(r.Context())
	})
	route := func(r *http.Request) string { return "/items/{id}" }

	h := TraceHandler(RouteTagger(route)(next), otelhttp.WithTracerProvider(tp))
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/items/1", nil))

	ended := spans.Ended()
	require.Len(t, ended, 1)
	assert.Equal(t, "GET /items/{id}", ended[0].Name())
	assert.Contains(t, ended[0].Attributes(), semconv.HTTPRoute("/items/{id}"))
	require.NotNil(t, labeler)
	assert.Contains(t, labeler.Get(), semconv.HTTPRoute("/items/{id}"))
}

func TestRouteTaggerNoRoute(t *testing.T) {
	spans := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spans))

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	route := func(r *http.Request) string { return "" }

	h := TraceHandler(RouteTagger(route)(next), otelhttp.WithTracerProvider(tp))
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/unknown", nil))

	ended := spans.Ended()
	require.Len(t, ended, 1)
	assert.Equal(t, "/unknown", ended[0].Name())
}
