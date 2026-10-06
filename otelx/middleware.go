package otelx

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
	"go.opentelemetry.io/otel/trace"
)

func TraceHandler(h http.Handler, opts ...otelhttp.Option) http.Handler {
	middlewareOpts := []otelhttp.Option{
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			return r.URL.Path
		}),
	}
	return otelhttp.NewHandler(h, "", append(middlewareOpts, opts...)...)
}

// RouteTagger names the request span after the route template and adds
// http.route to the span and to otelhttp's metrics, for routers that do not
// expose their pattern to otelhttp. It must run inside TraceHandler.
//
// route is called after next has served the request, once the router has
// finished matching. With chi:
//
//	r.Use(otelx.RouteTagger(func(r *http.Request) string {
//		return chi.RouteContext(r.Context()).RoutePattern()
//	}))
func RouteTagger(route func(*http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			ctx := r.Context()
			pattern := route(r)
			if pattern == "" {
				return
			}
			attr := semconv.HTTPRoute(pattern)
			// Metrics are recorded for unsampled requests too, so the labeler is
			// updated regardless of whether the span is recording.
			if labeler, ok := otelhttp.LabelerFromContext(ctx); ok {
				labeler.Add(attr)
			}
			if span := trace.SpanFromContext(ctx); span.IsRecording() {
				span.SetName(r.Method + " " + pattern)
				span.SetAttributes(attr)
			}
		})
	}
}
