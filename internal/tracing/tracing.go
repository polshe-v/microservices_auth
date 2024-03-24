package tracing

import (
	"github.com/uber/jaeger-client-go/config"
	/*	"context"
		"fmt"

		"go.opentelemetry.io/otel"
		"go.opentelemetry.io/otel/exporters/jaeger"
		"go.opentelemetry.io/otel/sdk/resource"
		tracesdk "go.opentelemetry.io/otel/sdk/trace"
		semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
		"go.opentelemetry.io/otel/trace"
	*/)

//var globalTracer trace.Tracer

const serviceName = "auth_service"

// Init creates global tracer object for tracing operations.
func Init() error {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			//LogSpans:           true,
			LocalAgentHostPort: "jaeger:6831",
		},
	}

	_, err := cfg.InitGlobalTracer(serviceName)
	if err != nil {
		return err
	}
	return nil
}

// Init creates global tracer object for tracing operations.
/*func Init() error {
	globalTracer = otel.Tracer(serviceName)
	return nil
}

// Start creates a span and a context.Context containing the newly-created span.
func Start(ctx context.Context, method string) (context.Context, trace.Span) {
	return globalTracer.Start(ctx, method)
}

// InitTracer creates global tracer object for tracing operations.
func InitTracer(jaegerURL string, serviceName string) error {
	exporter, err := newJaegerExporter(jaegerURL)
	if err != nil {
		return fmt.Errorf("initialize exporter: %w", err)
	}

	tp, err := newTraceProvider(exporter, serviceName)
	if err != nil {
		return fmt.Errorf("initialize provider: %w", err)
	}

	otel.SetTracerProvider(tp)

	globalTracer = tp.Tracer("main tracer")
	return nil
}

func newTraceProvider(exp tracesdk.SpanExporter, ServiceName string) (*tracesdk.TracerProvider, error) {
	// Ensure default SDK resources and the required service name are set.
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(ServiceName),
		),
	)
	if err != nil {
		return nil, err
	}

	return tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(r),
	), nil
}

// NewJaegerExporter creates new jaeger exporter.
//
//	URL example - http://localhost:14268/api/traces
func newJaegerExporter(url string) (tracesdk.SpanExporter, error) {
	return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
}
*/
