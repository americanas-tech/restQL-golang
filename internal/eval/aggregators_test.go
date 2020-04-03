package eval_test

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/eval"
	"reflect"
	"testing"
)

func TestApplyAggregators(t *testing.T) {
	tests := []struct {
		name      string
		query     domain.Query
		resources domain.Resources
		expected  domain.Resources
	}{
		{
			"should do nothing if there is no aggregator",
			domain.Query{Statements: []domain.Statement{{Resource: "hero"}, {Resource: "sidekick"}}},
			domain.Resources{
				"hero":     domain.DoneResource{Details: domain.Details{Success: true}, Result: nil},
				"sidekick": domain.DoneResource{Details: domain.Details{Success: true}, Result: nil},
			},
			domain.Resources{
				"hero":     domain.DoneResource{Details: domain.Details{Success: true}, Result: nil},
				"sidekick": domain.DoneResource{Details: domain.Details{Success: true}, Result: nil},
			},
		},
		{
			"should aggregate one resource inside other",
			domain.Query{Statements: []domain.Statement{
				{Resource: "hero"},
				{Resource: "sidekick", In: []string{"hero", "sidekick"}},
			}},
			domain.Resources{
				"hero":     domain.DoneResource{Details: domain.Details{Success: true}, Result: unmarshal(`{ "id": 1, "name": "batman" }`)},
				"sidekick": domain.DoneResource{Details: domain.Details{Success: true}, Result: unmarshal(`{ "id": 10, "name": "robin" }`)},
			},
			domain.Resources{
				"hero":     domain.DoneResource{Details: domain.Details{Success: true}, Result: unmarshal(`{ "id": 1, "name": "batman", "sidekick": { "id": 10, "name": "robin" } }`)},
				"sidekick": domain.DoneResource{Details: domain.Details{Success: true}, Result: nil},
			},
		},
		{
			"should aggregate one list resource inside other",
			domain.Query{Statements: []domain.Statement{
				{Resource: "hero"},
				{Resource: "sidekick", In: []string{"hero", "sidekick"}},
			}},
			domain.Resources{
				"hero":     domain.DoneResource{Details: domain.Details{Success: true}, Result: unmarshal(`{ "id": 1, "name": "batman" }`)},
				"sidekick": domain.DoneResource{Details: domain.Details{Success: true}, Result: unmarshal(`[{ "id": 10, "name": "robin" }, { "id": 11, "name": "batgirl" }]`)},
			},
			domain.Resources{
				"hero":     domain.DoneResource{Details: domain.Details{Success: true}, Result: unmarshal(`{ "id": 1, "name": "batman", "sidekick":  [{ "id": 10, "name": "robin" }, { "id": 11, "name": "batgirl" }]}`)},
				"sidekick": domain.DoneResource{Details: domain.Details{Success: true}, Result: nil},
			},
		},
		{
			"should aggregate one resource inside other multiplexed resource",
			domain.Query{Statements: []domain.Statement{
				{Resource: "hero"},
				{Resource: "sidekick", In: []string{"hero", "sidekick"}},
			}},
			domain.Resources{
				"hero": domain.DoneResources{
					domain.DoneResource{Details: domain.Details{Success: true}, Result: unmarshal(`{ "id": 1, "name": "batman" }`)},
					domain.DoneResource{Details: domain.Details{Success: true}, Result: unmarshal(`{ "id": 2, "name": "wonder woman" }`)},
				},
				"sidekick": domain.DoneResource{Details: domain.Details{Success: true}, Result: unmarshal(`{ "id": 10, "name": "robin" }`)},
			},
			domain.Resources{
				"hero": domain.DoneResources{
					domain.DoneResource{Details: domain.Details{Success: true}, Result: unmarshal(`{ "id": 1, "name": "batman", "sidekick": { "id": 10, "name": "robin" } }`)},
					domain.DoneResource{Details: domain.Details{Success: true}, Result: unmarshal(`{ "id": 2, "name": "wonder woman", "sidekick": { "id": 10, "name": "robin" } }`)},
				},
				"sidekick": domain.DoneResource{Details: domain.Details{Success: true}, Result: nil},
			},
		},
		{
			"should aggregate one resource inside every item on target result",
			domain.Query{Statements: []domain.Statement{
				{Resource: "hero"},
				{Resource: "sidekick", In: []string{"hero", "sidekick"}},
			}},
			domain.Resources{
				"hero":     domain.DoneResource{Details: domain.Details{Success: true}, Result: unmarshal(`[{ "id": 1, "name": "batman" }, { "id": 2, "name": "wonder woman" }]`)},
				"sidekick": domain.DoneResource{Details: domain.Details{Success: true}, Result: unmarshal(`{ "id": 10, "name": "robin" }`)},
			},
			domain.Resources{
				"hero":     domain.DoneResource{Details: domain.Details{Success: true}, Result: unmarshal(`[{ "id": 1, "name": "batman", "sidekick": { "id": 10, "name": "robin" } }, { "id": 2, "name": "wonder woman", "sidekick": { "id": 10, "name": "robin" } }]`)},
				"sidekick": domain.DoneResource{Details: domain.Details{Success: true}, Result: nil},
			},
		},
		{
			"should aggregate one multiplexed resource inside other",
			domain.Query{Statements: []domain.Statement{
				{Resource: "hero"},
				{Resource: "sidekick", In: []string{"hero", "sidekick"}},
			}},
			domain.Resources{
				"hero": domain.DoneResource{Details: domain.Details{Success: true}, Result: unmarshal(`{ "id": 1, "name": "batman" }`)},
				"sidekick": domain.DoneResources{
					domain.DoneResource{Details: domain.Details{Success: true}, Result: unmarshal(`{ "id": 10, "name": "robin" }`)},
					domain.DoneResource{Details: domain.Details{Success: true}, Result: unmarshal(`{ "id": 11, "name": "batgirl" }`)},
				},
			},
			domain.Resources{
				"hero": domain.DoneResource{
					Details: domain.Details{Success: true},
					Result:  unmarshal(`{ "id": 1, "name": "batman", "sidekick": [{ "id": 10, "name": "robin" }, { "id": 11, "name": "batgirl" }] }`),
				},
				"sidekick": domain.DoneResources{
					domain.DoneResource{Details: domain.Details{Success: true}, Result: nil},
					domain.DoneResource{Details: domain.Details{Success: true}, Result: nil},
				},
			},
		},
		{
			"should aggregate one multiplexed resource with other multiplexed resource",
			domain.Query{Statements: []domain.Statement{
				{Resource: "hero"},
				{Resource: "sidekick", In: []string{"hero", "sidekick"}},
			}},
			domain.Resources{
				"hero": domain.DoneResources{
					domain.DoneResource{Details: domain.Details{Success: true}, Result: unmarshal(`{ "id": 1, "name": "batman"}`)},
					domain.DoneResource{Details: domain.Details{Success: true}, Result: unmarshal(`{ "id": 2, "name": "wonder woman"}`)},
				},
				"sidekick": domain.DoneResources{
					domain.DoneResource{Details: domain.Details{Success: true}, Result: unmarshal(`{ "id": 10, "name": "robin" }`)},
					domain.DoneResource{Details: domain.Details{Success: true}, Result: unmarshal(`{ "id": 11, "name": "batgirl" }`)},
				},
			},
			domain.Resources{
				"hero": domain.DoneResources{
					domain.DoneResource{
						Details: domain.Details{Success: true},
						Result:  unmarshal(`{ "id": 1, "name": "batman", "sidekick": { "id": 10, "name": "robin" } }`),
					},
					domain.DoneResource{
						Details: domain.Details{Success: true},
						Result:  unmarshal(`{ "id": 2, "name": "wonder woman", "sidekick": { "id": 11, "name": "batgirl" } }`),
					},
				},
				"sidekick": domain.DoneResources{
					domain.DoneResource{Details: domain.Details{Success: true}, Result: nil},
					domain.DoneResource{Details: domain.Details{Success: true}, Result: nil},
				},
			},
		},
		{
			"should aggregate one list resource with other list resource",
			domain.Query{Statements: []domain.Statement{
				{Resource: "hero"},
				{Resource: "sidekick", In: []string{"hero", "sidekick"}},
			}},
			domain.Resources{
				"hero":     domain.DoneResource{Details: domain.Details{Success: true}, Result: unmarshal(`[{ "id": 1, "name": "batman" }, { "id": 2, "name": "wonder woman" }]`)},
				"sidekick": domain.DoneResource{Details: domain.Details{Success: true}, Result: unmarshal(`[{ "id": 10, "name": "robin" }, { "id": 11, "name": "batgirl" }]`)},
			},
			domain.Resources{
				"hero":     domain.DoneResource{Details: domain.Details{Success: true}, Result: unmarshal(`[{ "id": 1, "name": "batman", "sidekick": { "id": 10, "name": "robin" } }, { "id": 2, "name": "wonder woman", "sidekick": { "id": 11, "name": "batgirl" } }]`)},
				"sidekick": domain.DoneResource{Details: domain.Details{Success: true}, Result: nil},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := eval.ApplyAggregators(tt.query, tt.resources)

			if !reflect.DeepEqual(got, tt.expected) {
				t.Fatalf("ApplyAggregators = %+#v, want = %+#v", got, tt.expected)
			}
		})
	}
}
