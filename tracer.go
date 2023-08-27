package tracer

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	provider *tracesdk.TracerProvider
	tracer   trace.Tracer
)

func Init(ctx context.Context, name string, oo ...Option) error {
	if tracer != nil {
		return errors.New("tracer already initialized")
	}

	opts := BuildOptions(oo)

	if opts.IsNoop() {
		tracer = trace.NewNoopTracerProvider().Tracer(name)
		return nil
	}

	conn, err := grpc.DialContext(ctx, opts.GetTarget(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("grpc.DialContext: %w", err)
	}

	exp, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return fmt.Errorf("otlptracegrpc.New: %w", err)
	}

	provider = tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(name),
		)),
	)

	tracer = provider.Tracer(name)
	return nil
}

func Shutdown(ctx context.Context) error {
	return provider.Shutdown(ctx)
}

func TraceToContext(ctx context.Context, headerTraceId string) (context.Context, error) {
	traceId, err := trace.TraceIDFromHex(headerTraceId)
	if err != nil {
		return ctx, fmt.Errorf("trace.TraceIDFromHex: %w", err)
	}
	spanContext := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceId,
	})
	return trace.ContextWithSpanContext(ctx, spanContext), nil
}
