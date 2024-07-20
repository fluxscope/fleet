package fleet

import (
	"context"
	"github.com/fluxscope/fleet/pkg/log"
	"github.com/fluxscope/fleet/pkg/zap"
	configpb "github.com/fluxscope/fleet/proto/config"
	"github.com/oklog/run"
	"log/slog"
	"syscall"
	"time"
)

type App struct {
	id               string
	services         map[string]Service
	buildInfo        BuildInfo
	ctx              context.Context
	cancelCtx        context.CancelFunc
	shutdownTimeout  time.Duration
	logger           *slog.Logger
	beforeStartHooks []Hook
	afterStartHooks  []Hook
	beforeStopHooks  []Hook
}

type Option func(*App)

func WithLogger(logger *slog.Logger) Option {
	return func(a *App) {
		a.logger = logger
	}
}

func WithShutdownTimeout(timeout time.Duration) Option {
	return func(app *App) {
		app.shutdownTimeout = timeout
	}
}

func WithID(id string) Option {
	return func(app *App) {
		app.id = id
	}
}
func WithService(services ...Service) Option {
	return func(a *App) {
		for _, s := range services {
			a.services[s.ID()] = s
		}
	}
}

func WithBuildInfo(buildInfo BuildInfo) Option {
	return func(a *App) {
		a.buildInfo = buildInfo
	}
}

func WithBeforeStartHooks(hooks ...Hook) Option {
	return func(a *App) {
		a.beforeStartHooks = append(a.beforeStartHooks, hooks...)
	}
}

func WithAfterStartHooks(hooks ...Hook) Option {
	return func(a *App) {
		a.afterStartHooks = append(a.afterStartHooks, hooks...)
	}
}

func WithBeforeStopHooks(hooks ...Hook) Option {
	return func(a *App) {
		a.beforeStopHooks = append(a.beforeStopHooks, hooks...)
	}
}

func NewApp(opts ...Option) *App {
	app := &App{
		services: map[string]Service{},
	}
	app.ctx, app.cancelCtx = context.WithCancel(context.Background())

	for _, opt := range opts {
		opt(app)
	}

	if app.shutdownTimeout == 0 {
		app.shutdownTimeout = 60 * time.Second
	}

	return app
}

type BuildInfo struct {
	Name    string
	Version string
}

func (app *App) Run() {

	if app.logger == nil {
		slogger, err := zap.NewSLogger(&configpb.Logging{
			Level: "DEBUG",
			Zap: &configpb.Logging_Zap{
				Production:    true,
				ContextFields: []string{},
			},
		})
		if err != nil {
			panic(err)
		}
		slogger = slogger.With("id", app.id, "name", app.buildInfo.Name, "version", app.buildInfo.Version)

		app.logger = slogger
	}
	slog.SetDefault(app.logger)

	var group run.Group
	defers := []func() error{}

	for _, svc := range app.services {
		group.Add(func() error {
			ctx := log.Context(app.ctx, app.logger)
			return svc.Run(ctx)
		}, func(err error) {
			app.logger.Error("failed to run service", "error", err)

			ctx, cancelCtx := context.WithTimeout(context.Background(), app.shutdownTimeout)
			defer cancelCtx()
			if err := svc.Shutdown(ctx); err != nil {
				app.logger.ErrorContext(ctx, "failed to shutdown service", "error", err)
			}
		})
	}
	defer func() {
		for _, df := range defers {
			_ = df()
		}
	}()

	group.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))

	group.Add(func() error {
		for _, h := range app.afterStartHooks {
			ctx := log.Context(app.ctx, app.logger)
			if err := h(ctx); err != nil {
				return err
			}
		}
		select {}
	}, func(err error) {
		for _, h := range app.beforeStopHooks {
			_ = h(context.Background())
		}
	})
	for _, h := range app.beforeStartHooks {
		ctx := log.Context(app.ctx, app.logger)
		if err := h(ctx); err != nil {
			panic(err)
		}
	}

	err := group.Run()
	panic(err)
}
