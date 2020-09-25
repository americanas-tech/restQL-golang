package domain

import (
	"context"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
	"github.com/pkg/errors"
)

// ErrRequestTimeout is the error returned by HTTPClient
// when a HTTP call fails due to the request exceeding
// the timeout defined in HTTPRequest.
var ErrRequestTimeout = errors.New("request timed out")

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
	Do(ctx context.Context, request restql.HTTPRequest) (restql.HTTPResponse, error)
}
