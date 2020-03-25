package eval_test

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/eval"
	"reflect"
	"testing"
)

func TestHiddenFilter(t *testing.T) {
	query := domain.Query{Statements: []domain.Statement{
		{Resource: "hero", Hidden: true},
		{Resource: "sidekick"},
	}}

	resources := domain.Resources{
		"hero":     domain.DoneResource{Details: domain.Details{Success: true}, Result: nil},
		"sidekick": domain.DoneResource{Details: domain.Details{Success: true}, Result: nil},
	}

	expectedResources := domain.Resources{
		"sidekick": domain.DoneResource{Details: domain.Details{Success: true}, Result: nil},
	}

	got := eval.ApplyFilters(query, resources)

	if !reflect.DeepEqual(got, expectedResources) {
		t.Fatalf("ApplyFilters = %+#v, want = %+#v", got, expectedResources)
	}
}
