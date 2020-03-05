package parser_test

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
	. "github.com/b2wdigital/restQL-golang/internal/parser"
	"reflect"
	"testing"
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
			"from hero from sidekick",
		},
		{
			"Unique from statement and with parameters",
			domain.Query{Statements: []domain.Statement{{
				Method:   "from",
				Resource: "hero",
				With: domain.Params{
					"id":       1,
					"name":     "batman",
					"weapons":  []interface{}{"belt", "hands"},
					"family":   map[string]interface{}{"father": "Thomas Wayne"},
					"height":   10.5,
					"universe": domain.Variable{"universe"},
				},
			}}},
			`from hero with id = 1, name = "batman", weapons = ["belt", "hands"], family = { "father": "Thomas Wayne" }, height = 10.5, universe = $universe`,
		},
		{
			"Unique from statement and chained with parameters",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", With: domain.Params{"id": domain.Chain{"done-resource", "id"}}}}},
			"from hero with id = done-resource.id",
		},
		{
			"Unique from statement and list of chained with parameters",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", With: domain.Params{"id": []interface{}{domain.Chain{"done-resource", "id"}}}}}},
			"from hero with id = [done-resource.id]",
		},
		{
			"Unique from statement and parameterized chained with parameters",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", With: domain.Params{"id": domain.Chain{"done-resource", domain.Variable{"field"}, "id"}}}}},
			"from hero with id = done-resource.$field.id",
		},
		{
			"Unique from statement and only filters",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", Only: []interface{}{"name", "weapons"}}}},
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
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", Timeout: "some-time"}}},
			"from hero timeout $some-time",
		},
		{
			"Unique from statement and headers",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", Headers: map[string]interface{}{"X-Trace-Id": "12345"}}}},
			`from hero headers X-Trace-Id = "12345"`,
		},
		{
			"Unique from statement and max age",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", CacheControl: domain.CacheControl{MaxAge: 2000, SMaxAge: 4000}}}},
			"from hero max-age = 2000 s-max-age = 4000",
		},
		{
			"Unique from statement and flattened list parameters",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", With: domain.Params{"id": domain.Flatten{[]interface{}{1, 2}}}}}},
			"from hero with id = [1, 2] -> flatten",
		},
		{
			"Unique from statement and object parameter encoded as json",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", With: domain.Params{"id": domain.Json{map[string]interface{}{"internal": 1}}}}}},
			`from hero with id = { "internal": 1 } -> json`,
		},
		{
			"Unique from statement and parameter encoded as base64",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", With: domain.Params{"id": domain.Base64{"abdcef12345"}}}}},
			`from hero with id = "abdcef12345" -> base64`,
		},
		{
			"Unique from statement and only filters with match function",
			domain.Query{Statements: []domain.Statement{{Method: "from", Resource: "hero", Only: []interface{}{domain.Match{Target: "name", Arg: "^Super"}, "weapons"}}}},
			`from hero only name -> matches("^Super"), weapons`,
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
						With: map[string]interface{}{
							"id":      1,
							"name":    "batman",
							"weapons": []interface{}{"belt", "hands"},
							"family":  map[string]interface{}{"father": "Thomas Wayne"},
							"height":  10.5,
						},
						Only: []interface{}{"id", "name"},
					},
					{
						Method:   "from",
						Resource: "sidekick",
						Alias:    "s",
						Headers:  map[string]interface{}{"X-Trace-Id": "abcdef12345"},
						With: map[string]interface{}{
							"id":      1,
							"name":    "batman",
							"weapons": []interface{}{"belt", "hands"},
							"family":  map[string]interface{}{"father": "Thomas Wayne"},
							"height":  10.5,
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
			`use max-age = 600
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
		
				 from sidekick as s
					with
						id = 1
						name = "batman"
						weapons = ["belt", "hands"]
						family = { "father": "Thomas Wayne" }
						height = 10.5
					headers
						X-Trace-Id = "abcdef12345"
					hidden
					ignore-errors
		
				 from villain as v`,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := Parse(testCase.query)

			if err != nil {
				t.Errorf("An error occured during Parse\n%s", err)
			}

			if !reflect.DeepEqual(got, testCase.expected) {
				t.Errorf("Parse = %#v,\nwant %#v", got, testCase.expected)
			}
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

	for i := 0; i < b.N; i++ {
		_, err := Parse(query)

		if err != nil {
			b.Fatalf("An error ocurred when running the benchmark: %v", err)
		}
	}
}
