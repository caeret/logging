package logging

import "context"

var _ Logger = (*nopLogger)(nil)

var NopLogger Logger = nopLogger{}

type nopLogger struct{}

func (n nopLogger) Debug(message string, fields ...interface{}) {

}

func (n nopLogger) Info(message string, fields ...interface{}) {

}

func (n nopLogger) Warn(message string, fields ...interface{}) {

}

func (n nopLogger) Error(message string, fields ...interface{}) {

}

func (n nopLogger) Skip(i int) Logger {
	return n
}

func (n nopLogger) WithCallerPKG(s string) Logger {
	return n
}

func (n nopLogger) WithSkipPKG(s ...string) Logger {
	return n
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