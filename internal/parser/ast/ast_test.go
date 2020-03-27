package ast

import (
	"fmt"
	"reflect"
	"testing"
)

func String(s string) *string {
	return &s
}

func Int(i int) *int {
	return &i
}

func Float(f float64) *float64 {
	return &f
}

func TestAstParser(t *testing.T) {
	tests := []struct {
		name     string
		expected interface{}
		query    string
	}{
		{
			"Simple from resource query",
			Query{Blocks: []Block{{Method: FromMethod, Resource: "cart"}}},
			"from cart",
		},
		{
			"Query with two from statements",
			Query{Blocks: []Block{{Method: FromMethod, Resource: "cart"}, {Method: FromMethod, Resource: "hero"}}},
			"from cart from hero",
		},
		{
			"Simple from resource query with comment",
			Query{Blocks: []Block{{Method: FromMethod, Resource: "cart"}}},
			`
				// a comment
				from cart // some other comment
				`,
		},
		{
			"Simple from resource query with use modifier",
			Query{Use: []Use{{Key: "max-age", Value: UseValue{Int: Int(600)}}, {Key: "s-max-age", Value: UseValue{Int: Int(400)}}, {Key: "timeout", Value: UseValue{Int: Int(8000)}}}, Blocks: []Block{{Method: FromMethod, Resource: "cart"}}},
			`use max-age = 600 use s-max-age = 400 use timeout = 8000 from cart`,
		},
		{
			"Simple from resource query with resource name with hyphen",
			Query{Blocks: []Block{{Method: FromMethod, Resource: "shopping-cart"}}},
			"from shopping-cart",
		},
		{
			"Simple from resource query with alias",
			Query{Blocks: []Block{{Method: FromMethod, Resource: "cart", Alias: "shopping"}}},
			"from cart as shopping",
		},
		{
			"Get query with integer query parameters",
			Query{Blocks: []Block{{Method: FromMethod, Resource: "cart", Qualifiers: []Qualifier{{With: []WithItem{{Key: "id", Value: Value{Primitive: &Primitive{Int: Int(1)}}}}}}}}},
			`from cart with id = 1`,
		},
		{
			"Get query with string query parameters",
			Query{Blocks: []Block{{Method: FromMethod, Resource: "cart", Qualifiers: []Qualifier{{With: []WithItem{{Key: "name", Value: Value{Primitive: &Primitive{String: String("batman")}}}}}}}}},
			`from cart with name = "batman"`,
		},
		{
			"Get query with multiple parameters delimited by comma",
			Query{Blocks: []Block{{Method: FromMethod, Resource: "hero", Qualifiers: []Qualifier{{With: []WithItem{
				{Key: "id", Value: Value{Primitive: &Primitive{Int: Int(1)}}},
				{Key: "name", Value: Value{Primitive: &Primitive{String: String("batman")}}},
			}}}}}},
			`from hero with id = 1, name = "batman"`,
		},
		{
			"Get query with multiple parameters delimited by space",
			Query{Blocks: []Block{{Method: FromMethod, Resource: "hero", Qualifiers: []Qualifier{{With: []WithItem{
				{Key: "id", Value: Value{Primitive: &Primitive{Int: Int(1)}}},
				{Key: "name", Value: Value{Primitive: &Primitive{String: String("batman")}}},
			}}}}}},
			`from hero with id = 1 name = "batman"`,
		},
		{
			"Get query with multiple parameters delimited by new line",
			Query{Blocks: []Block{{Method: FromMethod, Resource: "hero", Qualifiers: []Qualifier{{With: []WithItem{
				{Key: "id", Value: Value{Primitive: &Primitive{Int: Int(1)}}},
				{Key: "name", Value: Value{Primitive: &Primitive{String: String("batman")}}},
			}}}}}},
			fmt.Sprintf("from hero with id = 1\nname = \"batman\""),
		},
		{
			"Get query with float query parameters",
			Query{Blocks: []Block{{Method: FromMethod, Resource: "hero", Qualifiers: []Qualifier{{With: []WithItem{{Key: "height", Value: Value{Primitive: &Primitive{Float: Float(10.5)}}}}}}}}},
			`from hero with height = 10.5`,
		},
		{
			"Get query with list query parameters",
			Query{Blocks: []Block{{Method: FromMethod, Resource: "hero", Qualifiers: []Qualifier{{With: []WithItem{{Key: "weapons", Value: Value{List: []*Value{{Primitive: &Primitive{String: String("sword")}}, {Primitive: &Primitive{String: String("shield")}}}}}}}}}}},
			`from hero with weapons = ["sword", "shield"]`,
		},
		{
			"Get query with nested list query parameters",
			Query{Blocks: []Block{{Method: FromMethod, Resource: "hero", Qualifiers: []Qualifier{{With: []WithItem{{Key: "weapons", Value: Value{List: []*Value{{List: []*Value{{Primitive: &Primitive{String: String("sword")}}}}, {List: []*Value{{Primitive: &Primitive{String: String("shield")}}}}}}}}}}}}},
			`from hero with weapons = [["sword"], ["shield"]]`,
		},
		{
			"Get query with list query parameters delimited by new line",
			Query{Blocks: []Block{{Method: FromMethod, Resource: "hero", Qualifiers: []Qualifier{{With: []WithItem{{Key: "weapons", Value: Value{List: []*Value{{Primitive: &Primitive{String: String("sword")}}, {Primitive: &Primitive{String: String("shield")}}}}}}}}}}},
			fmt.Sprintf("from hero with weapons = [\"sword\"\n\"shield\"]"),
		},
		{
			"Get query with object query parameters",
			Query{Blocks: []Block{{Method: FromMethod, Resource: "hero", Qualifiers: []Qualifier{{With: []WithItem{{Key: "id", Value: Value{Object: []ObjectEntry{{Key: "id", Value: ObjectValue{Primitive: &Primitive{String: String("1")}}}}}}}}}}}},
			`from hero with id = { "id": "1" }`,
		},
		{
			"Get query with object query parameters with multiple key/values",
			Query{Blocks: []Block{{Method: FromMethod, Resource: "hero", Qualifiers: []Qualifier{{With: []WithItem{{Key: "id", Value: Value{Object: []ObjectEntry{{Key: "id", Value: ObjectValue{Primitive: &Primitive{String: String("1")}}}, {Key: "name", Value: ObjectValue{Primitive: &Primitive{String: String("batman")}}}}}}}}}}}},
			`from hero with id = { "id": "1", "name": "batman" }`,
		},
		{
			"Get query with object query parameters with multiple key/values delimited by new line",
			Query{Blocks: []Block{{Method: FromMethod, Resource: "hero", Qualifiers: []Qualifier{{With: []WithItem{{Key: "id", Value: Value{Object: []ObjectEntry{{Key: "id", Value: ObjectValue{Primitive: &Primitive{String: String("1")}}}, {Key: "name", Value: ObjectValue{Primitive: &Primitive{String: String("batman")}}}}}}}}}}}},
			fmt.Sprintf("from hero with id = { \"id\": \"1\"\n\"name\": \"batman\" }"),
		},
		{
			"Get query with nested object query parameters",
			Query{Blocks: []Block{{Method: FromMethod, Resource: "hero", Qualifiers: []Qualifier{{With: []WithItem{{Key: "id", Value: Value{Object: []ObjectEntry{{Key: "id", Value: ObjectValue{Nested: []ObjectEntry{{Key: "internal", Value: ObjectValue{Primitive: &Primitive{String: String("1")}}}, {Key: "version", Value: ObjectValue{Primitive: &Primitive{Int: Int(10)}}}}}}}}}}}}}}},
			`from hero with id = { "id": { "internal": "1", "version": 10 } }`,
		},
		{
			"Get query with chained query parameters",
			Query{Blocks: []Block{{Method: FromMethod, Resource: "hero", Qualifiers: []Qualifier{{With: []WithItem{{Key: "id", Value: Value{Primitive: &Primitive{Chain: []Chained{{PathItem: "done-resource"}, {PathItem: "id"}}}}}}}}}}},
			`from hero with id = done-resource.id`,
		},
		{
			"Get query with variable query parameters",
			Query{Blocks: []Block{{Method: FromMethod, Resource: "hero", Qualifiers: []Qualifier{{With: []WithItem{{Key: "id", Value: Value{Variable: String("id")}}}}}}}},
			`from hero with id = $id`,
		},
		{
			"Get query with variable chained query parameters",
			Query{Blocks: []Block{{Method: FromMethod, Resource: "hero", Qualifiers: []Qualifier{{With: []WithItem{{Key: "id", Value: Value{Primitive: &Primitive{Chain: []Chained{{PathItem: "done-resource"}, {PathVariable: "path"}, {PathItem: "id"}}}}}}}}}}},
			`from hero with id = done-resource.$path.id`,
		},
		{
			"Get query with select filter",
			Query{Blocks: []Block{{Method: FromMethod, Resource: "hero", Qualifiers: []Qualifier{{Only: []Filter{{Field: []string{"name"}}, {Field: []string{"weapons"}}}}}}}},
			"from hero only name, weapons",
		},
		{
			"Get query with select filter delimited by new line",
			Query{Blocks: []Block{{Method: FromMethod, Resource: "hero", Qualifiers: []Qualifier{{Only: []Filter{{Field: []string{"name"}}, {Field: []string{"weapons"}}}}}}}},
			fmt.Sprintf("from hero only name\nweapons"),
		},
		{
			"Multiple gets with select filter",
			Query{Blocks: []Block{
				{Method: FromMethod, Resource: "hero", Qualifiers: []Qualifier{{Only: []Filter{{Field: []string{"name"}}, {Field: []string{"weapons"}}}}}},
				{Method: FromMethod, Resource: "sidekick"},
			}},
			"from hero only name, weapons from sidekick",
		},
		{
			"Get query with hidden",
			Query{Blocks: []Block{{Method: FromMethod, Resource: "hero", Qualifiers: []Qualifier{{Hidden: true}}}}},
			"from hero hidden",
		},
		{
			"Get query with ignore errors",
			Query{Blocks: []Block{{Method: FromMethod, Resource: "hero", Qualifiers: []Qualifier{{IgnoreErrors: true}}}}},
			"from hero ignore-errors",
		},
		{
			"Get query with integer timeout",
			Query{Blocks: []Block{{Method: FromMethod, Resource: "hero", Qualifiers: []Qualifier{{Timeout: &TimeoutValue{Int: Int(200)}}}}}},
			`from hero timeout 200`,
		},
		{
			"Get query with variable timeout",
			Query{Blocks: []Block{{Method: FromMethod, Resource: "hero", Qualifiers: []Qualifier{{Timeout: &TimeoutValue{Variable: String("some-time")}}}}}},
			`from hero timeout $some-time`,
		},
		{
			"Get query with headers",
			Query{Blocks: []Block{{Method: FromMethod, Resource: "hero", Qualifiers: []Qualifier{{Headers: []HeaderItem{
				{Key: "Authorization", Value: HeaderValue{String: String("abcdef12345")}},
				{Key: "X-Trace-Id", Value: HeaderValue{Variable: String("trace-id")}},
			}}}}}},
			`from hero headers Authorization = "abcdef12345", X-Trace-Id = $trace-id`,
		},
		{
			"Get query with max age",
			Query{Blocks: []Block{{Method: FromMethod, Resource: "hero", Qualifiers: []Qualifier{{MaxAge: &MaxAgeValue{Int: Int(2000)}}}}}},
			`from hero max-age = 2000`,
		},
		{
			"Get query with cache control",
			Query{Blocks: []Block{{Method: FromMethod, Resource: "hero", Qualifiers: []Qualifier{{MaxAge: &MaxAgeValue{Int: Int(2000)}}}}}},
			`from hero cache-control = 2000`,
		},
		{
			"Get query with cache control",
			Query{Blocks: []Block{{Method: FromMethod, Resource: "hero", Qualifiers: []Qualifier{{SMaxAge: &SMaxAgeValue{Int: Int(2000)}}}}}},
			`from hero s-max-age = 2000`,
		},
		{
			"Get query with list query parameters flattened",
			Query{Blocks: []Block{{
				Method:   FromMethod,
				Resource: "hero",
				Qualifiers: []Qualifier{
					{With: []WithItem{
						{
							Key: "weapons",
							Value: Value{List: []*Value{
								{Primitive: &Primitive{String: String("sword")}},
								{Primitive: &Primitive{String: String("shield")}},
							}},
							Flatten: true,
						}}},
				},
			}}},
			`from hero with weapons = ["sword", "shield"] -> flatten`,
		},
		{
			"Get query with chained query parameters flattened",
			Query{Blocks: []Block{{
				Method:   FromMethod,
				Resource: "hero",
				Qualifiers: []Qualifier{{
					With: []WithItem{
						{
							Key:     "id",
							Value:   Value{Primitive: &Primitive{Chain: []Chained{{PathItem: "done-resource"}, {PathItem: "id"}}}},
							Flatten: true,
						},
					},
				}},
			}}},
			`from hero with id = done-resource.id -> flatten`,
		},
		{
			"Get query with query parameters encoded in base64",
			Query{Blocks: []Block{{
				Method:   FromMethod,
				Resource: "hero",
				Qualifiers: []Qualifier{{
					With: []WithItem{
						{
							Key:    "id",
							Value:  Value{Primitive: &Primitive{String: String("abcdefg12345")}},
							Base64: true,
						},
					},
				}},
			}}},
			`from hero with id = "abcdefg12345" -> base64`,
		},
		{
			"Get query with object query parameters encoded as json",
			Query{Blocks: []Block{{
				Method:   FromMethod,
				Resource: "hero",
				Qualifiers: []Qualifier{{
					With: []WithItem{{
						Key:   "id",
						Value: Value{Object: []ObjectEntry{{Key: "id", Value: ObjectValue{Primitive: &Primitive{String: String("1")}}}}},
						Json:  true,
					}},
				}},
			}}},
			`from hero with id = { "id": "1" } -> json`,
		},
		{
			"Get query with select filters and match function",
			Query{Blocks: []Block{{
				Method:   FromMethod,
				Resource: "hero",
				Qualifiers: []Qualifier{
					{Only: []Filter{
						{Field: []string{"name"}, Match: "^Super"},
						{Field: []string{"weapons"}},
					}},
				},
			}}},
			`from hero
					only
						name -> matches("^Super")
						weapons`,
		},
		{
			"Full query",
			Query{Blocks: []Block{
				{
					Method:   FromMethod,
					Resource: "hero",
					Alias:    "h",
					Qualifiers: []Qualifier{
						{
							Headers: []HeaderItem{
								{Key: "X-Trace-Id", Value: HeaderValue{String: String("abcdef12345")}},
							},
						},
						{
							Only: []Filter{{Field: []string{"id"}}, {Field: []string{"name"}}},
						},
						{
							With: []WithItem{
								{Key: "id", Value: Value{Primitive: &Primitive{Int: Int(1)}}},
								{Key: "name", Value: Value{Primitive: &Primitive{String: String("batman")}}},
								{
									Key: "weapons",
									Value: Value{List: []*Value{
										{Primitive: &Primitive{String: String("belt")}}, {Primitive: &Primitive{String: String("hands")}},
									}},
								},
								{
									Key: "family",
									Value: Value{Object: []ObjectEntry{
										{Key: "father", Value: ObjectValue{Primitive: &Primitive{String: String("Thomas Wayne")}}},
									}},
								},
								{Key: "height", Value: Value{Primitive: &Primitive{Float: Float(10.5)}}},
								{Key: "var", Value: Value{Variable: String("myvar")}},
							},
						},
					},
				},
				{
					Method:   FromMethod,
					Resource: "sidekick",
					Alias:    "s",
					Qualifiers: []Qualifier{
						{
							With: []WithItem{
								{Key: "id", Value: Value{Primitive: &Primitive{Int: Int(1)}}},
								{Key: "name", Value: Value{Primitive: &Primitive{String: String("batman")}}},
								{
									Key: "weapons",
									Value: Value{List: []*Value{
										{Primitive: &Primitive{String: String("belt")}}, {Primitive: &Primitive{String: String("hands")}},
									}},
								},
								{
									Key: "family",
									Value: Value{Object: []ObjectEntry{
										{Key: "father", Value: ObjectValue{Primitive: &Primitive{String: String("Thomas Wayne")}}},
									}},
								},
								{Key: "height", Value: Value{Primitive: &Primitive{Float: Float(10.5)}}}},
						},
						{
							Headers: []HeaderItem{
								{Key: "X-Trace-Id", Value: HeaderValue{String: String("abcdef12345")}},
							},
						},
						{Hidden: true},
						{IgnoreErrors: true},
					},
				},
				{
					Method:   FromMethod,
					Resource: "villain",
					Alias:    "v",
				},
			}},
			`from hero as h
					headers
						X-Trace-Id = "abcdef12345"
					only
						id
						name
					with
						id = 1
						name = "batman"
						weapons = ["belt", "hands"]
						family = { "father": "Thomas Wayne" }
						height = 10.5
						var = $myvar
		
		
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
		{
			"Bug: words starting with `to` failed parser because they would match a method keyword",
			Query{Blocks: []Block{
				{
					Method:   FromMethod,
					Resource: "cart",
					Qualifiers: []Qualifier{
						{
							With: []WithItem{
								{Key: "id", Value: Value{Variable: String("id")}},
							},
						},
						{
							Only: []Filter{{Field: []string{"total"}}, {Field: []string{"opn"}}},
						},
					},
				},
				{
					Method:   FromMethod,
					Resource: "sku",
					Qualifiers: []Qualifier{
						{
							With: []WithItem{
								{Key: "sku", Value: Value{Primitive: &Primitive{Chain: []Chained{{PathItem: "cart"}, {PathItem: "lines"}, {PathItem: "productSku"}}}}},
							},
						},
					},
				},
			}},
			`
		from cart
		with 
			id = $id
		only	
			total
			opn
			
		
		from sku
		with
			sku = cart.lines.productSku
`,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := Parse(testCase.query)

			if err != nil {
				t.Errorf("Parse returned an unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(*got, testCase.expected) {
				t.Errorf("Parse = %#v,\nwant %#v", *got, testCase.expected)
			}
		})
	}
}
