package runner_test

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/runner"
	"reflect"
	"testing"
)

func TestResolveVariables(t *testing.T) {
	tests := []struct {
		name      string
		resources runner.Resources
		input     map[string]interface{}
		expected  runner.Resources
	}{
		{
			"resolve variable in timeout",
			runner.Resources{"hero": domain.Statement{Method: "from", Resource: "hero", Timeout: domain.Variable{"duration"}}},
			map[string]interface{}{"duration": "1000"},
			runner.Resources{"hero": domain.Statement{Method: "from", Resource: "hero", Timeout: 1000}},
		},
		{
			"resolve variable in with",
			runner.Resources{
				"hero": domain.Statement{
					Method:   "from",
					Resource: "hero",
					With: map[string]interface{}{
						"id":       "1234567890",
						"name":     domain.Variable{"name"},
						"weapons":  domain.Chain{"done-resource", domain.Variable{"field"}, "id"},
						"sidekick": []interface{}{domain.Variable{"sidekick"}},
						"places":   map[string]interface{}{"city": domain.Variable{"city"}},
					},
				}},
			map[string]interface{}{"name": "batman", "field": "weapon", "sidekick": "robbin", "city": "Gotham"},
			runner.Resources{
				"hero": domain.Statement{
					Method:   "from",
					Resource: "hero",
					With: map[string]interface{}{
						"id":       "1234567890",
						"name":     "batman",
						"weapons":  domain.Chain{"done-resource", "weapon", "id"},
						"sidekick": []interface{}{"robbin"},
						"places":   map[string]interface{}{"city": "Gotham"},
					},
				},
			},
		},
		{
			"resolve variable in max-age/s-max-age",
			runner.Resources{
				"hero": domain.Statement{Method: "from", Resource: "hero", CacheControl: domain.CacheControl{MaxAge: domain.Variable{"cache-control"}, SMaxAge: domain.Variable{"s-cache-control"}}},
			},
			map[string]interface{}{"cache-control": "200", "s-cache-control": "400"},
			runner.Resources{
				"hero": domain.Statement{Method: "from", Resource: "hero", CacheControl: domain.CacheControl{MaxAge: 200, SMaxAge: 400}},
			},
		},
		{
			"resolve variable in headers",
			runner.Resources{
				"hero": domain.Statement{
					Method:   "from",
					Resource: "hero",
					Headers: map[string]interface{}{
						"Authorization": domain.Variable{"auth"},
						"X-Some-Header": domain.Variable{"some-param"},
						"X-Id":          "1234567890",
					},
				},
			},
			map[string]interface{}{"auth": "abcdef0987", "some-param": "abc"},
			runner.Resources{
				"hero": domain.Statement{Method: "from", Resource: "hero", Headers: map[string]interface{}{"Authorization": "abcdef0987", "X-Id": "1234567890", "X-Some-Header": "abc"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runner.ResolveVariables(tt.resources, tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ResolveVariables = %#+v. Want %#+v", got, tt.expected)
			}
		})
	}
}
