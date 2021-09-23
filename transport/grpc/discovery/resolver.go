package discovery

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/lianmc123/app-frame/internal/endpoint"
	"github.com/lianmc123/app-frame/log"
	"github.com/lianmc123/app-frame/registry"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
	"time"
)

type discoveryResolver struct {
	insecure bool
	w        registry.Watcher
	cc       resolver.ClientConn

	ctx    context.Context
	cancel context.CancelFunc

	log *log.Helper
}

func (r *discoveryResolver) ResolveNow(options resolver.ResolveNowOptions) {}

func (r *discoveryResolver) Close() {
	r.cancel()
	err := r.w.Stop()
	if err != nil {
		r.log.ErrorF("[resolver] failed to watch top: %s", err)
	}
}

func (r *discoveryResolver) watch() {
	for {
		select {
		case <-r.ctx.Done():
			return
		default:
		}
		ins, err := r.w.Next()
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			r.log.ErrorF("[resolver] Failed to watch discovery endpoint: %v", err)
			time.Sleep(time.Second)
			continue
		}
		r.update(ins)
	}
}

func (r *discoveryResolver) update(ins []*registry.ServiceInstance) {
	var addrs []resolver.Address
	for _, in := range ins {
		ep, err := endpoint.ParseEndpoint(in.Endpoints, "grpc", !r.insecure)
		if err != nil {
			r.log.ErrorF("[resolver] Failed to parse discovery endpoint: %v", err)
			continue
		}
		if ep == "" {
			continue
		}
		addr := resolver.Address{
			ServerName: in.Name,
			Attributes: parseAttributes(in.Metadata),
			Addr:       ep,
		}
		addrs = append(addrs, addr)
	}
	if len(addrs) == 0 {
		r.log.WarningF("[resolver] Zero endpoint found,refused to write, instances: %v", ins)
		return
	}
	err := r.cc.UpdateState(resolver.State{Addresses: addrs})
	if err != nil {
		r.log.ErrorF("[resolver] failed to update state: %s", err)
	}
	b, _ := json.Marshal(ins)
	r.log.InfoF("[resolver] update instances: %s", b)
}

func parseAttributes(md map[string]string) *attributes.Attributes {
	pairs := make([]interface{}, 0, len(md))
	for k, v := range md {
		pairs = append(pairs, k, v)
	}
	return attributes.New(pairs...)
}
