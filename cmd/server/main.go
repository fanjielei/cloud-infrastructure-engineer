// Ref: https://opentelemetry.io/docs/languages/go/getting-started/
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/alfatraining/cloud-infrastructure-engineer/internal/otel"
	"github.com/alfatraining/cloud-infrastructure-engineer/internal/server"
)

const (
	// Address of your otel collector.
	otelAddr = "localhost:9876"
)

func main() {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Set up OpenTelemetry trace exporter
	// conn, err := grpc.NewClient(otelAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	panic(err)
	// }
	service := otel.Service{Name: "alfaview/training"}

	// otel.New configures which providers and exporters should be used.
	observer, err := otel.New(ctx, service,
		otel.WithZapLogger(),
		otel.WithPrometheusExporter(ctx, ":9090"),
		// otel.WithStdoutTracer(), # Enable if you'd like to experiment with traces
		// otel.WithOTLPTracer(ctx, conn), # Enable if you'd like to experiment with traces
	)
	if err != nil {
		return
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, observer.Shutdown(ctx))
	}()

	srvErr := make(chan error, 1)
	srv := server.New(ctx, observer, ":8080")
	go func() {
		observer.Info(ctx, fmt.Sprintf("starting status server on '%s'", srv.Addr))
		srvErr <- srv.ListenAndServe()
	}()

	// Wait for interruption.
	select {
	case err = <-srvErr:
		// Error when starting HTTP server.
		return
	case <-ctx.Done():
		// Wait for first CTRL+C.
		// Stop receiving signal notifications as soon as possible.
		stop()
	}

	// When Shutdown is called, ListenAndServe immediately returns ErrServerClosed.
	if err = srv.Shutdown(ctx); err != nil {
		panic(err)
	}
}
