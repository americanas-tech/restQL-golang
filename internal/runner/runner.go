package runner

import (
	"context"
	"github.com/b2wdigital/restQL-golang/internal/domain"
)

type Runner struct {
	config domain.Configuration
	log    domain.Logger
	client domain.HttpClient
}

func NewRunner(config domain.Configuration, httpClient domain.HttpClient, log domain.Logger) Runner {
	return Runner{
		config: config,
		log:    log,
		client: httpClient,
	}
}

func (r Runner) ExecuteQuery(ctx context.Context, query domain.Query, queryCtx QueryContext) interface{} {
	executor := Executor{mappings: queryCtx.Mappings, client: r.client, log: r.log}

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
				response := executor.DoStatement(ctx, statement)
				state.UpdateDone(resourceId, response)
			case []interface{}:
				responses := executor.DoMultiplexedStatement(ctx, statement)
				state.UpdateDone(resourceId, responses)
			}
		}
	}

	return state.Done()
}
