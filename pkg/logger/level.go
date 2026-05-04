package logger

import (
	"context"

	"go.uber.org/zap"
)

// Debug logs a debug-level message.
func (l *logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l.l.Debug(msg, fields...)
}

// Info logs an info-level message.
func (l *logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.l.Info(msg, fields...)
}

// Warn logs a warning-level message.
func (l *logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.l.Warn(msg, fields...)
}

// Error logs an error-level message.
func (l *logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.l.Error(msg, fields...)
}
