package logging

import (
	"context"
	"os"
	"sync/atomic"

	"github.com/caeret/zap"
	"github.com/caeret/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var defaultLogger atomic.Value

func init() {
	defaultLogger.Store(NewDefault().Skip(1))
}

func Default() Logger {
	return defaultLogger.Load().(Logger)
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

func SetLogger(logger Logger) {
	defaultLogger.Store(logger.Skip(1))
}

type Config struct {
	Level zapcore.Level
	Path  string
}

func NewRotator(conf Config) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   conf.Path,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
		LocalTime:  true,
	}
}

func NewLoggerLevel(conf Config) zap.AtomicLevel {
	return zap.NewAtomicLevelAt(conf.Level)
}

func NewDefault() Logger {
	encodeConf := zap.NewProductionEncoderConfig()
	encodeConf.EncodeTime = zapcore.ISO8601TimeEncoder
	encodeConf.TimeKey = "time"
	encodeConf.EncodeLevel = zapcore.CapitalColorLevelEncoder
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encodeConf), zapcore.AddSync(os.Stdout), zap.NewAtomicLevelAt(zap.DebugLevel))
	return &ZapLogger{
		logger: zap.New(core, zap.WithCaller(true), zap.AddCallerSkip(1)).Sugar(),
	}
}

func New(logger *lumberjack.Logger, level zap.AtomicLevel) Logger {
	encodeConf := zap.NewProductionEncoderConfig()
	encodeConf.EncodeTime = zapcore.ISO8601TimeEncoder
	encodeConf.TimeKey = "time"
	consoleConf := encodeConf
	consoleConf.EncodeLevel = zapcore.CapitalColorLevelEncoder
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encodeConf), zapcore.AddSync(logger), level),
		zapcore.NewCore(zapcore.NewConsoleEncoder(consoleConf), zapcore.AddSync(os.Stdout), level),
	)

	return &ZapLogger{
		logger: zap.New(core, zap.WithCaller(true), zap.AddCallerSkip(1)).Sugar(),
	}
}

type Logger interface {
	Debug(message string, fields ...interface{})
	Info(message string, fields ...interface{})
	Warn(message string, fields ...interface{})
	Error(message string, fields ...interface{})
	Skip(int) Logger
	WithCallerPKG(string) Logger
	WithSkipPKG(...string) Logger
	With(fields ...interface{}) Logger
	WithCtx(ctx context.Context) Logger
	Sync() error
}

type ZapLogger struct {
	logger *zap.SugaredLogger
}

func (l *ZapLogger) Debug(message string, fields ...interface{}) {
	l.logger.Debugw(message, fields...)
}

func (l *ZapLogger) Info(message string, fields ...interface{}) {
	l.logger.Infow(message, fields...)
}

func (l *ZapLogger) Warn(message string, fields ...interface{}) {
	l.logger.Warnw(message, fields...)
}

func (l *ZapLogger) Error(message string, fields ...interface{}) {
	l.logger.Errorw(message, fields...)
}

func (l *ZapLogger) Skip(n int) Logger {
	return &ZapLogger{logger: l.logger.Desugar().WithOptions(zap.AddCallerSkip(n)).Sugar()}
}

func (l *ZapLogger) WithCallerPKG(pkg string) Logger {
	return &ZapLogger{logger: l.logger.Desugar().WithOptions(zap.WithCallerPKG(pkg)).Sugar()}
}

func (l *ZapLogger) WithSkipPKG(pkg ...string) Logger {
	return &ZapLogger{logger: l.logger.Desugar().WithOptions(zap.WithSkipPKG(pkg...)).Sugar()}
}

func (l *ZapLogger) With(fields ...interface{}) Logger {
	return &ZapLogger{logger: l.logger.With(fields...)}
}

func (l *ZapLogger) Sync() error {
	return l.logger.Sync()
}
