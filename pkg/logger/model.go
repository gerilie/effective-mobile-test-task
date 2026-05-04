package logger

import (
	"context"
	"fmt"

	"github.com/yushafro/effective-mobile-tz/pkg/env"
	"go.uber.org/zap"
)

type Logger interface {
	Debug(ctx context.Context, msg string, fields ...zap.Field)
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Warn(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)

	With(fields ...zap.Field) Logger

	Zap() *zap.Logger

	Stop() error
}

type logger struct {
	l *zap.Logger
}

func NewWithConfig(cfg Config, environment string) (Logger, error) {
	var config zap.Config
	if environment == env.Dev {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	config.Level = cfg.Level
	l, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("build logger: %w", err)
	}

	return &logger{
		l: l,
	}, nil
}

func (l *logger) Zap() *zap.Logger {
	return l.l
}

func (l *logger) Stop() error {
	return l.l.Sync()
}
