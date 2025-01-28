package otel

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
)

// newTracer initializes a new TracerProvider.
func newTracer(r *resource.Resource) *trace.TracerProvider {
	return trace.NewTracerProvider(
		trace.WithResource(r),
		trace.WithSampler(trace.AlwaysSample()),
	)
}

// WithOTLPTracer registers a otlp tracer, sending traces to a gRPC endpoint.
func WithOTLPTracer(ctx context.Context, client *grpc.ClientConn) func(*Otel) error {
	return func(observer *Otel) error {
		if client == nil {
			return nil // allow to not configure any gRPC exporter
		}

		traceGRPC, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(client))
		if err != nil {
			return err
		}
		observer.RegisterSpanProcessor(trace.NewBatchSpanProcessor(traceGRPC, trace.WithBatchTimeout(time.Second)))
		return nil
	}
}

// WithStdoutTracer registers a stdout tracer.
func WithStdoutTracer() func(*Otel) error {
	return func(observer *Otel) error {
		traceStdout, err := stdouttrace.New(
			stdouttrace.WithPrettyPrint())
		if err != nil {
			return err
		}

		observer.TracerProvider.RegisterSpanProcessor(trace.NewBatchSpanProcessor(traceStdout, trace.WithBatchTimeout(time.Second)))
		return nil
	}
}
