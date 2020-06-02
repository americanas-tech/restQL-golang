package runner_test

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/runner"
	"github.com/b2wdigital/restQL-golang/test"
	"testing"
)

func TestUnwrapFlatten(t *testing.T) {
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
			domain.Resources{"hero": domain.Statement{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": domain.Flatten{Target: 1}}}}},
			domain.Resources{"hero": domain.Statement{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": 1}}}},
		},
		{
			"should unwrap flatten body in statement",
			domain.Resources{"hero": domain.Statement{Method: "from", Resource: "hero", With: domain.Params{Body: domain.Flatten{Target: map[string]interface{}{"id": 1}}}}},
			domain.Resources{"hero": domain.Statement{Method: "from", Resource: "hero", With: domain.Params{Body: map[string]interface{}{"id": 1}}}},
		},
		{
			"should unwrap flatten param in multiplexed statement",
			domain.Resources{"hero": []interface{}{
				[]interface{}{
					domain.Statement{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": domain.Flatten{Target: 1}}}},
					domain.Statement{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": domain.Flatten{Target: 2}}}},
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
			domain.Resources{"hero": domain.Statement{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"obj": map[string]interface{}{"id": domain.Flatten{Target: 1}}}}}},
			domain.Resources{"hero": domain.Statement{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"obj": map[string]interface{}{"id": 1}}}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runner.UnwrapFlatten(tt.input)
			test.Equal(t, got, tt.expected)
		})
	}
}
