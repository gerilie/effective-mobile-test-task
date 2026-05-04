package logger

import (
	"context"

	"go.uber.org/zap"
)

type contextKey struct{}

var loggerKey = contextKey{}

// WithLogger returns a new context containing the provided Logger.
func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// WithRequestID attaches a request ID to the logger stored in the context.
// If no logger is found in the context, it will return the context unchanged.
func WithRequestID(ctx context.Context, id string) context.Context {
	l := FromContext(ctx)
	if l == nil {
		return ctx
	}

	l = l.With(zap.String("request_id", id))

	return WithLogger(ctx, l)
}

// FromContext extracts Logger from the context.
// Returns nil if no logger is present.
func FromContext(ctx context.Context) Logger {
	l, ok := ctx.Value(loggerKey).(Logger)
	if !ok {
		return nil
	}

	return l
}
