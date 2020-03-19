package runner_test

import (
	"encoding/json"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/runner"
	"reflect"
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
			domain.Resources{"resource-name": domain.Statement{Method: "from", Resource: "resource-name", With: domain.Params{"id": "abcdef12345"}}},
			domain.Resources{
				"resource-name": domain.Statement{
					Method:   "from",
					Resource: "resource-name",
					With:     domain.Params{"id": "abcdef12345"},
				},
			},
			domain.Resources{"done-resource": domain.DoneResource{Details: domain.Details{Status: 200}, Result: unmarshal(`{"id": "abcdef12345"}`)}},
		},
		{
			"Returns a statement with EmptyChained as value if done-resource status code is not in 399 >= status >= 200",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": runner.EmptyChained}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": domain.Chain{"done-resource", "id"}}}},
			domain.Resources{"done-resource": domain.DoneResource{Details: domain.Details{Status: 404}, Result: unmarshal("{}")}},
		},
		{
			"Returns a statement with EmptyChained as value if done-resource status code is not in 399 >= status >= 200",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": []interface{}{runner.EmptyChained, "abcdef"}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": domain.Chain{"done-resource", "id"}}}},
			domain.Resources{"done-resource": domain.DoneResources{domain.DoneResource{Details: domain.Details{Status: 404}, Result: unmarshal("{}")}, domain.DoneResource{Details: domain.Details{Status: 200}, Result: unmarshal(`{"id": "abcdef"}`)}}},
		},
		{
			"Returns a statement with single done resource value",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": "abcdef"}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": domain.Chain{"done-resource", "id"}}}},
			domain.Resources{"done-resource": domain.DoneResource{Details: domain.Details{Status: 200}, Result: unmarshal(`{"id": "abcdef"}`)}},
		},
		{
			"Returns a statement with single done resource value",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": float64(1), "uuid": "12354-5656"}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": domain.Chain{"done-resource", "id"}, "uuid": "12354-5656"}}},
			domain.Resources{"done-resource": domain.DoneResource{Details: domain.Details{Status: 200}, Result: unmarshal(`{ "id": 1}`)}},
		},
		{
			"Returns a statement with single done resource value",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": 1, "name": []interface{}{"12354-5656"}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": 1, "name": domain.Chain{"done-resource", "name"}}}},
			domain.Resources{"done-resource": domain.DoneResource{Details: domain.Details{Status: 200}, Result: unmarshal(`{"name":["12354-5656"]}`)}},
		},
		{
			"Returns a statement with multiple done resource value",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": []interface{}{float64(1), float64(2)}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": domain.Chain{"done-resource", "id"}}}},
			domain.Resources{"done-resource": domain.DoneResources{
				domain.DoneResource{Details: domain.Details{Status: 200}, Result: unmarshal(`{"id": 1}`)},
				domain.DoneResource{Details: domain.Details{Status: 200}, Result: unmarshal(`{"id": 2}`)},
			}},
		},
		{
			"Returns a statement with single list value",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": []interface{}{float64(1), float64(2)}, "name": []interface{}{"a", "b"}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": domain.Chain{"done-resource", "id"}, "name": []interface{}{"a", "b"}}}},
			domain.Resources{"done-resource": domain.DoneResource{Details: domain.Details{Status: 200}, Result: unmarshal(`{"id": [1,2]}`)}},
		},
		{
			"Returns a statement with multiple list value",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": []interface{}{[]interface{}{float64(1), float64(2)}, []interface{}{float64(3), float64(4)}}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": domain.Chain{"done-resource", "sidekickId"}}}},
			domain.Resources{"done-resource": domain.DoneResource{
				Details: domain.Details{Status: 200},
				Result:  unmarshal(`[{"id":"A","sidekickId":[1,2]},{"id":"B","sidekickId":[3,4]}]`)}},
		},
		{
			"Returns a statement with multiple list value",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": []interface{}{[]interface{}{[]interface{}{float64(1), float64(2)}, []interface{}{float64(3), float64(4)}}}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": domain.Chain{"done-resource", "sidekickId"}}}},
			domain.Resources{"done-resource": domain.DoneResources{
				domain.DoneResource{Details: domain.Details{Status: 200}, Result: unmarshal(`[{"id":"A","sidekickId":[1,2]},{"id":"B","sidekickId":[3,4]}]`)},
			}},
		},
		{
			"Returns a statement with single list value",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": []interface{}{float64(1), nil}, "name": []interface{}{"a", "b"}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": domain.Chain{"done-resource", "id"}, "name": []interface{}{"a", "b"}}}},
			domain.Resources{"done-resource": domain.DoneResources{
				domain.DoneResource{Details: domain.Details{Status: 200}, Result: unmarshal(`{"id":1,"class":"rest"}`)},
				domain.DoneResource{Details: domain.Details{Status: 200}, Result: unmarshal(`{"id":null,"class":"rest"}`)},
			}},
		},
		{
			"Returns a statement with empty param",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": []interface{}{float64(1), nil}, "name": []interface{}{"a", "b"}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": domain.Chain{"done-resource", "id"}, "name": []interface{}{"a", "b"}}}},
			domain.Resources{"done-resource": domain.DoneResources{
				domain.DoneResource{Details: domain.Details{Status: 200}, Result: unmarshal(`{"id":1,"class":"rest"}`)},
				domain.DoneResource{Details: domain.Details{Status: 200}, Result: unmarshal(`{"class":"rest"}`)}}},
		},
		{
			"Resolve a statement with lists and nested values",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"name": []interface{}{"clojure", "java"}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"name": domain.Chain{"done-resource", "language", "id"}}}},
			domain.Resources{"done-resource": domain.DoneResource{Details: domain.Details{Status: 200}, Result: unmarshal(`{"language":{"id":["clojure","java"]}}`)}},
		},
		{
			"Resolve a statement with lists and nested values",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"name": "clojure"}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"name": domain.Chain{"done-resource", "language", "id"}}}},
			domain.Resources{"done-resource": domain.DoneResource{Details: domain.Details{Status: 200}, Result: unmarshal(`{"language":{"id":"clojure"}}`)}},
		},
		{
			"Resolve a statement with lists and nested values",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"name": []interface{}{"clojure"}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"name": domain.Chain{"done-resource", "language", "id"}}}},
			domain.Resources{"done-resource": domain.DoneResource{Details: domain.Details{Status: 200}, Result: unmarshal(`{"language":[{"id":"clojure"}]}`)}},
		},
		{
			"Resolve a statement with lists and nested values",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"name": []interface{}{"clojure", "java"}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"name": domain.Chain{"done-resource", "language", "id"}}}},
			domain.Resources{"done-resource": domain.DoneResource{Details: domain.Details{Status: 200}, Result: unmarshal(`{"language":[{"id":"clojure"},{"id":"java"}]}`)}},
		},
		{
			"Resolve a statement with lists and nested values",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"name": []interface{}{"python", "elixir"}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"name": domain.Chain{"done-resource", "language", "xpto", "id"}}}},
			domain.Resources{"done-resource": domain.DoneResource{Details: domain.Details{Status: 200}, Result: unmarshal(`{"language":[{"xpto":{"id":"python"}},{"xpto":{"id":"elixir"}}]}`)}},
		},
		{
			"Resolve a statement with lists and nested values",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"name": []interface{}{[]interface{}{"python", "123"}, []interface{}{"elixir", "345"}}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"name": domain.Chain{"done-resource", "language", "xpto", "id"}}}},
			domain.Resources{"done-resource": domain.DoneResource{Details: domain.Details{Status: 200}, Result: unmarshal(`{"language":[{"xpto":{"id":["python","123"]}},{"xpto":{"id":["elixir","345"]}}]}`)}},
		},
		{
			"Resolve a statement with lists and nested values",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"name": []interface{}{[]interface{}{"python", "123"}, []interface{}{"elixir", "345"}}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"name": domain.Chain{"done-resource", "language", "xpto", "asdf", "id"}}}},
			domain.Resources{"done-resource": domain.DoneResource{Details: domain.Details{Status: 200}, Result: unmarshal(`{"language":[{"xpto":{"asdf":[{"id":"python"},{"id":"123"}]}},{"xpto":{"asdf":[{"id":"elixir"},{"id":"345"}]}}]}`)}},
		},
		{
			"Resolve a statement with lists and nested values",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": []interface{}{[]interface{}{[]interface{}{"DAGGER"}, []interface{}{"GUN"}}, []interface{}{[]interface{}{"SWORD"}, []interface{}{"SHOTGUN"}}}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": domain.Chain{"done-resource", "weapons"}}}},
			domain.Resources{
				"done-resource": domain.DoneResources{
					domain.DoneResources{
						domain.DoneResource{Details: domain.Details{Status: 200}, Result: unmarshal(`{"name":"1", "weapons":["DAGGER"]}`)},
						domain.DoneResource{Details: domain.Details{Status: 200}, Result: unmarshal(`{"name":"2","weapons":["GUN"]}`)},
					},
					domain.DoneResources{
						domain.DoneResource{Details: domain.Details{Status: 200}, Result: unmarshal(`{"name":"3", "weapons":["SWORD"]}`)},
						domain.DoneResource{Details: domain.Details{Status: 200}, Result: unmarshal(`{"name":"4","weapons":["SHOTGUN"]}`)},
					},
				},
			},
		},
		{
			"Returns a statement with chained param inside list",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"weapon-class": []interface{}{"melee"}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"weapon-class": domain.Chain{"done-resource", "heroes", "weapons", "type", "class"}}}},
			domain.Resources{"done-resource": domain.DoneResource{Details: domain.Details{Status: 200}, Result: unmarshal(`{"heroes":[{"weapons":{"type":{"class":"melee"}}}]}`)}},
		},
		{
			"Returns a statement with chained param inside list",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"weapon-class": []interface{}{[]interface{}{"melee"}}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"weapon-class": []interface{}{domain.Chain{"done-resource", "heroes", "weapons", "type", "class"}}}}},
			domain.Resources{"done-resource": domain.DoneResource{Details: domain.Details{Status: 200}, Result: unmarshal(`{"heroes":[{"weapons":{"type":{"class":"melee"}}}]}`)}},
		},
		{
			"Returns a statement with chained param inside list",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"weapon-class": []interface{}{[]interface{}{}}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"weapon-class": []interface{}{domain.Chain{"done-resource", "heroes", "weapons", "type", "class"}}}}},
			domain.Resources{"done-resource": domain.DoneResource{Details: domain.Details{Status: 200}, Result: unmarshal(`{"heroes":[]}`)}},
		},
		{
			"Returns a statement with chained param inside object",
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"country": map[string]interface{}{"code": "USA"}}}},
			domain.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"country": map[string]interface{}{"code": domain.Chain{"done-resource", "hero", "origin"}}}}},
			domain.Resources{"done-resource": domain.DoneResource{Details: domain.Details{Status: 200}, Result: unmarshal(`{"hero": {"origin": "USA"}}`)}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := runner.ResolveChainedValues(tt.statementWithUnresolvedParam, tt.doneResources); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ResolveChainedValue = %+#v, want %+#v", got, tt.expected)
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
