package fleet

import (
	"context"
	"emperror.dev/emperror"
	"github.com/fluxscope/fleet/pkg/log"
	"github.com/fluxscope/fleet/pkg/zap"
	configpb "github.com/fluxscope/fleet/proto/config"
	"github.com/oklog/run"
	"log/slog"
	"syscall"
	"time"
)

type App struct {
	id                string
	services          map[string]Service
	buildInfo         BuildInfo
	ctx               context.Context
	cancelCtx         context.CancelFunc
	shutdownTimeout   time.Duration
	logger            *slog.Logger
	beforeStartHooks  []Hook
	onStartingHooks   []Hook
	onStoppingHooks   []Hook
	afterStoppedHooks []Hook
	configPath        string
}

type Option func(*App)

func WithConfigPath(path string) Option {
	return func(a *App) {
		a.configPath = path
	}
}

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

func WithOnStartingHooks(hooks ...Hook) Option {
	return func(a *App) {
		a.onStartingHooks = append(a.onStartingHooks, hooks...)
	}
}

func WithOnStoppingHooks(hooks ...Hook) Option {
	return func(a *App) {
		a.onStoppingHooks = append(a.onStoppingHooks, hooks...)
	}
}

func WithAfterStoppedHooks(hooks ...Hook) Option {
	return func(a *App) {
		a.afterStoppedHooks = append(a.afterStoppedHooks, hooks...)
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

func (app *App) runHook(h Hook) error {
	ctx := log.Context(app.ctx, app.logger)
	return h(ctx)
}

func (app *App) initLogger() {
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
}

func (app *App) loadConfig() error {
	return nil
}

func (app *App) shutdownCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), app.shutdownTimeout)
}

func (app *App) serviceHandler(svc Service) (execute func() error, interrupt func(error)) {
	ctx, cancel := context.WithCancel(app.ctx)
	return func() error {
			ctx := log.Context(ctx, app.logger)
			return svc.Run(ctx)
		}, func(err error) {
			app.logger.Warn("shutdown service", "cause", err)
			defer cancel()
			ctx, cancelCtx := app.shutdownCtx()
			defer cancelCtx()
			if err := svc.Shutdown(ctx); err != nil {
				app.logger.Error("failed to shutdown service", "error", err)
			}
		}
}

func (app *App) RunE() error {
	if err := app.loadConfig(); err != nil {
		return err
	}

	app.initLogger()

	var group run.Group
	for _, svc := range app.services {
		svc := svc
		group.Add(app.serviceHandler(svc))
	}

	group.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))

	group.Add(app.hookHandler())

	// hook
	for _, h := range app.beforeStartHooks {
		if err := app.runHook(h); err != nil {
			return err
		}
	}

	err := group.Run()
	app.cancelCtx()
	// hook
	for _, h := range app.afterStoppedHooks {
		if err := app.runHook(h); err != nil {
			return err
		}
	}
	return err
}

func (app *App) Run() {
	emperror.Panic(app.RunE())
}

func (app *App) hookHandler() (execute func() error, interrupt func(error)) {
	ctx, cancel := context.WithCancel(app.ctx)
	return func() error {
			for _, h := range app.onStartingHooks {
				if err := app.runHook(h); err != nil {
					return err
				}
			}
			select {
			case <-ctx.Done():
				return nil
			}
		}, func(err error) {
			ctx, cancelCtx := context.WithTimeout(context.Background(), app.shutdownTimeout)
			defer cancelCtx()
			for _, h := range app.onStoppingHooks {
				_ = h(log.Context(ctx, app.logger))
			}
			cancel()
		}
}
