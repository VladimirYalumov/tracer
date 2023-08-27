package tracer

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Span interface {
	Tag(key string, value interface{})
	HasSpanId() bool
	SpanId() string
	HasTraceId() bool
	TraceId() string
	End()
}

type span struct {
	s trace.Span
}

func (s *span) Tag(key string, value interface{}) {
	switch v := value.(type) {
	case int:
		s.s.SetAttributes(attribute.Int(key, v))
	case string:
		s.s.SetAttributes(attribute.String(key, v))
	case float64:
		s.s.SetAttributes(attribute.Float64(key, v))
	case int64:
		s.s.SetAttributes(attribute.Int64(key, v))
	case bool:
		s.s.SetAttributes(attribute.Bool(key, v))
	case []string:
		s.s.SetAttributes(attribute.StringSlice(key, v))
	case []int:
		s.s.SetAttributes(attribute.IntSlice(key, v))
	}
}

func (s *span) HasSpanId() bool {
	return s.s.SpanContext().HasSpanID()
}

func (s *span) SpanId() string {
	return s.s.SpanContext().SpanID().String()
}

func (s *span) HasTraceId() bool {
	return s.s.SpanContext().HasTraceID()
}

func (s *span) TraceId() string {
	return s.s.SpanContext().TraceID().String()
}

func (s *span) End() {
	s.s.End()
}

func StartSpan(ctx context.Context, name string) (context.Context, Span) {
	span := new(span)

	if tracer == nil {
		ctx, span.s = trace.NewNoopTracerProvider().Tracer("noop").Start(ctx, name)
	} else {
		ctx, span.s = tracer.Start(ctx, name)
	}

	return ctx, span
}

func SpanFromContext(ctx context.Context) Span {
	span := new(span)
	span.s = trace.SpanFromContext(ctx)

	return span
}
