package runner_test

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/runner"
	"reflect"
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
			"should make request with only url",
			domain.Statement{Resource: "hero"},
			domain.QueryContext{Mappings: map[string]domain.Mapping{"hero": {Schema: "http", Uri: "hero.io/api"}}},
			domain.HttpRequest{Schema: "http", Uri: "hero.io/api", Query: map[string]interface{}{}, Headers: map[string]string{}},
		},
		{
			"should make request with url and query params from statement",
			domain.Statement{Resource: "hero", With: map[string]interface{}{"id": "123456"}},
			domain.QueryContext{Mappings: map[string]domain.Mapping{"hero": {Schema: "http", Uri: "hero.io/api"}}},
			domain.HttpRequest{Schema: "http", Uri: "hero.io/api", Query: map[string]interface{}{"id": "123456"}, Headers: map[string]string{}},
		},
		{
			"should make request with url and header from statement",
			domain.Statement{Resource: "hero", Headers: map[string]interface{}{"X-TID": "1234567890"}},
			domain.QueryContext{Mappings: map[string]domain.Mapping{"hero": {Schema: "http", Uri: "hero.io/api"}}},
			domain.HttpRequest{Schema: "http", Uri: "hero.io/api", Query: map[string]interface{}{}, Headers: map[string]string{"X-TID": "1234567890"}},
		},
		{
			"should make request with url path params resolved",
			domain.Statement{Resource: "hero", With: map[string]interface{}{"id": "123456"}},
			domain.QueryContext{Mappings: map[string]domain.Mapping{"hero": {
				Schema:        "http",
				Uri:           "hero.io/api/:id",
				PathParams:    []string{"id"},
				PathParamsSet: map[string]struct{}{"id": {}},
			}}},
			domain.HttpRequest{Schema: "http", Uri: "hero.io/api/123456", Query: map[string]interface{}{}, Headers: map[string]string{}},
		},
		{
			"should make request with url, query params from statement and forward query params",
			domain.Statement{Resource: "hero", With: map[string]interface{}{"id": "123456"}},
			domain.QueryContext{
				Mappings: map[string]domain.Mapping{"hero": {Schema: "http", Uri: "hero.io/api"}},
				Input:    domain.QueryInput{Params: map[string]interface{}{"c_universe": "dc", "test": "test"}},
			},
			domain.HttpRequest{Schema: "http", Uri: "hero.io/api", Query: map[string]interface{}{"id": "123456", "c_universe": "dc"}, Headers: map[string]string{}},
		},
		{
			"should make request with url, header from statement and only allowed forward headers",
			domain.Statement{Resource: "hero", Headers: map[string]interface{}{"X-TID": "1234567890"}},
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
				Schema:  "http",
				Uri:     "hero.io/api",
				Query:   map[string]interface{}{},
				Headers: map[string]string{"X-TID": "1234567890", "Authorization": "Bearer abcdefgh"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runner.MakeRequest(tt.statement, tt.queryCtx)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("MakeRequest = %#+v, want = %#+v", got, tt.expected)
			}
		})
	}
}
