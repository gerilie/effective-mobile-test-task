package logger

import (
	"context"
	"fmt"

	"github.com/yushafro/effective-mobile-tz/pkg/env"
	"go.uber.org/zap"
)

// Logger defines a structured logging interface.
type Logger interface {
	Debug(ctx context.Context, msg string, fields ...zap.Field)
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Warn(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)

	// With returns a new Logger with additional structured fields.
	With(fields ...zap.Field) Logger

	// Zap returns the underlying zap.Logger instance.
	Zap() *zap.Logger

	// Stop flushes any buffered log entries.
	Stop() error
}

type logger struct {
	l *zap.Logger
}

// NewWithConfig creates a new Logger based on provided Config and environment.
// Uses development config for dev environment and production config otherwise.
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

// Zap returns the underlying zap.Logger.
func (l *logger) Zap() *zap.Logger {
	return l.l
}

// Stop flushes buffered logs.
func (l *logger) Stop() error {
	return l.l.Sync()
}
