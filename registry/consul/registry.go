package consul

import (
	"context"
	"github.com/hashicorp/consul/api"
	"github.com/lianmc123/app-frame/registry"
)

var (
	_ registry.Registrar = &Registry{}
	_ registry.Discovery = &Registry{}
)

type Option func(registry *Registry)

func WithHealthCheck(enable bool) Option {
	return func(o *Registry) {
		o.enableHealthCheck = enable
	}
}

type Registry struct {
	cli               *Client
	enableHealthCheck bool
}

// New creates consul registry
func New(cfg *api.Config, opts ...Option) (*Registry, error) {
	cli, err := NewClient(cfg)
	if err != nil {
		return nil, err
	}
	r := &Registry{
		cli: cli,
		//registry:          make(map[string]*serviceSet),
		enableHealthCheck: true,
	}
	for _, o := range opts {
		o(r)
	}
	return r, nil
}

func (r *Registry) Register(ctx context.Context, svc *registry.ServiceInstance) error {
	return r.cli.Register(ctx, svc, r.enableHealthCheck)
}

func (r *Registry) Deregister(ctx context.Context, svc *registry.ServiceInstance) error {
	return r.cli.Deregister(ctx, svc.ID)
}

func (r *Registry) GetService(ctx context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	panic("implement me")
}

func (r *Registry) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {

}
