package eval_test

import (
	"encoding/json"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/eval"
	"reflect"
	"testing"
)

func TestEvaluateSavedQuery(t *testing.T) {
	config := stubConfig{
		UnmarshalResult: map[string]interface{}{
			"queries": map[string]interface{}{
				"restQL": map[string][]string{"my-query": {"from hero"}},
			},
			"mappings": map[string]string{
				"hero": "http://heroes.com",
			},
		},
	}

	mr := eval.NewMappingReader(config, NoOpLogger{})
	qr := eval.NewQueryReader(config, NoOpLogger{})
	evaluator := eval.NewEvaluator(mr, qr, NoOpLogger{})

	t.Run("eval query found with no variables", func(t *testing.T) {
		options := eval.QueryOptions{
			Namespace: "restQL",
			Id:        "my-query",
			Revision:  1,
		}
		input := eval.QueryInput{}

		query, err := evaluator.SavedQuery(nil, options, input)
		if err != nil {
			t.Fatalf("evaluate.SavedQuery return an error : %v", err)
		}

		expected := domain.Query{
			Statements: []domain.Statement{
				{Method: "from", Resource: "hero"},
			},
		}

		if !reflect.DeepEqual(query, expected) {
			t.Fatalf("evaluate.SavedQuery return %+v, expected %+v", query, expected)
		}
	})

}

type NoOpLogger struct{}

func (n NoOpLogger) Panic(msg string, fields ...interface{})            {}
func (n NoOpLogger) Fatal(msg string, fields ...interface{})            {}
func (n NoOpLogger) Error(msg string, err error, fields ...interface{}) {}
func (n NoOpLogger) Warn(msg string, fields ...interface{})             {}
func (n NoOpLogger) Info(msg string, fields ...interface{})             {}
func (n NoOpLogger) Debug(msg string, fields ...interface{})            {}

type stubConfig struct {
	UnmarshalResult interface{}
	EnvResult       string
}

func (sc stubConfig) Build() string {
	return "test"
}

type stubFileSource struct {
	unmarshalResult interface{}
}

func (sfc stubFileSource) Unmarshal(target interface{}) error {
	bytes, _ := json.Marshal(sfc.unmarshalResult)
	json.Unmarshal(bytes, target)

	return nil
}

func (sc stubConfig) File() domain.FileSource {
	return stubFileSource{sc.UnmarshalResult}
}

type stubEnvSource struct {
	envResult string
}

func (sec stubEnvSource) GetString(key string) string {
	return sec.envResult
}

func (sc stubConfig) Env() domain.EnvSource {
	return stubEnvSource{sc.EnvResult}
}
