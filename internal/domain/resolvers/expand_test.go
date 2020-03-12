package resolvers_test

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/domain/resolvers"
	"reflect"
	"testing"
)

func TestMultiplexStatements(t *testing.T) {
	tests := []struct {
		name     string
		input    domain.Query
		expected domain.Query
	}{
		{
			"should change nothing if there is no with params",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero"}}},
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero"}}},
		},
		{
			"should change nothing if there is no list parameter",
			domain.Query{
				Use: map[string]interface{}{"max-age": 600},
				Statements: []domain.Statement{
					{
						Method:   "from",
						Resource: "hero",
						Alias:    "h",
						Headers:  map[string]interface{}{"X-Trace-Id": "abcdef12345"},
						With: map[string]interface{}{
							"id":     1,
							"name":   "batman",
							"family": map[string]interface{}{"father": "Thomas Wayne"},
							"height": 10.5,
						},
						Only: []interface{}{"id", "name"},
					},
					{
						Method:   "from",
						Resource: "sidekick",
						Alias:    "s",
						Headers:  map[string]interface{}{"X-Trace-Id": "abcdef12345"},
						With: map[string]interface{}{
							"id":     1,
							"name":   "batman",
							"family": map[string]interface{}{"father": "Thomas Wayne"},
							"height": 10.5,
						},
						Hidden:       true,
						IgnoreErrors: true,
					},
					{
						Method:   "from",
						Resource: "villain",
						Alias:    "v",
					},
				},
			},
			domain.Query{
				Use: map[string]interface{}{"max-age": 600},
				Statements: []domain.Statement{
					{
						Method:   "from",
						Resource: "hero",
						Alias:    "h",
						Headers:  map[string]interface{}{"X-Trace-Id": "abcdef12345"},
						With: map[string]interface{}{
							"id":     1,
							"name":   "batman",
							"family": map[string]interface{}{"father": "Thomas Wayne"},
							"height": 10.5,
						},
						Only: []interface{}{"id", "name"},
					},
					{
						Method:   "from",
						Resource: "sidekick",
						Alias:    "s",
						Headers:  map[string]interface{}{"X-Trace-Id": "abcdef12345"},
						With: map[string]interface{}{
							"id":     1,
							"name":   "batman",
							"family": map[string]interface{}{"father": "Thomas Wayne"},
							"height": 10.5,
						},
						Hidden:       true,
						IgnoreErrors: true,
					},
					{
						Method:   "from",
						Resource: "villain",
						Alias:    "v",
					},
				},
			},
		},
		{
			"should make a new statement for each list value",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", With: map[string]interface{}{"id": []interface{}{"12345", "67890"}}}}},
			domain.Query{Statements: []domain.Statement{
				{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "12345"}},
				{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "67890"}},
			}},
		},
		{
			"should make a new statement for each list value",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", With: map[string]interface{}{"id": []interface{}{"12345", "67890", "19283"}}}}},
			domain.Query{Statements: []domain.Statement{
				{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "12345"}},
				{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "67890"}},
				{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "19283"}},
			}},
		},
		{
			"should make a new statement for each list value and keep other params",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", With: map[string]interface{}{"id": []interface{}{"12345", "67890"}, "name": "batman", "age": 45}}}},
			domain.Query{Statements: []domain.Statement{
				{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "12345", "name": "batman", "age": 45}},
				{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "67890", "name": "batman", "age": 45}},
			}},
		},
		{
			"should make a new statement for each list value in each statement",
			domain.Query{Statements: []domain.Statement{
				{Method: "from", Resource: "hero", With: map[string]interface{}{"id": []interface{}{"12345", "67890"}}},
				{Method: "from", Resource: "villain", With: map[string]interface{}{"id": []interface{}{"abcdef", "ghijkl"}}},
			}},
			domain.Query{Statements: []domain.Statement{
				{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "12345"}},
				{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "67890"}},
				{Method: "from", Resource: "villain", With: map[string]interface{}{"id": "abcdef"}},
				{Method: "from", Resource: "villain", With: map[string]interface{}{"id": "ghijkl"}},
			}},
		},
		{
			"should not make a new statement if list is flattened",
			domain.Query{Statements: []domain.Statement{
				{Method: "from", Resource: "hero", With: map[string]interface{}{"id": domain.Flatten{[]interface{}{"12345", "67890"}}}},
				{Method: "from", Resource: "villain", With: map[string]interface{}{"id": []interface{}{"abcdef", "ghijkl"}}},
			}},
			domain.Query{Statements: []domain.Statement{
				{Method: "from", Resource: "hero", With: map[string]interface{}{"id": domain.Flatten{[]interface{}{"12345", "67890"}}}},
				{Method: "from", Resource: "villain", With: map[string]interface{}{"id": "abcdef"}},
				{Method: "from", Resource: "villain", With: map[string]interface{}{"id": "ghijkl"}},
			}},
		},
		{
			"should make a new statement for each value of the cartesian product of list params",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", With: map[string]interface{}{"id": []interface{}{"12345", "67890"}, "name": []interface{}{"batman", "superman"}}}}},
			domain.Query{Statements: []domain.Statement{
				{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "12345", "name": "batman"}},
				{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "67890", "name": "superman"}},
			}},
		},
		{
			"should make a new statement for each value of the cartesian product of list params",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", With: map[string]interface{}{"id": []interface{}{"12345", "67890"}, "name": []interface{}{"batman", "superman"}, "age": []interface{}{45, 38}}}}},
			domain.Query{Statements: []domain.Statement{
				{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "12345", "name": "batman", "age": 45}},
				{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "67890", "name": "superman", "age": 38}},
			}},
		},
		{
			"should make a new statement for each value of the cartesian product of list params",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", With: map[string]interface{}{"id": []interface{}{"12345", "67890"}, "name": []interface{}{"wonder woman", "batman", "superman"}}}}},
			domain.Query{Statements: []domain.Statement{
				{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "12345", "name": "wonder woman"}},
				{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "67890", "name": "batman"}},
			}},
		},
		{
			"should make a new statement for each value of the cartesian product of list params and keep other params",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", With: map[string]interface{}{"id": []interface{}{"12345", "67890"}, "name": []interface{}{"batman", "superman"}, "universe": "dc"}}}},
			domain.Query{Statements: []domain.Statement{
				{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "12345", "name": "batman", "universe": "dc"}},
				{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "67890", "name": "superman", "universe": "dc"}},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolvers.MultiplexStatements(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("MultiplexStatements = %#+v, want = %#+v", got, tt.expected)
			}
		})
	}
}

func BenchmarkMultiplexStatements(b *testing.B) {
	input := domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", With: map[string]interface{}{"id": []interface{}{"12345", "67890"}, "name": "batman", "age": 45}}}}

	for i := 0; i < b.N; i++ {
		resolvers.MultiplexStatements(input)
	}
}
