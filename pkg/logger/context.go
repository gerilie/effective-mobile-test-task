package logger

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"
)

type contextKey string

const (
	loggerKey    contextKey = "logger"
	requestIDKey contextKey = "request_id"
)

// WithLogger returns a new context containing the provided Logger.
func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext extracts Logger from the context.
// Returns nil if no logger is present.
func FromContext(ctx context.Context) Logger {
	l, ok := ctx.Value(loggerKey).(Logger)
	if !ok {
		fmt.Fprintln(os.Stderr, "WARNING: logger not found in context, using noop logger")

		return &noopLogger{}
	}

	return l
}

// WithRequestID attaches a request ID to the logger stored in the context
// and makes it available for direct retrieval via RequestIDFromContext.
//
// The request ID is added both to the logger's fields (for structured logging)
// and to the context directly (for easy extraction without accessing the logger).
func WithRequestID(ctx context.Context, id string) context.Context {
	l := FromContext(ctx)

	l = &logger{
		l: l.Zap().With(zap.String("request_id", id)),
	}

	ctx = WithLogger(ctx, l)
	ctx = context.WithValue(ctx, requestIDKey, id)

	return ctx
}

// RequestIDFromContext extracts the request ID from the context.
// It returns the request ID string if present, or an empty string if not found.
func RequestIDFromContext(ctx context.Context) string {
	id, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		return ""
	}

	return id
}
