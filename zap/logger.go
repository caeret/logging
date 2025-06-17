package zap

import (
	"context"

	"go.uber.org/zap"

	"github.com/caeret/logging"
)

var _ logging.Logger = (*Logger)(nil)

func New(logger *zap.Logger) logging.Logger {
	return &Logger{
		logger: logger.Sugar(),
	}
}

type Logger struct {
	logger *zap.SugaredLogger
}

func (l *Logger) Debug(message string, fields ...interface{}) {
	l.logger.Debugw(message, fields...)
}

func (l *Logger) Info(message string, fields ...interface{}) {
	l.logger.Infow(message, fields...)
}

func (l *Logger) Warn(message string, fields ...interface{}) {
	l.logger.Warnw(message, fields...)
}

func (l *Logger) Error(message string, fields ...interface{}) {
	l.logger.Errorw(message, fields...)
}

func (l *Logger) With(fields ...interface{}) logging.Logger {
	return &Logger{logger: l.logger.With(fields...)}
}

func (l *Logger) WithCtx(ctx context.Context) logging.Logger {
	return logging.WithCtx(ctx, l)
}

func (l *Logger) Sync() error {
	return l.logger.Sync()
}

func (l *Logger) Zap() *zap.Logger {
	return l.logger.Desugar()
}

type Zap interface {
	Zap() *zap.Logger
}
