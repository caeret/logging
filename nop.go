package logging

import (
	"context"

	"go.uber.org/zap"
)

var _ Logger = (*nopLogger)(nil)

var Nop Logger = nopLogger{}

type nopLogger struct{}

func (n nopLogger) Debug(message string, fields ...interface{}) {
}

func (n nopLogger) Info(message string, fields ...interface{}) {
}

func (n nopLogger) Warn(message string, fields ...interface{}) {
}

func (n nopLogger) Error(message string, fields ...interface{}) {
}

func (n nopLogger) With(fields ...interface{}) Logger {
	return n
}

func (n nopLogger) WithCtx(ctx context.Context) Logger {
	return n
}

func (n nopLogger) Sync() error {
	return nil
}

func (n nopLogger) Zap() *zap.Logger {
	return zap.NewNop()
}
