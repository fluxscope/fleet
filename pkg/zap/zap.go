package zap

import (
	"context"
	"github.com/fluxscope/fleet/pkg/zap/handler"
	configpb "github.com/fluxscope/fleet/proto/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log/slog"
)

func NewZapLogger(config *configpb.Logging) (*zap.Logger, error) {
	configLevel := "DEBUG"
	if config.Level != "" {
		configLevel = config.Level
	}
	prod := false
	if config.Zap != nil {
		prod = config.Zap.Production
	}
	var zapConfig zap.Config
	if prod {
		zapConfig = zap.NewProductionConfig()
	} else {
		zapConfig = zap.NewDevelopmentConfig()
	}
	zapConfig.DisableCaller = true
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapConfig.Sampling = nil
	level, err := zapcore.ParseLevel(configLevel)
	if err != nil {
		level = zapcore.DebugLevel
	}
	zapConfig.Level = zap.NewAtomicLevelAt(level)
	return zapConfig.Build()
}

func NewSLogger(config *configpb.Logging) (*slog.Logger, error) {
	zlog, err := NewZapLogger(config)
	if err != nil {
		return nil, err
	}

	return NewSLoggerFromZap(config, zlog)

}

func NewSLoggerFromZap(config *configpb.Logging, zlogger *zap.Logger) (*slog.Logger, error) {
	l := slog.LevelDebug
	if config.Level != "" {
		if err := l.UnmarshalText([]byte(config.Level)); err != nil {
			return nil, err
		}
	}
	fields := []string{}
	if config.Zap != nil && config.Zap.ContextFields != nil {
		fields = append(fields, config.Zap.ContextFields...)
	}
	logger := slog.New(handler.Option{Level: l, Logger: zlogger, Extractor: FieldsExtractor(fields...)}.NewZapHandler())
	//slog.SetDefault(logger)
	return logger, nil
}

func FieldsExtractor(fields ...string) handler.ContextExtractor {
	return func(ctx context.Context) map[string]interface{} {
		ctxValues := make(map[string]interface{})
		for _, field := range fields {
			ctxValues[field] = ctx.Value(field)
		}
		return ctxValues
	}
}
