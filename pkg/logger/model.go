package logger

import (
	"context"

	"go.uber.org/zap"
)

type key string

const (
	loggerKey key = "logger"
	RequestID key = "request_id"
)

type Config struct {
	Level string `mapstructure:"level"`
}

type Logger interface {
	Debug(ctx context.Context, msg string, fields ...zap.Field)
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Warn(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
	Fatal(ctx context.Context, msg string, fields ...zap.Field)
	Stop() error
}

type logger struct {
	l *zap.Logger
}

func NewBootstrap() Logger {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	l, _ := cfg.Build()

	return &logger{
		l: l,
	}
}

func NewWithConfig(cfg Config) Logger {
	_cfg := zap.NewProductionConfig()

	atomicLevel := zap.NewAtomicLevel()
	switch cfg.Level {
	case "debug":
		atomicLevel.SetLevel(zap.DebugLevel)
	case "info":
		atomicLevel.SetLevel(zap.InfoLevel)
	case "warn":
		atomicLevel.SetLevel(zap.WarnLevel)
	case "error":
		atomicLevel.SetLevel(zap.ErrorLevel)
	case "fatal":
		atomicLevel.SetLevel(zap.FatalLevel)
	default:
		atomicLevel.SetLevel(zap.InfoLevel)
	}
	_cfg.Level = atomicLevel

	l, err := _cfg.Build()
	if err != nil {
		return &logger{
			l: zap.L(),
		}
	}

	return &logger{
		l: l,
	}
}

func (l *logger) Stop() error {
	return l.l.Sync()
}
