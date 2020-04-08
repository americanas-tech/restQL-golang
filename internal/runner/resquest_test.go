package runner_test

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/runner"
	"github.com/b2wdigital/restQL-golang/test"
	"net/http"
	"testing"
)

func TestMakeRequest(t *testing.T) {
	tests := []struct {
		name      string
		statement domain.Statement
		queryCtx  domain.QueryContext
		expected  domain.HttpRequest
	}{
		{
			"should make get request with url",
			domain.Statement{Method: domain.FromMethod, Resource: "hero"},
			domain.QueryContext{Mappings: map[string]domain.Mapping{"hero": {Schema: "http", Uri: "hero.io/api"}}},
			domain.HttpRequest{Method: http.MethodGet, Schema: "http", Uri: "hero.io/api", Query: map[string]interface{}{}, Headers: map[string]string{"Content-Type": "application/json"}},
		},
		{
			"should make post request with url",
			domain.Statement{Method: domain.ToMethod, Resource: "hero", With: map[string]interface{}{"id": 1}},
			domain.QueryContext{Mappings: map[string]domain.Mapping{"hero": {Schema: "http", Uri: "hero.io/api"}}},
			domain.HttpRequest{Method: http.MethodPost, Schema: "http", Uri: "hero.io/api", Query: map[string]interface{}{}, Body: map[string]interface{}{"id": 1}, Headers: map[string]string{"Content-Type": "application/json"}},
		},
		{
			"should make put request with url",
			domain.Statement{Method: domain.UpdateMethod, Resource: "hero", With: map[string]interface{}{"id": 1}},
			domain.QueryContext{Mappings: map[string]domain.Mapping{"hero": {Schema: "http", Uri: "hero.io/api"}}},
			domain.HttpRequest{Method: http.MethodPut, Schema: "http", Uri: "hero.io/api", Query: map[string]interface{}{}, Body: map[string]interface{}{"id": 1}, Headers: map[string]string{"Content-Type": "application/json"}},
		},
		{
			"should make delete request with url",
			domain.Statement{Method: domain.DeleteMethod, Resource: "hero"},
			domain.QueryContext{Mappings: map[string]domain.Mapping{"hero": {Schema: "http", Uri: "hero.io/api"}}},
			domain.HttpRequest{Method: http.MethodDelete, Schema: "http", Uri: "hero.io/api", Query: map[string]interface{}{}, Headers: map[string]string{"Content-Type": "application/json"}},
		},
		{
			"should make request with url and query params from statement",
			domain.Statement{Method: domain.FromMethod, Resource: "hero", With: map[string]interface{}{"id": "123456"}},
			domain.QueryContext{Mappings: map[string]domain.Mapping{"hero": {Schema: "http", Uri: "hero.io/api"}}},
			domain.HttpRequest{Method: http.MethodGet, Schema: "http", Uri: "hero.io/api", Query: map[string]interface{}{"id": "123456"}, Headers: map[string]string{"Content-Type": "application/json"}},
		},
		{
			"should make request with url and header from statement",
			domain.Statement{Method: domain.FromMethod, Resource: "hero", Headers: map[string]interface{}{"X-TID": "1234567890"}},
			domain.QueryContext{Mappings: map[string]domain.Mapping{"hero": {Schema: "http", Uri: "hero.io/api"}}},
			domain.HttpRequest{Method: http.MethodGet, Schema: "http", Uri: "hero.io/api", Query: map[string]interface{}{}, Headers: map[string]string{"X-TID": "1234567890", "Content-Type": "application/json"}},
		},
		{
			"should make request with url path params resolved",
			domain.Statement{Method: domain.FromMethod, Resource: "hero", With: map[string]interface{}{"id": "123456"}},
			domain.QueryContext{Mappings: map[string]domain.Mapping{"hero": {
				Schema:        "http",
				Uri:           "hero.io/api/:id",
				PathParams:    []string{"id"},
				PathParamsSet: map[string]struct{}{"id": {}},
			}}},
			domain.HttpRequest{Method: http.MethodGet, Schema: "http", Uri: "hero.io/api/123456", Query: map[string]interface{}{}, Headers: map[string]string{"Content-Type": "application/json"}},
		},
		{
			"should make request with url, query params from statement and forward query params",
			domain.Statement{Method: domain.FromMethod, Resource: "hero", With: map[string]interface{}{"id": "123456"}},
			domain.QueryContext{
				Mappings: map[string]domain.Mapping{"hero": {Schema: "http", Uri: "hero.io/api"}},
				Input:    domain.QueryInput{Params: map[string]interface{}{"c_universe": "dc", "test": "test"}},
			},
			domain.HttpRequest{Method: http.MethodGet, Schema: "http", Uri: "hero.io/api", Query: map[string]interface{}{"id": "123456", "c_universe": "dc"}, Headers: map[string]string{"Content-Type": "application/json"}},
		},
		{
			"should make request with url, header from statement and only allowed forward headers",
			domain.Statement{Method: domain.FromMethod, Resource: "hero", Headers: map[string]interface{}{"X-TID": "1234567890"}},
			domain.QueryContext{
				Mappings: map[string]domain.Mapping{"hero": {Schema: "http", Uri: "hero.io/api"}},
				Input: domain.QueryInput{Headers: map[string]string{
					"Authorization":   "Bearer abcdefgh",
					"Host":            "http://hero.io/api",
					"Content-Type":    "application/json",
					"Content-Length":  "0",
					"Connection":      "keepalive",
					"Origin":          "www.test.com",
					"Accept-Encoding": "gzip, deflate",
				}},
			},
			domain.HttpRequest{
				Method:  http.MethodGet,
				Schema:  "http",
				Uri:     "hero.io/api",
				Query:   map[string]interface{}{},
				Headers: map[string]string{"X-TID": "1234567890", "Authorization": "Bearer abcdefgh", "Content-Type": "application/json"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runner.MakeRequest(tt.statement, tt.queryCtx)

			test.Equal(t, got, tt.expected)
		})
	}
}
