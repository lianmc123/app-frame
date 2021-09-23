package metadata

import (
	"context"
	"strings"
)

type Metadata map[string]string

// New creates an MD from a given key-values map.
func New(mds ...map[string]string) Metadata {
	md := Metadata{}
	for _, m := range mds {
		for k, v := range m {
			md.Set(k, v)
		}
	}
	return md
}

// Set stores the key-value pair.
func (m Metadata) Set(key string, value string) {
	if key == "" || value == "" {
		return
	}
	k := strings.ToLower(key)
	m[k] = value
}

// Get returns the value associated with the passed key.
func (m Metadata) Get(key string) string {
	k := strings.ToLower(key)
	return m[k]
}

// Clone returns a deep copy of Metadata
func (m Metadata) Clone() Metadata {
	md := Metadata{}
	for k, v := range m {
		md[k] = v
	}
	return md
}

type serverMetadataKey struct{}

// NewServerContext creates a new context with client md attached.
func NewServerContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, serverMetadataKey{}, md)
}

// FromServerContext returns the server metadata in ctx if it exists.
func FromServerContext(ctx context.Context) (Metadata, bool) {
	md, ok := ctx.Value(serverMetadataKey{}).(Metadata)
	return md, ok
}

type clientMetadataKey struct{}

// NewClientContext creates a new context with client md attached.
func NewClientContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, clientMetadataKey{}, md)
}

// FromClientContext returns the client metadata in ctx if it exists.
func FromClientContext(ctx context.Context) (Metadata, bool) {
	md, ok := ctx.Value(clientMetadataKey{}).(Metadata)
	return md, ok
}
