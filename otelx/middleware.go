package otelx

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func TraceHandler(h http.Handler, opts ...otelhttp.Option) http.Handler {
	middlewareOpts := []otelhttp.Option{
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			return r.URL.Path
		}),
	}
	return otelhttp.NewHandler(h, "", append(middlewareOpts, opts...)...)
}
