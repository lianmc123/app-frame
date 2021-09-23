package app_frame

import (
	"context"
	"errors"
	"github.com/lianmc123/app-frame/log"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type App struct {
	ctx    context.Context
	cancel context.CancelFunc
	opts   *options
}

func New(opt ...Option) *App {
	opts := &options{
		id:     uuid.NewV4().String(),
		logger: log.NewHelper(log.DefaultLogger),
		sigs:   []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		ctx:    context.Background(),
	}
	for _, o := range opt {
		o(opts)
	}
	ctx, cancel := context.WithCancel(opts.ctx)
	return &App{
		ctx:    ctx,
		cancel: cancel,
		opts:   opts,
	}
}

type AppInfo interface {
	ID() string
	Name() string
	Version() string
}

func (a *App) ID() string { return a.opts.id }

func (a *App) Name() string { return a.opts.name }

func (a *App) Version() string { return a.opts.version }

type AppKey struct{}

func NewContext(ctx context.Context, app AppInfo) context.Context {
	return context.WithValue(ctx, AppKey{}, app)
}

func FromContext(ctx context.Context) (app AppInfo, ok bool) {
	app, ok = ctx.Value(AppKey{}).(AppInfo)
	return
}

func (a *App) Run() error {
	ctx := NewContext(a.ctx, a)
	eg, ctx := errgroup.WithContext(ctx)

	wg := sync.WaitGroup{} // 确保所有服务都是已经启动了的
	for _, srv := range a.opts.servers {
		s := srv
		eg.Go(func() error {
			<-ctx.Done()
			return s.Stop(ctx)
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			return s.Start(ctx)
		})
	}
	wg.Wait()

	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opts.sigs...)
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				err := a.Stop()
				if err != nil {
					a.opts.logger.ErrorF("failed to app stop: %v", err)
				}
			}
		}
	})
	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func (a *App) Stop() error {
	if a.cancel != nil {
		a.cancel()
	}
	return nil
}
