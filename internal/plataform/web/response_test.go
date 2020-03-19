package web_test

import (
	"encoding/json"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/plataform/web"
	"reflect"
	"testing"
)

func TestMakeQueryResponse(t *testing.T) {
	tests := []struct {
		name        string
		queryResult domain.Resources
		expected    web.QueryResponse
	}{
		{
			"should make response for simple result",
			domain.Resources{
				"hero": domain.DoneResource{
					Details: domain.Details{Status: 200, Success: true},
					Result:  unmarshal(`{"id": "12345abcde"}`),
				},
			},
			web.QueryResponse{
				"hero": web.StatementResult{
					Details: web.StatementDetails{Status: 200, Success: true},
					Result:  unmarshal(`{"id": "12345abcde"}`),
				},
			},
		},
		{
			"should make response with metadata",
			domain.Resources{
				"hero": domain.DoneResource{
					Details: domain.Details{Status: 200, Success: true, IgnoreErrors: true},
					Result:  unmarshal(`{"id": "12345abcde"}`),
				},
			},
			web.QueryResponse{
				"hero": web.StatementResult{
					Details: web.StatementDetails{Status: 200, Success: true, Metadata: web.StatementMetadata{IgnoreErrors: "ignore"}},
					Result:  unmarshal(`{"id": "12345abcde"}`),
				},
			},
		},
		{
			"should make response with debugging",
			domain.Resources{
				"hero": domain.DoneResource{
					Details: domain.Details{Status: 200, Success: true, Debug: &domain.Debugging{
						Url:             "http://hero.io/api",
						RequestHeaders:  map[string]string{"X-Token": "abcabcacbabc"},
						ResponseHeaders: map[string]string{"X-New-Token": "efgefgefg"},
						Params:          map[string]interface{}{"filter": "no"},
						ResponseTime:    100,
					}},
					Result: unmarshal(`{"id": "12345abcde"}`),
				},
			},
			web.QueryResponse{
				"hero": web.StatementResult{
					Details: web.StatementDetails{Status: 200, Success: true, Debug: &web.StatementDebugging{
						Url:             "http://hero.io/api",
						RequestHeaders:  map[string]string{"X-Token": "abcabcacbabc"},
						ResponseHeaders: map[string]string{"X-New-Token": "efgefgefg"},
						Params:          map[string]interface{}{"filter": "no"},
						ResponseTime:    100,
					}},
					Result: unmarshal(`{"id": "12345abcde"}`),
				},
			},
		},
		{
			"should make response for multiplexed result",
			domain.Resources{
				"hero": domain.DoneResources{
					domain.DoneResource{
						Details: domain.Details{Status: 200, Success: true},
						Result:  unmarshal(`{"id": "12345abcde"}`),
					},
					domain.DoneResource{
						Details: domain.Details{Status: 200, Success: true},
						Result:  unmarshal(`{"id": "67890fghij"}`),
					},
				},
			},
			web.QueryResponse{
				"hero": web.StatementResult{
					Details: []interface{}{web.StatementDetails{Status: 200, Success: true}, web.StatementDetails{Status: 200, Success: true}},
					Result:  []interface{}{unmarshal(`{"id": "12345abcde"}`), unmarshal(`{"id": "67890fghij"}`)},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := web.MakeQueryResponse(tt.queryResult)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("MakeQueryResponse = %+#v, want = %+#v", got, tt.expected)
			}
		})
	}

}

func TestCalculateStatusCode(t *testing.T) {
	tests := []struct {
		name        string
		queryResult domain.Resources
		expected    int
	}{
		{
			"should return 200 when resources are successful",
			domain.Resources{
				"hero":     domain.DoneResource{Details: domain.Details{Status: 200}},
				"sidekick": domain.DoneResource{Details: domain.Details{Status: 204}},
				"villain":  domain.DoneResource{Details: domain.Details{Status: 201}},
			},
			200,
		},
		{
			"should return max status code",
			domain.Resources{
				"hero":     domain.DoneResource{Details: domain.Details{Status: 200}},
				"sidekick": domain.DoneResource{Details: domain.Details{Status: 500}},
				"villain":  domain.DoneResource{Details: domain.Details{Status: 408}},
			},
			500,
		},
		{
			"should return max status code",
			domain.Resources{
				"hero": domain.DoneResources{
					domain.DoneResources{
						domain.DoneResource{Details: domain.Details{Status: 200}},
						domain.DoneResource{Details: domain.Details{Status: 200}},
						domain.DoneResource{Details: domain.Details{Status: 408}},
					},
				},
				"sidekick": domain.DoneResource{Details: domain.Details{Status: 204}},
				"villain":  domain.DoneResource{Details: domain.Details{Status: 400}},
			},
			408,
		},
		{
			"should return max status code expect for result marked with ignore",
			domain.Resources{
				"hero":     domain.DoneResource{Details: domain.Details{Status: 200}},
				"sidekick": domain.DoneResource{Details: domain.Details{Status: 500, IgnoreErrors: true}},
				"villain":  domain.DoneResource{Details: domain.Details{Status: 400, IgnoreErrors: true}},
			},
			200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := web.CalculateStatusCode(tt.queryResult)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("CalculateStatusCode = %#+v, want = %#+v", got, tt.expected)
			}
		})
	}
}

func unmarshal(body string) interface{} {
	var f interface{}
	err := json.Unmarshal([]byte(body), &f)
	if err != nil {
		panic(err)
	}
	return f
}
