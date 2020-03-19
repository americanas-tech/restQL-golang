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

func (r Runner) ExecuteQuery(ctx context.Context, query domain.Query, queryCtx domain.QueryContext) (domain.Resources, error) {
	queryTimeout := parseQueryTimeout(query)
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	resources := initializeResources(query, queryCtx)

	state := NewState(resources)

	requestCh := make(chan request)
	outputCh := make(chan domain.Resources)
	errorCh := make(chan error)
	resultCh := make(chan result)

	stateWorker := &stateWorker{
		requestCh: requestCh,
		resultCh:  resultCh,
		outputCh:  outputCh,
		state:     state,
		ctx:       ctx,
	}

	requestWorker := &requestWorker{
		requestCh: requestCh,
		resultCh:  resultCh,
		executor:  r.executor,
		errorCh:   errorCh,
		queryCtx:  queryCtx,
		ctx:       ctx,
	}

	go stateWorker.Run()
	go requestWorker.Run()

	select {
	case output := <-outputCh:
		return output, nil
	case err := <-errorCh:
		cancel()
		return nil, err
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

func initializeResources(query domain.Query, queryCtx domain.QueryContext) domain.Resources {
	resources := domain.NewResources(query.Statements)

	resources = ApplyModifiers(query.Use, resources)
	resources = ResolveVariables(resources, queryCtx.Input.Params)
	resources = MultiplexStatements(resources)

	return resources
}

type request struct {
	ResourceIdentifier domain.ResourceId
	Statement          interface{}
}

type result struct {
	ResourceIdentifier domain.ResourceId
	Response           interface{}
}

type stateWorker struct {
	requestCh chan request
	resultCh  chan result
	outputCh  chan domain.Resources
	state     *State
	ctx       context.Context
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
	errorCh   chan error
	executor  Executor
	queryCtx  domain.QueryContext
	ctx       context.Context
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
					response, err := rw.executor.DoStatement(rw.ctx, statement, rw.queryCtx)
					if err != nil {
						rw.errorCh <- err
					}
					rw.resultCh <- result{ResourceIdentifier: resourceId, Response: response}
				}()
			case []interface{}:
				go func() {
					responses, err := rw.executor.DoMultiplexedStatement(rw.ctx, statement, rw.queryCtx)
					if err != nil {
						rw.errorCh <- err
					}
					rw.resultCh <- result{ResourceIdentifier: resourceId, Response: responses}
				}()
			}
		case <-rw.ctx.Done():
			return
		}
	}
}
