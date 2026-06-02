package logger

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"
)

type contextKey struct{}

//nolint:gochecknoglobals
var (
	loggerKey    = contextKey{}
	requestIDKey = contextKey{}
)

// WithLogger returns a new context containing the provided Logger.
func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// WithRequestID attaches a request ID to the logger stored in the context
// and makes it available for direct retrieval via RequestIDFromContext.
//
// If no logger is found in the context, the context is returned unchanged.
// The request ID is added both to the logger's fields (for structured logging)
// and to the context directly (for easy extraction without accessing the logger).
//
// This function is typically used in HTTP middleware to ensure every request
// has a unique identifier that appears in all log entries and can be retrieved
// by downstream handlers.
func WithRequestID(ctx context.Context, id string) context.Context {
	l := FromContext(ctx)
	if l == nil {
		return ctx
	}

	l = &logger{
		l: l.Zap().With(zap.String("request_id", id)),
	}

	ctx = WithLogger(ctx, l)
	ctx = context.WithValue(ctx, requestIDKey, id)

	return ctx
}

// RequestIDFromContext extracts the request ID from the context.
// It returns the request ID string if present, or an empty string if not found.
//
// The request ID is set by WithRequestID, typically called from HTTP middleware
// during request processing. It can be used for request tracing, logging
// correlation, and debugging across different parts of the application.
func RequestIDFromContext(ctx context.Context) string {
	id, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		return ""
	}

	return id
}

// FromContext extracts Logger from the context.
// Returns nil if no logger is present.
func FromContext(ctx context.Context) Logger {
	l, ok := ctx.Value(loggerKey).(Logger)
	if !ok {
		fmt.Fprintln(os.Stderr, "WARNING: logger not found in context, using loop")

		return &noopLogger{}
	}

	return l
}
