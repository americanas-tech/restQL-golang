package eval_test

import (
	"testing"

	"github.com/b2wdigital/restQL-golang/v6/internal/eval"
	"github.com/b2wdigital/restQL-golang/v6/pkg/restql"

	"github.com/b2wdigital/restQL-golang/v6/internal/domain"
	"github.com/b2wdigital/restQL-golang/v6/test"
)

func TestResolveVariables(t *testing.T) {
	tests := []struct {
		name      string
		resources domain.Query
		input     restql.QueryInput
		expected  domain.Query
	}{
		{
			"resolve variable in timeout from params",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", Timeout: domain.Variable{"duration"}}}},
			restql.QueryInput{Params: map[string]interface{}{"duration": "1000"}},
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", Timeout: 1000}}},
		},
		{
			"resolve variable in timeout from header",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", Timeout: domain.Variable{"duration"}}}},
			restql.QueryInput{Headers: map[string]string{"Duration": "1000"}},
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", Timeout: 1000}}},
		},
		{
			"resolve variable in timeout from body",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", Timeout: domain.Variable{"duration"}}}},
			restql.QueryInput{Body: map[string]interface{}{"duration": 1000}},
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", Timeout: 1000}}},
		},
		{
			"resolve variable in with from params",
			domain.Query{
				Statements: []domain.Statement{{
					Method:   "from",
					Resource: "hero",
					With: domain.Params{
						Body: domain.Variable{Target: "heroInfo"},
						Values: map[string]interface{}{
							"id":           "1234567890",
							"name":         domain.Variable{"name"},
							"affiliations": domain.NoMultiplex{Value: domain.Variable{"affiliations"}},
							"weapons":      domain.Chain{"done-resource", domain.Variable{"field"}, "id"},
							"sidekick":     []interface{}{[]interface{}{domain.Variable{"sidekick"}, map[string]interface{}{"mainHero": domain.Variable{"name"}}}},
							"places":       map[string]interface{}{"city": map[string]interface{}{"name": domain.Variable{"city"}, "neighborhoods": []interface{}{domain.Variable{"mainNeighborhood"}}}},
						},
					},
				}}},
			restql.QueryInput{Params: map[string]interface{}{
				"name":             "batman",
				"affiliations":     []string{"justice league", "batman family"},
				"field":            "weapon",
				"sidekick":         "robbin",
				"city":             "Gotham",
				"mainNeighborhood": "bowery",
				"heroInfo":         `{"id": "test"}`,
			}},
			domain.Query{
				Statements: []domain.Statement{{
					Method:   "from",
					Resource: "hero",
					With: domain.Params{
						Body: map[string]interface{}{"id": "test"},
						Values: map[string]interface{}{
							"id":           "1234567890",
							"name":         "batman",
							"affiliations": domain.NoMultiplex{Value: []string{"justice league", "batman family"}},
							"weapons":      domain.Chain{"done-resource", "weapon", "id"},
							"sidekick":     []interface{}{[]interface{}{"robbin", map[string]interface{}{"mainHero": "batman"}}},
							"places":       map[string]interface{}{"city": map[string]interface{}{"name": "Gotham", "neighborhoods": []interface{}{"bowery"}}},
						},
					},
				}},
			},
		},
		{
			"resolve variable in with from headers",
			domain.Query{
				Statements: []domain.Statement{{
					Method:   "from",
					Resource: "hero",
					With: domain.Params{Values: map[string]interface{}{
						"id":       "1234567890",
						"name":     domain.Variable{"name"},
						"weapons":  domain.Chain{"done-resource", domain.Variable{"field"}, "id"},
						"sidekick": []interface{}{[]interface{}{domain.Variable{"sidekick"}}},
						"places":   map[string]interface{}{"city": map[string]interface{}{"name": domain.Variable{"city"}}},
					}},
				}}},
			restql.QueryInput{Headers: map[string]string{"Name": "batman", "Field": "weapon", "Sidekick": "robbin", "City": "Gotham"}},
			domain.Query{
				Statements: []domain.Statement{{
					Method:   "from",
					Resource: "hero",
					With: domain.Params{Values: map[string]interface{}{
						"id":       "1234567890",
						"name":     "batman",
						"weapons":  domain.Chain{"done-resource", "weapon", "id"},
						"sidekick": []interface{}{[]interface{}{"robbin"}},
						"places":   map[string]interface{}{"city": map[string]interface{}{"name": "Gotham"}},
					}},
				}},
			},
		},
		{
			"resolve variable in with from body",
			domain.Query{
				Statements: []domain.Statement{{
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
					}},
				}},
			restql.QueryInput{Body: map[string]interface{}{
				"name":         "batman",
				"affiliations": []string{"justice league", "batman family"},
				"field":        "weapon",
				"sidekick":     "robbin",
				"city":         "Gotham",
				"heroInfo":     map[string]interface{}{"id": "test"},
			}},
			domain.Query{
				Statements: []domain.Statement{{
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
				}},
			},
		},
		{
			"resolve variable in max-age/s-max-age from params",
			domain.Query{
				Statements: []domain.Statement{{Method: "from", Resource: "hero", CacheControl: domain.CacheControl{MaxAge: domain.Variable{"cache-control"}, SMaxAge: domain.Variable{"s-cache-control"}}}},
			},
			restql.QueryInput{Params: map[string]interface{}{"cache-control": "200", "s-cache-control": "400"}},
			domain.Query{
				Statements: []domain.Statement{{Method: "from", Resource: "hero", CacheControl: domain.CacheControl{MaxAge: 200, SMaxAge: 400}}},
			},
		},
		{
			"resolve variable in max-age/s-max-age from headers",
			domain.Query{
				Statements: []domain.Statement{{Method: "from", Resource: "hero", CacheControl: domain.CacheControl{MaxAge: domain.Variable{"cache-control"}, SMaxAge: domain.Variable{"s-cache-control"}}}},
			},
			restql.QueryInput{Headers: map[string]string{"Cache-Control": "200", "S-Cache-Control": "400"}},
			domain.Query{
				Statements: []domain.Statement{{Method: "from", Resource: "hero", CacheControl: domain.CacheControl{MaxAge: 200, SMaxAge: 400}}},
			},
		},
		{
			"resolve variable in max-age/s-max-age from params",
			domain.Query{
				Statements: []domain.Statement{{Method: "from", Resource: "hero", CacheControl: domain.CacheControl{MaxAge: domain.Variable{"cache-control"}, SMaxAge: domain.Variable{"s-cache-control"}}}},
			},
			restql.QueryInput{Body: map[string]interface{}{"cache-control": "200", "s-cache-control": 400}},
			domain.Query{
				Statements: []domain.Statement{{Method: "from", Resource: "hero", CacheControl: domain.CacheControl{MaxAge: 200, SMaxAge: 400}}},
			},
		},
		{
			"resolve variable in headers from params",
			domain.Query{
				Statements: []domain.Statement{{
					Method:   "from",
					Resource: "hero",
					Headers: map[string]interface{}{
						"Authorization": domain.Chain{"done-resource", domain.Variable{"authField"}},
						"X-Some-Header": domain.Variable{"some-param"},
						"X-Id":          "1234567890",
					},
				}},
			},
			restql.QueryInput{Params: map[string]interface{}{"authField": "token", "some-param": "abc"}},
			domain.Query{
				Statements: []domain.Statement{{Method: "from", Resource: "hero", Headers: map[string]interface{}{"Authorization": domain.Chain{"done-resource", "token"}, "X-Id": "1234567890", "X-Some-Header": "abc"}}},
			},
		},
		{
			"resolve variable in headers from headers",
			domain.Query{
				Statements: []domain.Statement{{
					Method:   "from",
					Resource: "hero",
					Headers: map[string]interface{}{
						"Authorization": domain.Variable{"auth"},
						"X-Some-Header": domain.Variable{"some-param"},
						"X-Id":          "1234567890",
					},
				}},
			},
			restql.QueryInput{Headers: map[string]string{"Auth": "abcdef0987", "Some-Param": "abc"}},
			domain.Query{
				Statements: []domain.Statement{{Method: "from", Resource: "hero", Headers: map[string]interface{}{"Authorization": "abcdef0987", "X-Id": "1234567890", "X-Some-Header": "abc"}}},
			},
		},
		{
			"resolve variable in headers from body",
			domain.Query{
				Statements: []domain.Statement{{
					Method:   "from",
					Resource: "hero",
					Headers: map[string]interface{}{
						"Authorization": domain.Variable{"auth"},
						"X-Some-Header": domain.Variable{"some-param"},
						"X-Id":          "1234567890",
					},
				}},
			},
			restql.QueryInput{Body: map[string]interface{}{"auth": "abcdef0987", "some-param": "abc"}},
			domain.Query{
				Statements: []domain.Statement{{Method: "from", Resource: "hero", Headers: map[string]interface{}{"Authorization": "abcdef0987", "X-Id": "1234567890", "X-Some-Header": "abc"}}},
			},
		},
		{
			"resolve variable in function on only clause from params",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", Only: []interface{}{
				domain.NewFilterByRegex("name", domain.Variable{Target: "nameField"}, domain.Variable{Target: "namePattern"}),
				domain.Match{Value: "name", Args: []domain.Arg{{Name: domain.MatchArgRegex, Value: domain.Variable{Target: "heroName"}}}},
			}}}},
			restql.QueryInput{Params: map[string]interface{}{"nameField": "profile.name", "namePattern": "^Super", "heroName": "^Super"}},
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", Only: []interface{}{
				domain.NewFilterByRegex("name", "profile.name", "^Super"),
				domain.Match{Value: "name", Args: []domain.Arg{{Name: domain.MatchArgRegex, Value: "^Super"}}},
			}}}},
		},
		{
			"resolve variable in function on only clause from header",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", Only: []interface{}{
				domain.NewFilterByRegex("name", domain.Variable{Target: "nameField"}, domain.Variable{Target: "namePattern"}),
				domain.Match{Value: "name", Args: []domain.Arg{{Name: domain.MatchArgRegex, Value: domain.Variable{Target: "heroName"}}}},
			}}}},
			restql.QueryInput{Headers: map[string]string{"Namefield": "profile.name", "Namepattern": "^Super", "Heroname": "^Super"}},
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", Only: []interface{}{
				domain.NewFilterByRegex("name", "profile.name", "^Super"),
				domain.Match{Value: "name", Args: []domain.Arg{{Name: domain.MatchArgRegex, Value: "^Super"}}},
			}}}},
		},
		{
			"resolve variable in function on only clause from body",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", Only: []interface{}{
				domain.NewFilterByRegex("name", domain.Variable{Target: "nameField"}, domain.Variable{Target: "namePattern"}),
				domain.Match{Value: "name", Args: []domain.Arg{{Name: domain.MatchArgRegex, Value: domain.Variable{Target: "heroName"}}}},
			}}}},
			restql.QueryInput{Body: map[string]interface{}{"nameField": "profile.name", "namePattern": "^Super", "heroName": "^Super"}},
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", Only: []interface{}{
				domain.NewFilterByRegex("name", "profile.name", "^Super"),
				domain.Match{Value: "name", Args: []domain.Arg{{Name: domain.MatchArgRegex, Value: "^Super"}}},
			}}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := eval.ResolveVariables(tt.resources, tt.input)
			test.Equal(t, got, tt.expected)
		})
	}
}
