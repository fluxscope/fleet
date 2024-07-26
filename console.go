package fleet

import (
	"context"
	"emperror.dev/emperror"
	"github.com/fluxscope/fleet/pkg/log"
	"github.com/oklog/run"
	"syscall"
)

type Console struct {
	*App
	cmd Command
}

func (c *Console) Run(args []string) {
	emperror.Panic(c.RunE(args))
}

func (c *Console) RunE(args []string) error {
	if err := c.loadConfig(); err != nil {
		return err
	}

	c.initLogger()

	var group run.Group

	for _, svc := range c.services {
		svc := svc
		group.Add(c.serviceHandler(svc))
	}

	group.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))

	group.Add(func() error {
		ctx := log.Context(c.ctx, c.logger)
		return c.cmd.Run(ctx, args...)
	}, func(err error) {
	})

	err := group.Run()
	c.cancelCtx()
	return err
}

func NewConsole(cmd Command, opts ...Option) *Console {
	app := NewApp(opts...)
	return &Console{App: app, cmd: cmd}
}
