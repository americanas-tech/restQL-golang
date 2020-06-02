package ast_test

import (
	"fmt"
	"github.com/b2wdigital/restQL-golang/internal/parser/ast"
	"github.com/b2wdigital/restQL-golang/test"
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

func TestAstGenerator(t *testing.T) {
	tests := []struct {
		name     string
		expected interface{}
		query    string
	}{
		{
			"Simple from resource query",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "cart"}}},
			"from cart",
		},
		{
			"Query with two from statements",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "cart"}, {Method: ast.FromMethod, Resource: "hero"}}},
			"from cart from hero",
		},
		{
			"Simple from resource query with comment",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "cart"}}},
			`
				// a comment
				from cart // some other comment
				`,
		},
		{
			"Simple from resource query with use modifier",
			ast.Query{Use: []ast.Use{{Key: "max-age", Value: ast.UseValue{Int: Int(600)}}, {Key: "s-max-age", Value: ast.UseValue{Int: Int(400)}}, {Key: "timeout", Value: ast.UseValue{Int: Int(8000)}}}, Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "cart"}}},
			`use max-age 600 use s-max-age 400 use timeout 8000 from cart`,
		},
		{
			"Simple from resource query with resource name with hyphen",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "shopping-cart"}}},
			"from shopping-cart",
		},
		{
			"Simple from resource query with alias",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "cart", Alias: "shopping"}}},
			"from cart as shopping",
		},
		{
			"Get query with integer query parameters",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "cart", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{KeyValues: []ast.KeyValue{{Key: "id", Value: ast.Value{Primitive: &ast.Primitive{Int: Int(1)}}}}}}}}}},
			`from cart with id = 1`,
		},
		{
			"Get query with string query parameters",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "cart", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{KeyValues: []ast.KeyValue{{Key: "name", Value: ast.Value{Primitive: &ast.Primitive{String: String("batman")}}}}}}}}}},
			`from cart with name = "batman"`,
		},
		{
			"Get query with multiple parameters delimited by comma",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{KeyValues: []ast.KeyValue{
				{Key: "id", Value: ast.Value{Primitive: &ast.Primitive{Int: Int(1)}}},
				{Key: "name", Value: ast.Value{Primitive: &ast.Primitive{String: String("batman")}}},
			}}}}}}},
			`from hero with id = 1, name = "batman"`,
		},
		{
			"Get query with multiple parameters delimited by space",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{KeyValues: []ast.KeyValue{
				{Key: "id", Value: ast.Value{Primitive: &ast.Primitive{Int: Int(1)}}},
				{Key: "name", Value: ast.Value{Primitive: &ast.Primitive{String: String("batman")}}},
			}}}}}}},
			`from hero with id = 1 name = "batman"`,
		},
		{
			"Get query with default body value",
			ast.Query{Blocks: []ast.Block{{Method: ast.ToMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{Body: &ast.ParameterBody{Target: "hero"}}}}}}},
			`to hero with $hero`,
		},
		{
			"Get query with default body value and key-value parameter",
			ast.Query{Blocks: []ast.Block{{Method: ast.ToMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{
				Body: &ast.ParameterBody{Target: "hero"},
				KeyValues: []ast.KeyValue{
					{Key: "name", Value: ast.Value{Primitive: &ast.Primitive{String: String("batman")}}},
				},
			}}}}}},
			`to hero with $hero, name = "batman"`,
		},
		{
			"Get query with default body value flattened and key-value parameter",
			ast.Query{Blocks: []ast.Block{{Method: ast.ToMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{
				Body: &ast.ParameterBody{Target: "hero", Flatten: true},
				KeyValues: []ast.KeyValue{
					{Key: "name", Value: ast.Value{Primitive: &ast.Primitive{String: String("batman")}}},
				},
			}}}}}},
			`to hero with $hero -> flatten, name = "batman"`,
		},
		{
			"Get query with multiple parameters delimited by comma and using key with dots",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{KeyValues: []ast.KeyValue{
				{Key: "id", Value: ast.Value{Primitive: &ast.Primitive{Int: Int(1)}}},
				{Key: "name", Value: ast.Value{Primitive: &ast.Primitive{String: String("batman")}}},
			}}}}}}},
			fmt.Sprintf("from hero with id = 1\nname = \"batman\""),
		},
		{
			"Get query with string query parameters and key using dot",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "cart", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{KeyValues: []ast.KeyValue{
				{Key: "hero.name", Value: ast.Value{Primitive: &ast.Primitive{String: String("batman")}}},
			}}}}}}},
			`from cart with hero.name = "batman"`,
		},
		{
			"Get query with multiple parameters keys using dot",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{KeyValues: []ast.KeyValue{
				{Key: "hero.id", Value: ast.Value{Primitive: &ast.Primitive{Int: Int(1)}}},
				{Key: "hero.height", Value: ast.Value{Primitive: &ast.Primitive{Float: Float(10.5)}}},
			}}}}}}},
			`from hero with hero.id = 1, hero.height = 10.5`,
		},
		{
			"Get query with multiple parameters delimited by space and using keys with dots",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{KeyValues: []ast.KeyValue{
				{Key: "hero.id", Value: ast.Value{Primitive: &ast.Primitive{Int: Int(1)}}},
				{Key: "hero.name", Value: ast.Value{Primitive: &ast.Primitive{String: String("batman")}}},
			}}}}}}},
			`from hero with hero.id = 1 hero.name = "batman"`,
		},
		{
			"Get query with multiple parameters delimited by new line and using keys with dots",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{KeyValues: []ast.KeyValue{
				{Key: "hero.id", Value: ast.Value{Primitive: &ast.Primitive{Int: Int(1)}}},
				{Key: "hero.name", Value: ast.Value{Primitive: &ast.Primitive{String: String("batman")}}},
			}}}}}}},
			fmt.Sprintf("from hero with hero.id = 1\nhero.name = \"batman\""),
		},
		{
			"Get query with float query parameters",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{KeyValues: []ast.KeyValue{{Key: "height", Value: ast.Value{Primitive: &ast.Primitive{Float: Float(10.5)}}}}}}}}}},
			`from hero with height = 10.5`,
		},
		{
			"Get query with list query parameters",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{
				KeyValues: []ast.KeyValue{{Key: "weapons", Value: ast.Value{List: []*ast.Value{{Primitive: &ast.Primitive{String: String("sword")}}, {Primitive: &ast.Primitive{String: String("shield")}}}}}},
			}}}}}},
			`from hero with weapons = ["sword", "shield"]`,
		},
		{
			"Get query with nested list query parameters",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{
				KeyValues: []ast.KeyValue{{Key: "weapons", Value: ast.Value{List: []*ast.Value{{List: []*ast.Value{{Primitive: &ast.Primitive{String: String("sword")}}}}, {List: []*ast.Value{{Primitive: &ast.Primitive{String: String("shield")}}}}}}}},
			}}}}}},
			`from hero with weapons = [["sword"], ["shield"]]`,
		},
		{
			"Get query with list query parameters delimited by new line",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{
				KeyValues: []ast.KeyValue{{Key: "weapons", Value: ast.Value{List: []*ast.Value{{Primitive: &ast.Primitive{String: String("sword")}}, {Primitive: &ast.Primitive{String: String("shield")}}}}}},
			}}}}}},
			fmt.Sprintf("from hero with weapons = [\"sword\"\n\"shield\"]"),
		},
		{
			"Get query with object query parameters",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{
				KeyValues: []ast.KeyValue{{Key: "id", Value: ast.Value{Object: []ast.ObjectEntry{{Key: "id", Value: ast.ObjectValue{Primitive: &ast.Primitive{String: String("1")}}}}}}},
			}}}}}},
			`from hero with id = { "id": "1" }`,
		},
		{
			"Get query with object query parameters",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{
				KeyValues: []ast.KeyValue{{Key: "id", Value: ast.Value{Object: []ast.ObjectEntry{{Key: "id", Value: ast.ObjectValue{Primitive: &ast.Primitive{String: String("1")}}}}}}},
			}}}}}},
			`from hero with id = { id: "1" }`,
		},
		{
			"Get query with object query parameters using list value",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{
				KeyValues: []ast.KeyValue{{Key: "id", Value: ast.Value{Object: []ast.ObjectEntry{{Key: "names", Value: ast.ObjectValue{List: []*ast.ObjectValue{{Primitive: &ast.Primitive{String: String("batman")}}, {Primitive: &ast.Primitive{String: String("wonder woman")}}}}}}}}},
			}}}}}},
			`from hero with id = { names: ["batman", "wonder woman"] }`,
		},
		{
			"Get query with object query parameters using a variable value",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{
				KeyValues: []ast.KeyValue{{Key: "id", Value: ast.Value{Object: []ast.ObjectEntry{{Key: "id", Value: ast.ObjectValue{Variable: String("id")}}}}}},
			}}}}}},
			`from hero with id = { id: $id }`,
		},
		{
			"Get query with object query parameters with multiple key/values",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{
				KeyValues: []ast.KeyValue{{Key: "id", Value: ast.Value{Object: []ast.ObjectEntry{{Key: "id", Value: ast.ObjectValue{Primitive: &ast.Primitive{String: String("1")}}}, {Key: "name", Value: ast.ObjectValue{Primitive: &ast.Primitive{String: String("batman")}}}}}}},
			}}}}}},
			`from hero with id = { "id": "1", "name": "batman" }`,
		},
		{
			"Get query with object query parameters with multiple key/values delimited by new line",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{
				KeyValues: []ast.KeyValue{{Key: "id", Value: ast.Value{Object: []ast.ObjectEntry{{Key: "id", Value: ast.ObjectValue{Primitive: &ast.Primitive{String: String("1")}}}, {Key: "name", Value: ast.ObjectValue{Primitive: &ast.Primitive{String: String("batman")}}}}}}},
			}}}}}},
			fmt.Sprintf("from hero with id = { \"id\": \"1\"\n\"name\": \"batman\" }"),
		},
		{
			"Get query with nested object query parameters",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{
				KeyValues: []ast.KeyValue{{Key: "id", Value: ast.Value{Object: []ast.ObjectEntry{{Key: "id", Value: ast.ObjectValue{Nested: []ast.ObjectEntry{{Key: "internal", Value: ast.ObjectValue{Primitive: &ast.Primitive{String: String("1")}}}, {Key: "version", Value: ast.ObjectValue{Primitive: &ast.Primitive{Int: Int(10)}}}}}}}}}},
			}}}}}},
			`from hero with id = { "id": { "internal": "1", "version": 10 } }`,
		},
		{
			"Get query with chained query parameters",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{
				KeyValues: []ast.KeyValue{{Key: "id", Value: ast.Value{Primitive: &ast.Primitive{Chain: []ast.Chained{{PathItem: "done-resource"}, {PathItem: "id"}}}}}},
			}}}}}},
			`from hero with id = done-resource.id`,
		},
		{
			"Get query with variable query parameters",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{
				KeyValues: []ast.KeyValue{{Key: "id", Value: ast.Value{Variable: String("id")}}},
			}}}}}},
			`from hero with id = $id`,
		},
		{
			"Get query with variable chained query parameters",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{
				KeyValues: []ast.KeyValue{{Key: "id", Value: ast.Value{Primitive: &ast.Primitive{Chain: []ast.Chained{{PathItem: "done-resource"}, {PathVariable: "path"}, {PathItem: "id"}}}}}},
			}}}}}},
			`from hero with id = done-resource.$path.id`,
		},
		{
			"Get query with select filter",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{Only: []ast.Filter{{Field: []string{"name"}}, {Field: []string{"weapons"}}}}}}}},
			"from hero only name, weapons",
		},
		{
			"Get query with select filter delimited by new line",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{Only: []ast.Filter{{Field: []string{"name"}}, {Field: []string{"weapons"}}}}}}}},
			fmt.Sprintf("from hero only name\nweapons"),
		},
		{
			"Multiple gets with select filter",
			ast.Query{Blocks: []ast.Block{
				{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{Only: []ast.Filter{{Field: []string{"name"}}, {Field: []string{"weapons"}}}}}},
				{Method: ast.FromMethod, Resource: "sidekick"},
			}},
			"from hero only name, weapons from sidekick",
		},
		{
			"Get query with hidden",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{Hidden: true}}}}},
			"from hero hidden",
		},
		{
			"Get query with ignore errors",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{IgnoreErrors: true}}}}},
			"from hero ignore-errors",
		},
		{
			"Get query with integer timeout",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{Timeout: &ast.TimeoutValue{Int: Int(200)}}}}}},
			`from hero timeout 200`,
		},
		{
			"Get query with variable timeout",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{Timeout: &ast.TimeoutValue{Variable: String("some-time")}}}}}},
			`from hero timeout $some-time`,
		},
		{
			"Get query with headers",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{Headers: []ast.HeaderItem{
				{Key: "Authorization", Value: ast.HeaderValue{String: String("abcdef12345")}},
				{Key: "X-Trace-Id", Value: ast.HeaderValue{Variable: String("trace-id")}},
			}}}}}},
			`from hero headers Authorization = "abcdef12345", X-Trace-Id = $trace-id`,
		},
		{
			"Get query with max age",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{MaxAge: &ast.MaxAgeValue{Int: Int(2000)}}}}}},
			`from hero max-age 2000`,
		},
		{
			"Get query with cache control",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{SMaxAge: &ast.SMaxAgeValue{Int: Int(2000)}}}}}},
			`from hero s-max-age 2000`,
		},
		{
			"Get query with list query parameters flattened",
			ast.Query{Blocks: []ast.Block{{
				Method:   ast.FromMethod,
				Resource: "hero",
				Qualifiers: []ast.Qualifier{
					{
						With: &ast.Parameters{
							KeyValues: []ast.KeyValue{
								{
									Key: "weapons",
									Value: ast.Value{List: []*ast.Value{
										{Primitive: &ast.Primitive{String: String("sword")}},
										{Primitive: &ast.Primitive{String: String("shield")}},
									}},
									Flatten: true,
								},
							},
						},
					},
				},
			}}},
			`from hero with weapons = ["sword", "shield"] -> flatten`,
		},
		{
			"Get query with chained query parameters flattened",
			ast.Query{Blocks: []ast.Block{{
				Method:   ast.FromMethod,
				Resource: "hero",
				Qualifiers: []ast.Qualifier{{
					With: &ast.Parameters{
						KeyValues: []ast.KeyValue{
							{
								Key:     "id",
								Value:   ast.Value{Primitive: &ast.Primitive{Chain: []ast.Chained{{PathItem: "done-resource"}, {PathItem: "id"}}}},
								Flatten: true,
							},
						},
					},
				}},
			}}},
			`from hero with id = done-resource.id -> flatten`,
		},
		{
			"Get query with query parameters encoded in base64",
			ast.Query{Blocks: []ast.Block{{
				Method:   ast.FromMethod,
				Resource: "hero",
				Qualifiers: []ast.Qualifier{{
					With: &ast.Parameters{
						KeyValues: []ast.KeyValue{
							{
								Key:    "id",
								Value:  ast.Value{Primitive: &ast.Primitive{String: String("abcdefg12345")}},
								Base64: true,
							},
						},
					},
				}},
			}}},
			`from hero with id = "abcdefg12345" -> base64`,
		},
		{
			"Get query with object query parameters encoded as json",
			ast.Query{Blocks: []ast.Block{{
				Method:   ast.FromMethod,
				Resource: "hero",
				Qualifiers: []ast.Qualifier{{
					With: &ast.Parameters{
						KeyValues: []ast.KeyValue{{
							Key:   "id",
							Value: ast.Value{Object: []ast.ObjectEntry{{Key: "id", Value: ast.ObjectValue{Primitive: &ast.Primitive{String: String("1")}}}}},
							Json:  true,
						}},
					},
				}},
			}}},
			`from hero with id = { "id": "1" } -> json`,
		},
		{
			"Get query with select filters and match function",
			ast.Query{Blocks: []ast.Block{{
				Method:   ast.FromMethod,
				Resource: "hero",
				Qualifiers: []ast.Qualifier{
					{Only: []ast.Filter{
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
			"Simple from resource query with aggreation",
			ast.Query{Blocks: []ast.Block{
				{Method: ast.FromMethod, Resource: "hero"},
				{Method: ast.FromMethod, Resource: "sidekick", In: []string{"hero", "sidekick"}},
			}},
			"from hero from sidekick in hero.sidekick",
		},
		{
			"Full query",
			ast.Query{Blocks: []ast.Block{
				{
					Method:   ast.FromMethod,
					Resource: "hero",
					Alias:    "h",
					Qualifiers: []ast.Qualifier{
						{
							Headers: []ast.HeaderItem{
								{Key: "X-Trace-Id", Value: ast.HeaderValue{String: String("abcdef12345")}},
							},
						},
						{
							Only: []ast.Filter{{Field: []string{"id"}}, {Field: []string{"name"}}},
						},
						{
							With: &ast.Parameters{
								KeyValues: []ast.KeyValue{
									{Key: "id", Value: ast.Value{Primitive: &ast.Primitive{Int: Int(1)}}},
									{Key: "name", Value: ast.Value{Primitive: &ast.Primitive{String: String("batman")}}},
									{
										Key: "weapons",
										Value: ast.Value{List: []*ast.Value{
											{Primitive: &ast.Primitive{String: String("belt")}}, {Primitive: &ast.Primitive{String: String("hands")}},
										}},
									},
									{
										Key: "family",
										Value: ast.Value{Object: []ast.ObjectEntry{
											{Key: "father", Value: ast.ObjectValue{Primitive: &ast.Primitive{String: String("Thomas Wayne")}}},
										}},
									},
									{Key: "height", Value: ast.Value{Primitive: &ast.Primitive{Float: Float(10.5)}}},
									{Key: "var", Value: ast.Value{Variable: String("myvar")}},
								},
							},
						},
					},
				},
				{
					Method:   ast.FromMethod,
					Resource: "sidekick",
					Alias:    "s",
					In:       []string{"hero", "sidekick"},
					Qualifiers: []ast.Qualifier{
						{
							With: &ast.Parameters{
								KeyValues: []ast.KeyValue{
									{Key: "id", Value: ast.Value{Primitive: &ast.Primitive{Int: Int(1)}}},
									{Key: "name", Value: ast.Value{Primitive: &ast.Primitive{String: String("batman")}}},
									{
										Key: "weapons",
										Value: ast.Value{List: []*ast.Value{
											{Primitive: &ast.Primitive{String: String("belt")}}, {Primitive: &ast.Primitive{String: String("hands")}},
										}},
									},
									{
										Key: "family",
										Value: ast.Value{Object: []ast.ObjectEntry{
											{Key: "father", Value: ast.ObjectValue{Primitive: &ast.Primitive{String: String("Thomas Wayne")}}},
										}},
									},
									{Key: "height", Value: ast.Value{Primitive: &ast.Primitive{Float: Float(10.5)}}},
								},
							},
						},
						{
							Headers: []ast.HeaderItem{
								{Key: "X-Trace-Id", Value: ast.HeaderValue{String: String("abcdef12345")}},
							},
						},
						{Hidden: true},
						{IgnoreErrors: true},
					},
				},
				{
					Method:   ast.FromMethod,
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
		
		
				 from sidekick as s in hero.sidekick
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
			ast.Query{Blocks: []ast.Block{
				{
					Method:   ast.FromMethod,
					Resource: "cart",
					Qualifiers: []ast.Qualifier{
						{
							With: &ast.Parameters{
								KeyValues: []ast.KeyValue{
									{Key: "id", Value: ast.Value{Variable: String("id")}},
								},
							},
						},
						{
							Only: []ast.Filter{{Field: []string{"total"}}, {Field: []string{"opn"}}},
						},
					},
				},
				{
					Method:   ast.FromMethod,
					Resource: "sku",
					Qualifiers: []ast.Qualifier{
						{
							With: &ast.Parameters{
								KeyValues: []ast.KeyValue{
									{Key: "sku", Value: ast.Value{Primitive: &ast.Primitive{Chain: []ast.Chained{{PathItem: "cart"}, {PathItem: "lines"}, {PathItem: "productSku"}}}}},
								},
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

	generator, err := ast.New()

	test.VerifyError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generator.Parse(tt.query)

			test.VerifyError(t, err)
			test.Equal(t, *got, tt.expected)
		})
	}
}
