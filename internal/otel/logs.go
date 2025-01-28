// Ref: https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/bridges/otelzap/example_test.go
package otel

import (
	"context"
	"os"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger ensures that all log providers adhere to the same standard.
type Logger interface {
	Debug(context.Context, string)
	Info(context.Context, string)
	Error(context.Context, string)
	Fatal(context.Context, string)

	// Shutdown must ensure logs are being written before returning
	Shutdown(context.Context) error
}

var _ Logger = (*logger)(nil)

type logger struct {
	*zap.Logger
}

// WithZapLogger initialises a new zap logger, which adds trace and span IDs from context to the log output.
// TODO: This does not send logs to an OTEL collector yet.
func WithZapLogger() func(*Otel) error {
	return func(observer *Otel) error {
		core := zapcore.NewTee(
			zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig()), zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
			// otelzap.NewCore("alfaview/training", otelzap.WithLoggerProvider(provider)),
		)
		observer.logger = &logger{Logger: zap.New(core)}
		return nil
	}
}

func encoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "severity",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func (l *logger) Debug(ctx context.Context, msg string) {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		l.Logger.Debug(msg)
		return
	}

	l.Logger.Debug(msg,
		zap.String("traceID", span.SpanContext().TraceID().String()),
		zap.String("spanID", span.SpanContext().SpanID().String()),
	)
}
func (l *logger) Info(ctx context.Context, msg string) {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		l.Logger.Info(msg)
		return
	}

	l.Logger.Info(msg,
		zap.String("traceID", span.SpanContext().TraceID().String()),
		zap.String("spanID", span.SpanContext().SpanID().String()),
	)
}

func (l *logger) Error(ctx context.Context, msg string) {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		l.Logger.Error(msg)
		return
	}

	l.Logger.Error(msg,
		zap.String("traceID", span.SpanContext().TraceID().String()),
		zap.String("spanID", span.SpanContext().SpanID().String()),
	)
}
func (l *logger) Fatal(ctx context.Context, msg string) {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		l.Logger.Fatal(msg)
		return
	}

	l.Logger.Fatal(msg,
		zap.String("traceID", span.SpanContext().TraceID().String()),
		zap.String("spanID", span.SpanContext().SpanID().String()),
	)
}

func (l *logger) Shutdown(_ context.Context) error {
	return l.Sync()
}
