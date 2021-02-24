package runner_test

import (
	"testing"

	"github.com/b2wdigital/restQL-golang/v5/internal/domain"
	"github.com/b2wdigital/restQL-golang/v5/internal/runner"
	"github.com/b2wdigital/restQL-golang/v5/test"
)

func TestUnwrapNoMultiplex(t *testing.T) {
	tests := []struct {
		name     string
		input    domain.Resources
		expected domain.Resources
	}{
		{
			"should change nothing if there is no flatten param",
			domain.Resources{"hero": domain.Statement{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": 1}}}},
			domain.Resources{"hero": domain.Statement{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": 1}}}},
		},
		{
			"should unwrap flatten param in statement",
			domain.Resources{"hero": domain.Statement{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": domain.NoMultiplex{Value: 1}}}}},
			domain.Resources{"hero": domain.Statement{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": 1}}}},
		},
		{
			"should unwrap flatten body in statement",
			domain.Resources{"hero": domain.Statement{Method: "from", Resource: "hero", With: domain.Params{Body: domain.NoMultiplex{Value: map[string]interface{}{"id": 1}}}}},
			domain.Resources{"hero": domain.Statement{Method: "from", Resource: "hero", With: domain.Params{Body: map[string]interface{}{"id": 1}}}},
		},
		{
			"should unwrap flatten param in multiplexed statement",
			domain.Resources{"hero": []interface{}{
				[]interface{}{
					domain.Statement{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": domain.NoMultiplex{Value: 1}}}},
					domain.Statement{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": domain.NoMultiplex{Value: 2}}}},
				},
			}},
			domain.Resources{"hero": []interface{}{
				[]interface{}{
					domain.Statement{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": 1}}},
					domain.Statement{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": 2}}},
				},
			}},
		},
		{
			"should unwrap flatten param inside object",
			domain.Resources{"hero": domain.Statement{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"obj": map[string]interface{}{"id": domain.NoMultiplex{Value: 1}}}}}},
			domain.Resources{"hero": domain.Statement{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"obj": map[string]interface{}{"id": 1}}}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runner.UnwrapNoMultiplex(tt.input)
			test.Equal(t, got, tt.expected)
		})
	}
}
