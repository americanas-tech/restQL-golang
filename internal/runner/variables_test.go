package runner_test

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/runner"
	"github.com/b2wdigital/restQL-golang/test"
	"testing"
)

func TestResolveVariables(t *testing.T) {
	tests := []struct {
		name      string
		resources domain.Resources
		input     domain.QueryInput
		expected  domain.Resources
	}{
		{
			"resolve variable in timeout from params",
			domain.Resources{"hero": domain.Statement{Method: "from", Resource: "hero", Timeout: domain.Variable{"duration"}}},
			domain.QueryInput{Params: map[string]interface{}{"duration": "1000"}},
			domain.Resources{"hero": domain.Statement{Method: "from", Resource: "hero", Timeout: 1000}},
		},
		{
			"resolve variable in timeout from header",
			domain.Resources{"hero": domain.Statement{Method: "from", Resource: "hero", Timeout: domain.Variable{"duration"}}},
			domain.QueryInput{Headers: map[string]string{"duration": "1000"}},
			domain.Resources{"hero": domain.Statement{Method: "from", Resource: "hero", Timeout: 1000}},
		},
		{
			"resolve variable in timeout from body",
			domain.Resources{"hero": domain.Statement{Method: "from", Resource: "hero", Timeout: domain.Variable{"duration"}}},
			domain.QueryInput{Body: map[string]interface{}{"duration": 1000}},
			domain.Resources{"hero": domain.Statement{Method: "from", Resource: "hero", Timeout: 1000}},
		},
		{
			"resolve variable in with from params",
			domain.Resources{
				"hero": domain.Statement{
					Method:   "from",
					Resource: "hero",
					With: domain.Params{
						Body: domain.Variable{Target: "heroInfo"},
						Values: map[string]interface{}{
							"id":           "1234567890",
							"name":         domain.Variable{"name"},
							"affiliations": domain.Flatten{Target: domain.Variable{"affiliations"}},
							"weapons":      domain.Chain{"done-resource", domain.Variable{"field"}, "id"},
							"sidekick":     []interface{}{[]interface{}{domain.Variable{"sidekick"}}},
							"places":       map[string]interface{}{"city": map[string]interface{}{"name": domain.Variable{"city"}}},
						},
					},
				}},
			domain.QueryInput{Params: map[string]interface{}{
				"name":         "batman",
				"affiliations": []string{"justice league", "batman family"},
				"field":        "weapon",
				"sidekick":     "robbin",
				"city":         "Gotham",
				"heroInfo":     `{"id": "test"}`,
			}},
			domain.Resources{
				"hero": domain.Statement{
					Method:   "from",
					Resource: "hero",
					With: domain.Params{
						Body: map[string]interface{}{"id": "test"},
						Values: map[string]interface{}{
							"id":           "1234567890",
							"name":         "batman",
							"affiliations": domain.Flatten{Target: []string{"justice league", "batman family"}},
							"weapons":      domain.Chain{"done-resource", "weapon", "id"},
							"sidekick":     []interface{}{[]interface{}{"robbin"}},
							"places":       map[string]interface{}{"city": map[string]interface{}{"name": "Gotham"}},
						},
					},
				},
			},
		},
		{
			"resolve variable in with from headers",
			domain.Resources{
				"hero": domain.Statement{
					Method:   "from",
					Resource: "hero",
					With: domain.Params{Values: map[string]interface{}{
						"id":       "1234567890",
						"name":     domain.Variable{"name"},
						"weapons":  domain.Chain{"done-resource", domain.Variable{"field"}, "id"},
						"sidekick": []interface{}{[]interface{}{domain.Variable{"sidekick"}}},
						"places":   map[string]interface{}{"city": map[string]interface{}{"name": domain.Variable{"city"}}},
					}},
				}},
			domain.QueryInput{Headers: map[string]string{"name": "batman", "field": "weapon", "sidekick": "robbin", "city": "Gotham"}},
			domain.Resources{
				"hero": domain.Statement{
					Method:   "from",
					Resource: "hero",
					With: domain.Params{Values: map[string]interface{}{
						"id":       "1234567890",
						"name":     "batman",
						"weapons":  domain.Chain{"done-resource", "weapon", "id"},
						"sidekick": []interface{}{[]interface{}{"robbin"}},
						"places":   map[string]interface{}{"city": map[string]interface{}{"name": "Gotham"}},
					}},
				},
			},
		},
		{
			"resolve variable in with from body",
			domain.Resources{
				"hero": domain.Statement{
					Method:   "from",
					Resource: "hero",
					With: domain.Params{
						Body: domain.Variable{Target: "heroInfo"},
						Values: map[string]interface{}{
							"id":           "1234567890",
							"name":         domain.Variable{"name"},
							"affiliations": domain.Variable{"affiliations"},
							"weapons":      domain.Chain{"done-resource", domain.Variable{"field"}, "id"},
							"sidekick":     []interface{}{[]interface{}{domain.Variable{"sidekick"}}},
							"places":       map[string]interface{}{"city": map[string]interface{}{"name": domain.Variable{"city"}}},
						},
					},
				}},
			domain.QueryInput{Body: map[string]interface{}{
				"name":         "batman",
				"affiliations": []string{"justice league", "batman family"},
				"field":        "weapon",
				"sidekick":     "robbin",
				"city":         "Gotham",
				"heroInfo":     map[string]interface{}{"id": "test"},
			}},
			domain.Resources{
				"hero": domain.Statement{
					Method:   "from",
					Resource: "hero",
					With: domain.Params{
						Body: map[string]interface{}{"id": "test"},
						Values: map[string]interface{}{
							"id":           "1234567890",
							"name":         "batman",
							"affiliations": []string{"justice league", "batman family"},
							"weapons":      domain.Chain{"done-resource", "weapon", "id"},
							"sidekick":     []interface{}{[]interface{}{"robbin"}},
							"places":       map[string]interface{}{"city": map[string]interface{}{"name": "Gotham"}},
						},
					},
				},
			},
		},
		{
			"resolve variable in max-age/s-max-age from params",
			domain.Resources{
				"hero": domain.Statement{Method: "from", Resource: "hero", CacheControl: domain.CacheControl{MaxAge: domain.Variable{"cache-control"}, SMaxAge: domain.Variable{"s-cache-control"}}},
			},
			domain.QueryInput{Params: map[string]interface{}{"cache-control": "200", "s-cache-control": "400"}},
			domain.Resources{
				"hero": domain.Statement{Method: "from", Resource: "hero", CacheControl: domain.CacheControl{MaxAge: 200, SMaxAge: 400}},
			},
		},
		{
			"resolve variable in max-age/s-max-age from headers",
			domain.Resources{
				"hero": domain.Statement{Method: "from", Resource: "hero", CacheControl: domain.CacheControl{MaxAge: domain.Variable{"cache-control"}, SMaxAge: domain.Variable{"s-cache-control"}}},
			},
			domain.QueryInput{Headers: map[string]string{"cache-control": "200", "s-cache-control": "400"}},
			domain.Resources{
				"hero": domain.Statement{Method: "from", Resource: "hero", CacheControl: domain.CacheControl{MaxAge: 200, SMaxAge: 400}},
			},
		},
		{
			"resolve variable in max-age/s-max-age from params",
			domain.Resources{
				"hero": domain.Statement{Method: "from", Resource: "hero", CacheControl: domain.CacheControl{MaxAge: domain.Variable{"cache-control"}, SMaxAge: domain.Variable{"s-cache-control"}}},
			},
			domain.QueryInput{Body: map[string]interface{}{"cache-control": "200", "s-cache-control": 400}},
			domain.Resources{
				"hero": domain.Statement{Method: "from", Resource: "hero", CacheControl: domain.CacheControl{MaxAge: 200, SMaxAge: 400}},
			},
		},
		{
			"resolve variable in headers from params",
			domain.Resources{
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
			domain.QueryInput{Params: map[string]interface{}{"auth": "abcdef0987", "some-param": "abc"}},
			domain.Resources{
				"hero": domain.Statement{Method: "from", Resource: "hero", Headers: map[string]interface{}{"Authorization": "abcdef0987", "X-Id": "1234567890", "X-Some-Header": "abc"}},
			},
		},
		{
			"resolve variable in headers from headers",
			domain.Resources{
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
			domain.QueryInput{Headers: map[string]string{"auth": "abcdef0987", "some-param": "abc"}},
			domain.Resources{
				"hero": domain.Statement{Method: "from", Resource: "hero", Headers: map[string]interface{}{"Authorization": "abcdef0987", "X-Id": "1234567890", "X-Some-Header": "abc"}},
			},
		},
		{
			"resolve variable in headers from body",
			domain.Resources{
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
			domain.QueryInput{Body: map[string]interface{}{"auth": "abcdef0987", "some-param": "abc"}},
			domain.Resources{
				"hero": domain.Statement{Method: "from", Resource: "hero", Headers: map[string]interface{}{"Authorization": "abcdef0987", "X-Id": "1234567890", "X-Some-Header": "abc"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runner.ResolveVariables(tt.resources, tt.input)
			test.Equal(t, got, tt.expected)
		})
	}
}
