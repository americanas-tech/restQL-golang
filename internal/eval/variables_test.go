package eval_test

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/eval"
	"reflect"
	"testing"
)

func TestResolveVariables(t *testing.T) {
	tests := []struct {
		name     string
		query    domain.Query
		input    eval.QueryInput
		expected domain.Query
	}{
		{
			"resolve variable in timeout",
			domain.Query{Statements: []domain.Statement{
				{Method: "from", Resource: "hero", Timeout: domain.Variable{"duration"}},
			}},
			eval.QueryInput{Params: map[string]interface{}{"duration": "1000"}},
			domain.Query{Statements: []domain.Statement{
				{Method: "from", Resource: "hero", Timeout: 1000},
			}},
		},
		{
			"resolve variable in with",
			domain.Query{Statements: []domain.Statement{
				{
					Method:   "from",
					Resource: "hero",
					With: map[string]interface{}{
						"id":      "1234567890",
						"name":    domain.Variable{"name"},
						"weapons": domain.Chain{"done-resource", domain.Variable{"field"}, "id"},
					},
				},
			}},
			eval.QueryInput{Params: map[string]interface{}{"name": "batman", "field": "weapon"}},
			domain.Query{Statements: []domain.Statement{
				{Method: "from", Resource: "hero", With: map[string]interface{}{
					"id":      "1234567890",
					"name":    "batman",
					"weapons": domain.Chain{"done-resource", "weapon", "id"},
				}},
			}},
		},
		{
			"resolve variable in max-age/s-max-age",
			domain.Query{Statements: []domain.Statement{
				{Method: "from", Resource: "hero", CacheControl: domain.CacheControl{MaxAge: domain.Variable{"cache-control"}, SMaxAge: domain.Variable{"s-cache-control"}}},
			}},
			eval.QueryInput{Params: map[string]interface{}{"cache-control": "200", "s-cache-control": "400"}},
			domain.Query{Statements: []domain.Statement{
				{Method: "from", Resource: "hero", CacheControl: domain.CacheControl{MaxAge: 200, SMaxAge: 400}},
			}},
		},
		{
			"resolve variable in headers",
			domain.Query{Statements: []domain.Statement{
				{
					Method:   "from",
					Resource: "hero",
					Headers: map[string]interface{}{
						"Authorization": domain.Variable{"auth"},
						"X-Some-Header": domain.Variable{"some-param"},
						"X-Id":          "1234567890",
					},
				},
			}},
			eval.QueryInput{Params: map[string]interface{}{"auth": "abcdef0987", "some-param": "abc"}},
			domain.Query{Statements: []domain.Statement{
				{Method: "from", Resource: "hero", Headers: map[string]interface{}{"Authorization": "abcdef0987", "X-Id": "1234567890", "X-Some-Header": "abc"}},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := eval.ResolveVariables(tt.query, tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ResolveVariables = %#+v. Want %#+v", got, tt.expected)
			}
		})
	}
}
