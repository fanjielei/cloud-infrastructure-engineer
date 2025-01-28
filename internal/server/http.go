package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/alfatraining/cloud-infrastructure-engineer/internal/otel"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var (
	ErrInvalidStatusCode = "invalid http status code"
)

// New creates a new http server.
func New(ctx context.Context, observer otel.Observer, addr string) *http.Server {
	srv := &http.Server{
		Addr:         addr,
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      httpHandler(observer),
	}
	return srv
}

func httpHandler(observer otel.Observer) http.Handler {
	mux := http.NewServeMux()
	status := statusHandler{
		statusCode: 200,
		observer:   observer,
	}

	otel.HandlerFunc("/status", mux, http.HandlerFunc(status.Get))
	otel.HandlerFunc("/status/{code}", mux, http.HandlerFunc(status.Post))
	otel.HandlerFunc("/flaky", mux, http.HandlerFunc(status.Flaky))

	// Add HTTP instrumentation for the whole server.
	handler := otelhttp.NewHandler(mux, "/",
		otelhttp.WithPropagators(observer),
		otelhttp.WithMeterProvider(observer),
		otelhttp.WithTracerProvider(observer),
	)
	return handler
}
