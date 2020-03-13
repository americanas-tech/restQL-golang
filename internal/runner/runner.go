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

func (r Runner) ExecuteQuery(ctx context.Context, query domain.Query, mappings map[string]domain.Mapping) interface{} {
	executor := Executor{mappings: mappings, client: r.client, log: r.log}

	state := NewState(query)

	for !state.HasFinished() {
		availableResources := state.Available()

		for resourceId := range availableResources {
			state.SetAsRequest(resourceId)
		}

		for resourceId, statement := range availableResources {
			switch statement := statement.(type) {
			case domain.Statement:
				response := executor.DoStatement(ctx, statement)
				state.UpdateDone(resourceId, DoneRequest(response))
			case []domain.Statement:
				responses := executor.DoMultiplexedStatement(ctx, statement)
				state.UpdateDone(resourceId, DoneRequests(responses))
			}
		}
	}

	return state.Done()
}
