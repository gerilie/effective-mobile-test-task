package logger

import "go.uber.org/zap"

// With returns a new logger instance enriched with the provided fields.
func (l *logger) With(fields ...zap.Field) Logger {
	return &logger{
		l: l.l.With(fields...),
	}
}
