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
		expected                     runner.Resources
		statementWithUnresolvedParam runner.Resources
		doneResources                runner.DoneResources
	}{
		{
			"Do nothing if there is no with chained",
			runner.Resources{"resource-name": domain.Statement{Method: "from", Resource: "resource-name", With: domain.Params{"id": "abcdef12345"}}},
			runner.Resources{
				"resource-name": domain.Statement{
					Method:   "from",
					Resource: "resource-name",
					With:     domain.Params{"id": "abcdef12345"},
				},
			},
			runner.DoneResources{"done-resource": runner.DoneRequest{StatusCode: 200, Body: unmarshal(`{"id": "abcdef12345"}`)}},
		},
		{
			"Returns a statement with EmptyChained as value if done-resource status code is not in 399 >= status >= 200",
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": runner.EmptyChained}}},
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": domain.Chain{"done-resource", "id"}}}},
			runner.DoneResources{"done-resource": runner.DoneRequest{StatusCode: 404, Body: unmarshal("{}")}},
		},
		{
			"Returns a statement with EmptyChained as value if done-resource status code is not in 399 >= status >= 200",
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": []interface{}{runner.EmptyChained, "abcdef"}}}},
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": domain.Chain{"done-resource", "id"}}}},
			runner.DoneResources{"done-resource": runner.DoneRequests{runner.DoneRequest{StatusCode: 404, Body: unmarshal("{}")}, runner.DoneRequest{StatusCode: 200, Body: unmarshal(`{"id": "abcdef"}`)}}},
		},
		{
			"Returns a statement with single done resource value",
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": "abcdef"}}},
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": domain.Chain{"done-resource", "id"}}}},
			runner.DoneResources{"done-resource": runner.DoneRequest{StatusCode: 200, Body: unmarshal(`{"id": "abcdef"}`)}},
		},
		{
			"Returns a statement with single done resource value",
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": float64(1), "uuid": "12354-5656"}}},
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": domain.Chain{"done-resource", "id"}, "uuid": "12354-5656"}}},
			runner.DoneResources{"done-resource": runner.DoneRequest{StatusCode: 200, Body: unmarshal(`{ "id": 1}`)}},
		},
		{
			"Returns a statement with single done resource value",
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": 1, "name": []interface{}{"12354-5656"}}}},
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": 1, "name": domain.Chain{"done-resource", "name"}}}},
			runner.DoneResources{"done-resource": runner.DoneRequest{StatusCode: 200, Body: unmarshal(`{"name":["12354-5656"]}`)}},
		},
		{
			"Returns a statement with multiple done resource value",
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": []interface{}{float64(1), float64(2)}}}},
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": domain.Chain{"done-resource", "id"}}}},
			runner.DoneResources{"done-resource": runner.DoneRequests{
				runner.DoneRequest{StatusCode: 200, Body: unmarshal(`{"id": 1}`)},
				runner.DoneRequest{StatusCode: 200, Body: unmarshal(`{"id": 2}`)},
			}},
		},
		{
			"Returns a statement with single list value",
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": []interface{}{float64(1), float64(2)}, "name": []interface{}{"a", "b"}}}},
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": domain.Chain{"done-resource", "id"}, "name": []interface{}{"a", "b"}}}},
			runner.DoneResources{"done-resource": runner.DoneRequest{StatusCode: 200, Body: unmarshal(`{"id": [1,2]}`)}},
		},
		{
			"Returns a statement with multiple list value",
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": []interface{}{[]interface{}{float64(1), float64(2)}, []interface{}{float64(3), float64(4)}}}}},
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": domain.Chain{"done-resource", "sidekickId"}}}},
			runner.DoneResources{"done-resource": runner.DoneRequest{
				StatusCode: 200,
				Body:       unmarshal(`[{"id":"A","sidekickId":[1,2]},{"id":"B","sidekickId":[3,4]}]`)}},
		},
		{
			"Returns a statement with multiple list value",
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": []interface{}{[]interface{}{[]interface{}{float64(1), float64(2)}, []interface{}{float64(3), float64(4)}}}}}},
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": domain.Chain{"done-resource", "sidekickId"}}}},
			runner.DoneResources{"done-resource": runner.DoneRequests{
				runner.DoneRequest{StatusCode: 200, Body: unmarshal(`[{"id":"A","sidekickId":[1,2]},{"id":"B","sidekickId":[3,4]}]`)},
			}},
		},
		{
			"Returns a statement with single list value",
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": []interface{}{float64(1), nil}, "name": []interface{}{"a", "b"}}}},
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": domain.Chain{"done-resource", "id"}, "name": []interface{}{"a", "b"}}}},
			runner.DoneResources{"done-resource": runner.DoneRequests{
				runner.DoneRequest{StatusCode: 200, Body: unmarshal(`{"id":1,"class":"rest"}`)},
				runner.DoneRequest{StatusCode: 200, Body: unmarshal(`{"id":null,"class":"rest"}`)},
			}},
		},
		{
			"Returns a statement with empty param",
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": []interface{}{float64(1), nil}, "name": []interface{}{"a", "b"}}}},
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": domain.Chain{"done-resource", "id"}, "name": []interface{}{"a", "b"}}}},
			runner.DoneResources{"done-resource": runner.DoneRequests{
				runner.DoneRequest{StatusCode: 200, Body: unmarshal(`{"id":1,"class":"rest"}`)},
				runner.DoneRequest{StatusCode: 200, Body: unmarshal(`{"class":"rest"}`)}}},
		},
		{
			"Resolve a statement with lists and nested values",
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"name": []interface{}{"clojure", "java"}}}},
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"name": domain.Chain{"done-resource", "language", "id"}}}},
			runner.DoneResources{"done-resource": runner.DoneRequest{StatusCode: 200, Body: unmarshal(`{"language":{"id":["clojure","java"]}}`)}},
		},
		{
			"Resolve a statement with lists and nested values",
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"name": "clojure"}}},
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"name": domain.Chain{"done-resource", "language", "id"}}}},
			runner.DoneResources{"done-resource": runner.DoneRequest{StatusCode: 200, Body: unmarshal(`{"language":{"id":"clojure"}}`)}},
		},
		{
			"Resolve a statement with lists and nested values",
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"name": []interface{}{"clojure"}}}},
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"name": domain.Chain{"done-resource", "language", "id"}}}},
			runner.DoneResources{"done-resource": runner.DoneRequest{StatusCode: 200, Body: unmarshal(`{"language":[{"id":"clojure"}]}`)}},
		},
		{
			"Resolve a statement with lists and nested values",
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"name": []interface{}{"clojure", "java"}}}},
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"name": domain.Chain{"done-resource", "language", "id"}}}},
			runner.DoneResources{"done-resource": runner.DoneRequest{StatusCode: 200, Body: unmarshal(`{"language":[{"id":"clojure"},{"id":"java"}]}`)}},
		},
		{
			"Resolve a statement with lists and nested values",
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"name": []interface{}{"python", "elixir"}}}},
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"name": domain.Chain{"done-resource", "language", "xpto", "id"}}}},
			runner.DoneResources{"done-resource": runner.DoneRequest{StatusCode: 200, Body: unmarshal(`{"language":[{"xpto":{"id":"python"}},{"xpto":{"id":"elixir"}}]}`)}},
		},
		{
			"Resolve a statement with lists and nested values",
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"name": []interface{}{[]interface{}{"python", "123"}, []interface{}{"elixir", "345"}}}}},
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"name": domain.Chain{"done-resource", "language", "xpto", "id"}}}},
			runner.DoneResources{"done-resource": runner.DoneRequest{StatusCode: 200, Body: unmarshal(`{"language":[{"xpto":{"id":["python","123"]}},{"xpto":{"id":["elixir","345"]}}]}`)}},
		},
		{
			"Resolve a statement with lists and nested values",
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"name": []interface{}{[]interface{}{"python", "123"}, []interface{}{"elixir", "345"}}}}},
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"name": domain.Chain{"done-resource", "language", "xpto", "asdf", "id"}}}},
			runner.DoneResources{"done-resource": runner.DoneRequest{StatusCode: 200, Body: unmarshal(`{"language":[{"xpto":{"asdf":[{"id":"python"},{"id":"123"}]}},{"xpto":{"asdf":[{"id":"elixir"},{"id":"345"}]}}]}`)}},
		},
		{
			"Resolve a statement with lists and nested values",
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": []interface{}{[]interface{}{[]interface{}{"DAGGER"}, []interface{}{"GUN"}}, []interface{}{[]interface{}{"SWORD"}, []interface{}{"SHOTGUN"}}}}}},
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"id": domain.Chain{"done-resource", "weapons"}}}},
			runner.DoneResources{
				"done-resource": runner.DoneRequests{
					runner.DoneRequests{
						runner.DoneRequest{StatusCode: 200, Body: unmarshal(`{"name":"1", "weapons":["DAGGER"]}`)},
						runner.DoneRequest{StatusCode: 200, Body: unmarshal(`{"name":"2","weapons":["GUN"]}`)},
					},
					runner.DoneRequests{
						runner.DoneRequest{StatusCode: 200, Body: unmarshal(`{"name":"3", "weapons":["SWORD"]}`)},
						runner.DoneRequest{StatusCode: 200, Body: unmarshal(`{"name":"4","weapons":["SHOTGUN"]}`)},
					},
				},
			},
		},
		{
			"Returns a statement with chained param inside list",
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"weapon-class": []interface{}{"melee"}}}},
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"weapon-class": domain.Chain{"done-resource", "heroes", "weapons", "type", "class"}}}},
			runner.DoneResources{"done-resource": runner.DoneRequest{StatusCode: 200, Body: unmarshal(`{"heroes":[{"weapons":{"type":{"class":"melee"}}}]}`)}},
		},
		{
			"Returns a statement with chained param inside list",
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"weapon-class": []interface{}{[]interface{}{"melee"}}}}},
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"weapon-class": []interface{}{domain.Chain{"done-resource", "heroes", "weapons", "type", "class"}}}}},
			runner.DoneResources{"done-resource": runner.DoneRequest{StatusCode: 200, Body: unmarshal(`{"heroes":[{"weapons":{"type":{"class":"melee"}}}]}`)}},
		},
		{
			"Returns a statement with chained param inside list",
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"weapon-class": []interface{}{[]interface{}{}}}}},
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"weapon-class": []interface{}{domain.Chain{"done-resource", "heroes", "weapons", "type", "class"}}}}},
			runner.DoneResources{"done-resource": runner.DoneRequest{StatusCode: 200, Body: unmarshal(`{"heroes":[]}`)}},
		},
		{
			"Returns a statement with chained param inside object",
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"country": map[string]interface{}{"code": "USA"}}}},
			runner.Resources{"resource-name": domain.Statement{Resource: "resource-name", With: domain.Params{"country": map[string]interface{}{"code": domain.Chain{"done-resource", "hero", "origin"}}}}},
			runner.DoneResources{"done-resource": runner.DoneRequest{StatusCode: 200, Body: unmarshal(`{"hero": {"origin": "USA"}}`)}},
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
