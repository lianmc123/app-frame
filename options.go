package app_frame

import (
	"context"
	"github.com/lianmc123/app-frame/log"
	"github.com/lianmc123/app-frame/transport"
	"os"
)

type Option func(o *options)

type options struct {
	id      string
	name    string
	version string

	servers []transport.Server
	logger  *log.Helper
	sigs    []os.Signal
	ctx     context.Context
}

func ID(id string) Option { return func(o *options) { o.id = id } }

func Name(name string) Option { return func(o *options) { o.name = name } }

func Version(version string) Option { return func(o *options) { o.version = version } }

func Logger(logger log.Logger) Option { return func(o *options) { o.logger = log.NewHelper(logger) } }

func Service(srv ...transport.Server) Option { return func(o *options) { o.servers = srv } }

func Signal(sigs ...os.Signal) Option { return func(o *options) { o.sigs = sigs } }

func Context(ctx context.Context) Option { return func(o *options) { o.ctx = ctx } }
