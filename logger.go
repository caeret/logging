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
	defaultLogger.Store(NewDefault())
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

func With(fields ...interface{}) Logger {
	return Default().With(fields...)
}

func SetLogger(logger Logger) {
	defaultLogger.Store(logger)
}

type Config struct {
	Level  zapcore.Level
	Path   string
	MaxAge int
}

func NewRotator(conf Config) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:  conf.Path,
		Compress:  false,
		LocalTime: true,
		MaxAge:    conf.MaxAge,
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

func New(logger *lumberjack.Logger, level zap.AtomicLevel, callerPKG string, skipPKG ...string) Logger {
	encodeConf := zap.NewProductionEncoderConfig()
	encodeConf.EncodeTime = zapcore.ISO8601TimeEncoder
	encodeConf.TimeKey = "time"
	consoleConf := encodeConf
	consoleConf.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleConf.ConsoleSeparator = "  "
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encodeConf), zapcore.AddSync(logger), level),
		zapcore.NewCore(zapcore.NewConsoleEncoder(consoleConf), zapcore.AddSync(os.Stdout), level),
	)

	return &ZapLogger{
		logger: zap.New(core, zap.WithCaller(true), zap.AddCallerSkip(1), zap.WithCallerPKG(callerPKG), zap.WithSkipPKG(skipPKG...)).Sugar(),
	}
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

func (l *ZapLogger) With(fields ...interface{}) Logger {
	return &ZapLogger{logger: l.logger.With(fields...)}
}

func (l *ZapLogger) Sync() error {
	return l.logger.Sync()
}

func (l *ZapLogger) Zap() *zap.Logger {
	return l.logger.Desugar()
}
