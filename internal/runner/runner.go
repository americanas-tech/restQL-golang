package runner

import (
	"context"
	"github.com/b2wdigital/restQL-golang/internal/domain"
)

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

func (r Runner) ExecuteQuery(ctx context.Context, query domain.Query, queryCtx QueryContext) interface{} {
	resources := NewResources(query.Statements)
	resources = ResolveVariables(resources, queryCtx.Input.Params)
	resources = MultiplexStatements(resources)

	state := NewState(resources)

	for !state.HasFinished() {
		availableResources := state.Available()
		for resourceId := range availableResources {
			state.SetAsRequest(resourceId)
		}

		availableResources = ResolveChainedValues(availableResources, state.Done())
		availableResources = MultiplexStatements(availableResources)

		for resourceId, statement := range availableResources {
			switch statement := statement.(type) {
			case domain.Statement:
				response := r.executor.DoStatement(ctx, statement, queryCtx)
				state.UpdateDone(resourceId, response)
			case []interface{}:
				responses := r.executor.DoMultiplexedStatement(ctx, statement, queryCtx)
				state.UpdateDone(resourceId, responses)
			}
		}
	}

	return state.Done()
}
