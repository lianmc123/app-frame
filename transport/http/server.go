package http

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/gin-gonic/gin"
	app_err "github.com/lianmc123/app-frame/errors"
	"github.com/lianmc123/app-frame/internal/endpoint"
	"github.com/lianmc123/app-frame/internal/host"
	"github.com/lianmc123/app-frame/log"
	"github.com/lianmc123/app-frame/middleware"
	"github.com/lianmc123/app-frame/transport"
	"golang.org/x/net/http2"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type ServerOption func(*Server)

type ginOption struct {
	handler     *gin.Engine
	recoverFunc gin.RecoveryFunc
}
type Server struct {
	recoverFunc gin.RecoveryFunc
	//gin *ginOption
	addr         string
	readTimeout  time.Duration
	writeTimeout time.Duration
	log          *log.Helper
	tlsCfg       *tls.Config
	http2Enable  bool
	mw           []middleware.Middleware
	httpServer   *http.Server
	endpoint     *url.URL
	lis          net.Listener

	once sync.Once
	err  error
}

func NewServer(addr string, routerConfigure func(gin.IRouter), release bool, opts ...ServerOption) *Server {
	if release {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	engine.Use(gin.Logger())
	srv := &Server{
		addr: addr,
		httpServer: &http.Server{
			Addr:    addr,
			Handler: engine,
		},
		log: log.NewHelper(log.DefaultLogger),
	}

	for _, opt := range opts {
		opt(srv)
	}
	if srv.recoverFunc != nil {
		engine.Use(gin.CustomRecovery(srv.recoverFunc))
	} else {
		engine.Use(gin.Recovery())
	}
	engine.Use(srv.serverMiddleware())

	if routerConfigure != nil {
		routerConfigure(engine)
	}
	return srv
}
func (s *Server) Endpoint() (*url.URL, error) {
	s.once.Do(func() {
		if s.endpoint != nil {
			return
		}
		lis, err := net.Listen("tcp", s.addr)
		if err != nil {
			s.err = err
			return
		}
		addr, err := host.Extract(s.addr, lis)
		if err != nil {
			lis.Close()
			s.err = err
			return
		}
		s.lis = lis
		s.endpoint = endpoint.NewEndpoint("http", addr, s.tlsCfg != nil)
	})
	if s.err != nil {
		return nil, s.err
	}
	return s.endpoint, nil
}

func (s *Server) Start(_ context.Context) error {
	if _, err := s.Endpoint(); err != nil {
		return err
	}
	s.log.InfoF("[HTTP] transport listening on: %s", s.lis.Addr())
	var err error
	if s.tlsCfg != nil {
		if s.http2Enable {
			err = http2.ConfigureServer(s.httpServer, &http2.Server{})
			if err != nil {
				return err
			}
		}
		err = s.httpServer.ServeTLS(s.lis, "", "")
		//err = s.httpServer.ListenAndServeTLS("", "")
	}
	err = s.httpServer.Serve(s.lis)
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func ReadTimeout(readTimeout time.Duration) ServerOption {
	return func(s *Server) { s.readTimeout = readTimeout; s.httpServer.ReadTimeout = readTimeout }
}

func WriteTimeout(writeTimeout time.Duration) ServerOption {
	return func(s *Server) {
		s.writeTimeout = writeTimeout
		s.httpServer.WriteTimeout = writeTimeout
	}
}

func Logger(logger log.Logger) ServerOption {
	return func(s *Server) {
		s.log = log.NewHelper(logger)
	}
}

func TlsConfig(tlsCfg *tls.Config) ServerOption {
	return func(s *Server) {
		s.tlsCfg = tlsCfg
		s.httpServer.TLSConfig = tlsCfg
	}
}

func EnableHttp2(tlsCfg *tls.Config) ServerOption {
	return func(s *Server) {
		s.http2Enable = true
		s.tlsCfg = tlsCfg
		s.httpServer.TLSConfig = tlsCfg
	}
}

func CustomRecover(rf gin.RecoveryFunc) ServerOption {
	return func(s *Server) {
		s.recoverFunc = rf
	}
}

func Middleware(mw ...middleware.Middleware) ServerOption {
	return func(s *Server) {
		s.mw = mw
	}
}

/**
返回app-frame/errors.啥啥啥

*/
func (s *Server) serverMiddleware() func(*gin.Context) {
	return func(c *gin.Context) {
		next := func(ctx context.Context, req interface{}) (interface{}, error) {
			c.Next()
			var err error
			if c.Writer.Status() >= 400 {
				err = app_err.Errorf(c.Writer.Status(), app_err.UnknownReason, app_err.UnknownReason)
			}

			//err := c.Errors.Last()
			//if err != nil {
			//	b := app_err.FromError(err)
			//	if b != nil {
			//		c.JSON(int(b.Code), b)
			//	}else {
			//		c.JSON(app_err.UnknownCode, app_err.Errorf(app_err.UnknownCode, app_err.UnknownReason, app_err.UnknownReason))
			//	}
			//	/*_ = app_err.Errorf(http.StatusInternalServerError, app_err.UnknownReason, app_err.UnknownReason)
			//	err = app_err.un*/
			//	return c.Writer, err
			//}
			return c.Writer, err
		}
		if len(s.mw) > 0 {
			next = middleware.Chain(s.mw...)(next)
		}
		ctx := transport.NewServerContext(c.Request.Context(), &Transport{
			endpoint:    s.endpoint.String(),
			operation:   c.FullPath(),
			reqHeader:   headerCarrier(c.Request.Header),
			replyHeader: headerCarrier(c.Writer.Header()),
			request:     c.Request,
		})
		c.Request = c.Request.WithContext(ctx)
		_, _ = next(ctx, c.Request)
	}
}
