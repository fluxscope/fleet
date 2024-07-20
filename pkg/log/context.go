package log

import (
	"context"
	"log/slog"
)

type logContext struct{}

var ctxLogger = logContext{}

func Context(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxLogger, logger)
}

func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(ctxLogger).(*slog.Logger); ok {
		return logger
	} else {
		return slog.Default()
	}
}
