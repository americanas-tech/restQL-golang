package domain

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

type EnvSource interface {
	GetString(key string) string
	GetAll() map[string]string
}

type HttpClient interface {
	Do(ctx context.Context, request HttpRequest) (HttpResponse, error)
}

var ErrRequestTimeout = errors.New("request timed out")

type Headers map[string]string
type Body interface{}

type HttpRequest struct {
	Method  string
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
