package slog

import (
	"context"
	"log/slog"

	"github.com/caeret/logging"
)

var _ logging.Logger = (*logger)(nil)

func New(l *slog.Logger) logging.Logger {
	return &logger{l: l}
}

type logger struct {
	l *slog.Logger
}

func (l *logger) Debug(msg string, args ...any) {
	l.l.Debug(msg, args...)
}

func (l *logger) Info(msg string, args ...any) {
	l.l.Info(msg, args...)
}

func (l *logger) Warn(msg string, args ...any) {
	l.l.Warn(msg, args...)
}

func (l *logger) Error(msg string, args ...any) {
	l.l.Error(msg, args...)
}

func (l *logger) With(fields ...interface{}) logging.Logger {
	return &logger{l: l.l.With(fields...)}
}

func (l *logger) WithCtx(ctx context.Context) logging.Logger {
	return logging.WithCtx(ctx, l)
}

func (l *logger) Sync() error {
	return nil
}
