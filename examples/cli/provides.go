package main

import (
	"context"
	"github.com/fluxscope/fleet"
	"github.com/fluxscope/fleet/pkg/log"
	"github.com/google/wire"
	"time"
)

var ProvideSet = wire.NewSet(NewApp)

func NewApp() *fleet.App {
	return fleet.NewApp(
		fleet.WithCommand(func(ctx context.Context, args ...string) error {
			logger := log.FromContext(ctx)
			logger.InfoContext(ctx, "run hello cli", "args", args)
			time.Sleep(5 * time.Second)
			logger.InfoContext(ctx, "run hello cli success")
			return nil
		}, 1*time.Second),
		fleet.WithBeforeStartHooks(func(ctx context.Context) error {
			logger := log.FromContext(ctx)
			logger.Info("before start hook")
			return nil
		}),
		fleet.WithOnStartingHooks(func(ctx context.Context) error {
			logger := log.FromContext(ctx)
			logger.Info("on starting hook")
			return nil
		}),
		fleet.WithOnStoppingHooks(func(ctx context.Context) error {
			logger := log.FromContext(ctx)
			logger.Info("on stopping hook")
			return nil
		}),
		fleet.WithAfterStoppedHooks(func(ctx context.Context) error {
			logger := log.FromContext(ctx)
			logger.Info("after stopped hook")
			return nil
		}),
	)
}
