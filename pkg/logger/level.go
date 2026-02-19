package logger

import (
	"context"

	"go.uber.org/zap"
)

func (l *logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	fields = addID(ctx, fields)
	l.l.Debug(msg, fields...)
}

func (l *logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	fields = addID(ctx, fields)
	l.l.Info(msg, fields...)
}

func (l *logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	fields = addID(ctx, fields)
	l.l.Warn(msg, fields...)
}

func (l *logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	fields = addID(ctx, fields)
	l.l.Error(msg, fields...)
}

func (l *logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	fields = addID(ctx, fields)
	l.l.Fatal(msg, fields...)
}
