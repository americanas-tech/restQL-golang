package plugins_test

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/platform/plugins"
	"github.com/b2wdigital/restQL-golang/test"
	"net/http"
	"testing"
)

func TestDecodeQueryResult(t *testing.T) {
	tests := []struct {
		name        string
		queryResult domain.Resources
		expected    map[string]interface{}
	}{
		{
			"should make response for simple result",
			domain.Resources{
				"hero": domain.DoneResource{
					Status:          200,
					Success:         true,
					Method:          http.MethodGet,
					Url:             "http://hero.io/api",
					RequestHeaders:  map[string]string{"X-Token": "abcabcacbabc"},
					RequestParams:   map[string]interface{}{"filter": "no"},
					RequestBody:     map[string]interface{}{},
					ResponseTime:    100,
					ResponseHeaders: map[string]string{"X-New-Token": "efgefgefg"},
					ResponseBody:    test.Unmarshal(`{"id": "12345abcde"}`),
				},
			},
			map[string]interface{}{
				"hero": map[string]interface{}{
					"details": map[string]interface{}{
						"status":   200,
						"success":  true,
						"metadata": map[string]interface{}{},
						"debugging": map[string]interface{}{
							"method":          http.MethodGet,
							"url":             "http://hero.io/api",
							"requestHeaders":  map[string]string{"X-Token": "abcabcacbabc"},
							"params":          map[string]interface{}{"filter": "no"},
							"requestBody":     map[string]interface{}{},
							"responseHeaders": map[string]string{"X-New-Token": "efgefgefg"},
							"responseTime":    int64(100),
						},
					},
					"result": test.Unmarshal(`{"id": "12345abcde"}`),
				},
			},
		},
		{
			"should make response with metadata",
			domain.Resources{
				"hero": domain.DoneResource{
					Status:          200,
					Success:         true,
					Method:          http.MethodGet,
					IgnoreErrors:    true,
					Url:             "http://hero.io/api",
					RequestHeaders:  map[string]string{"X-Token": "abcabcacbabc"},
					RequestParams:   map[string]interface{}{"filter": "no"},
					RequestBody:     map[string]interface{}{},
					ResponseTime:    100,
					ResponseHeaders: map[string]string{"X-New-Token": "efgefgefg"},
					ResponseBody:    test.Unmarshal(`{"id": "12345abcde"}`),
				},
			},
			map[string]interface{}{
				"hero": map[string]interface{}{
					"details": map[string]interface{}{
						"status":   200,
						"success":  true,
						"metadata": map[string]interface{}{"ignore-errors": true},
						"debugging": map[string]interface{}{
							"method":          http.MethodGet,
							"url":             "http://hero.io/api",
							"requestHeaders":  map[string]string{"X-Token": "abcabcacbabc"},
							"params":          map[string]interface{}{"filter": "no"},
							"requestBody":     map[string]interface{}{},
							"responseHeaders": map[string]string{"X-New-Token": "efgefgefg"},
							"responseTime":    int64(100),
						},
					},
					"result": test.Unmarshal(`{"id": "12345abcde"}`),
				},
			},
		},
		{
			"should make response for multiplexed result",
			domain.Resources{
				"hero": domain.DoneResources{
					domain.DoneResource{
						Status:          200,
						Success:         true,
						Method:          http.MethodGet,
						Url:             "http://hero.io/api",
						RequestHeaders:  map[string]string{"X-Token": "abcabcacbabc"},
						RequestParams:   map[string]interface{}{"filter": "no"},
						RequestBody:     nil,
						ResponseTime:    100,
						ResponseHeaders: map[string]string{"X-New-Token": "efgefgefg"},
						ResponseBody:    test.Unmarshal(`{"id": "12345abcde"}`),
					},
					domain.DoneResource{
						Status:          200,
						Success:         true,
						Method:          http.MethodGet,
						Url:             "http://hero.io/api",
						RequestHeaders:  map[string]string{"X-Token": "abcabcacbabc"},
						RequestParams:   map[string]interface{}{"filter": "no"},
						RequestBody:     nil,
						ResponseTime:    100,
						ResponseHeaders: map[string]string{"X-New-Token": "kgkgkgkgkgkgkg"},
						ResponseBody:    test.Unmarshal(`{"id": "98765lkjhgf"}`),
					},
				},
			},
			map[string]interface{}{
				"hero": map[string]interface{}{
					"details": []interface{}{
						map[string]interface{}{
							"status":   200,
							"success":  true,
							"metadata": map[string]interface{}{},
							"debugging": map[string]interface{}{
								"method":          http.MethodGet,
								"url":             "http://hero.io/api",
								"requestHeaders":  map[string]string{"X-Token": "abcabcacbabc"},
								"params":          map[string]interface{}{"filter": "no"},
								"requestBody":     nil,
								"responseHeaders": map[string]string{"X-New-Token": "efgefgefg"},
								"responseTime":    int64(100),
							},
						},
						map[string]interface{}{
							"status":   200,
							"success":  true,
							"metadata": map[string]interface{}{},
							"debugging": map[string]interface{}{
								"method":          http.MethodGet,
								"url":             "http://hero.io/api",
								"requestHeaders":  map[string]string{"X-Token": "abcabcacbabc"},
								"params":          map[string]interface{}{"filter": "no"},
								"requestBody":     nil,
								"responseHeaders": map[string]string{"X-New-Token": "kgkgkgkgkgkgkg"},
								"responseTime":    int64(100),
							},
						},
					},
					"result": []interface{}{
						test.Unmarshal(`{"id": "12345abcde"}`),
						test.Unmarshal(`{"id": "98765lkjhgf"}`),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := plugins.DecodeQueryResult(tt.queryResult)
			test.Equal(t, got, tt.expected)
		})
	}

}
