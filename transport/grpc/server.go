package grpc

import (
	"context"
	"crypto/tls"
	ic "github.com/lianmc123/app-frame/internal/context"
	"github.com/lianmc123/app-frame/internal/endpoint"
	"github.com/lianmc123/app-frame/internal/host"
	"github.com/lianmc123/app-frame/log"
	"github.com/lianmc123/app-frame/middleware"
	"github.com/lianmc123/app-frame/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"net"
	"net/url"
	"sync"
	"time"
)

type ServerOption func(o *Server)

func Logger(logger log.Logger) ServerOption {
	return func(s *Server) {
		s.log = log.NewHelper(logger)
	}
}

func Timeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = timeout
	}
}

func TLSConfig(c *tls.Config) ServerOption {
	return func(s *Server) {
		s.tlsConf = c
	}
}

func Options(opts ...grpc.ServerOption) ServerOption {
	return func(s *Server) {
		s.grpcOpts = opts
	}
}

func UnaryInterceptor(interceptors ...grpc.UnaryServerInterceptor) ServerOption {
	return func(s *Server) {
		s.interceptors = interceptors
	}
}

func Address(addr string) ServerOption {
	return func(s *Server) {
		s.address = addr
	}
}

func Middleware(m ...middleware.Middleware) ServerOption {
	return func(s *Server) {
		s.middleware = m
	}
}

type Server struct {
	ctx     context.Context
	Server  *grpc.Server
	address string

	tlsConf         *tls.Config
	timeout         time.Duration
	grpcOpts        []grpc.ServerOption
	log             *log.Helper
	once            sync.Once
	endpoint        *url.URL
	interceptors    []grpc.UnaryServerInterceptor
	middleware      []middleware.Middleware

	health *health.Server

	err error
	lis net.Listener
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		log:     log.NewHelper(log.DefaultLogger),
		health:  health.NewServer(),
		timeout: 1 * time.Second,
		address: ":0",
	}

	for _, o := range opts {
		o(srv)
	}

	var interceptors = []grpc.UnaryServerInterceptor{
		srv.unaryServerInterceptor(),
	}
	if len(srv.interceptors) > 0 {
		interceptors = append(interceptors, srv.interceptors...)
	}

	var grpcOpts = []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(interceptors...),
	}
	if srv.tlsConf != nil {
		grpcOpts = append(grpcOpts, grpc.Creds(credentials.NewTLS(srv.tlsConf)))
	}
	if len(srv.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, srv.grpcOpts...)
	}
	srv.Server = grpc.NewServer(grpcOpts...)
	grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)
	return srv
}

func (s *Server) unaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		//fmt.Println("..1..")
		ctx, cancel := ic.Merge(ctx, s.ctx)
		defer cancel()
		md, _ := metadata.FromIncomingContext(ctx)
		replyHeader := metadata.MD{}
		ctx = transport.NewServerContext(ctx, &Transport{
			endpoint:    s.endpoint.String(),
			operation:   info.FullMethod,
			reqHeader:   headerCarrier(md),
			replyHeader: headerCarrier(replyHeader),
		})
		if s.timeout > 0 {
			ctx, cancel = context.WithTimeout(ctx, s.timeout)
			defer cancel()
		}
		h := func(ctx context.Context, req interface{}) (interface{}, error) {
			return handler(ctx, req)
		}
		if len(s.middleware) > 0 {
			h = middleware.Chain(s.middleware...)(h)
		}
		reply, err := h(ctx, req)
		if len(replyHeader) > 0 {
			_ = grpc.SetHeader(ctx, replyHeader)
		}
		return reply, err
	}
}

func (s *Server) Start(ctx context.Context) error {
	if _, err := s.Endpoint(); err != nil {
		return err
	}
	s.ctx = ctx
	s.log.InfoF("[gRPC] transport listening on: %s", s.lis.Addr().String())
	s.health.Resume()
	return s.Server.Serve(s.lis)
}

func (s *Server) Stop(ctx context.Context) error {
	s.Server.GracefulStop()
	s.health.Shutdown()
	s.log.Info("[gRPC] transport stopping")
	return nil
}

func (s *Server) Endpoint() (*url.URL, error) {
	s.once.Do(func() {
		lis, err := net.Listen("tcp", s.address)
		if err != nil {
			s.err = err
			return
		}
		addr, err := host.Extract(s.address, lis)
		if err != nil {
			_ = lis.Close()
			s.err = err
			return
		}
		s.lis = lis
		s.endpoint = endpoint.NewEndpoint("grpc", addr, s.tlsConf != nil)
	})
	if s.err != nil {
		return nil, s.err
	}
	return s.endpoint, nil
}
