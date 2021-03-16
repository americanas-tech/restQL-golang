package runner

import (
	"context"
	"sync"
	"time"

	"github.com/b2wdigital/restQL-golang/v5/internal/domain"
	"github.com/b2wdigital/restQL-golang/v5/pkg/restql"
	"github.com/pkg/errors"
)

var (
	// ErrQueryTimedOut represents the event of a query that
	// exceed the maximum execution time for it.
	ErrQueryTimedOut = errors.New("query timed out")

	// ErrMaxQueryDenied represents the event when restQL reaches the maximum
	// number of concurrent queries and cannot process anymore.
	ErrMaxQueryDenied = errors.New("max concurrent query reached: query execution denied")

	// ErrMaxGoroutineDenied represents the event when restQL reaches the maximum
	// number of concurrent goroutines and cannot process anymore.
	ErrMaxGoroutineDenied = errors.New("max concurrent goroutine reached: statement execution denied")
)

// Options wraps all configuration parameters for the Runner
type Options struct {
	GlobalQueryTimeout      time.Duration
	MaxConcurrentQueries    int
	MaxConcurrentGoroutines int
}

// Runner process a query into a Resource collection
// with the results in the most efficient way.
// All statements that can be executed in parallel,
// hence not having co-dependency, are done so.
type Runner struct {
	log              restql.Logger
	executor         Executor
	queryLimiter     *limiter
	goroutineLimiter *limiter
	options          Options
}

// NewRunner returns a Runner instance.
func NewRunner(log restql.Logger, executor Executor, options Options) Runner {
	return Runner{
		log:              log,
		executor:         executor,
		queryLimiter:     newLimiter(int32(options.MaxConcurrentQueries)),
		goroutineLimiter: newLimiter(int32(options.MaxConcurrentGoroutines)),
		options:          options,
	}
}

// ExecuteQuery process a query into a Resource collection.
func (r Runner) ExecuteQuery(ctx context.Context, query domain.Query, queryCtx restql.QueryContext) (domain.Resources, error) {
	log := restql.GetLogger(ctx)

	success := r.queryLimiter.Acquire()
	if !success {
		return nil, ErrMaxQueryDenied
	}
	defer r.queryLimiter.Release()

	var cancel context.CancelFunc
	queryTimeout, ok := r.parseQueryTimeout(query)
	if ok {
		ctx, cancel = context.WithTimeout(ctx, queryTimeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	defer cancel()

	resources, err := r.initializeResources(query)
	if err != nil {
		return nil, err
	}

	state := NewState(resources)

	requestCh := make(chan request, 10)
	resultCh := make(chan result, 10)
	outputCh := make(chan domain.Resources)
	errorCh := make(chan error)

	stateWorker := &stateWorker{
		log:              log,
		requestCh:        requestCh,
		resultCh:         resultCh,
		outputCh:         outputCh,
		errorCh:          errorCh,
		state:            state,
		ctx:              ctx,
		goroutineLimiter: r.goroutineLimiter,
	}

	requestWorker := &requestWorker{
		requestCh:        requestCh,
		resultCh:         resultCh,
		errorCh:          errorCh,
		executor:         r.executor,
		queryCtx:         queryCtx,
		ctx:              ctx,
		goroutineLimiter: r.goroutineLimiter,
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
		return r.options.GlobalQueryTimeout, false
	}

	duration, ok := timeout.(int)
	if !ok {
		return r.options.GlobalQueryTimeout, false
	}

	return time.Millisecond * time.Duration(duration), true
}

func (r Runner) initializeResources(query domain.Query) (domain.Resources, error) {
	resources := domain.NewResources(query.Statements)

	err := ValidateDependsOnTarget(resources)
	if err != nil {
		return nil, err
	}

	err = ValidateChainedValues(resources)
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
	log              restql.Logger
	requestCh        chan request
	resultCh         chan result
	outputCh         chan domain.Resources
	errorCh          chan error
	state            *State
	ctx              context.Context
	goroutineLimiter *limiter
}

func (sw *stateWorker) Run() {
	for !sw.state.HasFinished() {
		availableResources := sw.state.Available()
		for resourceID := range availableResources {
			sw.state.SetAsRequest(resourceID)
		}

		availableResources = ResolveChainedValues(availableResources, sw.state.Done())
		availableResources = ResolveDependsOn(availableResources, sw.state.Done())
		availableResources = ApplyEncoders(availableResources, sw.log)
		availableResources = MultiplexStatements(availableResources)
		availableResources = UnwrapNoMultiplex(availableResources)

		for resourceID, stmt := range availableResources {
			resourceID, stmt := resourceID, stmt
			success := sw.goroutineLimiter.Acquire()
			if !success {
				select {
				case sw.errorCh <- ErrMaxGoroutineDenied:
				case <-sw.ctx.Done():
				}

				return
			}

			go func() {
				select {
				case sw.requestCh <- request{ResourceIdentifier: resourceID, Statement: stmt}:
				case <-sw.ctx.Done():
				}

				sw.goroutineLimiter.Release()
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
	requestCh        chan request
	resultCh         chan result
	errorCh          chan error
	executor         Executor
	queryCtx         restql.QueryContext
	ctx              context.Context
	goroutineLimiter *limiter
}

func (rw *requestWorker) Run() {
	for {
		select {
		case req := <-rw.requestCh:
			resourceID := req.ResourceIdentifier
			statement := req.Statement

			success := rw.goroutineLimiter.Acquire()
			if !success {
				select {
				case rw.errorCh <- ErrMaxGoroutineDenied:
				case <-rw.ctx.Done():
				}

				return
			}

			switch statement := statement.(type) {
			case domain.Statement:
				go func() {
					response := rw.executor.DoStatement(rw.ctx, statement, rw.queryCtx)
					writeResult(rw.ctx, rw.resultCh, result{ResourceIdentifier: resourceID, Response: response})
					rw.goroutineLimiter.Release()
				}()
			case []interface{}:
				go func() {
					result := rw.runMultiplexedStatement(statement, resourceID)
					writeResult(rw.ctx, rw.resultCh, result)
					rw.goroutineLimiter.Release()
				}()
			default:
				// Should never reach this point
				rw.goroutineLimiter.Release()
			}
		case <-rw.ctx.Done():
			return
		}
	}
}

func (rw *requestWorker) runMultiplexedStatement(statements []interface{}, resourceID domain.ResourceID) result {
	responseChans := make([]chan interface{}, len(statements))
	for i := range responseChans {
		responseChans[i] = make(chan interface{}, 1)
	}

	var wg sync.WaitGroup

	wg.Add(len(statements))
	for i, stmt := range statements {
		select {
		case <-rw.ctx.Done():
			return result{}
		default:
		}

		i, stmt := i, stmt
		ch := responseChans[i]

		success := rw.goroutineLimiter.Acquire()
		if !success {
			select {
			case rw.errorCh <- ErrMaxGoroutineDenied:
			case <-rw.ctx.Done():
			}

			return result{}
		}

		switch stmt := stmt.(type) {
		case domain.Statement:
			go func() {
				response := rw.executor.DoStatement(rw.ctx, stmt, rw.queryCtx)
				ch <- response
				wg.Done()
				rw.goroutineLimiter.Release()
			}()
		case []interface{}:
			go func() {
				subResult := rw.runMultiplexedStatement(stmt, resourceID)
				ch <- subResult.Response
				wg.Done()
				rw.goroutineLimiter.Release()
			}()
		}
	}

	wg.Wait()
	responses := make(restql.DoneResources, len(statements))
	for i, ch := range responseChans {
		responses[i] = <-ch
	}

	return result{ResourceIdentifier: resourceID, Response: responses}
}

func writeResult(ctx context.Context, out chan result, r result) {
	select {
	case out <- r:
	case <-ctx.Done():
	}
}

type limiter struct {
	mu sync.Mutex

	limit  int32
	bucket int32
}

func newLimiter(limit int32) *limiter {
	return &limiter{limit: limit, bucket: limit}
}

// Acquire will try to reserve a token on limiter.
// If the limiter bucket is full (default case on select), it will fail.
func (l *limiter) Acquire() bool {
	if l.limit <= 0 {
		return true
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.bucket <= 0 {
		l.bucket = 0
		return false
	}

	l.bucket = l.bucket - 1
	return true
}

// Release will return a token to the bucket with a non-blocking read.
func (l *limiter) Release() {
	if l.limit <= 0 {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.bucket >= l.limit {
		l.bucket = l.limit
		return
	}

	l.bucket = l.bucket + 1
	return
}
