package runner_test

import (
	"fmt"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/runner"
	"github.com/b2wdigital/restQL-golang/test"
	"testing"
)

func TestResolveChainedValues(t *testing.T) {
	tests := []struct {
		name                         string
		expected                     domain.Resources
		statementWithUnresolvedParam domain.Resources
		doneResources                domain.Resources
	}{
		{
			"Do nothing if there is no with chained",
			domain.Resources{"resource-name": domain.Statement{Method: "from", Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": "abcdef12345"}}}},
			domain.Resources{
				"resource-name": domain.Statement{
					Method:   "from",
					Resource: "resource-name",
					With:     domain.Params{Values: map[string]interface{}{"id": "abcdef12345"}},
				},
			},
			domain.Resources{"done-resource": domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"id": "abcdef12345"}`)}},
		},
		{
			"Returns a statement with EmptyChained as value if done-resource status code is not in 399 >= status >= 200",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": runner.EmptyChained}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": domain.Chain{"done-resource", "id"}}}}},
			domain.Resources{"done-resource": domain.DoneResource{Status: 404, ResponseBody: test.Unmarshal("{}")}},
		},
		{
			"Returns a statement with EmptyChained as value if done-resource status code is not in 399 >= status >= 200",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": []interface{}{runner.EmptyChained, "abcdef"}}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": domain.Chain{"done-resource", "id"}}}}},
			domain.Resources{"done-resource": domain.DoneResources{domain.DoneResource{Status: 404, ResponseBody: test.Unmarshal("{}")}, domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"id": "abcdef"}`)}}},
		},
		{
			"Returns a statement with single done resource value",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": "abcdef"}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": domain.Chain{"done-resource", "id"}}}}},
			domain.Resources{"done-resource": domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"id": "abcdef"}`)}},
		},
		{
			"Returns a statement with single done resource value",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": float64(1), "uuid": "12354-5656"}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": domain.Chain{"done-resource", "id"}, "uuid": "12354-5656"}}}},
			domain.Resources{"done-resource": domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{ "id": 1}`)}},
		},
		{
			"Returns a statement with single done resource value",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": 1, "name": []interface{}{"12354-5656"}}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": 1, "name": domain.Chain{"done-resource", "name"}}}}},
			domain.Resources{"done-resource": domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"name":["12354-5656"]}`)}},
		},
		{
			"Returns a statement with multiple done resource value",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": []interface{}{float64(1), float64(2)}}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": domain.Chain{"done-resource", "id"}}}}},
			domain.Resources{"done-resource": domain.DoneResources{
				domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"id": 1}`)},
				domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"id": 2}`)},
			}},
		},
		{
			"Returns a statement with single list value 01",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": []interface{}{float64(1), float64(2)}, "name": []interface{}{"a", "b"}}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": domain.Chain{"done-resource", "id"}, "name": []interface{}{"a", "b"}}}}},
			domain.Resources{"done-resource": domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"id": [1,2]}`)}},
		},
		{
			"Returns a statement with single list value",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": []interface{}{float64(1), nil}, "name": []interface{}{"a", "b"}}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": domain.Chain{"done-resource", "id"}, "name": []interface{}{"a", "b"}}}}},
			domain.Resources{"done-resource": domain.DoneResources{
				domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"id":1,"class":"rest"}`)},
				domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"id":null,"class":"rest"}`)},
			}},
		},
		{
			"Returns a statement with multiple list value",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": []interface{}{[]interface{}{float64(1), float64(2)}, []interface{}{float64(3), float64(4)}}}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": domain.Chain{"done-resource", "sidekickId"}}}}},
			domain.Resources{"done-resource": domain.DoneResource{
				Status:       200,
				ResponseBody: test.Unmarshal(`[{"id":"A","sidekickId":[1,2]},{"id":"B","sidekickId":[3,4]}]`)}},
		},
		{
			"Returns a statement with multiple list value",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": []interface{}{[]interface{}{[]interface{}{float64(1), float64(2)}, []interface{}{float64(3), float64(4)}}}}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": domain.Chain{"done-resource", "sidekickId"}}}}},
			domain.Resources{"done-resource": domain.DoneResources{
				domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`[{"id":"A","sidekickId":[1,2]},{"id":"B","sidekickId":[3,4]}]`)},
			}},
		},
		{
			"Returns a statement with empty param",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": []interface{}{float64(1), nil}, "name": []interface{}{"a", "b"}}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": domain.Chain{"done-resource", "id"}, "name": []interface{}{"a", "b"}}}}},
			domain.Resources{"done-resource": domain.DoneResources{
				domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"id":1,"class":"rest"}`)},
				domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"class":"rest"}`)}}},
		},
		{
			"Resolve a statement with lists and nested values",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"name": []interface{}{"clojure", "java"}}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"name": domain.Chain{"done-resource", "language", "id"}}}}},
			domain.Resources{"done-resource": domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"language":{"id":["clojure","java"]}}`)}},
		},
		{
			"Resolve a statement with lists and nested values",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"name": "clojure"}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"name": domain.Chain{"done-resource", "language", "id"}}}}},
			domain.Resources{"done-resource": domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"language":{"id":"clojure"}}`)}},
		},
		{
			"Resolve a statement with lists and nested values",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"name": []interface{}{"clojure"}}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"name": domain.Chain{"done-resource", "language", "id"}}}}},
			domain.Resources{"done-resource": domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"language":[{"id":"clojure"}]}`)}},
		},
		{
			"Resolve a statement with lists and nested values",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"name": []interface{}{"clojure", "java"}}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"name": domain.Chain{"done-resource", "language", "id"}}}}},
			domain.Resources{"done-resource": domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"language":[{"id":"clojure"},{"id":"java"}]}`)}},
		},
		{
			"Resolve a statement with lists and nested values",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"name": []interface{}{"python", "elixir"}}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"name": domain.Chain{"done-resource", "language", "xpto", "id"}}}}},
			domain.Resources{"done-resource": domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"language":[{"xpto":{"id":"python"}},{"xpto":{"id":"elixir"}}]}`)}},
		},
		{
			"Resolve a statement with lists and nested values",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"name": []interface{}{[]interface{}{"python", "123"}, []interface{}{"elixir", "345"}}}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"name": domain.Chain{"done-resource", "language", "xpto", "id"}}}}},
			domain.Resources{"done-resource": domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"language":[{"xpto":{"id":["python","123"]}},{"xpto":{"id":["elixir","345"]}}]}`)}},
		},
		{
			"Resolve a statement with lists and nested values",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"name": []interface{}{[]interface{}{"python", "123"}, []interface{}{"elixir", "345"}}}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"name": domain.Chain{"done-resource", "language", "xpto", "asdf", "id"}}}}},
			domain.Resources{"done-resource": domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"language":[{"xpto":{"asdf":[{"id":"python"},{"id":"123"}]}},{"xpto":{"asdf":[{"id":"elixir"},{"id":"345"}]}}]}`)}},
		},
		{
			"Resolve a statement with lists and nested values",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": []interface{}{[]interface{}{[]interface{}{"DAGGER"}, []interface{}{"GUN"}}, []interface{}{[]interface{}{"SWORD"}, []interface{}{"SHOTGUN"}}}}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": domain.Chain{"done-resource", "weapons"}}}}},
			domain.Resources{
				"done-resource": domain.DoneResources{
					domain.DoneResources{
						domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"name":"1", "weapons":["DAGGER"]}`)},
						domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"name":"2","weapons":["GUN"]}`)},
					},
					domain.DoneResources{
						domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"name":"3", "weapons":["SWORD"]}`)},
						domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"name":"4","weapons":["SHOTGUN"]}`)},
					},
				},
			},
		},
		{
			"Returns a statement with chained param inside list",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"weapon-class": []interface{}{"melee"}}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"weapon-class": domain.Chain{"done-resource", "heroes", "weapons", "type", "class"}}}}},
			domain.Resources{"done-resource": domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"heroes":[{"weapons":{"type":{"class":"melee"}}}]}`)}},
		},
		{
			"Returns a statement with chained param inside list",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"weapon-class": []interface{}{[]interface{}{"melee"}}}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"weapon-class": []interface{}{domain.Chain{"done-resource", "heroes", "weapons", "type", "class"}}}}}},
			domain.Resources{"done-resource": domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"heroes":[{"weapons":{"type":{"class":"melee"}}}]}`)}},
		},
		{
			"Returns a statement with chained param inside list",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"weapon-class": []interface{}{[]interface{}{}}}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"weapon-class": []interface{}{domain.Chain{"done-resource", "heroes", "weapons", "type", "class"}}}}}},
			domain.Resources{"done-resource": domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"heroes":[]}`)}},
		},
		{
			"Returns a statement with chained param inside object",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"country": map[string]interface{}{"code": "USA"}}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"country": map[string]interface{}{"code": domain.Chain{"done-resource", "hero", "origin"}}}}}},
			domain.Resources{"done-resource": domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"hero": {"origin": "USA"}}`)}},
		},
		{
			"Returns a statement with flattened chained value",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": domain.NoMultiplex{Value: []interface{}{"abcdef", "12345"}}}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": domain.NoMultiplex{Value: domain.Chain{"done-resource", "id"}}}}}},
			domain.Resources{"done-resource": domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"id": ["abcdef", "12345"]}`)}},
		},
		{
			"Returns a statement with single flattened list value",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": domain.NoMultiplex{Value: []interface{}{float64(1), float64(2)}}, "name": []interface{}{"a", "b"}}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": domain.NoMultiplex{Value: domain.Chain{"done-resource", "id"}}, "name": []interface{}{"a", "b"}}}}},
			domain.Resources{"done-resource": domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"id": [1,2]}`)}},
		},
		{
			"Returns a statement with single done resource value resolved from header",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": "abcdef", "xtid": "12345678"}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": domain.Chain{"done-resource", "location"}, "xtid": domain.Chain{"done-resource", "x-tid"}}}}},
			domain.Resources{"done-resource": domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"id": "abcdef"}`), ResponseHeaders: map[string]string{"location": "abcdef", "X-TID": "12345678"}}},
		},
		{
			"Returns a multiplexed statement with single done resource value",
			domain.Resources{"resource-name": []interface{}{
				domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": "abcdef"}}},
				domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": "abcdef"}}},
			}},
			domain.Resources{"resource-name": []interface{}{
				domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": domain.Chain{"done-resource", "id"}}}},
				domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"id": domain.Chain{"done-resource", "id"}}}},
			}},
			domain.Resources{"done-resource": domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"id": "abcdef"}`)}},
		},
		{
			"Returns a statement with chained value in header resolved",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", Headers: map[string]interface{}{"x-id": "abcdef"}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", Headers: map[string]interface{}{"x-id": domain.Chain{"done-resource", "id"}}}},
			domain.Resources{"done-resource": domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"id": "abcdef"}`)}},
		},
		{
			"Returns a statement with chained value in header resolved to a complex value",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", Headers: map[string]interface{}{"x-id": `["abcdef","ghijkl"]`}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", Headers: map[string]interface{}{"x-id": domain.Chain{"done-resource", "tokens"}}}},
			domain.Resources{"done-resource": domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"tokens": ["abcdef","ghijkl"]}`)}},
		},
		{
			"Returns a statement with object param with resolved list values exploded",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"info": domain.NoMultiplex{Value: []interface{}{map[string]interface{}{"weapon": "batarang"}, map[string]interface{}{"weapon": "batbelt"}}}}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{Values: map[string]interface{}{"info": domain.NoMultiplex{Value: map[string]interface{}{"weapon": domain.Chain{"done-resource", "hero", "weapons"}}}}}}},
			domain.Resources{"done-resource": domain.DoneResource{Status: 200, ResponseBody: test.Unmarshal(`{"hero": {"weapons": ["batarang", "batbelt"]}}`)}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runner.ResolveChainedValues(tt.statementWithUnresolvedParam, tt.doneResources)
			test.Equal(t, got, tt.expected)
		})
	}
}

func TestValidateChainedValues(t *testing.T) {
	tests := []struct {
		name      string
		expected  error
		resources domain.Resources
	}{
		{
			"Do nothing if there is no with chained",
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
			"Fail validation if chained parameter target unknown resource",
			fmt.Errorf("%w : done-resource.id", runner.ErrInvalidChainedParameter),
			domain.Resources{
				"resource-name": domain.Statement{
					Method:   "from",
					Resource: "resource-name",
					With:     domain.Params{Values: map[string]interface{}{"id": domain.Chain{"done-resource", "id"}}},
				},
			},
		},
		{
			"Fail validation if chained parameter inside list target unknown resource",
			fmt.Errorf("%w : done-resource.id", runner.ErrInvalidChainedParameter),
			domain.Resources{
				"resource-name": domain.Statement{
					Method:   "from",
					Resource: "resource-name",
					With:     domain.Params{Values: map[string]interface{}{"id": []interface{}{[]interface{}{domain.Chain{"done-resource", "id"}}}}},
				},
			},
		},
		{
			"Fail validation if chained parameter inside object target unknown resource",
			fmt.Errorf("%w : done-resource.id", runner.ErrInvalidChainedParameter),
			domain.Resources{
				"resource-name": domain.Statement{
					Method:   "from",
					Resource: "resource-name",
					With: domain.Params{Values: map[string]interface{}{
						"id": map[string]interface{}{
							"universal": map[string]interface{}{
								"value": domain.Chain{"done-resource", "id"},
							},
						},
					}},
				},
			},
		},
		{
			"Fail validation if chained parameter inside flatten target unknown resource",
			fmt.Errorf("%w : done-resource.id", runner.ErrInvalidChainedParameter),
			domain.Resources{
				"resource-name": domain.Statement{
					Method:   "from",
					Resource: "resource-name",
					With: domain.Params{Values: map[string]interface{}{
						"id": domain.NoMultiplex{Value: domain.Chain{"done-resource", "id"}},
					}},
				},
			},
		},
		{
			"Fail validation if chained parameter inside base64 target unknown resource",
			fmt.Errorf("%w : done-resource.id", runner.ErrInvalidChainedParameter),
			domain.Resources{
				"resource-name": domain.Statement{
					Method:   "from",
					Resource: "resource-name",
					With: domain.Params{Values: map[string]interface{}{
						"id": domain.Base64{Value: domain.Chain{"done-resource", "id"}},
					}},
				},
			},
		},
		{
			"Fail validation if chained parameter inside json target unknown resource",
			fmt.Errorf("%w : done-resource.id", runner.ErrInvalidChainedParameter),
			domain.Resources{
				"resource-name": domain.Statement{
					Method:   "from",
					Resource: "resource-name",
					With: domain.Params{Values: map[string]interface{}{
						"id": domain.Json{Value: domain.Chain{"done-resource", "id"}},
					}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runner.ValidateChainedValues(tt.resources)
			test.Equal(t, fmt.Sprintf("%s", got), fmt.Sprintf("%s", tt.expected))
		})
	}
}
