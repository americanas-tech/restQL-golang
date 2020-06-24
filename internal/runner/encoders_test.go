package runner_test

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/runner"
	"github.com/b2wdigital/restQL-golang/test"
	"testing"
)

func TestApplyEncoders(t *testing.T) {
	tests := []struct {
		name      string
		resources domain.Resources
		expected  domain.Resources
	}{
		{
			"should do nothing if there is no encode",
			domain.Resources{"hero": domain.Statement{Method: "from", Resource: "hero"}},
			domain.Resources{"hero": domain.Statement{Method: "from", Resource: "hero"}},
		},
		{
			"should not apply encoders to chained value",
			domain.Resources{"hero": domain.Statement{
				Method:   "from",
				Resource: "hero",
				With: domain.Params{Values: map[string]interface{}{
					"base64": domain.Base64{Value: domain.Chain{"done-resource", "id"}},
					"json":   domain.Json{Value: domain.Chain{"done-resource", "id"}},
				}},
			}},
			domain.Resources{"hero": domain.Statement{
				Method:   "from",
				Resource: "hero",
				With: domain.Params{Values: map[string]interface{}{
					"base64": domain.Base64{Value: domain.Chain{"done-resource", "id"}},
					"json":   domain.Json{Value: domain.Chain{"done-resource", "id"}},
				}},
			}},
		},
		{
			"should apply base64 encoder to with value",
			domain.Resources{"hero": domain.Statement{
				Method:   "from",
				Resource: "hero",
				With: domain.Params{Values: map[string]interface{}{
					"id": domain.Base64{Value: "12345abcdef"},
				}},
			}},
			domain.Resources{"hero": domain.Statement{
				Method:   "from",
				Resource: "hero",
				With: domain.Params{Values: map[string]interface{}{
					"id": "MTIzNDVhYmNkZWY=",
				}},
			}},
		},
		{
			"should apply base64 encoder to with body",
			domain.Resources{"hero": domain.Statement{
				Method:   "from",
				Resource: "hero",
				With: domain.Params{
					Body: domain.Base64{Value: map[string]interface{}{"id": "test"}},
					Values: map[string]interface{}{
						"id": domain.Base64{Value: "12345abcdef"},
					},
				},
			}},
			domain.Resources{"hero": domain.Statement{
				Method:   "from",
				Resource: "hero",
				With: domain.Params{
					Body: "bWFwW2lkOnRlc3Rd",
					Values: map[string]interface{}{
						"id": "MTIzNDVhYmNkZWY=",
					},
				},
			}},
		},
		{
			"should apply multiple encoders to with body",
			domain.Resources{"hero": domain.Statement{
				Method:   "from",
				Resource: "hero",
				With: domain.Params{
					Body: domain.NoMultiplex{Value: domain.Json{Value: map[string]interface{}{"id": "test"}}},
					Values: map[string]interface{}{
						"id": domain.Base64{Value: "12345abcdef"},
					},
				},
			}},
			domain.Resources{"hero": domain.Statement{
				Method:   "from",
				Resource: "hero",
				With: domain.Params{
					Body: domain.NoMultiplex{Value: map[string]interface{}{"id": "test"}},
					Values: map[string]interface{}{
						"id": "MTIzNDVhYmNkZWY=",
					},
				},
			}},
		},
		{
			"should unwrap json encoder in with body",
			domain.Resources{"hero": domain.Statement{
				Method:   "from",
				Resource: "hero",
				With: domain.Params{
					Body: domain.Json{Value: map[string]interface{}{"id": "test"}},
					Values: map[string]interface{}{
						"id": domain.Base64{Value: "12345abcdef"},
					},
				},
			}},
			domain.Resources{"hero": domain.Statement{
				Method:   "from",
				Resource: "hero",
				With: domain.Params{
					Body: map[string]interface{}{"id": "test"},
					Values: map[string]interface{}{
						"id": "MTIzNDVhYmNkZWY=",
					},
				},
			}},
		},
		{
			"should apply json encoder to with value",
			domain.Resources{"hero": domain.Statement{
				Method:   "from",
				Resource: "hero",
				With: domain.Params{Values: map[string]interface{}{
					"weapons": domain.Json{Value: map[string]interface{}{"id": 1, "name": "sword"}},
				}},
			}},
			domain.Resources{"hero": domain.Statement{
				Method:   "from",
				Resource: "hero",
				With: domain.Params{Values: map[string]interface{}{
					"weapons": `{"id":1,"name":"sword"}`,
				}},
			}},
		},
		{
			"should apply nested encoders to with value",
			domain.Resources{"hero": domain.Statement{
				Method:   "from",
				Resource: "hero",
				With: domain.Params{Values: map[string]interface{}{
					"weapons": domain.Base64{Value: domain.Json{Value: map[string]interface{}{"id": 1, "name": "sword"}}},
				}},
			}},
			domain.Resources{"hero": domain.Statement{
				Method:   "from",
				Resource: "hero",
				With: domain.Params{Values: map[string]interface{}{
					"weapons": "eyJpZCI6MSwibmFtZSI6InN3b3JkIn0=",
				}},
			}},
		},
		{
			"should apply encoder nested with flatten to with value",
			domain.Resources{"hero": domain.Statement{
				Method:   "from",
				Resource: "hero",
				With: domain.Params{Values: map[string]interface{}{
					"weapons": domain.NoMultiplex{Value: domain.Json{Value: []interface{}{"id", "name"}}},
				}},
			}},
			domain.Resources{"hero": domain.Statement{
				Method:   "from",
				Resource: "hero",
				With: domain.Params{Values: map[string]interface{}{
					"weapons": domain.NoMultiplex{Value: `["id","name"]`},
				}},
			}},
		},
		{
			"should apply encoders inside nested data structures",
			domain.Resources{"hero": domain.Statement{
				Method:   "from",
				Resource: "hero",
				With: domain.Params{Values: map[string]interface{}{
					"map": map[string]interface{}{
						"nested": map[string]interface{}{
							"json":   domain.Json{Value: map[string]interface{}{"id": 1, "name": "sword"}},
							"base64": domain.Base64{Value: "12345abcdef"},
						},
					},
					"list": []interface{}{
						[]interface{}{
							domain.Json{Value: map[string]interface{}{"id": 1, "name": "sword"}},
							domain.Base64{Value: "12345abcdef"},
						},
					},
				}},
			}},
			domain.Resources{"hero": domain.Statement{
				Method:   "from",
				Resource: "hero",
				With: domain.Params{Values: map[string]interface{}{
					"map": map[string]interface{}{
						"nested": map[string]interface{}{
							"json":   `{"id":1,"name":"sword"}`,
							"base64": "MTIzNDVhYmNkZWY=",
						},
					},
					"list": []interface{}{
						[]interface{}{
							`{"id":1,"name":"sword"}`,
							"MTIzNDVhYmNkZWY=",
						},
					},
				}},
			}},
		},
	}

	logger := noOpLogger{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runner.ApplyEncoders(tt.resources, logger)
			test.Equal(t, got, tt.expected)
		})
	}
}

type noOpLogger struct{}

func (n noOpLogger) Panic(msg string, fields ...interface{})            {}
func (n noOpLogger) Fatal(msg string, fields ...interface{})            {}
func (n noOpLogger) Error(msg string, err error, fields ...interface{}) {}
func (n noOpLogger) Warn(msg string, fields ...interface{})             {}
func (n noOpLogger) Info(msg string, fields ...interface{})             {}
func (n noOpLogger) Debug(msg string, fields ...interface{})            {}
