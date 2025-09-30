package otelx

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"
)

type CloudLoggingHandler struct {
	handler   slog.Handler
	projectID string
}

// CloudLoggingHandler is a slog.Handler which adds attributes from the
// span context.
func NewCloudLoggingHandler(projectID string, opts *slog.HandlerOptions) *CloudLoggingHandler {
	if projectID == "" {
		panic("projectID is required")
	}

	return &CloudLoggingHandler{
		handler:   slog.NewJSONHandler(os.Stdout, opts),
		projectID: projectID,
	}
}

func (h *CloudLoggingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// Handle overrides slog.Handler's Handle method. This adds attributes from the
// span context to the slog.Record.
func (h *CloudLoggingHandler) Handle(ctx context.Context, record slog.Record) error {
	if s := trace.SpanContextFromContext(ctx); s.IsValid() {
		// Add trace context attributes following Cloud Logging structured log format described
		// in https://cloud.google.com/logging/docs/structured-logging#special-payload-fields
		record.AddAttrs(
			slog.String("logging.googleapis.com/trace", h.gcpTraceID(s.TraceID())),
			slog.Any("logging.googleapis.com/spanId", s.SpanID()),
			slog.Bool("logging.googleapis.com/trace_sampled", s.TraceFlags().IsSampled()),
		)
	}
	return h.handler.Handle(ctx, record)
}

func (h *CloudLoggingHandler) gcpTraceID(traceID trace.TraceID) string {
	return fmt.Sprintf("projects/%s/traces/%s", h.projectID, traceID.String())
}

func (h *CloudLoggingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CloudLoggingHandler{
		handler:   h.handler.WithAttrs(attrs),
		projectID: h.projectID,
	}
}

func (h *CloudLoggingHandler) WithGroup(name string) slog.Handler {
	return &CloudLoggingHandler{
		handler:   h.handler.WithGroup(name),
		projectID: h.projectID,
	}
}

func GCPReplacer(groups []string, a slog.Attr) slog.Attr {
	// Rename attribute keys to match Cloud Logging structured log format
	switch a.Key {
	case slog.LevelKey:
		a.Key = "severity"
		// Map slog.Level string values to Cloud Logging LogSeverity
		// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogSeverity
		if level := a.Value.Any().(slog.Level); level == slog.LevelWarn {
			a.Value = slog.StringValue("WARNING")
		}
	case slog.TimeKey:
		a.Key = "timestamp"
	case slog.MessageKey:
		a.Key = "message"
	case slog.SourceKey:
		a.Key = "logging.googleapis.com/sourceLocation"
	}
	return a
}
