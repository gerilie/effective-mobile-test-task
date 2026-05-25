package logger

import "go.uber.org/zap"

type noopLogger struct{}

func (l *noopLogger) Debug(_ string, _ ...zap.Field) {}
func (l *noopLogger) Info(_ string, _ ...zap.Field)  {}
func (l *noopLogger) Warn(_ string, _ ...zap.Field)  {}
func (l *noopLogger) Error(_ string, _ ...zap.Field) {}

func (l *noopLogger) Zap() *zap.Logger { return zap.NewNop() }
func (l *noopLogger) Stop() error      { return nil }
