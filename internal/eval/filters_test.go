package eval_test

import (
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
	"regexp"
	"testing"

	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
	"github.com/b2wdigital/restQL-golang/v4/internal/eval"
	"github.com/b2wdigital/restQL-golang/v4/test"
)

func TestHiddenFilter(t *testing.T) {
	query := domain.Query{Statements: []domain.Statement{
		{Resource: "hero", Hidden: true},
		{Resource: "sidekick"},
	}}

	resources := restql.Resources{
		"hero":     restql.DoneResource{ResponseBody: nil},
		"sidekick": restql.DoneResource{ResponseBody: nil},
	}

	expectedResources := restql.Resources{
		"sidekick": restql.DoneResource{ResponseBody: nil},
	}

	got := eval.ApplyHidden(query, resources)

	test.Equal(t, got, expectedResources)
}

func TestOnlyFilters(t *testing.T) {
	tests := []struct {
		name      string
		query     domain.Query
		resources restql.Resources
		expected  restql.Resources
	}{
		{
			"should do nothing if there is no filter",
			domain.Query{Statements: []domain.Statement{{Resource: "hero"}, {Resource: "sidekick"}}},
			restql.Resources{
				"hero":     restql.DoneResource{ResponseBody: nil},
				"sidekick": restql.DoneResource{ResponseBody: nil},
			},
			restql.Resources{
				"hero":     restql.DoneResource{ResponseBody: nil},
				"sidekick": restql.DoneResource{ResponseBody: nil},
			},
		},
		{
			"should do nothing if there is resource result is a primitive",
			domain.Query{Statements: []domain.Statement{{Resource: "auth"}}},
			restql.Resources{
				"auth": restql.DoneResource{ResponseBody: "1234567890abcdefg"},
			},
			restql.Resources{
				"auth": restql.DoneResource{ResponseBody: "1234567890abcdefg"},
			},
		},
		{
			"should bring only the given fields",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{[]string{"name"}, []string{"age"}},
			}}},
			restql.Resources{
				"hero": restql.DoneResource{
					ResponseBody: test.Unmarshal(`{ "id": "12345", "name": "batman", "age": 42 }`),
				},
			},
			restql.Resources{
				"hero": restql.DoneResource{
					ResponseBody: test.Unmarshal(`{ "name": "batman", "age": 42 }`),
				},
			},
		},
		{
			"should bring only the given fields and not return field not present",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{[]string{"name"}, []string{"age"}},
			}}},
			restql.Resources{
				"hero": restql.DoneResource{
					ResponseBody: test.Unmarshal(`{ "id": "12345", "name": "batman" }`),
				},
			},
			restql.Resources{
				"hero": restql.DoneResource{
					ResponseBody: test.Unmarshal(`{ "name": "batman" }`),
				},
			},
		},
		{
			"should bring multiple nested fields",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{[]string{"city", "name"}, []string{"city", "population"}},
			}}},
			restql.Resources{
				"hero": restql.DoneResource{
					ResponseBody: test.Unmarshal(`{ "id": "12345", "name": "batman", "age": 42, "city": { "name": "gotham", "population": 10000000 } }`),
				},
			},
			restql.Resources{
				"hero": restql.DoneResource{
					ResponseBody: test.Unmarshal(`{ "city": { "name": "gotham", "population": 10000000 } }`),
				},
			},
		},
		{
			"should bring only the given nested fields",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{[]string{"id"}, []string{"nested", "some-field", "even-more-nested"}},
			}}},
			restql.Resources{
				"hero": restql.DoneResource{
					ResponseBody: test.Unmarshal(`{ "id": "12345", "nested": {"some-field": {"even-more-nested": "abcdef", "other-field": 1} } }`),
				},
			},
			restql.Resources{
				"hero": restql.DoneResource{
					ResponseBody: test.Unmarshal(`{ "id": "12345", "nested": {"some-field": {"even-more-nested": "abcdef"} } }`),
				},
			},
		},
		{
			"should bring only the given fields when field is in a list",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{[]string{"weapons", "name"}},
			}}},
			restql.Resources{
				"hero": restql.DoneResource{
					ResponseBody: test.Unmarshal(`{ "id": "12345", "weapons": [{"id": 1, "name": "belt"}, {"id": 2, "name": "batarang"}] }`),
				},
			},
			restql.Resources{
				"hero": restql.DoneResource{
					ResponseBody: test.Unmarshal(`{"weapons": [{"name": "belt"}, {"name": "batarang"}] }`),
				},
			},
		},
		{
			"should bring only the given fields when field is a list of primitivies",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{[]string{"weapons"}},
			}}},
			restql.Resources{
				"hero": restql.DoneResource{
					ResponseBody: test.Unmarshal(`{ "id": "12345", "weapons": ["belt", "batarang"] }`),
				},
			},
			restql.Resources{
				"hero": restql.DoneResource{
					ResponseBody: test.Unmarshal(`{"weapons": ["belt", "batarang"] }`),
				},
			},
		},
		{
			"should bring only the given fields when field is nested and in a list",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{[]string{"weapons", "nested", "some-field", "even-more-nested"}},
			}}},
			restql.Resources{
				"hero": restql.DoneResource{
					ResponseBody: test.Unmarshal(`{ "id": "12345", "weapons": [{"id": 1, "name": "belt", "nested": {"some-field": {"even-more-nested": "abcdef"} }}, {"id": 2, "name": "batarang", "nested": {"some-field": {"even-more-nested": "abcdef"} }}] }`),
				},
			},
			restql.Resources{
				"hero": restql.DoneResource{
					ResponseBody: test.Unmarshal(`{"weapons": [{"nested": {"some-field": {"even-more-nested": "abcdef"} }}, {"nested": {"some-field": {"even-more-nested": "abcdef"} }}] }`),
				},
			},
		},
		{
			"should bring fields in deep nested lists",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{[]string{"weapons", "name"}, []string{"weapons", "properties", "name"}, []string{"weapons", "properties", "value"}},
			}}},
			restql.Resources{
				"hero": restql.DoneResource{
					ResponseBody: test.Unmarshal(`{ "id": "12345", "weapons": [{"id": 1, "name": "belt", "properties": [{"name": "color", "value": "yellow"}, {"name": "weight", "value": "10"}]}, {"id": 2, "name": "batarang", "properties": [{"name": "color", "value": "black"}, {"name": "weight", "value": "1"}]}] }`),
				},
			},
			restql.Resources{
				"hero": restql.DoneResource{
					ResponseBody: test.Unmarshal(`{"weapons": [{"name": "belt", "properties": [{"name": "color", "value": "yellow"}, {"name": "weight", "value": "10"}]}, {"name": "batarang", "properties": [{"name": "color", "value": "black"}, {"name": "weight", "value": "1"}]}] }`),
				},
			},
		},
		{
			"should bring only the given fields when body is a list",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{[]string{"name"}, []string{"age"}},
			}}},
			restql.Resources{
				"hero": restql.DoneResource{
					ResponseBody: test.Unmarshal(`[{ "id": "12345", "name": "batman", "age": 42 },{ "id": "67890", "name": "wonder woman", "age": 35 }]`),
				},
			},
			restql.Resources{
				"hero": restql.DoneResource{
					ResponseBody: test.Unmarshal(`[{ "name": "batman", "age": 42 },{ "name": "wonder woman", "age": 35 }]`),
				},
			},
		},
		{
			"should bring only the given fields when resource is multiplexed",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{[]string{"name"}, []string{"age"}},
			}}},
			restql.Resources{
				"hero": restql.DoneResources{
					restql.DoneResource{
						ResponseBody: test.Unmarshal(`{ "id": "12345", "name": "batman", "age": 42 }`),
					},
					restql.DoneResource{
						ResponseBody: test.Unmarshal(`{ "id": "56789", "name": "wonder woman", "age": 35 }`),
					},
				},
			},
			restql.Resources{
				"hero": restql.DoneResources{
					restql.DoneResource{
						ResponseBody: test.Unmarshal(`{ "name": "batman", "age": 42 }`),
					},
					restql.DoneResource{
						ResponseBody: test.Unmarshal(`{ "name": "wonder woman", "age": 35 }`),
					},
				},
			},
		},
		{
			"should bring only the given fields that matches arg",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only: []interface{}{
					domain.Match{Value: []string{"id"}, Arg: regexp.MustCompile("56789")},
					domain.Match{Value: []string{"name"}, Arg: regexp.MustCompile("batman")},
					domain.Match{Value: []string{"city"}, Arg: "Gotham"},
					[]string{"age"},
				},
			}}},
			restql.Resources{
				"hero": restql.DoneResource{
					ResponseBody: test.Unmarshal(`{ "id": "12345", "name": "batman", "age": 42, "city": "Gotham" }`),
				},
			},
			restql.Resources{
				"hero": restql.DoneResource{
					ResponseBody: test.Unmarshal(`{ "name": "batman", "age": 42, "city": "Gotham" }`),
				},
			},
		},
		{
			"should bring only the given fields that matches regex arg",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{domain.Match{Value: []string{"id"}, Arg: regexp.MustCompile("9$")}, domain.Match{Value: []string{"name"}, Arg: regexp.MustCompile("^b")}, domain.Match{Value: []string{"age"}, Arg: regexp.MustCompile("42")}},
			}}},
			restql.Resources{
				"hero": restql.DoneResource{
					ResponseBody: test.Unmarshal(`{ "id": "12345", "name": "batman", "age": 42 }`),
				},
			},
			restql.Resources{
				"hero": restql.DoneResource{
					ResponseBody: test.Unmarshal(`{ "name": "batman", "age": 42 }`),
				},
			},
		},
		{
			"should bring only the list elements that matches arg",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{domain.Match{Value: []string{"weapons"}, Arg: regexp.MustCompile("^b")}},
			}}},
			restql.Resources{
				"hero": restql.DoneResource{
					ResponseBody: test.Unmarshal(`{ "id": "12345", "weapons": ["belt", "batarang", "katana"] }`),
				},
			},
			restql.Resources{
				"hero": restql.DoneResource{
					ResponseBody: test.Unmarshal(`{ "weapons": ["belt", "batarang"] }`),
				},
			},
		},
		{
			"should bring everything",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{[]string{"*"}},
			}}},
			restql.Resources{
				"hero": restql.DoneResource{
					ResponseBody: test.Unmarshal(`{ "id": "12345", "name": "batman", "age": 42 }`),
				},
			},
			restql.Resources{
				"hero": restql.DoneResource{
					ResponseBody: test.Unmarshal(`{ "id": "12345", "name": "batman", "age": 42 }`),
				},
			},
		},
		{
			"should bring everything except non matching field",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{[]string{"*"}, domain.Match{Value: []string{"name"}, Arg: regexp.MustCompile("^c")}},
			}}},
			restql.Resources{
				"hero": restql.DoneResource{
					ResponseBody: test.Unmarshal(`{ "id": "12345", "name": "batman", "age": 42 }`),
				},
			},
			restql.Resources{
				"hero": restql.DoneResource{
					ResponseBody: test.Unmarshal(`{ "id": "12345", "age": 42 }`),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := eval.ApplyFilters(test.NoOpLogger{}, tt.query, tt.resources)

			test.VerifyError(t, err)
			test.Equal(t, got, tt.expected)
		})
	}
}
