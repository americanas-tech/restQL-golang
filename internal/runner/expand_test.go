package runner_test

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/runner"
	"reflect"
	"testing"
)

func TestMultiplexStatements(t *testing.T) {
	tests := []struct {
		name     string
		input    domain.Resources
		expected domain.Resources
	}{
		{
			"should change nothing if there is no with params",
			domain.Resources{"hero": domain.Statement{Method: "from", Resource: "hero"}},
			domain.Resources{"hero": domain.Statement{Method: "from", Resource: "hero"}},
		},
		{
			"should change nothing if there is no list parameter",
			domain.Resources{
				"h": domain.Statement{
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
				"s": domain.Statement{
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
				"v": domain.Statement{
					Method:   "from",
					Resource: "villain",
					Alias:    "v",
				},
			},
			domain.Resources{
				"h": domain.Statement{
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
				"s": domain.Statement{
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
				"v": domain.Statement{
					Method:   "from",
					Resource: "villain",
					Alias:    "v",
				},
			},
		},
		{
			"should make nested new statements for each list value",
			domain.Resources{
				"hero": domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": []interface{}{"12345", "67890"}}},
			},
			domain.Resources{
				"hero": []interface{}{
					domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "12345"}},
					domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "67890"}},
				},
			},
		},
		{
			"should make deep nested new statements for each list value",
			domain.Resources{
				"hero": domain.Statement{
					Method:   "from",
					Resource: "hero",
					With:     map[string]interface{}{"id": []interface{}{[]interface{}{"12345"}, []interface{}{"67890"}}},
				},
			},
			domain.Resources{
				"hero": []interface{}{
					[]interface{}{domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "12345"}}},
					[]interface{}{domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "67890"}}},
				},
			},
		},
		{
			"should make a new statement for each list value",
			domain.Resources{
				"hero": domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": []interface{}{"12345", "67890", "19283"}}},
			},
			domain.Resources{
				"hero": []interface{}{
					domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "12345"}},
					domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "67890"}},
					domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "19283"}},
				},
			},
		},
		{
			"should make a new statement for each list value and keep other params",
			domain.Resources{
				"hero": domain.Statement{
					Method:   "from",
					Resource: "hero",
					With:     map[string]interface{}{"id": []interface{}{"12345", "67890"}, "name": "batman", "age": 45},
				},
			},
			domain.Resources{
				"hero": []interface{}{
					domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "12345", "name": "batman", "age": 45}},
					domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "67890", "name": "batman", "age": 45}},
				},
			},
		},
		{
			"should make a new statement for each list value only on statement with list param",
			domain.Resources{
				"hero":    domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "12345"}},
				"villain": domain.Statement{Method: "from", Resource: "villain", With: map[string]interface{}{"id": []interface{}{"abcdef", "ghijkl"}}},
			},
			domain.Resources{
				"hero": domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "12345"}},
				"villain": []interface{}{
					domain.Statement{Method: "from", Resource: "villain", With: map[string]interface{}{"id": "abcdef"}},
					domain.Statement{Method: "from", Resource: "villain", With: map[string]interface{}{"id": "ghijkl"}},
				},
			},
		},
		{
			"should make a new statement for each list value in each statement",
			domain.Resources{
				"hero":    domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": []interface{}{"12345", "67890"}}},
				"villain": domain.Statement{Method: "from", Resource: "villain", With: map[string]interface{}{"id": []interface{}{"abcdef", "ghijkl"}}},
			},
			domain.Resources{
				"hero": []interface{}{
					domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "12345"}},
					domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "67890"}},
				},
				"villain": []interface{}{
					domain.Statement{Method: "from", Resource: "villain", With: map[string]interface{}{"id": "abcdef"}},
					domain.Statement{Method: "from", Resource: "villain", With: map[string]interface{}{"id": "ghijkl"}},
				},
			},
		},
		{
			"should not make a new statement if list is flattened",
			domain.Resources{
				"hero":    domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": domain.Flatten{[]interface{}{"12345", "67890"}}}},
				"villain": domain.Statement{Method: "from", Resource: "villain", With: map[string]interface{}{"id": []interface{}{"abcdef", "ghijkl"}}},
			},
			domain.Resources{
				"hero": domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": domain.Flatten{[]interface{}{"12345", "67890"}}}},
				"villain": []interface{}{
					domain.Statement{Method: "from", Resource: "villain", With: map[string]interface{}{"id": "abcdef"}},
					domain.Statement{Method: "from", Resource: "villain", With: map[string]interface{}{"id": "ghijkl"}},
				},
			},
		},
		{
			"should make a new statement for each value of the cartesian product of list params",
			domain.Resources{
				"hero": domain.Statement{
					Method:   "from",
					Resource: "hero",
					With:     map[string]interface{}{"id": []interface{}{"12345", "67890"}, "name": []interface{}{"batman", "superman"}},
				},
			},
			domain.Resources{
				"hero": []interface{}{
					domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "12345", "name": "batman"}},
					domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "67890", "name": "superman"}},
				},
			},
		},
		{
			"should make a new statement for each value of the cartesian product of list params",
			domain.Resources{
				"hero": domain.Statement{
					Method:   "from",
					Resource: "hero",
					With:     map[string]interface{}{"id": []interface{}{"12345", "67890"}, "name": []interface{}{"batman", "superman"}, "age": []interface{}{45, 38}},
				},
			},
			domain.Resources{
				"hero": []interface{}{
					domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "12345", "name": "batman", "age": 45}},
					domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "67890", "name": "superman", "age": 38}},
				},
			},
		},
		{
			"should make a new statement for each value of the cartesian product of list params",
			domain.Resources{
				"hero": domain.Statement{
					Method:   "from",
					Resource: "hero",
					With:     map[string]interface{}{"id": []interface{}{"12345", "67890"}, "name": []interface{}{"wonder woman", "batman", "superman"}},
				},
			},
			domain.Resources{
				"hero": []interface{}{
					domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "12345", "name": "wonder woman"}},
					domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "67890", "name": "batman"}},
				},
			},
		},
		{
			"should make a new statement for each value of the cartesian product of list params and keep other params",
			domain.Resources{
				"hero": domain.Statement{
					Method:   "from",
					Resource: "hero",
					With:     map[string]interface{}{"id": []interface{}{"12345", "67890"}, "name": []interface{}{"batman", "superman"}, "universe": "dc"},
				},
			},
			domain.Resources{
				"hero": []interface{}{
					domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "12345", "name": "batman", "universe": "dc"}},
					domain.Statement{Method: "from", Resource: "hero", With: map[string]interface{}{"id": "67890", "name": "superman", "universe": "dc"}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runner.MultiplexStatements(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("MultiplexStatements = %#+v, want = %#+v", got, tt.expected)
			}
		})
	}
}
