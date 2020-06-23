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

func Boolean(b bool) *bool {
	return &b
}

func TestAstGenerator(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected ast.Query
	}{
		{
			"Simple from resource query",
			"from cart",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "cart"}}},
		},
		{
			"Simple from resource query with comment",
			`// a comment
						  from cart // some other comment`,
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "cart"}}},
		},
		{
			"Simple from resource query with use modifier",
			`
							use max-age 600
							use s-max-age 400
							use timeout 8000
							
							from cart
					`,
			ast.Query{
				Use: []ast.Use{
					{Key: ast.MaxAgeKeyword, Value: ast.UseValue{Int: Int(600)}},
					{Key: ast.SmaxAgeKeyword, Value: ast.UseValue{Int: Int(400)}},
					{Key: ast.TimeoutKeyword, Value: ast.UseValue{Int: Int(8000)}},
				},
				Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "cart"}},
			},
		},
		{
			"Query with two from statements",
			`
							from cart
							from hero
					`,
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "cart"}, {Method: ast.FromMethod, Resource: "hero"}}},
		},
		{
			"Simple from resource query with resource name with hyphen",
			"from shopping-cart",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "shopping-cart"}}},
		},
		{
			"Simple from resource query with alias",
			"from cart as shopping",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "cart", Alias: "shopping"}}},
		},
		{
			"Get query with integer query parameters",
			`from hero with id = 1`,
			ast.Query{
				Blocks: []ast.Block{
					{
						Method:   ast.FromMethod,
						Resource: "hero",
						Qualifiers: []ast.Qualifier{
							{
								With: &ast.Parameters{
									KeyValues: []ast.KeyValue{{Key: "id", Value: ast.Value{Primitive: &ast.Primitive{Int: Int(1)}}}},
								},
							},
						},
					},
				},
			},
		},
		{
			"Get query with string query parameters",
			`from cart with name = "batman"`,
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "cart", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{KeyValues: []ast.KeyValue{{Key: "name", Value: ast.Value{Primitive: &ast.Primitive{String: String("batman")}}}}}}}}}},
		},
		{
			"Get query with float query parameters",
			`from hero with height = 10.5`,
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{KeyValues: []ast.KeyValue{{Key: "height", Value: ast.Value{Primitive: &ast.Primitive{Float: Float(10.5)}}}}}}}}}},
		},
		{
			"Get query with chained query parameters",
			`from hero with id = done-resource.id`,
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{
				KeyValues: []ast.KeyValue{{Key: "id", Value: ast.Value{Primitive: &ast.Primitive{Chain: []ast.Chained{{PathItem: "done-resource"}, {PathItem: "id"}}}}}},
			}}}}}},
		},
		{
			"Get query with variable chained query parameters",
			`from hero with id = done-resource.$path.id`,
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{
				KeyValues: []ast.KeyValue{{Key: "id", Value: ast.Value{Primitive: &ast.Primitive{Chain: []ast.Chained{{PathItem: "done-resource"}, {PathVariable: "path"}, {PathItem: "id"}}}}}},
			}}}}}},
		},
		{
			"Get query with variable query parameters",
			`from hero with id = $id`,
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{
				KeyValues: []ast.KeyValue{{Key: "id", Value: ast.Value{Variable: String("id")}}},
			}}}}}},
		},
		{
			"Get query with null query parameters",
			`from hero with id = null`,
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{
				KeyValues: []ast.KeyValue{{Key: "id", Value: ast.Value{Primitive: &ast.Primitive{Null: true}}}},
			}}}}}},
		},
		{
			"Get query with boolean query parameters",
			`
							from hero with marvel = true
							from sidekick with marvel = false
					`,
			ast.Query{Blocks: []ast.Block{
				{
					Method:   ast.FromMethod,
					Resource: "hero",
					Qualifiers: []ast.Qualifier{{
						With: &ast.Parameters{
							KeyValues: []ast.KeyValue{{Key: "marvel", Value: ast.Value{Primitive: &ast.Primitive{Boolean: Boolean(true)}}}},
						},
					}},
				},
				{
					Method:   ast.FromMethod,
					Resource: "sidekick",
					Qualifiers: []ast.Qualifier{{
						With: &ast.Parameters{
							KeyValues: []ast.KeyValue{{Key: "marvel", Value: ast.Value{Primitive: &ast.Primitive{Boolean: Boolean(false)}}}},
						},
					}},
				},
			}},
		},
		{
			"Get query with list query parameters",
			`from hero with weapons = ["sword", "shield"]`,
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{
				KeyValues: []ast.KeyValue{{Key: "weapons", Value: ast.Value{List: []ast.Value{{Primitive: &ast.Primitive{String: String("sword")}}, {Primitive: &ast.Primitive{String: String("shield")}}}}}},
			}}}}}},
		},
		{
			"Get query with nested list query parameters",
			`from hero with weapons = [["sword"], ["shield"]]`,
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{
				KeyValues: []ast.KeyValue{{Key: "weapons", Value: ast.Value{List: []ast.Value{{List: []ast.Value{{Primitive: &ast.Primitive{String: String("sword")}}}}, {List: []ast.Value{{Primitive: &ast.Primitive{String: String("shield")}}}}}}}},
			}}}}}},
		},
		{
			"Get query with list query parameters delimited by new line",
			fmt.Sprintf("from hero with weapons = [\"sword\"\n\"shield\"]"),
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{
				KeyValues: []ast.KeyValue{{Key: "weapons", Value: ast.Value{List: []ast.Value{{Primitive: &ast.Primitive{String: String("sword")}}, {Primitive: &ast.Primitive{String: String("shield")}}}}}},
			}}}}}},
		},
		{
			"Get query with multiple parameters delimited by comma",
			`from hero with id = 1, name = "batman"`,
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{KeyValues: []ast.KeyValue{
				{Key: "id", Value: ast.Value{Primitive: &ast.Primitive{Int: Int(1)}}},
				{Key: "name", Value: ast.Value{Primitive: &ast.Primitive{String: String("batman")}}},
			}}}}}}},
		},
		{
			"Get query with object query parameters",
			`from hero with id = { "id": "1" }`,
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{
				KeyValues: []ast.KeyValue{{Key: "id", Value: ast.Value{Object: []ast.ObjectEntry{{Key: "id", Value: ast.Value{Primitive: &ast.Primitive{String: String("1")}}}}}}},
			}}}}}},
		},
		{
			"Get query with object query parameters",
			`from hero with id = { id: "1" }`,
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{
				KeyValues: []ast.KeyValue{{Key: "id", Value: ast.Value{Object: []ast.ObjectEntry{{Key: "id", Value: ast.Value{Primitive: &ast.Primitive{String: String("1")}}}}}}},
			}}}}}},
		},
		{
			"Get query with object query parameters using list value",
			`from hero with id = { names: ["batman", "wonder woman"] }`,
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{
				KeyValues: []ast.KeyValue{{Key: "id", Value: ast.Value{Object: []ast.ObjectEntry{{Key: "names", Value: ast.Value{List: []ast.Value{{Primitive: &ast.Primitive{String: String("batman")}}, {Primitive: &ast.Primitive{String: String("wonder woman")}}}}}}}}},
			}}}}}},
		},
		{
			"Get query with object query parameters using a variable value",
			`from hero with id = { id: $id }`,
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{
				KeyValues: []ast.KeyValue{{Key: "id", Value: ast.Value{Object: []ast.ObjectEntry{{Key: "id", Value: ast.Value{Variable: String("id")}}}}}},
			}}}}}},
		},
		{
			"Get query with object query parameters with multiple key/values",
			`from hero with id = { "id": "1", "name": "batman" }`,
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{
				KeyValues: []ast.KeyValue{{Key: "id", Value: ast.Value{Object: []ast.ObjectEntry{{Key: "id", Value: ast.Value{Primitive: &ast.Primitive{String: String("1")}}}, {Key: "name", Value: ast.Value{Primitive: &ast.Primitive{String: String("batman")}}}}}}},
			}}}}}},
		},
		{
			"Get query with object query parameters with multiple key/values delimited by new line",
			fmt.Sprintf("from hero with id = { \"id\": \"1\",\n\"name\": \"batman\" }"),
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{
				KeyValues: []ast.KeyValue{{Key: "id", Value: ast.Value{Object: []ast.ObjectEntry{{Key: "id", Value: ast.Value{Primitive: &ast.Primitive{String: String("1")}}}, {Key: "name", Value: ast.Value{Primitive: &ast.Primitive{String: String("batman")}}}}}}},
			}}}}}},
		},
		{
			"Get query with nested object query parameters",
			`from hero with id = { "id": { "internal": "1", "version": 10 } }`,
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{
				KeyValues: []ast.KeyValue{{Key: "id", Value: ast.Value{Object: []ast.ObjectEntry{{Key: "id", Value: ast.Value{Object: []ast.ObjectEntry{{Key: "internal", Value: ast.Value{Primitive: &ast.Primitive{String: String("1")}}}, {Key: "version", Value: ast.Value{Primitive: &ast.Primitive{Int: Int(10)}}}}}}}}}},
			}}}}}},
		},
		{
			"Get query with parameter using keyword as key",
			`from hero with from = "5m"`,
			ast.Query{
				Blocks: []ast.Block{
					{
						Method:   ast.FromMethod,
						Resource: "hero",
						Qualifiers: []ast.Qualifier{
							{
								With: &ast.Parameters{
									KeyValues: []ast.KeyValue{
										{Key: "from", Value: ast.Value{Primitive: &ast.Primitive{String: String("5m")}}},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			"Get query with string query parameters and key using dot",
			`from cart with hero.name = "batman"`,
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "cart", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{KeyValues: []ast.KeyValue{
				{Key: "hero.name", Value: ast.Value{Primitive: &ast.Primitive{String: String("batman")}}},
			}}}}}}},
		},
		{
			"Get query with multiple parameters keys using dot",
			`from hero with hero.id = 1, hero.height = 10.5`,
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{KeyValues: []ast.KeyValue{
				{Key: "hero.id", Value: ast.Value{Primitive: &ast.Primitive{Int: Int(1)}}},
				{Key: "hero.height", Value: ast.Value{Primitive: &ast.Primitive{Float: Float(10.5)}}},
			}}}}}}},
		},
		{
			"Get query with multiple parameters delimited by space and using keys with dots",
			`from hero with hero.id = 1, hero.name = "batman"`,
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{KeyValues: []ast.KeyValue{
				{Key: "hero.id", Value: ast.Value{Primitive: &ast.Primitive{Int: Int(1)}}},
				{Key: "hero.name", Value: ast.Value{Primitive: &ast.Primitive{String: String("batman")}}},
			}}}}}}},
		},
		{
			"Get query with multiple parameters delimited by new line and using keys with dots",
			fmt.Sprintf("from hero with hero.id = 1\nhero.name = \"batman\""),
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{With: &ast.Parameters{KeyValues: []ast.KeyValue{
				{Key: "hero.id", Value: ast.Value{Primitive: &ast.Primitive{Int: Int(1)}}},
				{Key: "hero.name", Value: ast.Value{Primitive: &ast.Primitive{String: String("batman")}}},
			}}}}}}},
		},
		{
			"Get query with list query parameters flattened",
			`from hero with weapons = ["sword", "shield"] -> flatten`,
			ast.Query{Blocks: []ast.Block{{
				Method:   ast.FromMethod,
				Resource: "hero",
				Qualifiers: []ast.Qualifier{
					{
						With: &ast.Parameters{
							KeyValues: []ast.KeyValue{
								{
									Key: "weapons",
									Value: ast.Value{List: []ast.Value{
										{Primitive: &ast.Primitive{String: String("sword")}},
										{Primitive: &ast.Primitive{String: String("shield")}},
									}},
									Functions: []string{"flatten"},
								},
							},
						},
					},
				},
			}}},
		},
		{
			"Get query with chained query parameters flattened",
			`from hero with id = done-resource.id -> flatten`,
			ast.Query{Blocks: []ast.Block{{
				Method:   ast.FromMethod,
				Resource: "hero",
				Qualifiers: []ast.Qualifier{{
					With: &ast.Parameters{
						KeyValues: []ast.KeyValue{
							{
								Key:       "id",
								Value:     ast.Value{Primitive: &ast.Primitive{Chain: []ast.Chained{{PathItem: "done-resource"}, {PathItem: "id"}}}},
								Functions: []string{"flatten"},
							},
						},
					},
				}},
			}}},
		},
		{
			"Get query with query parameters encoded in base64",
			`from hero with id = "abcdefg12345" -> base64`,
			ast.Query{Blocks: []ast.Block{{
				Method:   ast.FromMethod,
				Resource: "hero",
				Qualifiers: []ast.Qualifier{{
					With: &ast.Parameters{
						KeyValues: []ast.KeyValue{
							{
								Key:       "id",
								Value:     ast.Value{Primitive: &ast.Primitive{String: String("abcdefg12345")}},
								Functions: []string{"base64"},
							},
						},
					},
				}},
			}}},
		},
		{
			"Get query with object query parameters encoded as json",
			`from hero with id = { "id": "1" } -> json`,
			ast.Query{Blocks: []ast.Block{{
				Method:   ast.FromMethod,
				Resource: "hero",
				Qualifiers: []ast.Qualifier{{
					With: &ast.Parameters{
						KeyValues: []ast.KeyValue{{
							Key:       "id",
							Value:     ast.Value{Object: []ast.ObjectEntry{{Key: "id", Value: ast.Value{Primitive: &ast.Primitive{String: String("1")}}}}},
							Functions: []string{"json"},
						}},
					},
				}},
			}}},
		},
		{
			"Get query with parameter using multiple functions",
			`from hero with id = [1,2,3] -> flatten -> json`,
			ast.Query{Blocks: []ast.Block{{
				Method:   ast.FromMethod,
				Resource: "hero",
				Qualifiers: []ast.Qualifier{{
					With: &ast.Parameters{
						KeyValues: []ast.KeyValue{{
							Key:       "id",
							Value:     ast.Value{List: []ast.Value{{Primitive: &ast.Primitive{Int: Int(1)}}, {Primitive: &ast.Primitive{Int: Int(2)}}, {Primitive: &ast.Primitive{Int: Int(3)}}}},
							Functions: []string{"flatten", "json"},
						}},
					},
				}},
			}}},
		},
		{
			"Get query with dynamic body parameter and multiple statements",
			`from hero
							with
								$update
								id = 1
								from = "5m"
								timeout = 100
		
							from sidekick
								with
									id = 1
									from = "5m"
									timeout = 100`,
			ast.Query{
				Blocks: []ast.Block{
					{
						Method:   ast.FromMethod,
						Resource: "hero",
						Qualifiers: []ast.Qualifier{
							{
								With: &ast.Parameters{
									Body: &ast.ParameterBody{Target: "update"},
									KeyValues: []ast.KeyValue{
										{Key: "id", Value: ast.Value{Primitive: &ast.Primitive{Int: Int(1)}}},
										{Key: "from", Value: ast.Value{Primitive: &ast.Primitive{String: String("5m")}}},
										{Key: "timeout", Value: ast.Value{Primitive: &ast.Primitive{Int: Int(100)}}},
									},
								},
							},
						},
					},
					{
						Method:   ast.FromMethod,
						Resource: "sidekick",
						Qualifiers: []ast.Qualifier{
							{
								With: &ast.Parameters{
									KeyValues: []ast.KeyValue{
										{Key: "id", Value: ast.Value{Primitive: &ast.Primitive{Int: Int(1)}}},
										{Key: "from", Value: ast.Value{Primitive: &ast.Primitive{String: String("5m")}}},
										{Key: "timeout", Value: ast.Value{Primitive: &ast.Primitive{Int: Int(100)}}},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			"Get query with dynamic body parameter flattened",
			`from hero with $body -> flatten`,
			ast.Query{
				Blocks: []ast.Block{
					{
						Method:   ast.FromMethod,
						Resource: "hero",
						Qualifiers: []ast.Qualifier{
							{
								With: &ast.Parameters{
									Body: &ast.ParameterBody{
										Target:    "body",
										Functions: []string{"flatten"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			"Get query with select filter",
			`from hero only *, name, 
						to,
         		      
						weapons`,
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{Only: []ast.Filter{{Field: []string{"*"}}, {Field: []string{"name"}}, {Field: []string{"to"}}, {Field: []string{"weapons"}}}}}}}},
		},
		{
			"Get query with select filter delimited by new line",
			fmt.Sprintf("from hero only name\nweapons"),
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{Only: []ast.Filter{{Field: []string{"name"}}, {Field: []string{"weapons"}}}}}}}},
		},
		{
			"Multiple gets with select filter",
			`from hero
		only
			name
			weapons
		
		from sidekick`,
			ast.Query{Blocks: []ast.Block{
				{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{Only: []ast.Filter{{Field: []string{"name"}}, {Field: []string{"weapons"}}}}}},
				{Method: ast.FromMethod, Resource: "sidekick"},
			}},
		},
		{
			"Get query with select filters and match function",
			`from hero
								only
										name -> matches("^Super")
										weapons`,
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
		},
		{
			"Get query with hidden",
			"from hero hidden",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{Hidden: true}}}}},
		},
		{
			"Get query with ignore errors",
			"from hero ignore-errors",
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{IgnoreErrors: true}}}}},
		},
		{
			"Get query with integer timeout",
			`from hero timeout 200`,
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{Timeout: &ast.TimeoutValue{Int: Int(200)}}}}}},
		},
		{
			"Get query with variable timeout",
			`from hero timeout $some-time`,
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{Timeout: &ast.TimeoutValue{Variable: String("some-time")}}}}}},
		},
		{
			"Get query with headers",
			`from hero headers Authorization = "abcdef12345", X-Trace-Id = $trace-id`,
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{Headers: []ast.HeaderItem{
				{Key: "Authorization", Value: ast.HeaderValue{String: String("abcdef12345")}},
				{Key: "X-Trace-Id", Value: ast.HeaderValue{Variable: String("trace-id")}},
			}}}}}},
		},
		{
			"Get query with max age",
			`from hero max-age 2000`,
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{MaxAge: &ast.MaxAgeValue{Int: Int(2000)}}}}}},
		},
		{
			"Get query with cache control",
			`from hero s-max-age 2000`,
			ast.Query{Blocks: []ast.Block{{Method: ast.FromMethod, Resource: "hero", Qualifiers: []ast.Qualifier{{SMaxAge: &ast.SMaxAgeValue{Int: Int(2000)}}}}}},
		},
		{
			"Simple from resource query with aggregation",
			`
							from hero
							from sidekick in hero.sidekick
						`,
			ast.Query{Blocks: []ast.Block{
				{Method: ast.FromMethod, Resource: "hero"},
				{Method: ast.FromMethod, Resource: "sidekick", In: []string{"hero", "sidekick"}},
			}},
		},
		{
			"Full query",
			`from hero as h
							timeout 200
							headers
								X-Trace-Id = "abcdef12345"
							with
								id = 1,   
								name = "batman",    
								weapons = ["belt", "hands"],   		
								family = { "father": "Thomas Wayne" },		
								height = 10.5,		
								var = $myvar
							only 
								from
								name
						 from sidekick as s in hero.sidekick
							headers
								X-Trace-Id = "abcdef12345"
							timeout $timeout
							with
								id = 1
								hero.name = "batman"
								weapons = ["belt", "hands"]
								family = { "father": "Thomas Wayne" }
								height = 10.5
								chain = hero.with.foo
							hidden
							ignore-errors

						from villain as v
							only
								ignore-errors
								name
							ignore-errors`,
			ast.Query{Blocks: []ast.Block{
				{
					Method:   ast.FromMethod,
					Resource: "hero",
					Alias:    "h",
					Qualifiers: []ast.Qualifier{
						{
							Timeout: &ast.TimeoutValue{Int: Int(200)},
						},
						{
							Headers: []ast.HeaderItem{
								{Key: "X-Trace-Id", Value: ast.HeaderValue{String: String("abcdef12345")}},
							},
						},
						{
							With: &ast.Parameters{
								KeyValues: []ast.KeyValue{
									{Key: "id", Value: ast.Value{Primitive: &ast.Primitive{Int: Int(1)}}},
									{Key: "name", Value: ast.Value{Primitive: &ast.Primitive{String: String("batman")}}},
									{
										Key: "weapons",
										Value: ast.Value{List: []ast.Value{
											{Primitive: &ast.Primitive{String: String("belt")}}, {Primitive: &ast.Primitive{String: String("hands")}},
										}},
									},
									{
										Key: "family",
										Value: ast.Value{Object: []ast.ObjectEntry{
											{Key: "father", Value: ast.Value{Primitive: &ast.Primitive{String: String("Thomas Wayne")}}},
										}},
									},
									{Key: "height", Value: ast.Value{Primitive: &ast.Primitive{Float: Float(10.5)}}},
									{Key: "var", Value: ast.Value{Variable: String("myvar")}},
								},
							},
						},
						{
							Only: []ast.Filter{{Field: []string{"from"}}, {Field: []string{"name"}}},
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
							Headers: []ast.HeaderItem{
								{Key: "X-Trace-Id", Value: ast.HeaderValue{String: String("abcdef12345")}},
							},
						},
						{
							Timeout: &ast.TimeoutValue{Variable: String("timeout")},
						},
						{
							With: &ast.Parameters{
								KeyValues: []ast.KeyValue{
									{Key: "id", Value: ast.Value{Primitive: &ast.Primitive{Int: Int(1)}}},
									{Key: "hero.name", Value: ast.Value{Primitive: &ast.Primitive{String: String("batman")}}},
									{
										Key: "weapons",
										Value: ast.Value{List: []ast.Value{
											{Primitive: &ast.Primitive{String: String("belt")}}, {Primitive: &ast.Primitive{String: String("hands")}},
										}},
									},
									{
										Key: "family",
										Value: ast.Value{Object: []ast.ObjectEntry{
											{Key: "father", Value: ast.Value{Primitive: &ast.Primitive{String: String("Thomas Wayne")}}},
										}},
									},
									{Key: "height", Value: ast.Value{Primitive: &ast.Primitive{Float: Float(10.5)}}},
									{Key: "chain", Value: ast.Value{Primitive: &ast.Primitive{Chain: []ast.Chained{{PathItem: "hero"}, {PathItem: "with"}, {PathItem: "foo"}}}}},
								},
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
					Qualifiers: []ast.Qualifier{
						{
							Only: []ast.Filter{
								{Field: []string{"ignore-errors"}},
								{Field: []string{"name"}},
							},
						},
						{IgnoreErrors: true},
					},
				},
			}},
		},
		{
			"Bug: words starting with `to` failed parser because they would match a method keyword",
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
