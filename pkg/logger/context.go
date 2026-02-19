package logger

import (
	"context"

	"go.uber.org/zap"
)

func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func FromContext(ctx context.Context) Logger {
	l, ok := ctx.Value(loggerKey).(Logger)
	if !ok {
		l := NewBootstrap()

		l.Warn(ctx, "logger not found in context, creating new one")

		return l
	}

	return l
}

func addID(ctx context.Context, fields []zap.Field) []zap.Field {
	id, ok := ctx.Value(RequestID).(string)
	if ok {
		fields = append(fields, zap.String("request_id", id))
	}

	return fields
}
