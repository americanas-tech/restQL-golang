package runner_test

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/runner"
	"github.com/b2wdigital/restQL-golang/test"
	"testing"
	"time"
)

func TestNewDoneResource(t *testing.T) {
	tests := []struct {
		name     string
		request  domain.HttpRequest
		response domain.HttpResponse
		options  runner.DoneResourceOptions
		expected domain.DoneResource
	}{
		{
			"should create done resource for successful execution",
			domain.HttpRequest{},
			domain.HttpResponse{StatusCode: 200, Body: nil},
			runner.DoneResourceOptions{},
			domain.DoneResource{Status: 200, Success: true, IgnoreErrors: false, ResponseBody: nil},
		},
		{
			"should create done resource for failed execution",
			domain.HttpRequest{},
			domain.HttpResponse{StatusCode: 400, Body: nil},
			runner.DoneResourceOptions{},
			domain.DoneResource{Status: 400, Success: false, IgnoreErrors: false, ResponseBody: nil},
		},
		{
			"should create done resource with debug",
			domain.HttpRequest{
				Schema:  "http",
				Host:    "hero.io",
				Path:    "/api",
				Query:   map[string]interface{}{"id": "123456"},
				Headers: map[string]string{"X-TID": "12345abdef"},
			},
			domain.HttpResponse{
				Url:        "http://hero.io/api",
				StatusCode: 200,
				Body:       nil,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Duration:   100 * time.Millisecond,
			},
			runner.DoneResourceOptions{},
			domain.DoneResource{
				Status:          200,
				Success:         true,
				IgnoreErrors:    false,
				Url:             "http://hero.io/api",
				RequestHeaders:  map[string]string{"X-TID": "12345abdef"},
				ResponseHeaders: map[string]string{"Content-Type": "application/json"},
				RequestParams:   map[string]interface{}{"id": "123456"},
				ResponseTime:    100,
				ResponseBody:    nil,
			},
		},
		{
			"should create done resource with ignore errors",
			domain.HttpRequest{},
			domain.HttpResponse{StatusCode: 200, Body: nil},
			runner.DoneResourceOptions{IgnoreErrors: true},
			domain.DoneResource{
				Status:       200,
				Success:      true,
				IgnoreErrors: true,
				ResponseBody: nil,
			},
		},
		{
			"should create done resource with cache control information returned by resource",
			domain.HttpRequest{},
			domain.HttpResponse{StatusCode: 200, Body: nil, Headers: map[string]string{"Cache-Control": "max-age=400, s-maxage=600"}},
			runner.DoneResourceOptions{},
			domain.DoneResource{
				Status:  200,
				Success: true,
				CacheControl: domain.ResourceCacheControl{
					MaxAge:  domain.ResourceCacheControlValue{Exist: true, Time: 400},
					SMaxAge: domain.ResourceCacheControlValue{Exist: true, Time: 600},
				},
				ResponseHeaders: map[string]string{"Cache-Control": "max-age=400, s-maxage=600"},
				IgnoreErrors:    false,
				ResponseBody:    nil,
			},
		},
		{
			"should create done resource with cache control information returned by resource",
			domain.HttpRequest{},
			domain.HttpResponse{StatusCode: 200, Body: nil, Headers: map[string]string{"Cache-Control": "no-cache"}},
			runner.DoneResourceOptions{},
			domain.DoneResource{
				Status:  200,
				Success: true,
				CacheControl: domain.ResourceCacheControl{
					NoCache: true,
				},
				ResponseHeaders: map[string]string{"Cache-Control": "no-cache"},
				IgnoreErrors:    false,
				ResponseBody:    nil,
			},
		},
		{
			"should create done resource with cache control information defined in statement if not returned by resource",
			domain.HttpRequest{},
			domain.HttpResponse{StatusCode: 200, Body: nil},
			runner.DoneResourceOptions{MaxAge: 100, SMaxAge: 300},
			domain.DoneResource{
				Status:  200,
				Success: true,
				CacheControl: domain.ResourceCacheControl{
					MaxAge:  domain.ResourceCacheControlValue{Exist: true, Time: 100},
					SMaxAge: domain.ResourceCacheControlValue{Exist: true, Time: 300},
				},
				IgnoreErrors: false,
				ResponseBody: nil,
			},
		},
		{
			"should create done resource with cache control information defined in statement if not returned by resource",
			domain.HttpRequest{},
			domain.HttpResponse{StatusCode: 200, Body: nil},
			runner.DoneResourceOptions{MaxAge: 100},
			domain.DoneResource{
				Status:  200,
				Success: true,
				CacheControl: domain.ResourceCacheControl{
					MaxAge: domain.ResourceCacheControlValue{Exist: true, Time: 100},
				},
				IgnoreErrors: false,
				ResponseBody: nil,
			},
		},
		{
			"should create done resource with minimum cache control information between the returned by resource and the defined in statement",
			domain.HttpRequest{},
			domain.HttpResponse{StatusCode: 200, Body: nil, Headers: map[string]string{"Cache-Control": "max-age=100, s-maxage=600"}},
			runner.DoneResourceOptions{MaxAge: 400, SMaxAge: 300},
			domain.DoneResource{
				Status:  200,
				Success: true,
				CacheControl: domain.ResourceCacheControl{
					MaxAge:  domain.ResourceCacheControlValue{Exist: true, Time: 100},
					SMaxAge: domain.ResourceCacheControlValue{Exist: true, Time: 300},
				},
				ResponseHeaders: map[string]string{"Cache-Control": "max-age=100, s-maxage=600"},
				IgnoreErrors:    false,
				ResponseBody:    nil,
			},
		},
		{
			"should create done resource with minimum cache control information between the returned by resource and the defined in statement",
			domain.HttpRequest{},
			domain.HttpResponse{StatusCode: 200, Body: nil, Headers: map[string]string{"Cache-Control": "no-cache"}},
			runner.DoneResourceOptions{MaxAge: 400, SMaxAge: 300},
			domain.DoneResource{
				Status:  200,
				Success: true,
				CacheControl: domain.ResourceCacheControl{
					NoCache: true,
				},
				ResponseHeaders: map[string]string{"Cache-Control": "no-cache"},
				IgnoreErrors:    false,
				ResponseBody:    nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runner.NewDoneResource(tt.request, tt.response, tt.options)

			test.Equal(t, got, tt.expected)
		})
	}

}

func TestNewTimeoutResponse(t *testing.T) {
	timeoutErr := domain.ErrRequestTimeout

	tests := []struct {
		name     string
		request  domain.HttpRequest
		response domain.HttpResponse
		options  runner.DoneResourceOptions
		expected domain.DoneResource
	}{
		{
			"should create response for time outed execution",
			domain.HttpRequest{
				Schema:  "http",
				Host:    "hero.io",
				Path:    "/api",
				Query:   map[string]interface{}{"id": "123456"},
				Headers: map[string]string{"X-TID": "12345abdef"},
			},
			domain.HttpResponse{
				StatusCode: 408,
				Url:        "http://hero.io/api",
				Duration:   100 * time.Millisecond,
			},
			runner.DoneResourceOptions{},
			domain.DoneResource{
				Status:         408,
				Success:        false,
				IgnoreErrors:   false,
				Url:            "http://hero.io/api",
				RequestHeaders: map[string]string{"X-TID": "12345abdef"},
				RequestParams:  map[string]interface{}{"id": "123456"},
				ResponseTime:   100,
				ResponseBody:   timeoutErr.Error(),
			},
		},
		{
			"should create response for time outed execution with ignore errors",
			domain.HttpRequest{
				Schema:  "http",
				Host:    "hero.io",
				Path:    "/api",
				Query:   map[string]interface{}{"id": "123456"},
				Headers: map[string]string{"X-TID": "12345abdef"},
			},
			domain.HttpResponse{
				StatusCode: 408,
				Url:        "http://hero.io/api",
				Duration:   100 * time.Millisecond,
			},
			runner.DoneResourceOptions{IgnoreErrors: true},
			domain.DoneResource{
				Status:         408,
				Success:        false,
				IgnoreErrors:   true,
				Url:            "http://hero.io/api",
				RequestHeaders: map[string]string{"X-TID": "12345abdef"},
				RequestParams:  map[string]interface{}{"id": "123456"},
				ResponseTime:   100,
				ResponseBody:   timeoutErr.Error(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runner.NewErrorResponse(timeoutErr, tt.request, tt.response, tt.options)

			test.Equal(t, got, tt.expected)
		})
	}
}

func TestNewEmptyChainedResponse(t *testing.T) {
	t.Run("should create response for single empty chained param", func(t *testing.T) {
		params := []string{"id"}
		options := runner.DoneResourceOptions{}

		expected := domain.DoneResource{
			Status:       400,
			Success:      false,
			IgnoreErrors: false,
			ResponseBody: "The request was skipped due to missing { :id } param value",
		}

		got := runner.NewEmptyChainedResponse(params, options)

		test.Equal(t, got, expected)
	})

	t.Run("should create response for multiple empty chained param", func(t *testing.T) {
		params := []string{"id", "name", "city"}
		options := runner.DoneResourceOptions{}

		expected := domain.DoneResource{
			Status:       400,
			Success:      false,
			IgnoreErrors: false,
			ResponseBody: "The request was skipped due to missing { :id :name :city } param value",
		}

		got := runner.NewEmptyChainedResponse(params, options)

		test.Equal(t, got, expected)
	})

	t.Run("should create response for empty chained statement with ignore errors", func(t *testing.T) {
		params := []string{"id"}
		options := runner.DoneResourceOptions{IgnoreErrors: true}

		expected := domain.DoneResource{
			Status:       400,
			Success:      false,
			IgnoreErrors: true,
			ResponseBody: "The request was skipped due to missing { :id } param value",
		}

		got := runner.NewEmptyChainedResponse(params, options)

		test.Equal(t, got, expected)
	})
}

func TestGetEmptyChainedParams(t *testing.T) {
	tests := []struct {
		name      string
		statement domain.Statement
		expected  []string
	}{
		{
			"should return nothing if there is no empty chained param",
			domain.Statement{With: domain.Params{Values: map[string]interface{}{"id": "12345"}}},
			nil,
		},
		{
			"should return name of empty chained param",
			domain.Statement{With: domain.Params{Values: map[string]interface{}{"id": "12345", "name": runner.EmptyChained}}},
			[]string{"name"},
		},
		{
			"should return name of empty chained param inside list",
			domain.Statement{With: domain.Params{Values: map[string]interface{}{"id": "12345", "name": []interface{}{runner.EmptyChained}}}},
			[]string{"name"},
		},
		{
			"should return name of empty chained param inside map",
			domain.Statement{With: domain.Params{Values: map[string]interface{}{"id": "12345", "name": map[string]interface{}{"first": runner.EmptyChained}}}},
			[]string{"name"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runner.GetEmptyChainedParams(tt.statement)

			test.Equal(t, got, tt.expected)
		})
	}
}
