package logging

import (
	"context"
	"fmt"
)

type contextKey string

type ctxKey struct {
	key contextKey
}

const (
	contextKeyLogging = "logging_context"
)

func (l *ZapLogger) WithCtx(ctx context.Context) Logger {
	var args []interface{}
	if v := ctx.Value(ctxKey{contextKeyLogging}); v != nil {
		for k, v := range v.(map[string]interface{}) {
			args = append(args, k, v)
		}
	}
	if len(args) > 0 {
		return l.With(args...)
	}
	return l
}

func SetCtx(ctx context.Context, pairs ...interface{}) context.Context {
	var m map[string]interface{}
	if v := ctx.Value(ctxKey{contextKeyLogging}); v == nil {
		m = map[string]interface{}{}
		ctx = context.WithValue(ctx, ctxKey{contextKeyLogging}, m)
	} else {
		m = v.(map[string]interface{})
	}
	for i := 0; i < len(pairs); i += 2 {
		var key string
		switch k := pairs[i].(type) {
		case string:
			key = k
		case fmt.Stringer:
			key = k.String()
		default:
			panic("invalid api code data")
		}
		m[key] = pairs[i+1]
	}
	return ctx
}
