package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var globalTracer trace.Tracer

// InitGlobalTracer creates global tracer object for tracing operations.
func InitGlobalTracer(serviceName string) error {
	globalTracer = otel.Tracer(serviceName)
	return nil
}

// Start creates a span and a context containing the newly-created span.
func Start(ctx context.Context, method string) (context.Context, trace.Span) {
	return globalTracer.Start(ctx, method)
}
