package consul

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/lianmc123/app-frame/registry"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	cfg    *api.Config
	cli    *api.Client
	ctx    context.Context
	cancel context.CancelFunc
}

// NewClient creates consul client
func NewClient(cfg *api.Config) (*Client, error) {
	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	c := &Client{cfg: cfg, cli: client}
	c.ctx, c.cancel = context.WithCancel(context.Background())
	return c, nil
}

func (c *Client) Register(ctx context.Context, svc *registry.ServiceInstance, enableHealthCheck bool) error {
	addresses := make(map[string]api.ServiceAddress)
	var addr string
	var port uint64
	for _, endpoint := range svc.Endpoints {
		raw, err := url.Parse(endpoint)
		if err != nil {
			return err
		}
		addr = raw.Hostname()
		port, _ = strconv.ParseUint(raw.Port(), 10, 16)
		addresses[raw.Scheme] = api.ServiceAddress{Address: endpoint, Port: int(port)}
	}
	asr := &api.AgentServiceRegistration{
		ID:              svc.ID,
		Name:            svc.Name,
		Meta:            svc.Metadata,
		Tags:            []string{fmt.Sprintf("version=%s", svc.Version)},
		TaggedAddresses: addresses,
		Address:         addr,
		Port:            int(port),
		Checks: []*api.AgentServiceCheck{
			{
				CheckID: "service:" + svc.ID,
				TTL:     "50s",
				Status:  "passing",
			},
		},
	}
	if enableHealthCheck {
		asr.Checks = append(asr.Checks, &api.AgentServiceCheck{
			GRPC:     fmt.Sprintf("%s:%d", addr, port),
			Interval: "20s",
			Timeout:  "3s",
		})
	}
	err := c.cli.Agent().ServiceRegister(asr)
	if err != nil {
		return err
	}
	go func() {
		ticker := time.NewTicker(time.Second * 20)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := c.cli.Agent().UpdateTTL("service:"+svc.ID, "pass", "pass"); err != nil {
					// todo
				}
			case <-c.ctx.Done():
				return
			}
		}
	}()
	return nil
}

func (c *Client) Deregister(ctx context.Context, serviceID string) error {
	c.cancel()
	return c.cli.Agent().ServiceDeregister(serviceID)
}
