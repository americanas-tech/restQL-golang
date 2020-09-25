package runner

import (
	"context"
	"sync"
	"time"

	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
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
func (e Executor) DoStatement(ctx context.Context, statement domain.Statement, queryCtx restql.QueryContext) domain.DoneResource {
	log := restql.GetLogger(ctx)

	drOptions := DoneResourceOptions{
		IgnoreErrors: statement.IgnoreErrors,
		MaxAge:       statement.CacheControl.MaxAge,
		SMaxAge:      statement.CacheControl.SMaxAge,
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

// DoMultiplexedStatement process multiplexed statements into a result by executing the relevant HTTP calls to the upstream dependency.
func (e Executor) DoMultiplexedStatement(ctx context.Context, statements []interface{}, queryCtx restql.QueryContext) domain.DoneResources {
	responseChans := make([]chan interface{}, len(statements))
	for i := range responseChans {
		responseChans[i] = make(chan interface{}, 1)
	}

	var wg sync.WaitGroup

	wg.Add(len(statements))
	for i, stmt := range statements {
		i, stmt := i, stmt
		ch := responseChans[i]

		go func() {
			response := e.doCurrentStatement(ctx, stmt, queryCtx)
			ch <- response
			wg.Done()
		}()
	}

	wg.Wait()

	responses := make(domain.DoneResources, len(statements))
	for i, ch := range responseChans {
		responses[i] = <-ch
	}

	return responses
}

func (e Executor) doCurrentStatement(ctx context.Context, stmt interface{}, queryCtx restql.QueryContext) interface{} {
	switch stmt := stmt.(type) {
	case domain.Statement:
		return e.DoStatement(ctx, stmt, queryCtx)
	case []interface{}:
		return e.DoMultiplexedStatement(ctx, stmt, queryCtx)
	default:
		return nil
	}
}
