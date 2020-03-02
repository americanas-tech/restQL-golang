package parser

import (
	"errors"
	"github.com/b2wdigital/restQL-golang/internal/parser/ast"
)

var (
	ErrDuplicatedQualifiers                        = errors.New("duplicated qualifiers found in query")
	ErrInvalidOnlyAndHiddenKeywordsInSameStatement = errors.New("it is not allowed to use only and hidden keywords together")
)

func validateQuery(q *ast.Query) error {
	for _, block := range q.Blocks {
		err := validateBlock(block)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateBlock(block ast.Block) error {
	occurrences := reduceBlockToOccurrences(block)

	for _, occurrence := range occurrences {
		if occurrence > 1 {
			return ErrDuplicatedQualifiers
		}
	}

	if occurrences[ast.OnlyKeyword] == 1 && occurrences[ast.HiddenKeyword] == 1 {
		return ErrInvalidOnlyAndHiddenKeywordsInSameStatement
	}

	return nil
}

func reduceBlockToOccurrences(block ast.Block) map[string]int {
	occurrences := newOccurrenceMap()

	for _, qualifier := range block.Qualifiers {
		if qualifier.With != nil {
			occurrences[ast.WithKeyword] += 1
		}

		if qualifier.Only != nil {
			occurrences[ast.OnlyKeyword] += 1
		}

		if qualifier.Timeout != nil {
			occurrences[ast.TimeoutKeyword] += 1
		}

		if qualifier.Headers != nil {
			occurrences[ast.HeadersKeyword] += 1
		}

		if qualifier.MaxAge != nil {
			occurrences[ast.MaxAgeKeyword] += 1
		}

		if qualifier.SMaxAge != nil {
			occurrences[ast.SmaxAgeKeyword] += 1
		}

		if qualifier.Hidden == true {
			occurrences[ast.HiddenKeyword] += 1
		}

		if qualifier.IgnoreErrors == true {
			occurrences[ast.IgnoreErrorsKeyword] += 1
		}
	}

	return occurrences
}

func newOccurrenceMap() map[string]int {
	occurrences := map[string]int{}

	occurrences[ast.WithKeyword] = 0
	occurrences[ast.OnlyKeyword] = 0
	occurrences[ast.TimeoutKeyword] = 0
	occurrences[ast.HeadersKeyword] = 0
	occurrences[ast.HiddenKeyword] = 0
	occurrences[ast.MaxAgeKeyword] = 0
	occurrences[ast.SmaxAgeKeyword] = 0
	occurrences[ast.IgnoreErrorsKeyword] = 0

	return occurrences
}
