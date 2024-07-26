package main

import (
	"context"
	"github.com/fluxscope/fleet"
	"github.com/fluxscope/fleet/pkg/log"
)

var _ fleet.Command = new(Command)

func NewCommand() fleet.Command {
	return &Command{}
}

type Command struct {
}

func (c *Command) Run(ctx context.Context, args ...string) error {
	logger := log.FromContext(ctx)
	logger.Info("hello world", "args", args)
	return nil
}
