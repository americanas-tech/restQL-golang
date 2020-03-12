package runner_test

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/runner"
	"reflect"
	"testing"
)

func TestRequestsThatCanBeRequested(t *testing.T) {
	tests := []struct {
		name     string
		input    domain.Query
		expected []domain.Statement
	}{
		{
			"should return statement with static parameters",
			domain.Query{
				Statements: []domain.Statement{
					{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "123456"}},
				},
			},
			[]domain.Statement{
				{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "123456"}},
			},
		},
		{
			"should not return statement with chained parameter",
			domain.Query{
				Statements: []domain.Statement{
					{Method: "from", Resource: "hero", With: map[string]interface{}{"id": domain.Chain{"done-resource", "id"}}},
				},
			},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := runner.NewRunnerState(tt.input)
			got := state.GetAvailableRequests()

			if !reflect.DeepEqual(got, tt.expected) {
				t.Fatalf("GetAvailableRequests = %#+v, want = %#+v", got, tt.expected)
			}
		})
	}
}
