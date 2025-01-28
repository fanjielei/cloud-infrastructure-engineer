// Ref: https://opentelemetry.io/docs/languages/go/getting-started/
package otel

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

// Observer orchestrates logs, metrics and traces.
type Observer interface {
	Logger
	metric.MeterProvider
	trace.TracerProvider
	propagation.TextMapPropagator
	Shutdown(context.Context) error
}

// Ensure Otel satisfies Observer.
var _ Observer = (*Otel)(nil)

// Otel is the main object to act as Observer.
type Otel struct {
	*logger
	*sdkmetric.MeterProvider
	*sdktrace.TracerProvider
	propagation.TextMapPropagator

	*resource.Resource
	prometheusServer *http.Server
	shutdownFuncs    []func(context.Context) error
}

// Service stores otel specific attributes and attaches them to otel resources.
type Service struct {
	Name string
}

// Options configures Otel options.
type Options func(*Otel) error

// New bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func New(ctx context.Context, s Service, options ...Options) (*Otel, error) {
	r, err := resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(s.Name)),
	)
	if err != nil {
		return nil, err
	}

	observer := &Otel{
		Resource:          r,
		TracerProvider:    newTracer(r),
		TextMapPropagator: newPropagator(),
		shutdownFuncs:     make([]func(context.Context) error, 0),
	}

	for _, fn := range options {
		if err := fn(observer); err != nil {
			return nil, err
		}
	}

	if observer.logger != nil {
		observer.shutdownFuncs = append(observer.shutdownFuncs, observer.logger.Shutdown)
	}
	if observer.TracerProvider != nil {
		observer.shutdownFuncs = append(observer.shutdownFuncs, observer.TracerProvider.Shutdown)
	}
	if observer.MeterProvider != nil {
		observer.shutdownFuncs = append(observer.shutdownFuncs, observer.MeterProvider.Shutdown)
	}

	go func() {
		if err := observer.PrometheusListenAndServe(ctx); err != nil {
			if errors.Is(err, io.EOF) {
				observer.logger.Error(ctx, fmt.Sprintf("serving metrics %s", err))
			}
		}
	}()
	return observer, nil
}

// PrometheusListenAndServe starts a webserver serving prometheus metrics.
func (o *Otel) PrometheusListenAndServe(ctx context.Context) error {
	if o.prometheusServer == nil {
		return nil
	}
	o.shutdownFuncs = append(o.shutdownFuncs, o.prometheusServer.Shutdown)
	o.logger.Info(ctx, fmt.Sprintf("starting metrics server on '%s'", o.prometheusServer.Addr))
	return o.prometheusServer.ListenAndServe()
}

// Shutdown ensures all observer stats are flushed before returning.
func (o *Otel) Shutdown(ctx context.Context) error {
	var err error
	for _, f := range o.shutdownFuncs {
		err = errors.Join(err, f(ctx))
	}
	return err
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}
