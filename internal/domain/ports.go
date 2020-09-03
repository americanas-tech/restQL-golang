package domain

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

// ErrRequestTimeout is the error returned by HTTPClient
// when a HTTP call fails due to the request exceeding
// the timeout defined in HTTPRequest.
var ErrRequestTimeout = errors.New("request timed out")

// ErrMappingsNotFound is the error returned when
// the resource mappings is not found anywhere
var ErrMappingsNotFound = errors.New("mappings not found")

// ErrQueryNotFound is the error returned when
// the query text is not found anywhere
var ErrQueryNotFound = errors.New("query not found")

// EnvSource expose access to environment variables.
type EnvSource interface {
	GetString(key string) string
	GetAll() map[string]string
}

// HTTPClient is the interface that wrap the method Do
//
// Do takes an HTTPRequest and execute it respecting
// the cancellation signal from the given Context.
type HTTPClient interface {
	Do(ctx context.Context, request HTTPRequest) (HTTPResponse, error)
}

// Headers represents all HTTP header in a request or response.
type Headers map[string]string

// Body represents a HTTP body in a request or response.
type Body interface{}

// HTTPRequest describe a HTTP call to be made by HTTPClient.
type HTTPRequest struct {
	Method  string
	Schema  string
	Host    string
	Path    string
	Query   map[string]interface{}
	Body    Body
	Headers Headers
	Timeout time.Duration
}

// HTTPResponse describe a HTTP response returned by HTTPClient.
type HTTPResponse struct {
	URL        string
	StatusCode int
	Body       Body
	Headers    Headers
	Duration   time.Duration
}
