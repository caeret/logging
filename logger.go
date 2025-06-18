package logging

import (
	"context"
	"sync/atomic"
)

var defaultLogger atomic.Value

func init() {
	defaultLogger.Store(&wrapper{Nop})
}

func SetDefault(logger Logger) {
	defaultLogger.Store(&wrapper{logger})
}

func Default() Logger {
	return defaultLogger.Load().(*wrapper).Logger
}

type wrapper struct {
	Logger
}

func Debug(message string, fields ...interface{}) {
	Default().Debug(message, fields...)
}

func Info(message string, fields ...interface{}) {
	Default().Info(message, fields...)
}

func Warn(message string, fields ...interface{}) {
	Default().Warn(message, fields...)
}

func Error(message string, fields ...interface{}) {
	Default().Error(message, fields...)
}

func With(fields ...interface{}) Logger {
	return Default().With(fields...)
}

type Logger interface {
	Debug(message string, fields ...interface{})
	Info(message string, fields ...interface{})
	Warn(message string, fields ...interface{})
	Error(message string, fields ...interface{})
	With(fields ...interface{}) Logger
	WithCtx(ctx context.Context) Logger
	Sync() error
}
