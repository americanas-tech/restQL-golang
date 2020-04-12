package restql

import (
	"context"
	"github.com/b2wdigital/restQL-golang/internal/domain"
)

type Plugin interface {
	Name() string
	BeforeQuery(ctx context.Context, query string, queryCtx QueryContext)
	AfterQuery(ctx context.Context, query string, result map[string]interface{})
	BeforeRequest(ctx context.Context, request HttpRequest)
	AfterRequest(ctx context.Context, request HttpRequest, response HttpResponse, err error)
}

type QueryInput = domain.QueryInput
type QueryOptions = domain.QueryOptions
type QueryContext = domain.QueryContext

type HttpRequest = domain.HttpRequest
type HttpResponse = domain.HttpResponse
