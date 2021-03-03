package parser_test

import (
	"regexp"
	"testing"

	"github.com/b2wdigital/restQL-golang/v5/internal/domain"
	"github.com/b2wdigital/restQL-golang/v5/internal/parser"
	"github.com/b2wdigital/restQL-golang/v5/test"
)

func TestQueryParser(t *testing.T) {
	tests := []struct {
		name     string
		expected interface{}
		query    string
	}{
		{

			"Unique from statement",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero"}}},
			"from hero",
		},
		{
			"Multiple from statement",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero"}, {Method: "from", Resource: "sidekick"}}},
			`from hero
					from sidekick`,
		},
		{
			"Unique from statement and with parameters",
			domain.Query{Statements: []domain.Statement{{
				Method:   "from",
				Resource: "hero",
				With: domain.Params{
					Values: map[string]interface{}{
						"id":      1,
						"name":    "batman",
						"weapons": []interface{}{"belt", "hands"},
						"family": map[string]interface{}{
							"father":       "Thomas Wayne",
							"grandparents": []interface{}{"John", "William"},
							"familyName":   domain.Variable{"familyName"},
						},
						"height":    10.5,
						"universe":  domain.Variable{"universe"},
						"sorted":    true,
						"timestamp": nil,
					},
				},
			}}},
			`from hero with id = 1, name = "batman", weapons = ["belt", "hands"], family = { "father": "Thomas Wayne", "grandparents": ["John", "William"], "familyName": $familyName }, height = 10.5, universe = $universe, sorted = true, timestamp = null`,
		},
		{
			"Unique from statement and chained with parameters",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": domain.Chain{"done-resource", "id"}}}}}},
			"from hero with id = done-resource.id",
		},
		{
			"Unique from statement and list of chained with parameters",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": []interface{}{domain.Chain{"done-resource", "id"}}}}}}},
			"from hero with id = [done-resource.id]",
		},
		{
			"Unique from statement and parameterized chained with parameters",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": domain.Chain{"done-resource", domain.Variable{"field"}, "id"}}}}}},
			"from hero with id = done-resource.$field.id",
		},
		{
			"Unique from statement and only filters",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", Only: []interface{}{[]string{"name"}, []string{"weapons"}}}}},
			"from hero only name, weapons",
		},
		{
			"Unique from statement and hidden filter",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", Hidden: true}}},
			"from hero hidden",
		},
		{
			"Unique from statement and ignore errors flag",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", IgnoreErrors: true}}},
			"from hero ignore-errors",
		},
		{
			"Unique from statement and fixed timeout",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", Timeout: 2000}}},
			"from hero timeout 2000",
		},
		{
			"Unique from statement and variable timeout",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", Timeout: domain.Variable{"some-time"}}}},
			"from hero timeout $some-time",
		},
		{
			"Unique from statement and headers",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", Headers: map[string]interface{}{"X-Trace-Id": "12345"}}}},
			`from hero headers X-Trace-Id = "12345"`,
		},
		{
			"Unique from statement and variable header",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", Headers: map[string]interface{}{"X-Trace-Id": domain.Variable{"traceId"}}}}},
			`from hero headers X-Trace-Id = $traceId`,
		},
		{
			"Unique from statement and chained header",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", Headers: map[string]interface{}{"X-Trace-Id": domain.Chain{"done-resource", "traceId"}}}}},
			`from hero headers X-Trace-Id = done-resource.traceId`,
		},
		{
			"Unique from statement and max age",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", CacheControl: domain.CacheControl{MaxAge: 2000, SMaxAge: 4000}}}},
			"from hero max-age 2000 s-max-age 4000",
		},
		{
			"Unique from statement and variable max age",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", CacheControl: domain.CacheControl{MaxAge: domain.Variable{"maxAge"}, SMaxAge: domain.Variable{"sMaxAge"}}}}},
			"from hero max-age $maxAge s-max-age $sMaxAge",
		},
		{
			"Unique from statement and flattened list parameters",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": domain.NoMultiplex{[]interface{}{1, 2}}}}}}},
			"from hero with id = [1, 2] -> no-multiplex",
		},
		{
			"Unique from statement with multiple functions applied to parameter",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": domain.NoMultiplex{domain.JSON{[]interface{}{1, 2}}}}}}}},
			"from hero with id = [1, 2] -> json -> no-multiplex",
		},
		{
			"Unique from statement and object parameter encoded as json",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": domain.JSON{map[string]interface{}{"internal": 1}}}}}}},
			`from hero with id = { "internal": 1 } -> json`,
		},
		{
			"Unique from statement and parameter encoded as base64",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": domain.Base64{"abdcef12345"}}}}}},
			`from hero with id = "abdcef12345" -> base64`,
		},
		{
			"Unique from statement and parameter defined as body",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": domain.AsBody{Value: []interface{}{map[string]interface{}{"registryNumber": "abdcef12345"}}}}}}}},
			`from hero with id = [{"registryNumber": "abdcef12345"}] -> as-body`,
		},
		{
			"Unique from statement and parameter flattened",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": domain.Flatten{[]interface{}{[]interface{}{1}, []interface{}{2}, []interface{}{3}}}}}}}},
			`from hero with id = [[1], [2], [3]] -> flatten`,
		},
		{
			"Unique to statement with default body value and custom parameter",
			domain.Query{Statements: []domain.Statement{{Method: "to", Resource: "hero", With: domain.Params{Body: domain.Variable{"hero"}, Values: map[string]interface{}{"name": "batman"}}}}},
			`to hero with $hero, name = "batman"`,
		},
		{
			"Unique to statement with default body flattened value and custom parameter",
			domain.Query{Statements: []domain.Statement{{Method: "to", Resource: "hero", With: domain.Params{Body: domain.NoMultiplex{Value: domain.Variable{"hero"}}, Values: map[string]interface{}{"name": "batman"}}}}},
			`to hero with $hero -> no-multiplex, name = "batman"`,
		},
		{
			"Unique from statement and only filters with match function",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", Only: []interface{}{domain.Match{Value: []string{"name"}, Arg: regexp.MustCompile("^Super")}, []string{"weapons"}}}}},
			`from hero only name -> matches("^Super"), weapons`,
		},
		{
			"Unique from statement and only filters with match function using variable as argument",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", Only: []interface{}{domain.Match{Value: []string{"name"}, Arg: domain.Variable{Target: "heroName"}}, []string{"weapons"}}}}},
			`from hero only name -> matches($heroName), weapons`,
		},
		{
			"Unique from statement with aggregation",
			domain.Query{Statements: []domain.Statement{
				{Method: "from", Resource: "hero"},
				{Method: "from", Resource: "sidekick", In: []string{"hero", "sidekick"}},
			}},
			`
					from hero
					from sidekick in hero.sidekick
			`,
		},
		{
			"Full query",
			domain.Query{
				Use: map[string]interface{}{"max-age": 600},
				Statements: []domain.Statement{
					{
						Method:   "from",
						Resource: "hero",
						Alias:    "h",
						Headers:  map[string]interface{}{"X-Trace-Id": "abcdef12345"},
						With: domain.Params{Values: map[string]interface{}{
							"id":      1,
							"name":    "batman",
							"weapons": []interface{}{"belt", "hands"},
							"family":  map[string]interface{}{"father": "Thomas Wayne"},
							"height":  10.5,
						}},
						Only: []interface{}{[]string{"id"}, []string{"name"}},
					},
					{
						Method:   "from",
						Resource: "sidekick",
						Alias:    "s",
						In:       []string{"hero", "sidekick"},
						Headers:  map[string]interface{}{"X-Trace-Id": "abcdef12345"},
						With: domain.Params{Values: map[string]interface{}{
							"id":      1,
							"name":    "batman",
							"weapons": []interface{}{"belt", "hands"},
							"family":  map[string]interface{}{"father": "Thomas Wayne"},
							"height":  10.5,
						}},
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
			`use max-age 600
				from hero as h
					headers
						X-Trace-Id = "abcdef12345"
					with
						id = 1
						name = "batman"
						weapons = ["belt", "hands"]
						family = { "father": "Thomas Wayne" }
						height = 10.5
					only
						id
						name
		
				 from sidekick as s in hero.sidekick
					headers
						X-Trace-Id = "abcdef12345"
					with
						id = 1
						name = "batman"
						weapons = ["belt", "hands"]
						family = { "father": "Thomas Wayne" }
						height = 10.5
					hidden
					ignore-errors
		
				 from villain as v`,
		},
		{
			"Unique from statement with no explode applied to parameter",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", With: domain.Params{Values: map[string]interface{}{"info": domain.NoExplode{Value: map[string]interface{}{"weapons": []interface{}{"batrang", "batbelt"}}}}}}}},
			`from hero with info = {weapons: ["batrang", "batbelt"]} -> no-explode`,
		},
		{
			"Multiple statements with second using depends on first",
			domain.Query{Statements: []domain.Statement{
				{Method: "from", Resource: "hero"},
				{Method: "from", Resource: "sidekick", DependsOn: domain.DependsOn{Target: "hero"}},
			}},
			`
					from hero
					from sidekick
						depends-on hero
			`,
		},
	}

	queryParser, err := parser.New()
	test.VerifyError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := queryParser.Parse(tt.query)

			test.VerifyError(t, err)
			test.Equal(t, got, tt.expected)
		})
	}
}

func BenchmarkParse(b *testing.B) {
	query := `
from hero as h
	with id = $id

from sidekick as s
	with
		heroId = hero.id
		ageA = 18

from hero as h
	with id = $id


from hero as h
	with id = $id
`

	queryParser, err := parser.New()
	if err != nil {
		b.Fatalf("failed to compile the queryParser : %v", err)
	}

	for i := 0; i < b.N; i++ {
		_, err := queryParser.Parse(query)

		if err != nil {
			b.Fatalf("An error occurred when running the benchmark: %v", err)
		}
	}
}
