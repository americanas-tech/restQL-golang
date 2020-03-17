package domain

import "context"

type EnvSource interface {
	GetString(key string) string
}

type FileSource interface {
	Unmarshal(target interface{}) error
}

type Configuration interface {
	Env() EnvSource
	File() FileSource
	Build() string
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
	Do(ctx context.Context, request Request) (Response, error)
}

type Headers map[string]string
type Body interface{}

type Request struct {
	Schema  string
	Uri     string
	Query   map[string]string
	Body    Body
	Headers Headers
}

type Response struct {
	StatusCode int
	Body       Body
	Headers    Headers
}
