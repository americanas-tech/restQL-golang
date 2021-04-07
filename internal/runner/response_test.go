package runner_test

import (
	"github.com/b2wdigital/restQL-golang/v6/pkg/restql"
	"testing"
	"time"

	"github.com/b2wdigital/restQL-golang/v6/internal/domain"
	"github.com/b2wdigital/restQL-golang/v6/internal/runner"
	"github.com/b2wdigital/restQL-golang/v6/test"
)

func TestNewDoneResource(t *testing.T) {
	tests := []struct {
		name     string
		request  restql.HTTPRequest
		response restql.HTTPResponse
		options  runner.DoneResourceOptions
		expected restql.DoneResource
	}{
		{
			"should create done resource for successful execution",
			restql.HTTPRequest{},
			restql.HTTPResponse{StatusCode: 200, Body: nil},
			runner.DoneResourceOptions{},
			restql.DoneResource{Status: 200, Success: true, IgnoreErrors: false, ResponseBody: nil},
		},
		{
			"should create done resource for failed execution",
			restql.HTTPRequest{},
			restql.HTTPResponse{StatusCode: 400, Body: nil},
			runner.DoneResourceOptions{},
			restql.DoneResource{Status: 400, Success: false, IgnoreErrors: false, ResponseBody: nil},
		},
		{
			"should create done resource with debug",
			restql.HTTPRequest{
				Schema:  "http",
				Host:    "hero.io",
				Path:    "/api",
				Query:   map[string]interface{}{"id": "123456"},
				Headers: map[string]string{"X-TID": "12345abdef"},
			},
			restql.HTTPResponse{
				URL:        "http://hero.io/api",
				StatusCode: 200,
				Body:       nil,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Duration:   100 * time.Millisecond,
			},
			runner.DoneResourceOptions{},
			restql.DoneResource{
				Status:          200,
				Success:         true,
				IgnoreErrors:    false,
				URL:             "http://hero.io/api",
				RequestHeaders:  map[string]string{"X-TID": "12345abdef"},
				ResponseHeaders: map[string]string{"Content-Type": "application/json"},
				RequestParams:   map[string]interface{}{"id": "123456"},
				ResponseTime:    100,
				ResponseBody:    nil,
			},
		},
		{
			"should create done resource with ignore errors",
			restql.HTTPRequest{},
			restql.HTTPResponse{StatusCode: 200, Body: nil},
			runner.DoneResourceOptions{IgnoreErrors: true},
			restql.DoneResource{
				Status:       200,
				Success:      true,
				IgnoreErrors: true,
				ResponseBody: nil,
			},
		},
		{
			"should create done resource with cache control information returned by resource",
			restql.HTTPRequest{},
			restql.HTTPResponse{StatusCode: 200, Body: nil, Headers: map[string]string{"Cache-Control": "max-age=400, s-maxage=600"}},
			runner.DoneResourceOptions{},
			restql.DoneResource{
				Status:  200,
				Success: true,
				CacheControl: restql.ResourceCacheControl{
					MaxAge:  restql.ResourceCacheControlValue{Exist: true, Time: 400},
					SMaxAge: restql.ResourceCacheControlValue{Exist: true, Time: 600},
				},
				ResponseHeaders: map[string]string{"Cache-Control": "max-age=400, s-maxage=600"},
				IgnoreErrors:    false,
				ResponseBody:    nil,
			},
		},
		{
			"should create done resource with cache control information returned by resource",
			restql.HTTPRequest{},
			restql.HTTPResponse{StatusCode: 200, Body: nil, Headers: map[string]string{"Cache-Control": "no-cache"}},
			runner.DoneResourceOptions{},
			restql.DoneResource{
				Status:  200,
				Success: true,
				CacheControl: restql.ResourceCacheControl{
					NoCache: true,
				},
				ResponseHeaders: map[string]string{"Cache-Control": "no-cache"},
				IgnoreErrors:    false,
				ResponseBody:    nil,
			},
		},
		{
			"should create done resource with cache control information defined in statement if not returned by resource",
			restql.HTTPRequest{},
			restql.HTTPResponse{StatusCode: 200, Body: nil},
			runner.DoneResourceOptions{MaxAge: 100, SMaxAge: 300},
			restql.DoneResource{
				Status:  200,
				Success: true,
				CacheControl: restql.ResourceCacheControl{
					MaxAge:  restql.ResourceCacheControlValue{Exist: true, Time: 100},
					SMaxAge: restql.ResourceCacheControlValue{Exist: true, Time: 300},
				},
				IgnoreErrors: false,
				ResponseBody: nil,
			},
		},
		{
			"should create done resource with cache control information defined in statement if not returned by resource",
			restql.HTTPRequest{},
			restql.HTTPResponse{StatusCode: 200, Body: nil},
			runner.DoneResourceOptions{MaxAge: 100},
			restql.DoneResource{
				Status:  200,
				Success: true,
				CacheControl: restql.ResourceCacheControl{
					MaxAge: restql.ResourceCacheControlValue{Exist: true, Time: 100},
				},
				IgnoreErrors: false,
				ResponseBody: nil,
			},
		},
		{
			"should create done resource with minimum cache control information between the returned by resource and the defined in statement",
			restql.HTTPRequest{},
			restql.HTTPResponse{StatusCode: 200, Body: nil, Headers: map[string]string{"Cache-Control": "max-age=100, s-maxage=600"}},
			runner.DoneResourceOptions{MaxAge: 400, SMaxAge: 300},
			restql.DoneResource{
				Status:  200,
				Success: true,
				CacheControl: restql.ResourceCacheControl{
					MaxAge:  restql.ResourceCacheControlValue{Exist: true, Time: 100},
					SMaxAge: restql.ResourceCacheControlValue{Exist: true, Time: 300},
				},
				ResponseHeaders: map[string]string{"Cache-Control": "max-age=100, s-maxage=600"},
				IgnoreErrors:    false,
				ResponseBody:    nil,
			},
		},
		{
			"should create done resource with minimum cache control information between the returned by resource and the defined in statement",
			restql.HTTPRequest{},
			restql.HTTPResponse{StatusCode: 200, Body: nil, Headers: map[string]string{"Cache-Control": "no-cache"}},
			runner.DoneResourceOptions{MaxAge: 400, SMaxAge: 300},
			restql.DoneResource{
				Status:  200,
				Success: true,
				CacheControl: restql.ResourceCacheControl{
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
		request  restql.HTTPRequest
		response restql.HTTPResponse
		options  runner.DoneResourceOptions
		expected restql.DoneResource
	}{
		{
			"should create response for time outed execution",
			restql.HTTPRequest{
				Schema:  "http",
				Host:    "hero.io",
				Path:    "/api",
				Query:   map[string]interface{}{"id": "123456"},
				Headers: map[string]string{"X-TID": "12345abdef"},
			},
			restql.HTTPResponse{
				StatusCode: 408,
				URL:        "http://hero.io/api",
				Duration:   100 * time.Millisecond,
			},
			runner.DoneResourceOptions{},
			restql.DoneResource{
				Status:         408,
				Success:        false,
				IgnoreErrors:   false,
				URL:            "http://hero.io/api",
				RequestHeaders: map[string]string{"X-TID": "12345abdef"},
				RequestParams:  map[string]interface{}{"id": "123456"},
				ResponseTime:   100,
				ResponseBody:   restql.NewResponseBodyFromValue(test.NoOpLogger, timeoutErr.Error()),
			},
		},
		{
			"should create response for time outed execution with ignore errors",
			restql.HTTPRequest{
				Schema:  "http",
				Host:    "hero.io",
				Path:    "/api",
				Query:   map[string]interface{}{"id": "123456"},
				Headers: map[string]string{"X-TID": "12345abdef"},
			},
			restql.HTTPResponse{
				StatusCode: 408,
				URL:        "http://hero.io/api",
				Duration:   100 * time.Millisecond,
			},
			runner.DoneResourceOptions{IgnoreErrors: true},
			restql.DoneResource{
				Status:         408,
				Success:        false,
				IgnoreErrors:   true,
				URL:            "http://hero.io/api",
				RequestHeaders: map[string]string{"X-TID": "12345abdef"},
				RequestParams:  map[string]interface{}{"id": "123456"},
				ResponseTime:   100,
				ResponseBody:   restql.NewResponseBodyFromValue(test.NoOpLogger, timeoutErr.Error()),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runner.NewErrorResponse(test.NoOpLogger, timeoutErr, tt.request, tt.response, tt.options)

			test.Equal(t, got, tt.expected)
		})
	}
}

func TestNewEmptyChainedResponse(t *testing.T) {
	t.Run("should create response for single empty chained param", func(t *testing.T) {
		params := []string{"id"}
		options := runner.DoneResourceOptions{}

		expected := restql.DoneResource{
			Status:       400,
			Success:      false,
			IgnoreErrors: false,
			ResponseBody: restql.NewResponseBodyFromValue(
				test.NoOpLogger,
				"The request was skipped due to missing { :id } param value",
			),
		}

		got := runner.NewEmptyChainedResponse(test.NoOpLogger, params, options)

		test.Equal(t, got, expected)
	})

	t.Run("should create response for multiple empty chained param", func(t *testing.T) {
		params := []string{"id", "name", "city"}
		options := runner.DoneResourceOptions{}

		expected := restql.DoneResource{
			Status:       400,
			Success:      false,
			IgnoreErrors: false,
			ResponseBody: restql.NewResponseBodyFromValue(
				test.NoOpLogger,
				"The request was skipped due to missing { :id :name :city } param value",
			),
		}

		got := runner.NewEmptyChainedResponse(test.NoOpLogger, params, options)

		test.Equal(t, got, expected)
	})

	t.Run("should create response for empty chained statement with ignore errors", func(t *testing.T) {
		params := []string{"id"}
		options := runner.DoneResourceOptions{IgnoreErrors: true}

		expected := restql.DoneResource{
			Status:       400,
			Success:      false,
			IgnoreErrors: true,
			ResponseBody: restql.NewResponseBodyFromValue(
				test.NoOpLogger,
				"The request was skipped due to missing { :id } param value",
			),
		}

		got := runner.NewEmptyChainedResponse(test.NoOpLogger, params, options)

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

func TestNewDependsOnUnresolvedResponse(t *testing.T) {
	t.Run("should create response for unresolved statement", func(t *testing.T) {
		statement := domain.Statement{
			DependsOn: domain.DependsOn{
				Target:   "failed-resource",
				Resolved: false,
			},
		}
		options := runner.DoneResourceOptions{}

		expected := restql.DoneResource{
			Status:       400,
			Success:      false,
			IgnoreErrors: false,
			ResponseBody: restql.NewResponseBodyFromValue(
				test.NoOpLogger,
				"The request was skipped due to unresolved dependency { failed-resource }",
			),
		}

		got := runner.NewNewDependsOnUnresolvedResponse(test.NoOpLogger, statement, options)

		test.Equal(t, got, expected)
	})

	t.Run("should create response for empty chained statement with ignore errors", func(t *testing.T) {
		statement := domain.Statement{
			DependsOn: domain.DependsOn{
				Target:   "failed-resource",
				Resolved: false,
			},
		}
		options := runner.DoneResourceOptions{IgnoreErrors: true}

		expected := restql.DoneResource{
			Status:       400,
			Success:      false,
			IgnoreErrors: true,
			ResponseBody: restql.NewResponseBodyFromValue(
				test.NoOpLogger,
				"The request was skipped due to unresolved dependency { failed-resource }",
			),
		}

		got := runner.NewNewDependsOnUnresolvedResponse(test.NoOpLogger, statement, options)

		test.Equal(t, got, expected)
	})
}
