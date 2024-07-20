package handler

import (
	"context"
	slogcommon "github.com/samber/slog-common"
	slogzap "github.com/samber/slog-zap/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log/slog"
	"runtime"
)

type Option struct {
	// log level (default: debug)
	Level slog.Leveler

	// optional: zap logger (default: zap.L())
	Logger *zap.Logger

	// optional: customize json payload builder
	Converter slogzap.Converter

	// optional: see slog.HandlerOptions
	AddSource   bool
	ReplaceAttr func(groups []string, a slog.Attr) slog.Attr

	Extractor ContextExtractor
}

func (o Option) NewZapHandler() slog.Handler {
	if o.Level == nil {
		o.Level = slog.LevelDebug
	}

	if o.Logger == nil {
		// should be selected lazily ?
		o.Logger = zap.L()
	}

	return &ZapHandler{
		option: o,
		attrs:  []slog.Attr{},
		groups: []string{},
	}
}

var _ slog.Handler = (*ZapHandler)(nil)

type ZapHandler struct {
	option Option
	attrs  []slog.Attr
	groups []string
}

func (h *ZapHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.option.Level.Level()
}

func (h *ZapHandler) Handle(ctx context.Context, record slog.Record) error {
	converter := slogzap.DefaultConverter
	if h.option.Converter != nil {
		converter = h.option.Converter
	}

	level := slogzap.LogLevels[record.Level]
	fields := converter(h.option.AddSource, h.option.ReplaceAttr, h.attrs, h.groups, &record)

	extractor := NoOpExtractor()
	if h.option.Extractor != nil {
		extractor = h.option.Extractor

	}
	ctxValues := extractor(ctx)
	for k, v := range ctxValues {
		fields = append(fields, zap.Any(k, v))
	}

	checked := h.option.Logger.Check(level, record.Message)
	if checked != nil {
		if h.option.AddSource {
			frame, _ := runtime.CallersFrames([]uintptr{record.PC}).Next()
			checked.Caller = zapcore.NewEntryCaller(0, frame.File, frame.Line, true)
			checked.Stack = "" //@TODO
		} else {
			checked.Caller = zapcore.EntryCaller{}
			checked.Stack = ""
		}
		checked.Write(fields...)
		return nil
	} else {
		h.option.Logger.Log(level, record.Message, fields...)
	}

	return nil
}

func (h *ZapHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ZapHandler{
		option: h.option,
		attrs:  slogcommon.AppendAttrsToGroup(h.groups, h.attrs, attrs...),
		groups: h.groups,
	}
}

func (h *ZapHandler) WithGroup(name string) slog.Handler {
	return &ZapHandler{
		option: h.option,
		attrs:  h.attrs,
		groups: append(h.groups, name),
	}
}
