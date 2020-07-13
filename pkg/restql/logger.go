package restql

import (
	"context"
)

type Logger interface {
	Panic(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
	Error(msg string, err error, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
	With(key string, value interface{}) Logger
}

const loggerCtxKey = "__LOGGER__"

func WithLogger(ctx context.Context, l Logger) context.Context {
	if lp, ok := ctx.Value(loggerCtxKey).(Logger); ok {
		if lp == l {
			return ctx
		}
	}
	return context.WithValue(ctx, loggerCtxKey, l)
}

func GetLogger(ctx context.Context) Logger {
	if l, ok := ctx.Value(loggerCtxKey).(Logger); ok {
		return l
	}
	return noOpLogger{}
}

type noOpLogger struct{}

func (n noOpLogger) Panic(msg string, fields ...interface{})            {}
func (n noOpLogger) Fatal(msg string, fields ...interface{})            {}
func (n noOpLogger) Error(msg string, err error, fields ...interface{}) {}
func (n noOpLogger) Warn(msg string, fields ...interface{})             {}
func (n noOpLogger) Info(msg string, fields ...interface{})             {}
func (n noOpLogger) Debug(msg string, fields ...interface{})            {}
func (n noOpLogger) With(key string, value interface{}) Logger          { return n }
