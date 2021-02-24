package runner

import (
	"context"
	"time"

	"github.com/b2wdigital/restQL-golang/v5/internal/domain"
	"github.com/b2wdigital/restQL-golang/v5/pkg/restql"
	"github.com/pkg/errors"
)

// ErrQueryTimedOut represents the event of a query that
// exceed the maximum execution time for it.
var ErrQueryTimedOut = errors.New("query timed out")

// Runner process a query into a Resource collection
// with the results in the most efficient way.
// All statements that can be executed in parallel,
// hence not having co-dependency, are done so.
type Runner struct {
	log                restql.Logger
	executor           Executor
	globalQueryTimeout time.Duration
}

// NewRunner returns a Runner instance.
func NewRunner(log restql.Logger, executor Executor, globalQueryTimeout time.Duration) Runner {
	return Runner{
		log:                log,
		executor:           executor,
		globalQueryTimeout: globalQueryTimeout,
	}
}

// ExecuteQuery process a query into a Resource collection.
func (r Runner) ExecuteQuery(ctx context.Context, query domain.Query, queryCtx restql.QueryContext) (domain.Resources, error) {
	log := restql.GetLogger(ctx)

	var cancel context.CancelFunc
	queryTimeout, ok := r.parseQueryTimeout(query)
	if ok {
		ctx, cancel = context.WithTimeout(ctx, queryTimeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	defer cancel()

	resources, err := r.initializeResources(query, queryCtx)
	if err != nil {
		return nil, err
	}

	state := NewState(resources)

	requestCh := make(chan request, 10)
	resultCh := make(chan result, 10)
	outputCh := make(chan domain.Resources)
	errorCh := make(chan error)

	stateWorker := &stateWorker{
		log:       log,
		requestCh: requestCh,
		resultCh:  resultCh,
		outputCh:  outputCh,
		state:     state,
		ctx:       ctx,
	}

	requestWorker := &requestWorker{
		requestCh: requestCh,
		resultCh:  resultCh,
		errorCh:   errorCh,
		executor:  r.executor,
		queryCtx:  queryCtx,
		ctx:       ctx,
	}

	go stateWorker.Run()
	go requestWorker.Run()

	select {
	case output := <-outputCh:
		return output, nil
	case err := <-errorCh:
		log.Debug("an error occurred when running the query", "error", err)
		return nil, err
	case <-ctx.Done():
		log.Debug("query timed out")
		return nil, ErrQueryTimedOut
	}
}

func (r Runner) parseQueryTimeout(query domain.Query) (time.Duration, bool) {
	timeout, found := query.Use["timeout"]
	if !found {
		return r.globalQueryTimeout, false
	}

	duration, ok := timeout.(int)
	if !ok {
		return r.globalQueryTimeout, false
	}

	return time.Millisecond * time.Duration(duration), true
}

func (r Runner) initializeResources(query domain.Query, queryCtx restql.QueryContext) (domain.Resources, error) {
	resources := domain.NewResources(query.Statements)

	err := ValidateChainedValues(resources)
	if err != nil {
		return nil, err
	}

	resources = ApplyModifiers(resources, query.Use)
	resources = ApplyEncoders(resources, r.log)
	resources = MultiplexStatements(resources)

	return resources, nil
}

type request struct {
	ResourceIdentifier domain.ResourceID
	Statement          interface{}
}

type result struct {
	ResourceIdentifier domain.ResourceID
	Response           interface{}
}

type stateWorker struct {
	log       restql.Logger
	requestCh chan request
	resultCh  chan result
	outputCh  chan domain.Resources
	state     *State
	ctx       context.Context
}

func (sw *stateWorker) Run() {
	for !sw.state.HasFinished() {
		availableResources := sw.state.Available()
		for resourceID := range availableResources {
			sw.state.SetAsRequest(resourceID)
		}

		availableResources = ResolveChainedValues(availableResources, sw.state.Done())
		availableResources = ApplyEncoders(availableResources, sw.log)
		availableResources = MultiplexStatements(availableResources)
		availableResources = UnwrapNoMultiplex(availableResources)

		for resourceID, stmt := range availableResources {
			resourceID, stmt := resourceID, stmt
			go func() {
				select {
				case sw.requestCh <- request{ResourceIdentifier: resourceID, Statement: stmt}:
				case <-sw.ctx.Done():
				}
			}()
		}

		select {
		case result := <-sw.resultCh:
			sw.state.UpdateDone(result.ResourceIdentifier, result.Response)
		case <-sw.ctx.Done():
			return
		}
	}

	select {
	case sw.outputCh <- sw.state.Done():
	case <-sw.ctx.Done():
	}
}

type requestWorker struct {
	requestCh chan request
	resultCh  chan result
	errorCh   chan error
	executor  Executor
	queryCtx  restql.QueryContext
	ctx       context.Context
}

func (rw *requestWorker) Run() {
	for {
		select {
		case req := <-rw.requestCh:
			resourceID := req.ResourceIdentifier
			statement := req.Statement

			switch statement := statement.(type) {
			case domain.Statement:
				go func() {
					response := rw.executor.DoStatement(rw.ctx, statement, rw.queryCtx)
					writeResult(rw.ctx, rw.resultCh, result{ResourceIdentifier: resourceID, Response: response})
				}()
			case []interface{}:
				go func() {
					responses := rw.executor.DoMultiplexedStatement(rw.ctx, statement, rw.queryCtx)
					writeResult(rw.ctx, rw.resultCh, result{ResourceIdentifier: resourceID, Response: responses})
				}()
			}
		case <-rw.ctx.Done():
			return
		}
	}
}

func writeResult(ctx context.Context, out chan result, r result) {
	select {
	case out <- r:
	case <-ctx.Done():
	}
}
