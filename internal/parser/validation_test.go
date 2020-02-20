package parser

import (
	"github.com/b2wdigital/restQL-golang/internal/parser/ast"
	"testing"
)

func TestParseValidation(t *testing.T) {
	t.Run("should have only one qualifier of each type (with, only, timeout, headers, etc)", func(t *testing.T) {
		err := ValidateQuery(&ast.Query{
			Blocks: []ast.Block{
				{
					Method:   "from",
					Resource: "hero",
					Qualifiers: []ast.Qualifier{
						{
							With: []ast.WithItem{{Key: "id", Value: ast.Value{Primitive: &ast.Primitive{Int: Int(1)}}}},
						},
						{
							With: []ast.WithItem{{Key: "name", Value: ast.Value{Primitive: &ast.Primitive{String: String("batman")}}}},
						},
					},
				},
			},
		})

		if err == nil {
			t.Fatalf("Expected an error from Parse but didn't got one")
		}

		if err != ErrDuplicatedQualifiers {
			t.Fatalf("Expected error %v got %v", ErrDuplicatedQualifiers, err)
		}
	})

	t.Run("should not have only and hidden keywords in the same block", func(t *testing.T) {
		err := ValidateQuery(&ast.Query{
			Blocks: []ast.Block{
				{
					Method:   "from",
					Resource: "hero",
					Qualifiers: []ast.Qualifier{
						{
							Only: []ast.Filter{{Field: "name"}},
						},
						{
							Hidden: true,
						},
					},
				},
			},
		})

		if err == nil {
			t.Fatalf("Expected an error from Parse but didn't got one")
		}

		if err != ErrInvalidOnlyAndHiddenKeywordsInSameStatement {
			t.Fatalf("Expected error %v got %v", ErrDuplicatedQualifiers, err)
		}
	})
}

func String(s string) *string {
	return &s
}

func Int(i int) *int {
	return &i
}
