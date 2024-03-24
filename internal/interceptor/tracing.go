package interceptor

import (
	"context"

	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/polshe-v/microservices_auth/internal/tracing"
)

const traceIDKey = "x-trace-id"

// TracingInterceptor creates traces for fucntion calls.
func TracingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ctx, span := tracing.Start(ctx, info.FullMethod)
	defer span.End()

	traceID := span.SpanContext().TraceID().String()
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(traceIDKey, traceID))

	header := metadata.New(map[string]string{traceIDKey: traceID})
	err := grpc.SendHeader(ctx, header)
	if err != nil {
		return nil, err
	}

	res, err := handler(ctx, req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}

	return res, err
}
