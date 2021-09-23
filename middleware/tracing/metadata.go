package tracing

import (
	"context"
	app_frame "github.com/lianmc123/app-frame"
	"github.com/lianmc123/app-frame/metadata"
	"go.opentelemetry.io/otel/propagation"
)

const serviceHeader = "x-md-service-name"

// Metadata is tracing metadata propagator
type Metadata struct{}

var _ propagation.TextMapPropagator = Metadata{}

// Inject sets metadata key-values from ctx into the carrier.
func (b Metadata) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	app, ok := app_frame.FromContext(ctx)
	if ok {
		carrier.Set(serviceHeader, app.Name())
	}
}

// Extract returns a copy of parent with the metadata from the carrier added.
func (b Metadata) Extract(parent context.Context, carrier propagation.TextMapCarrier) context.Context {
	name := carrier.Get(serviceHeader)
	if name != "" {
		if md, ok := metadata.FromServerContext(parent); ok {
			md.Set(serviceHeader, name)
		} else {
			md := metadata.New()
			md.Set(serviceHeader, name)
			parent = metadata.NewServerContext(parent, md)
		}
	}

	return parent
}

// Fields returns the keys who's values are set with Inject.
func (b Metadata) Fields() []string {
	return []string{serviceHeader}
}