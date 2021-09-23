package tracing

import (
	"context"
	"github.com/lianmc123/app-frame/middleware"
	"github.com/lianmc123/app-frame/transport"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type options struct {
	tracerProvider trace.TracerProvider
	propagator     propagation.TextMapPropagator
}

// Option specifies instrumentation configuration options.
type Option func(*options)

// WithPropagators specifies propagators to use for extracting
// information from the HTTP requests. If none are specified, global
// ones will be used.
func WithPropagators(propagators propagation.TextMapPropagator) Option {
	return func(cfg *options) {
		cfg.propagator = propagators
	}
}

// WithTracerProvider specifies a tracer provider to use for creating a tracer.
// If none is specified, the global provider is used.
func WithTracerProvider(provider trace.TracerProvider) Option {
	return func(cfg *options) {
		cfg.tracerProvider = provider
	}
}

// Server returns a new server middleware for OpenTelemetry.
func Server(opts ...Option) middleware.Middleware {
	tracer := NewTracer(trace.SpanKindServer, opts...)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				var span trace.Span
				ctx, span = tracer.Start(ctx, tr.Operation(), tr.RequestHeader())
				setServerSpan(ctx, span, req)
				defer func() { tracer.End(ctx, span, reply, err) }()
			}
			return handler(ctx, req)
		}
	}
}

// Client returns a new client middleware for OpenTelemetry.
func Client(opts ...Option) middleware.Middleware {
	tracer := NewTracer(trace.SpanKindClient, opts...)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromClientContext(ctx); ok {
				var span trace.Span
				ctx, span = tracer.Start(ctx, tr.Operation(), tr.RequestHeader())
				setClientSpan(ctx, span, req)
				defer func() { tracer.End(ctx, span, reply, err) }()
			}
			return handler(ctx, req)
		}
	}
}
