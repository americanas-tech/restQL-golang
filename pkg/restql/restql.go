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
	AfterTransaction(ctx context.Context, tr TransactionResponse)
	BeforeQuery(ctx context.Context, query string, queryCtx QueryContext)
	AfterQuery(ctx context.Context, query string, result map[string]interface{})
	BeforeRequest(ctx context.Context, request HttpRequest)
	AfterRequest(ctx context.Context, request HttpRequest, response HttpResponse, err error)
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
