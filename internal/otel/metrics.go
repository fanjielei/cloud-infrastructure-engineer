package otel

import (
	"context"
	"net"
	"net/http"
	"time"

	goprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	otelprometheus "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

// WithPrometheusExporter initializes a new prometheus exporter and registers it.
func WithPrometheusExporter(ctx context.Context, addr string) func(*Otel) error {
	return func(observer *Otel) error {
		registry := goprometheus.NewRegistry()
		registerer := goprometheus.WrapRegistererWith(map[string]string{}, registry)

		exporter, err := otelprometheus.New(
			otelprometheus.WithRegisterer(registerer),
			otelprometheus.WithoutScopeInfo(),
			otelprometheus.WithoutCounterSuffixes(),
		)
		if err != nil {
			return err
		}

		observer.MeterProvider = metric.NewMeterProvider(
			metric.WithResource(observer.Resource),
			metric.WithReader(exporter),
		)
		mux := http.NewServeMux()
		handler := promhttp.InstrumentMetricHandler(registerer, promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registerer}))
		HandlerFunc("/metrics", mux, handler)

		observer.prometheusServer = &http.Server{
			Addr:         addr,
			BaseContext:  func(_ net.Listener) context.Context { return ctx },
			ReadTimeout:  time.Second,
			WriteTimeout: 1 * time.Second,
			Handler: otelhttp.NewHandler(mux, "/",
				otelhttp.WithPropagators(observer.TextMapPropagator),
				otelhttp.WithMeterProvider(observer.MeterProvider),
				otelhttp.WithTracerProvider(observer.TracerProvider),
			),
		}
		return nil
	}
}

// HandlerFunc adds the route as otel tag.
func HandlerFunc(pattern string, mux *http.ServeMux, f http.Handler) {
	handler := otelhttp.WithRouteTag(pattern, f)
	mux.Handle(pattern, handler)
}
