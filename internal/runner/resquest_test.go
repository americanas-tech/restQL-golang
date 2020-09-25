package runner_test

import (
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
	"net/http"
	"testing"

	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
	"github.com/b2wdigital/restQL-golang/v4/internal/runner"
	"github.com/b2wdigital/restQL-golang/v4/test"
)

func TestMakeRequest(t *testing.T) {
	tests := []struct {
		name      string
		statement domain.Statement
		queryCtx  restql.QueryContext
		expected  restql.HTTPRequest
	}{
		{
			"should make get request with url",
			domain.Statement{Method: domain.FromMethod, Resource: "hero"},
			restql.QueryContext{Mappings: map[string]restql.Mapping{"hero": mapping(t, "http://hero.io/api")}},
			restql.HTTPRequest{Method: http.MethodGet, Schema: "http", Host: "hero.io", Path: "/api", Query: map[string]interface{}{}, Headers: map[string]string{"Content-Type": "application/json"}},
		},
		{
			"should make post request with url",
			domain.Statement{Method: domain.ToMethod, Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": 1}}},
			restql.QueryContext{Mappings: map[string]restql.Mapping{"hero": mapping(t, "http://hero.io/api")}},
			restql.HTTPRequest{Method: http.MethodPost, Schema: "http", Host: "hero.io", Path: "/api", Query: map[string]interface{}{}, Body: map[string]interface{}{"id": 1}, Headers: map[string]string{"Content-Type": "application/json"}},
		},
		{
			"should make patch request with url",
			domain.Statement{Method: domain.UpdateMethod, Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": 1}}},
			restql.QueryContext{Mappings: map[string]restql.Mapping{"hero": mapping(t, "http://hero.io/api")}},
			restql.HTTPRequest{Method: http.MethodPatch, Schema: "http", Host: "hero.io", Path: "/api", Query: map[string]interface{}{}, Body: map[string]interface{}{"id": 1}, Headers: map[string]string{"Content-Type": "application/json"}},
		},
		{
			"should make put request with url",
			domain.Statement{Method: domain.IntoMethod, Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": 1}}},
			restql.QueryContext{Mappings: map[string]restql.Mapping{"hero": mapping(t, "http://hero.io/api")}},
			restql.HTTPRequest{Method: http.MethodPut, Schema: "http", Host: "hero.io", Path: "/api", Query: map[string]interface{}{}, Body: map[string]interface{}{"id": 1}, Headers: map[string]string{"Content-Type": "application/json"}},
		},
		{
			"should make delete request with url",
			domain.Statement{Method: domain.DeleteMethod, Resource: "hero"},
			restql.QueryContext{Mappings: map[string]restql.Mapping{"hero": mapping(t, "http://hero.io/api")}},
			restql.HTTPRequest{Method: http.MethodDelete, Schema: "http", Host: "hero.io", Path: "/api", Query: map[string]interface{}{}, Headers: map[string]string{"Content-Type": "application/json"}},
		},
		{
			"should make request with url and query params from statement",
			domain.Statement{Method: domain.FromMethod, Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": "123456"}}},
			restql.QueryContext{Mappings: map[string]restql.Mapping{"hero": mapping(t, "http://hero.io/api")}},
			restql.HTTPRequest{Method: http.MethodGet, Schema: "http", Host: "hero.io", Path: "/api", Query: map[string]interface{}{"id": "123456"}, Headers: map[string]string{"Content-Type": "application/json"}},
		},
		{
			"should make request with url and header from statement",
			domain.Statement{Method: domain.FromMethod, Resource: "hero", Headers: map[string]interface{}{"X-TID": "1234567890"}},
			restql.QueryContext{Mappings: map[string]restql.Mapping{"hero": mapping(t, "http://hero.io/api")}},
			restql.HTTPRequest{Method: http.MethodGet, Schema: "http", Host: "hero.io", Path: "/api", Query: map[string]interface{}{}, Headers: map[string]string{"X-TID": "1234567890", "Content-Type": "application/json"}},
		},
		{
			"should make request with url path params resolved",
			domain.Statement{Method: domain.FromMethod, Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": "123456"}}},
			restql.QueryContext{Mappings: map[string]restql.Mapping{"hero": mapping(t, " http://hero.io/api/:id")}},
			restql.HTTPRequest{Method: http.MethodGet, Schema: "http", Host: "hero.io", Path: "/api/123456", Query: map[string]interface{}{}, Headers: map[string]string{"Content-Type": "application/json"}},
		},
		{
			"should make request with url, query params from statement and forward query params",
			domain.Statement{Method: domain.FromMethod, Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": "123456"}}},
			restql.QueryContext{
				Mappings: map[string]restql.Mapping{"hero": mapping(t, "http://hero.io/api")},
				Input:    restql.QueryInput{Params: map[string]interface{}{"c_universe": "dc", "test": "test"}},
			},
			restql.HTTPRequest{Method: http.MethodGet, Schema: "http", Host: "hero.io", Path: "/api", Query: map[string]interface{}{"id": "123456", "c_universe": "dc"}, Headers: map[string]string{"Content-Type": "application/json"}},
		},
		{
			"should make request with url, header from statement and only allowed forward headers",
			domain.Statement{Method: domain.FromMethod, Resource: "hero", Headers: map[string]interface{}{"X-TID": "1234567890"}},
			restql.QueryContext{
				Mappings: map[string]restql.Mapping{"hero": mapping(t, "http://hero.io/api")},
				Input: restql.QueryInput{Headers: map[string]string{
					"Authorization":   "Bearer abcdefgh",
					"host":            "http://hero.io/api",
					"Content-Type":    "application/json",
					"Content-Length":  "0",
					"Connection":      "keepalive",
					"Origin":          "www.test.com",
					"Accept-Encoding": "gzip, deflate",
				}},
			},
			restql.HTTPRequest{
				Method:  http.MethodGet,
				Schema:  "http",
				Host:    "hero.io",
				Path:    "/api",
				Query:   map[string]interface{}{},
				Headers: map[string]string{"X-TID": "1234567890", "Authorization": "Bearer abcdefgh", "Content-Type": "application/json"},
			},
		},
		{
			"should make post request with parameter as body",
			domain.Statement{Method: domain.ToMethod, Resource: "hero", With: domain.Params{Values: map[string]interface{}{"id": domain.AsBody{Value: []interface{}{"1", "2", "3"}}}}},
			restql.QueryContext{Mappings: map[string]restql.Mapping{"hero": mapping(t, "http://hero.io/api")}},
			restql.HTTPRequest{Method: http.MethodPost, Schema: "http", Host: "hero.io", Path: "/api", Query: map[string]interface{}{}, Body: []interface{}{"1", "2", "3"}, Headers: map[string]string{"Content-Type": "application/json"}},
		},
	}

	forwardPrefix := "c_"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runner.MakeRequest(0, forwardPrefix, tt.statement, tt.queryCtx)

			test.Equal(t, got, tt.expected)
		})
	}
}

func mapping(t *testing.T, url string) restql.Mapping {
	m, err := restql.NewMapping("test-resource", url)
	if err != nil {
		t.Fatal("failed to create stub mapping", err)
	}

	return m
}
