package restql

import (
	"context"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"net/http"
	"net/url"
)

type Plugin interface {
	Name() string
	BeforeTransaction(ctx context.Context, tr TransactionRequest) context.Context
	AfterTransaction(ctx context.Context, tr TransactionResponse) context.Context
	BeforeQuery(ctx context.Context, query string, queryCtx QueryContext) context.Context
	AfterQuery(ctx context.Context, query string, result map[string]interface{}) context.Context
	BeforeRequest(ctx context.Context, request HttpRequest) context.Context
	AfterRequest(ctx context.Context, request HttpRequest, response HttpResponse, err error) context.Context
}

type TransactionRequest struct {
	Url    *url.URL
	Method string
	Header http.Header
}

type TransactionResponse struct {
	Status int
	Header []byte
	Body   []byte
}

type QueryInput = domain.QueryInput
type QueryOptions = domain.QueryOptions
type QueryContext = domain.QueryContext

type HttpRequest = domain.HttpRequest
type HttpResponse = domain.HttpResponse
