package eval_test

import (
	"encoding/json"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/eval"
	"reflect"
	"testing"
)

func TestHiddenFilter(t *testing.T) {
	query := domain.Query{Statements: []domain.Statement{
		{Resource: "hero", Hidden: true},
		{Resource: "sidekick"},
	}}

	resources := domain.Resources{
		"hero":     domain.DoneResource{Details: domain.Details{Success: true}, Result: nil},
		"sidekick": domain.DoneResource{Details: domain.Details{Success: true}, Result: nil},
	}

	expectedResources := domain.Resources{
		"sidekick": domain.DoneResource{Details: domain.Details{Success: true}, Result: nil},
	}

	got, _ := eval.ApplyFilters(query, resources)

	if !reflect.DeepEqual(got, expectedResources) {
		t.Fatalf("ApplyFilters = %+#v, want = %+#v", got, expectedResources)
	}
}

func TestOnlyFilters(t *testing.T) {
	tests := []struct {
		name      string
		query     domain.Query
		resources domain.Resources
		expected  domain.Resources
	}{
		{
			"should do nothing if there is no filter",
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
			"should bring only the given fields",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{"name", "age"},
			}}},
			domain.Resources{
				"hero": domain.DoneResource{
					Details: domain.Details{Success: true},
					Result:  unmarshal(`{ "id": "12345", "name": "batman", "age": 42 }`),
				},
			},
			domain.Resources{
				"hero": domain.DoneResource{
					Details: domain.Details{Success: true},
					Result:  unmarshal(`{ "name": "batman", "age": 42 }`),
				},
			},
		},
		{
			"should bring multiple nested fields",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{"city.name", "city.population"},
			}}},
			domain.Resources{
				"hero": domain.DoneResource{
					Details: domain.Details{Success: true},
					Result:  unmarshal(`{ "id": "12345", "name": "batman", "age": 42, "city": { "name": "gotham", "population": 10000000 } }`),
				},
			},
			domain.Resources{
				"hero": domain.DoneResource{
					Details: domain.Details{Success: true},
					Result:  unmarshal(`{ "city": { "name": "gotham", "population": 10000000 } }`),
				},
			},
		},
		{
			"should bring only the given nested fields",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{"id", "nested.some-field.even-more-nested"},
			}}},
			domain.Resources{
				"hero": domain.DoneResource{
					Details: domain.Details{Success: true},
					Result:  unmarshal(`{ "id": "12345", "nested": {"some-field": {"even-more-nested": "abcdef", "other-field": 1} } }`),
				},
			},
			domain.Resources{
				"hero": domain.DoneResource{
					Details: domain.Details{Success: true},
					Result:  unmarshal(`{ "id": "12345", "nested": {"some-field": {"even-more-nested": "abcdef"} } }`),
				},
			},
		},
		{
			"should bring only the given fields when field is in a list",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{"weapons.name"},
			}}},
			domain.Resources{
				"hero": domain.DoneResource{
					Details: domain.Details{Success: true},
					Result:  unmarshal(`{ "id": "12345", "weapons": [{"id": 1, "name": "belt"}, {"id": 2, "name": "batarang"}] }`),
				},
			},
			domain.Resources{
				"hero": domain.DoneResource{
					Details: domain.Details{Success: true},
					Result:  unmarshal(`{"weapons": [{"name": "belt"}, {"name": "batarang"}] }`),
				},
			},
		},
		{
			"should bring only the given fields when field is a list of primitivies",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{"weapons"},
			}}},
			domain.Resources{
				"hero": domain.DoneResource{
					Details: domain.Details{Success: true},
					Result:  unmarshal(`{ "id": "12345", "weapons": ["belt", "batarang"] }`),
				},
			},
			domain.Resources{
				"hero": domain.DoneResource{
					Details: domain.Details{Success: true},
					Result:  unmarshal(`{"weapons": ["belt", "batarang"] }`),
				},
			},
		},
		{
			"should bring only the given fields when field is nested and in a list",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{"weapons.nested.some-field.even-more-nested"},
			}}},
			domain.Resources{
				"hero": domain.DoneResource{
					Details: domain.Details{Success: true},
					Result:  unmarshal(`{ "id": "12345", "weapons": [{"id": 1, "name": "belt", "nested": {"some-field": {"even-more-nested": "abcdef"} }}, {"id": 2, "name": "batarang", "nested": {"some-field": {"even-more-nested": "abcdef"} }}] }`),
				},
			},
			domain.Resources{
				"hero": domain.DoneResource{
					Details: domain.Details{Success: true},
					Result:  unmarshal(`{"weapons": [{"nested": {"some-field": {"even-more-nested": "abcdef"} }}, {"nested": {"some-field": {"even-more-nested": "abcdef"} }}] }`),
				},
			},
		},
		{
			"should bring fields in deep nested lists",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{"weapons.name", "weapons.properties.name", "weapons.properties.value"},
			}}},
			domain.Resources{
				"hero": domain.DoneResource{
					Details: domain.Details{Success: true},
					Result:  unmarshal(`{ "id": "12345", "weapons": [{"id": 1, "name": "belt", "properties": [{"name": "color", "value": "yellow"}, {"name": "weight", "value": "10"}]}, {"id": 2, "name": "batarang", "properties": [{"name": "color", "value": "black"}, {"name": "weight", "value": "1"}]}] }`),
				},
			},
			domain.Resources{
				"hero": domain.DoneResource{
					Details: domain.Details{Success: true},
					Result:  unmarshal(`{"weapons": [{"name": "belt", "properties": [{"name": "color", "value": "yellow"}, {"name": "weight", "value": "10"}]}, {"name": "batarang", "properties": [{"name": "color", "value": "black"}, {"name": "weight", "value": "1"}]}] }`),
				},
			},
		},
		{
			"should bring only the given fields when body is a list",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{"name", "age"},
			}}},
			domain.Resources{
				"hero": domain.DoneResource{
					Details: domain.Details{Success: true},
					Result:  unmarshal(`[{ "id": "12345", "name": "batman", "age": 42 },{ "id": "67890", "name": "wonder woman", "age": 35 }]`),
				},
			},
			domain.Resources{
				"hero": domain.DoneResource{
					Details: domain.Details{Success: true},
					Result:  unmarshal(`[{ "name": "batman", "age": 42 },{ "name": "wonder woman", "age": 35 }]`),
				},
			},
		},
		{
			"should bring only the given fields when resource is multiplexed",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{"name", "age"},
			}}},
			domain.Resources{
				"hero": domain.DoneResources{
					domain.DoneResource{
						Details: domain.Details{Success: true},
						Result:  unmarshal(`{ "id": "12345", "name": "batman", "age": 42 }`),
					},
					domain.DoneResource{
						Details: domain.Details{Success: true},
						Result:  unmarshal(`{ "id": "56789", "name": "wonder woman", "age": 35 }`),
					},
				},
			},
			domain.Resources{
				"hero": domain.DoneResources{
					domain.DoneResource{
						Details: domain.Details{Success: true},
						Result:  unmarshal(`{ "name": "batman", "age": 42 }`),
					},
					domain.DoneResource{
						Details: domain.Details{Success: true},
						Result:  unmarshal(`{ "name": "wonder woman", "age": 35 }`),
					},
				},
			},
		},
		{
			"should bring only the given fields that matches arg",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{domain.Match{Target: "id", Arg: "56789"}, domain.Match{Target: "name", Arg: "batman"}, "age"},
			}}},
			domain.Resources{
				"hero": domain.DoneResource{
					Details: domain.Details{Success: true},
					Result:  unmarshal(`{ "id": "12345", "name": "batman", "age": 42 }`),
				},
			},
			domain.Resources{
				"hero": domain.DoneResource{
					Details: domain.Details{Success: true},
					Result:  unmarshal(`{ "name": "batman", "age": 42 }`),
				},
			},
		},
		{
			"should bring only the given fields that matches regex arg",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{domain.Match{Target: "id", Arg: "9$"}, domain.Match{Target: "name", Arg: "^b"}, domain.Match{Target: "age", Arg: "42"}},
			}}},
			domain.Resources{
				"hero": domain.DoneResource{
					Details: domain.Details{Success: true},
					Result:  unmarshal(`{ "id": "12345", "name": "batman", "age": 42 }`),
				},
			},
			domain.Resources{
				"hero": domain.DoneResource{
					Details: domain.Details{Success: true},
					Result:  unmarshal(`{ "name": "batman", "age": 42 }`),
				},
			},
		},
		{
			"should bring only the list elements that matches arg",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{domain.Match{Target: "weapons", Arg: "^b"}},
			}}},
			domain.Resources{
				"hero": domain.DoneResource{
					Details: domain.Details{Success: true},
					Result:  unmarshal(`{ "id": "12345", "weapons": ["belt", "batarang", "katana"] }`),
				},
			},
			domain.Resources{
				"hero": domain.DoneResource{
					Details: domain.Details{Success: true},
					Result:  unmarshal(`{ "weapons": ["belt", "batarang"] }`),
				},
			},
		},
		{
			"should bring everything",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{"*"},
			}}},
			domain.Resources{
				"hero": domain.DoneResource{
					Details: domain.Details{Success: true},
					Result:  unmarshal(`{ "id": "12345", "name": "batman", "age": 42 }`),
				},
			},
			domain.Resources{
				"hero": domain.DoneResource{
					Details: domain.Details{Success: true},
					Result:  unmarshal(`{ "id": "12345", "name": "batman", "age": 42 }`),
				},
			},
		},
		{
			"should bring everything except non matching field",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{"*", domain.Match{Target: "name", Arg: "^c"}},
			}}},
			domain.Resources{
				"hero": domain.DoneResource{
					Details: domain.Details{Success: true},
					Result:  unmarshal(`{ "id": "12345", "name": "batman", "age": 42 }`),
				},
			},
			domain.Resources{
				"hero": domain.DoneResource{
					Details: domain.Details{Success: true},
					Result:  unmarshal(`{ "id": "12345", "age": 42 }`),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := eval.ApplyFilters(tt.query, tt.resources)
			if err != nil {
				t.Fatalf("ApplyFilters returned unexpected error: %s", err)
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Fatalf("ApplyFilters = %+#v, want = %+#v", got, tt.expected)
			}
		})
	}
}

func unmarshal(body string) interface{} {
	var f interface{}
	err := json.Unmarshal([]byte(body), &f)
	if err != nil {
		panic(err)
	}
	return f
}
