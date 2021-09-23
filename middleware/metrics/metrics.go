package metrics

import (
	"context"
	"github.com/lianmc123/app-frame/errors"
	"github.com/lianmc123/app-frame/metrics"
	"github.com/lianmc123/app-frame/middleware"
	"github.com/lianmc123/app-frame/transport"
	"strconv"
	"time"
)

// Option is metrics option.
type Option func(*options)

// WithRequests with requests counter.
func WithRequests(c metrics.Counter) Option {
	return func(o *options) {
		o.requests = c
	}
}

// WithSeconds with seconds histogram.
func WithSeconds(c metrics.Observer) Option {
	return func(o *options) {
		o.seconds = c
	}
}

type options struct {
	// counter: <client/transport>_requests_code_total{kind, operation, code, reason}
	requests metrics.Counter
	// histogram: <client/transport>_requests_seconds_bucket{kind, operation}
	seconds metrics.Observer
}

// Server is middleware transport-side metrics.
func Server(opts ...Option) middleware.Middleware {
	options := options{}
	for _, o := range opts {
		o(&options)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var (
				code      int
				reason    string
				kind      string
				operation string
			)
			startTime := time.Now()
			if info, ok := transport.FromServerContext(ctx); ok {
				kind = info.Kind().String()
				operation = info.Operation()
			}
			reply, err := handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = int(se.Code)
				reason = se.Reason
			}
			if options.requests != nil {
				options.requests.With(kind, operation, strconv.Itoa(code), reason).Inc()
			}
			if options.seconds != nil {
				options.seconds.With(kind, operation).Observe(time.Since(startTime).Seconds() * 1000)
			}
			return reply, err
		}
	}
}

// Client is middleware client-side metrics.
func Client(opts ...Option) middleware.Middleware {
	options := options{}
	for _, o := range opts {
		o(&options)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var (
				code      int
				reason    string
				kind      string
				operation string
			)
			startTime := time.Now()
			if info, ok := transport.FromClientContext(ctx); ok {
				kind = info.Kind().String()
				operation = info.Operation()
			}
			reply, err := handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = int(se.Code)
				reason = se.Reason
			}
			if options.requests != nil {
				options.requests.With(kind, operation, strconv.Itoa(code), reason).Inc()
			}
			if options.seconds != nil {
				options.seconds.With(kind, operation).Observe(time.Since(startTime).Seconds())
			}
			return reply, err
		}
	}
}
