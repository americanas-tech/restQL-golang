package eval_test

import (
	"regexp"
	"testing"

	"github.com/b2wdigital/restQL-golang/v6/internal/domain"
	"github.com/b2wdigital/restQL-golang/v6/internal/eval"
	"github.com/b2wdigital/restQL-golang/v6/pkg/restql"
	"github.com/b2wdigital/restQL-golang/v6/test"
)

func TestHiddenFilter(t *testing.T) {
	query := domain.Query{Statements: []domain.Statement{
		{Resource: "hero", Hidden: true},
		{Resource: "sidekick"},
	}}

	resources := domain.Resources{
		"hero":     restql.DoneResource{ResponseBody: nil},
		"sidekick": restql.DoneResource{ResponseBody: nil},
	}

	expectedResources := domain.Resources{
		"sidekick": restql.DoneResource{ResponseBody: nil},
	}

	got := eval.ApplyHidden(query, resources)

	test.Equal(t, got, expectedResources)
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
				"hero":     restql.DoneResource{ResponseBody: nil},
				"sidekick": restql.DoneResource{ResponseBody: nil},
			},
			domain.Resources{
				"hero":     restql.DoneResource{ResponseBody: nil},
				"sidekick": restql.DoneResource{ResponseBody: nil},
			},
		},
		{
			"should do nothing if there is resource result is a primitive",
			domain.Query{Statements: []domain.Statement{{Resource: "auth"}}},
			domain.Resources{
				"auth": restql.DoneResource{ResponseBody: restql.NewResponseBodyFromValue(test.NoOpLogger, "1234567890abcdefg")},
			},
			domain.Resources{
				"auth": restql.DoneResource{ResponseBody: restql.NewResponseBodyFromValue(test.NoOpLogger, "1234567890abcdefg")},
			},
		},
		{
			"should bring only the given fields",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{[]string{"name"}, []string{"age"}},
			}}},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(test.NoOpLogger, test.Unmarshal(`{ "id": "12345", "name": "batman", "age": 42 }`)),
				},
			},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(test.NoOpLogger, test.Unmarshal(`{ "name": "batman", "age": 42 }`)),
				},
			},
		},
		{
			"should bring only the given fields and not return field not present",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{[]string{"name"}, []string{"age"}},
			}}},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(test.NoOpLogger, test.Unmarshal(`{ "id": "12345", "name": "batman" }`)),
				},
			},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(test.NoOpLogger, test.Unmarshal(`{ "name": "batman" }`)),
				},
			},
		},
		{
			"should bring multiple nested fields",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{[]string{"city", "name"}, []string{"city", "population"}},
			}}},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": "12345", "name": "batman", "age": 42, "city": { "name": "gotham", "population": 10000000 } }`),
					),
				},
			},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(test.NoOpLogger, test.Unmarshal(`{ "city": { "name": "gotham", "population": 10000000 } }`)),
				},
			},
		},
		{
			"should bring only the given nested fields",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{[]string{"id"}, []string{"nested", "some-field", "even-more-nested"}},
			}}},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": "12345", "nested": {"some-field": {"even-more-nested": "abcdef", "other-field": 1} } }`),
					),
				},
			},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": "12345", "nested": {"some-field": {"even-more-nested": "abcdef"} } }`),
					),
				},
			},
		},
		{
			"should bring only the given fields when field is in a list",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{[]string{"weapons", "name"}},
			}}},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": "12345", "weapons": [{"id": 1, "name": "belt"}, {"id": 2, "name": "batarang"}] }`),
					),
				},
			},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{"weapons": [{"name": "belt"}, {"name": "batarang"}] }`),
					),
				},
			},
		},
		{
			"should bring only the given fields when field is a list of primitivies",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{[]string{"weapons"}},
			}}},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(test.NoOpLogger, test.Unmarshal(`{ "id": "12345", "weapons": ["belt", "batarang"] }`)),
				},
			},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(test.NoOpLogger, test.Unmarshal(`{"weapons": ["belt", "batarang"] }`)),
				},
			},
		},
		{
			"should bring only the given fields when field is nested and in a list",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{[]string{"weapons", "nested", "some-field", "even-more-nested"}},
			}}},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": "12345", "weapons": [{"id": 1, "name": "belt", "nested": {"some-field": {"even-more-nested": "abcdef"} }}, {"id": 2, "name": "batarang", "nested": {"some-field": {"even-more-nested": "abcdef"} }}] }`),
					),
				},
			},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{"weapons": [{"nested": {"some-field": {"even-more-nested": "abcdef"} }}, {"nested": {"some-field": {"even-more-nested": "abcdef"} }}] }`),
					),
				},
			},
		},
		{
			"should bring fields in deep nested lists",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{[]string{"weapons", "name"}, []string{"weapons", "properties", "name"}, []string{"weapons", "properties", "value"}},
			}}},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": "12345", "weapons": [{"id": 1, "name": "belt", "properties": [{"name": "color", "value": "yellow"}, {"name": "weight", "value": "10"}]}, {"id": 2, "name": "batarang", "properties": [{"name": "color", "value": "black"}, {"name": "weight", "value": "1"}]}] }`),
					),
				},
			},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{"weapons": [{"name": "belt", "properties": [{"name": "color", "value": "yellow"}, {"name": "weight", "value": "10"}]}, {"name": "batarang", "properties": [{"name": "color", "value": "black"}, {"name": "weight", "value": "1"}]}] }`),
					),
				},
			},
		},
		{
			"should bring only the given fields when body is a list",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{[]string{"name"}, []string{"age"}},
			}}},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`[{ "id": "12345", "name": "batman", "age": 42 },{ "id": "67890", "name": "wonder woman", "age": 35 }]`),
					),
				},
			},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`[{ "name": "batman", "age": 42 },{ "name": "wonder woman", "age": 35 }]`),
					),
				},
			},
		},
		{
			"should bring only the given fields when resource is multiplexed",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{[]string{"name"}, []string{"age"}},
			}}},
			domain.Resources{
				"hero": restql.DoneResources{
					restql.DoneResource{
						ResponseBody: restql.NewResponseBodyFromValue(
							test.NoOpLogger,
							test.Unmarshal(`{ "id": "12345", "name": "batman", "age": 42 }`),
						),
					},
					restql.DoneResource{
						ResponseBody: restql.NewResponseBodyFromValue(
							test.NoOpLogger,
							test.Unmarshal(`{ "id": "56789", "name": "wonder woman", "age": 35 }`),
						),
					},
				},
			},
			domain.Resources{
				"hero": restql.DoneResources{
					restql.DoneResource{
						ResponseBody: restql.NewResponseBodyFromValue(
							test.NoOpLogger,
							test.Unmarshal(`{ "name": "batman", "age": 42 }`),
						),
					},
					restql.DoneResource{
						ResponseBody: restql.NewResponseBodyFromValue(
							test.NoOpLogger,
							test.Unmarshal(`{ "name": "wonder woman", "age": 35 }`),
						),
					},
				},
			},
		},
		{
			"should bring only the given fields that matches arg",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only: []interface{}{
					domain.Match{Value: []string{"id"}, Args: []domain.Arg{{Name: domain.MatchArgRegex, Value: regexp.MustCompile("56789")}}},
					domain.Match{Value: []string{"name"}, Args: []domain.Arg{{Name: domain.MatchArgRegex, Value: regexp.MustCompile("batman")}}},
					domain.Match{Value: []string{"city"}, Args: []domain.Arg{{Name: domain.MatchArgRegex, Value: "Gotham"}}},
					[]string{"age"},
				},
			}}},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": "12345", "name": "batman", "age": 42, "city": "Gotham" }`),
					),
				},
			},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "name": "batman", "age": 42, "city": "Gotham" }`),
					),
				},
			},
		},
		{
			"should bring only the given fields that matches regex arg",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only: []interface{}{
					domain.Match{Value: []string{"id"}, Args: []domain.Arg{{Name: domain.MatchArgRegex, Value: regexp.MustCompile("9$")}}},
					domain.Match{Value: []string{"name"}, Args: []domain.Arg{{Name: domain.MatchArgRegex, Value: regexp.MustCompile("^b")}}},
					domain.Match{Value: []string{"age"}, Args: []domain.Arg{{Name: domain.MatchArgRegex, Value: regexp.MustCompile("42")}}},
				},
			}}},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": "12345", "name": "batman", "age": 42 }`),
					),
				},
			},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "name": "batman", "age": 42 }`),
					),
				},
			},
		},
		{
			"should bring only the list elements that matches arg",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{domain.Match{Value: []string{"weapons"}, Args: []domain.Arg{{Name: domain.MatchArgRegex, Value: regexp.MustCompile("^b")}}}},
			}}},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": "12345", "weapons": ["belt", "batarang", "katana"] }`),
					),
				},
			},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "weapons": ["belt", "batarang"] }`),
					),
				},
			},
		},
		{
			"should bring everything",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only:     []interface{}{[]string{"*"}},
			}}},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": "12345", "name": "batman", "age": 42 }`),
					),
				},
			},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": "12345", "name": "batman", "age": 42 }`),
					),
				},
			},
		},
		{
			"should bring everything except non matching field",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only: []interface{}{
					[]string{"*"},
					domain.Match{Value: []string{"name"}, Args: []domain.Arg{{Name: domain.MatchArgRegex, Value: regexp.MustCompile("^c")}}},
				},
			}}},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": "12345", "name": "batman", "age": 42 }`),
					),
				},
			},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": "12345", "age": 42 }`),
					),
				},
			},
		},
		{
			"should bring only the given list items that pass filterByRegex",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only: []interface{}{
					domain.NewFilterByRegex([]string{"weapons"}, "profile.type", regexp.MustCompile("attack|buff")),
					[]string{"id"},
				},
			}}},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": "12345", "weapons": [{"id": 1, "profile": {"type": "attack"}, "stats": {"damage": 2}},{"id": 2, "profile": {"type": "defense"}, "stats": {"damage": 2}}, {"id": 3, "profile": {"type": "buff"}, "stats": {"damage": 2}}] }`),
					),
				},
			},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": "12345", "weapons": [{"id": 1, "profile": {"type": "attack"}, "stats": {"damage": 2}}, {"id": 3, "profile": {"type": "buff"}, "stats": {"damage": 2}}] }`),
					),
				},
			},
		},
		{
			"should bring only the given list items that pass filterByRegex",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only: []interface{}{
					domain.NewFilterByRegex([]string{"weapons"}, "profile.type", "attack|buff"),
					[]string{"id"},
				},
			}}},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": "12345", "weapons": [{"id": 1, "profile": {"type": "attack"}, "stats": {"damage": 2}},{"id": 2, "profile": {"type": "defense"}, "stats": {"damage": 2}}, {"id": 3, "profile": {"type": "buff"}, "stats": {"damage": 2}}] }`),
					),
				},
			},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": "12345", "weapons": [{"id": 1, "profile": {"type": "attack"}, "stats": {"damage": 2}}, {"id": 3, "profile": {"type": "buff"}, "stats": {"damage": 2}}] }`),
					),
				},
			},
		},
		{
			"should bring all elements and filter a nested list of objects using filterByRegex",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only: []interface{}{
					domain.NewFilterByRegex([]string{"weapons", "bases"}, "name", "^two-handed-axe$"),
					[]string{"*"},
				},
			}}},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": "12345", "weapons": [{"id": 1, "profile": {"type": "attack"}, "bases": [{"name":"one-handed-axe","damage":2},{"name":"two-handed-axe","damage":5}]},{"id": 2, "profile": {"type": "defense"}, "stats": {"damage": 2}}, {"id": 3, "profile": {"type": "buff"}, "stats": {"damage": 2}}] }`),
					),
				},
			},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": "12345", "weapons": [{"id": 1, "profile": {"type": "attack"}, "bases": [{"name":"two-handed-axe","damage":5}]},{"id": 2, "profile": {"type": "defense"}, "stats": {"damage": 2}}, {"id": 3, "profile": {"type": "buff"}, "stats": {"damage": 2}}] }`),
					),
				},
			},
		},
		{
			"should not apply filterByRegex if target value is not a list",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only: []interface{}{
					domain.NewFilterByRegex([]string{"stats"}, "profile.type", "attack|buff"),
					[]string{"id"},
				},
			}}},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": "12345", "stats": {"profile": {"name": "connan", "type": "warrior"}}, "weapons": [{"id": 1, "profile": {"type": "attack"}, "stats": {"damage": 2}},{"id": 2, "profile": {"type": "defense"}, "stats": {"damage": 2}}, {"id": 3, "profile": {"type": "buff"}, "stats": {"damage": 2}}] }`),
					),
				},
			},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": "12345", "stats": {"profile": {"name": "connan", "type": "warrior"}}}`),
					),
				},
			},
		},
		{
			"should not apply filterByRegex if arguments are not resolved",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only: []interface{}{
					domain.NewFilterByRegex([]string{"weapons"}, domain.Variable{"profile.type"}, "attack|buff"),
					[]string{"id"},
				},
			}}},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": "12345", "weapons": [{"id": 1, "profile": {"type": "attack"}, "stats": {"damage": 2}},{"id": 2, "profile": {"type": "defense"}, "stats": {"damage": 2}}, {"id": 3, "profile": {"type": "buff"}, "stats": {"damage": 2}}] }`),
					),
				},
			},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": "12345", "weapons": [{"id": 1, "profile": {"type": "attack"}, "stats": {"damage": 2}},{"id": 2, "profile": {"type": "defense"}, "stats": {"damage": 2}}, {"id": 3, "profile": {"type": "buff"}, "stats": {"damage": 2}}] }`),
					),
				},
			},
		},
		{
			"should not apply filterByRegex if arguments are not resolved",
			domain.Query{Statements: []domain.Statement{{
				Resource: "hero",
				Only: []interface{}{
					domain.NewFilterByRegex([]string{"weapons"}, "profile.type", domain.Variable{"typePattern"}),
					[]string{"id"},
				},
			}}},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": "12345", "weapons": [{"id": 1, "profile": {"type": "attack"}, "stats": {"damage": 2}},{"id": 2, "profile": {"type": "defense"}, "stats": {"damage": 2}}, {"id": 3, "profile": {"type": "buff"}, "stats": {"damage": 2}}] }`),
					),
				},
			},
			domain.Resources{
				"hero": restql.DoneResource{
					ResponseBody: restql.NewResponseBodyFromValue(
						test.NoOpLogger,
						test.Unmarshal(`{ "id": "12345", "weapons": [{"id": 1, "profile": {"type": "attack"}, "stats": {"damage": 2}},{"id": 2, "profile": {"type": "defense"}, "stats": {"damage": 2}}, {"id": 3, "profile": {"type": "buff"}, "stats": {"damage": 2}}] }`),
					),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := eval.ApplyFilters(test.NoOpLogger, tt.query, tt.resources)

			test.VerifyError(t, err)
			test.Equal(t, got, tt.expected)
		})
	}
}
