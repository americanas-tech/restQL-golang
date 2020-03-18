package runner

import (
	"context"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/pkg/errors"
	"time"
)

var ErrQueryTimedOut = errors.New("query timed out")

type Runner struct {
	config   domain.Configuration
	log      domain.Logger
	client   domain.HttpClient
	executor Executor
}

func NewRunner(config domain.Configuration, httpClient domain.HttpClient, log domain.Logger) Runner {
	return Runner{
		config:   config,
		log:      log,
		client:   httpClient,
		executor: Executor{client: httpClient, log: log},
	}
}

func (r Runner) ExecuteQuery(ctx context.Context, query domain.Query, queryCtx QueryContext) (interface{}, error) {
	queryTimeout := parseQueryTimeout(query)
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	resources := initializeResources(query, queryCtx)

	state := NewState(resources)

	requestCh := make(chan request)
	outputCh := make(chan Resources)
	resultCh := make(chan result)

	stateWorker := &stateWorker{
		requestCh: requestCh,
		resultCh:  resultCh,
		outputCh:  outputCh,
		ctx:       ctx,
		state:     state,
	}

	requestWorker := &requestWorker{
		requestCh: requestCh,
		resultCh:  resultCh,
		executor:  r.executor,
		ctx:       ctx,
		queryCtx:  queryCtx,
	}

	go stateWorker.Run()
	go requestWorker.Run()

	select {
	case output := <-outputCh:
		return output, nil
	case <-ctx.Done():
		r.log.Debug("query timed out")
		return nil, ErrQueryTimedOut
	}
}

func parseQueryTimeout(query domain.Query) time.Duration {
	defaultQueryTimeout := 5 * time.Second

	timeout, found := query.Use["timeout"]
	if !found {
		return defaultQueryTimeout
	}

	duration, ok := timeout.(int)
	if !ok {
		return defaultQueryTimeout
	}

	return time.Millisecond * time.Duration(duration)
}

func initializeResources(query domain.Query, queryCtx QueryContext) Resources {
	resources := NewResources(query.Statements)

	resources = ApplyModifiers(query.Use, resources)
	resources = ResolveVariables(resources, queryCtx.Input.Params)
	resources = MultiplexStatements(resources)

	return resources
}

type request struct {
	ResourceIdentifier ResourceId
	Statement          interface{}
}

type result struct {
	ResourceIdentifier ResourceId
	Response           interface{}
}

type stateWorker struct {
	requestCh chan request
	resultCh  chan result
	outputCh  chan Resources
	ctx       context.Context
	state     *State
}

func (sw *stateWorker) Run() {
	for !sw.state.HasFinished() {
		availableResources := sw.state.Available()
		for resourceId := range availableResources {
			sw.state.SetAsRequest(resourceId)
		}

		availableResources = ResolveChainedValues(availableResources, sw.state.Done())
		availableResources = MultiplexStatements(availableResources)

		for resourceId, stmt := range availableResources {
			resourceId, stmt := resourceId, stmt
			go func() {
				sw.requestCh <- request{ResourceIdentifier: resourceId, Statement: stmt}
			}()
		}

		select {
		case result := <-sw.resultCh:
			sw.state.UpdateDone(result.ResourceIdentifier, result.Response)
		case <-sw.ctx.Done():
			return
		}
	}

	sw.outputCh <- sw.state.Done()
}

type requestWorker struct {
	requestCh chan request
	resultCh  chan result
	executor  Executor
	ctx       context.Context
	queryCtx  QueryContext
}

func (rw *requestWorker) Run() {
	for {
		select {
		case req := <-rw.requestCh:
			resourceId := req.ResourceIdentifier
			statement := req.Statement

			switch statement := statement.(type) {
			case domain.Statement:
				go func() {
					response := rw.executor.DoStatement(rw.ctx, statement, rw.queryCtx)
					rw.resultCh <- result{ResourceIdentifier: resourceId, Response: response}
				}()
			case []interface{}:
				go func() {
					responses := rw.executor.DoMultiplexedStatement(rw.ctx, statement, rw.queryCtx)
					rw.resultCh <- result{ResourceIdentifier: resourceId, Response: responses}
				}()
			}
		case <-rw.ctx.Done():
			return
		}
	}
}
