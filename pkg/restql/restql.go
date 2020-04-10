package restql

import "github.com/b2wdigital/restQL-golang/internal/domain"

type Plugin interface {
	BeforeQuery(query string, queryCtx QueryContext)
	AfterQuery(query string, result map[string]interface{})
	BeforeRequest(request HttpRequest)
	AfterRequest(response HttpResponse, err error)
}

type QueryInput = domain.QueryInput
type QueryOptions = domain.QueryOptions
type QueryContext = domain.QueryContext

type HttpRequest = domain.HttpRequest
type HttpResponse = domain.HttpResponse
