package slog

import (
	"context"
	"log/slog"

	"github.com/caeret/logging"
)

var _ logging.Logger = (*Logger)(nil)

func New(l *slog.Logger) logging.Logger {
	return &Logger{l: l}
}

type Logger struct {
	l *slog.Logger
}

func (l *Logger) Debug(msg string, args ...any) {
	l.l.Debug(msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.l.Info(msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.l.Warn(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.l.Error(msg, args...)
}

func (l *Logger) With(fields ...interface{}) logging.Logger {
	return &Logger{l: l.l.With(fields...)}
}

func (l *Logger) WithCtx(ctx context.Context) logging.Logger {
	return logging.WithCtx(ctx, l)
}

func (l *Logger) Sync() error {
	return nil
}
