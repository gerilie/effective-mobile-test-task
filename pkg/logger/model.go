package logger

import (
	"context"

	"github.com/yushafro/effective-mobile-tz/pkg/env"
	"go.uber.org/zap"
)

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

func NewBootstrap(environment string) Logger {
	var cfg zap.Config
	if environment == env.Dev {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	l, err := cfg.Build()
	if err != nil {
		return &logger{
			l: zap.L(),
		}
	}

	return &logger{
		l: l,
	}
}

func NewWithConfig(cfg Config, environment string) Logger {
	var config zap.Config
	if environment == env.Dev {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

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
	config.Level = atomicLevel

	l, err := config.Build()
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
