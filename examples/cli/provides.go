package main

import (
	"context"
	"github.com/fluxscope/fleet"
	"github.com/fluxscope/fleet/pkg/log"
	"github.com/fluxscope/fleet/pkg/runx"
	"github.com/google/wire"
)

//go:generate wire

var ProvideSet = wire.NewSet(NewApp, NewCommand)

func NewApp(cmd fleet.Command) *fleet.Console {
	return fleet.NewConsole(
		cmd,
		fleet.WithService(&service{}),
	)
}

var _ fleet.Service = new(service)

type service struct {
}

func (s *service) ID() string {
	return "hello-test"
}

func (s *service) Run(ctx context.Context) error {
	logger := log.FromContext(ctx)
	logger.ErrorContext(ctx, "hello service run")
	return runx.RunForever(ctx)
}

func (s *service) Shutdown(ctx context.Context) error {
	return nil
}
