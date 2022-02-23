package tracing

import (
	"context"
	"github.com/gogo/protobuf/proto"
	"github.com/wyy-go/wgo/nrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func Trace(service string, opts ...Option) nrpc.Middleware {
	cfg := config{}
	for _, opt := range opts {
		opt.apply(&cfg)
	}
	if cfg.TracerProvider == nil {
		cfg.TracerProvider = otel.GetTracerProvider()
	}
	tracer := cfg.TracerProvider.Tracer(
		"nrpc",
		trace.WithInstrumentationVersion(SemVersion()),
	)
	if cfg.Propagators == nil {
		cfg.Propagators = otel.GetTextMapPropagator()
	}

	return func(handler nrpc.Handler) nrpc.Handler {
		return func(ctx context.Context) (rsp proto.Message, err error) {
			opts := []trace.SpanStartOption{
				trace.WithSpanKind(trace.SpanKindServer),
			}
			tCtx, span := tracer.Start(ctx, service, opts...)
			defer span.End()
			// serve the request to the next middleware
			rsp, err = handler(tCtx)
			if err != nil {
				span.SetAttributes(attribute.String("nrpc.errors", err.Error()))
			}

			return
		}
	}
}

func FromTraceId(ctx context.Context) string {
	if span := trace.SpanContextFromContext(ctx); span.HasTraceID() {
		return span.SpanID().String()
	}
	return ""
}

func GetTraceId(ctx context.Context) string {
	return FromTraceId(ctx)
}
