package domain

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

type EnvSource interface {
	GetString(key string) string
}

type Logger interface {
	Panic(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
	Error(msg string, err error, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
}

type HttpClient interface {
	Do(ctx context.Context, request HttpRequest) (HttpResponse, error)
}

var ErrRequestTimeout = errors.New("request timed out")

type Headers map[string]string
type Body interface{}

type HttpRequest struct {
	Schema  string
	Uri     string
	Query   map[string]interface{}
	Body    Body
	Headers Headers
}

type HttpResponse struct {
	Url        string
	StatusCode int
	Body       Body
	Headers    Headers
	Duration   time.Duration
}
