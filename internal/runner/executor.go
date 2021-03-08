package runner

import (
	"context"
	"time"

	"github.com/b2wdigital/restQL-golang/v5/internal/domain"
	"github.com/b2wdigital/restQL-golang/v5/pkg/restql"
)

// Executor process statements into a result
// by executing the relevant HTTP calls to
// the upstream dependency.
type Executor struct {
	client          domain.HTTPClient
	log             restql.Logger
	resourceTimeout time.Duration
	forwardPrefix   string
}

// NewExecutor constructs an instance of Executor.
func NewExecutor(log restql.Logger, client domain.HTTPClient, resourceTimeout time.Duration, forwardPrefix string) Executor {
	return Executor{client: client, log: log, resourceTimeout: resourceTimeout, forwardPrefix: forwardPrefix}
}

// DoStatement process a single statement into a result by executing the relevant HTTP calls to the upstream dependency.
func (e Executor) DoStatement(ctx context.Context, statement domain.Statement, queryCtx restql.QueryContext) restql.DoneResource {
	log := restql.GetLogger(ctx)

	drOptions := DoneResourceOptions{
		IgnoreErrors: statement.IgnoreErrors,
		MaxAge:       statement.CacheControl.MaxAge,
		SMaxAge:      statement.CacheControl.SMaxAge,
	}

	if !statement.DependsOn.Resolved {
		failedDependsOnResponse := NewNewDependsOnUnresolvedResponse(log, statement, drOptions)
		log.Debug("request execution skipped due to unresolved dependency", "resource", statement.Resource, "method", statement.Method)
		return failedDependsOnResponse
	}

	emptyChainedParams := GetEmptyChainedParams(statement)
	if len(emptyChainedParams) > 0 {
		emptyChainedResponse := NewEmptyChainedResponse(log, emptyChainedParams, drOptions)
		log.Debug("request execution skipped due to empty chained parameters", "resource", statement.Resource, "method", statement.Method)
		return emptyChainedResponse
	}

	request := MakeRequest(e.resourceTimeout, e.forwardPrefix, statement, queryCtx)

	log.Debug("executing request for statement", "resource", statement.Resource, "method", statement.Method, "request", request)

	response, err := e.client.Do(ctx, request)
	if err != nil {
		errorResponse := NewErrorResponse(log, err, request, response, drOptions)
		log.Debug("request execution failed", "error", err, "resource", statement.Resource, "method", statement.Method, "response", errorResponse)
		return errorResponse
	}

	dr := NewDoneResource(request, response, drOptions)

	log.Debug("request execution done", "resource", statement.Resource, "method", statement.Method, "response", dr)

	return dr
}
