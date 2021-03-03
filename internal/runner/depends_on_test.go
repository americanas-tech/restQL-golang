package runner_test

import (
	"fmt"
	"github.com/b2wdigital/restQL-golang/v5/internal/domain"
	"github.com/b2wdigital/restQL-golang/v5/internal/runner"
	"github.com/b2wdigital/restQL-golang/v5/pkg/restql"
	"github.com/b2wdigital/restQL-golang/v5/test"
	"testing"
)

func TestResolveDependsOn(t *testing.T) {
	tests := []struct {
		name                string
		expected            domain.Resources
		statementUnresolved domain.Resources
		doneResources       domain.Resources
	}{
		{
			"Should set all resources to resolved if there is no depends on",
			domain.Resources{
				"hero":     domain.Statement{Method: "from", Resource: "hero", DependsOn: domain.DependsOn{Resolved: true}, With: domain.Params{Values: map[string]interface{}{"id": "abcdef12345"}}},
				"sidekick": domain.Statement{Method: "from", Resource: "sidekick", DependsOn: domain.DependsOn{Resolved: true}, With: domain.Params{Values: map[string]interface{}{"id": "zyxv98765"}}},
			},
			domain.Resources{
				"hero":     domain.Statement{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": "abcdef12345"}}},
				"sidekick": domain.Statement{Method: "from", Resource: "sidekick", With: domain.Params{Values: map[string]interface{}{"id": "zyxv98765"}}},
			},
			domain.Resources{"done-resource": restql.DoneResource{Status: 200, ResponseBody: restql.NewResponseBodyFromValue(test.NoOpLogger, test.Unmarshal(`{"id": "abcdef12345"}`))}},
		},
		{
			"Returns a statement with depends on resolved if done-resource status code is in 200 >= status >= 399",
			domain.Resources{
				"hero":     domain.Statement{Method: "from", Resource: "hero", DependsOn: domain.DependsOn{Resolved: true}, With: domain.Params{Values: map[string]interface{}{"id": "abcdef12345"}}},
				"sidekick": domain.Statement{Method: "from", Resource: "sidekick", DependsOn: domain.DependsOn{Target: "hero", Resolved: true}, With: domain.Params{Values: map[string]interface{}{"id": "zyxv98765"}}},
			},
			domain.Resources{
				"hero":     domain.Statement{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": "abcdef12345"}}},
				"sidekick": domain.Statement{Method: "from", Resource: "sidekick", DependsOn: domain.DependsOn{Target: "hero"}, With: domain.Params{Values: map[string]interface{}{"id": "zyxv98765"}}},
			},
			domain.Resources{"hero": restql.DoneResource{Status: 201, ResponseBody: restql.NewResponseBodyFromValue(test.NoOpLogger, test.Unmarshal("{}"))}},
		},
		{
			"Returns a statement with depends on not resolved if done-resource status code is not successful",
			domain.Resources{
				"hero":     domain.Statement{Method: "from", Resource: "hero", DependsOn: domain.DependsOn{Resolved: true}, With: domain.Params{Values: map[string]interface{}{"id": "abcdef12345"}}},
				"sidekick": domain.Statement{Method: "from", Resource: "sidekick", DependsOn: domain.DependsOn{Target: "hero", Resolved: false}, With: domain.Params{Values: map[string]interface{}{"id": "zyxv98765"}}},
			},
			domain.Resources{
				"hero":     domain.Statement{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": "abcdef12345"}}},
				"sidekick": domain.Statement{Method: "from", Resource: "sidekick", DependsOn: domain.DependsOn{Target: "hero"}, With: domain.Params{Values: map[string]interface{}{"id": "zyxv98765"}}},
			},
			domain.Resources{"hero": restql.DoneResource{Status: 408, ResponseBody: restql.NewResponseBodyFromValue(test.NoOpLogger, test.Unmarshal("{}"))}},
		},
		{
			"Returns a statement with depends on resolved if done-resource status code is not successful but has ignore-errors",
			domain.Resources{
				"hero":     domain.Statement{Method: "from", Resource: "hero", DependsOn: domain.DependsOn{Resolved: true}, With: domain.Params{Values: map[string]interface{}{"id": "abcdef12345"}}},
				"sidekick": domain.Statement{Method: "from", Resource: "sidekick", DependsOn: domain.DependsOn{Target: "hero", Resolved: true}, With: domain.Params{Values: map[string]interface{}{"id": "zyxv98765"}}},
			},
			domain.Resources{
				"hero":     domain.Statement{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": "abcdef12345"}}},
				"sidekick": domain.Statement{Method: "from", Resource: "sidekick", DependsOn: domain.DependsOn{Target: "hero"}, With: domain.Params{Values: map[string]interface{}{"id": "zyxv98765"}}},
			},
			domain.Resources{"hero": restql.DoneResource{Status: 408, IgnoreErrors: true, ResponseBody: restql.NewResponseBodyFromValue(test.NoOpLogger, test.Unmarshal("{}"))}},
		},
		{
			"Returns a statements with depends on resolved if at least one of the multiplexed done-resource status code is in 200 >= status >= 399",
			domain.Resources{
				"hero":     domain.Statement{Method: "from", Resource: "hero", DependsOn: domain.DependsOn{Resolved: true}, With: domain.Params{Values: map[string]interface{}{"id": "abcdef12345"}}},
				"sidekick": domain.Statement{Method: "from", Resource: "sidekick", DependsOn: domain.DependsOn{Target: "hero", Resolved: true}, With: domain.Params{Values: map[string]interface{}{"id": "zyxv98765"}}},
			},
			domain.Resources{
				"hero":     domain.Statement{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": "abcdef12345"}}},
				"sidekick": domain.Statement{Method: "from", Resource: "sidekick", DependsOn: domain.DependsOn{Target: "hero"}, With: domain.Params{Values: map[string]interface{}{"id": "zyxv98765"}}},
			},
			domain.Resources{"hero": restql.DoneResources{
				restql.DoneResource{Status: 500, ResponseBody: restql.NewResponseBodyFromValue(test.NoOpLogger, test.Unmarshal("{}"))},
				restql.DoneResource{Status: 201, ResponseBody: restql.NewResponseBodyFromValue(test.NoOpLogger, test.Unmarshal("{}"))},
				restql.DoneResource{Status: 408, ResponseBody: restql.NewResponseBodyFromValue(test.NoOpLogger, test.Unmarshal("{}"))},
			}},
		},
		{
			"Returns multiplexed statement with depends on resolved if done-resource status code is in 200 >= status >= 399",
			domain.Resources{
				"hero": domain.Statement{Method: "from", Resource: "hero", DependsOn: domain.DependsOn{Resolved: true}, With: domain.Params{Values: map[string]interface{}{"id": "abcdef12345"}}},
				"sidekick": []interface{}{
					domain.Statement{Method: "from", Resource: "sidekick", DependsOn: domain.DependsOn{Target: "hero", Resolved: true}, With: domain.Params{Values: map[string]interface{}{"id": "zyxv98765"}}},
					domain.Statement{Method: "from", Resource: "sidekick", DependsOn: domain.DependsOn{Target: "hero", Resolved: true}, With: domain.Params{Values: map[string]interface{}{"id": "5678xpot"}}},
				},
			},
			domain.Resources{
				"hero": domain.Statement{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": "abcdef12345"}}},
				"sidekick": []interface{}{
					domain.Statement{Method: "from", Resource: "sidekick", DependsOn: domain.DependsOn{Target: "hero"}, With: domain.Params{Values: map[string]interface{}{"id": "zyxv98765"}}},
					domain.Statement{Method: "from", Resource: "sidekick", DependsOn: domain.DependsOn{Target: "hero"}, With: domain.Params{Values: map[string]interface{}{"id": "5678xpot"}}},
				},
			},
			domain.Resources{"hero": restql.DoneResource{Status: 201, ResponseBody: restql.NewResponseBodyFromValue(test.NoOpLogger, test.Unmarshal("{}"))}},
		},
		{
			"Returns multiplexed statement with depends on resolved if multiplexed done-resource status code is in 200 >= status >= 399",
			domain.Resources{
				"hero": domain.Statement{Method: "from", Resource: "hero", DependsOn: domain.DependsOn{Resolved: true}, With: domain.Params{Values: map[string]interface{}{"id": "abcdef12345"}}},
				"sidekick": []interface{}{
					domain.Statement{Method: "from", Resource: "sidekick", DependsOn: domain.DependsOn{Target: "hero", Resolved: true}, With: domain.Params{Values: map[string]interface{}{"id": "zyxv98765"}}},
					domain.Statement{Method: "from", Resource: "sidekick", DependsOn: domain.DependsOn{Target: "hero", Resolved: true}, With: domain.Params{Values: map[string]interface{}{"id": "5678xpot"}}},
				},
			},
			domain.Resources{
				"hero": domain.Statement{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": "abcdef12345"}}},
				"sidekick": []interface{}{
					domain.Statement{Method: "from", Resource: "sidekick", DependsOn: domain.DependsOn{Target: "hero"}, With: domain.Params{Values: map[string]interface{}{"id": "zyxv98765"}}},
					domain.Statement{Method: "from", Resource: "sidekick", DependsOn: domain.DependsOn{Target: "hero"}, With: domain.Params{Values: map[string]interface{}{"id": "5678xpot"}}},
				},
			},
			domain.Resources{"hero": restql.DoneResources{
				restql.DoneResource{Status: 500, ResponseBody: restql.NewResponseBodyFromValue(test.NoOpLogger, test.Unmarshal("{}"))},
				restql.DoneResource{Status: 201, ResponseBody: restql.NewResponseBodyFromValue(test.NoOpLogger, test.Unmarshal("{}"))},
				restql.DoneResource{Status: 408, ResponseBody: restql.NewResponseBodyFromValue(test.NoOpLogger, test.Unmarshal("{}"))},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runner.ResolveDependsOn(tt.statementUnresolved, tt.doneResources)
			test.Equal(t, got, tt.expected)
		})
	}
}

func TestValidateDependsOnTarget(t *testing.T) {
	tests := []struct {
		name      string
		expected  error
		resources domain.Resources
	}{
		{
			"Should do nothing if there is no depends on",
			nil,
			domain.Resources{
				"resource-name": domain.Statement{
					Method:   "from",
					Resource: "resource-name",
					With:     domain.Params{Values: map[string]interface{}{"id": "abcdef12345"}},
				},
			},
		},
		{
			"Fail validation if depends on target unknown resource",
			fmt.Errorf("%w: unknown", runner.ErrInvalidDependsOnTarget),
			domain.Resources{
				"resource-name": domain.Statement{
					Method:    "from",
					Resource:  "resource-name",
					DependsOn: domain.DependsOn{Target: "unknown"},
					With:      domain.Params{Values: map[string]interface{}{"id": "abcdef12345"}},
				},
			},
		},
		{
			"Pass validation if depends on target known resource",
			nil,
			domain.Resources{
				"known-resource": domain.Statement{
					Method:   "from",
					Resource: "known-resource",
					With:     domain.Params{Values: map[string]interface{}{"id": "abcdef12345"}},
				},
				"resource-name": domain.Statement{
					Method:    "from",
					Resource:  "resource-name",
					DependsOn: domain.DependsOn{Target: "known-resource"},
					With:      domain.Params{Values: map[string]interface{}{"id": "abcdef12345"}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runner.ValidateDependsOnTarget(tt.resources)
			test.Equal(t, fmt.Sprintf("%s", got), fmt.Sprintf("%s", tt.expected))
		})
	}
}
