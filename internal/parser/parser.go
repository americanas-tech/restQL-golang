package parser

import (
	"errors"
	"github.com/b2wdigital/restQL-golang/v6/internal/domain"
	"github.com/b2wdigital/restQL-golang/v6/internal/parser/ast"
)

// ErrInvalidQuery represents a given query that not comply with the restQL syntax
var ErrInvalidQuery = errors.New("invalid query")

// Parser is the interface implemented by types that
// can transform a query string into an internal representation.
type Parser interface {
	Parse(queryStr string) (domain.Query, error)
}

type parser struct {
	astGenerator ast.Generator
}

// New returns an instance of a Parser.
func New() (Parser, error) {
	generator, err := ast.New()
	if err != nil {
		return parser{}, err
	}

	return parser{astGenerator: generator}, nil
}

func (p parser) Parse(queryStr string) (domain.Query, error) {
	query, err := p.astGenerator.Parse(queryStr)
	if err != nil {
		return domain.Query{}, err
	}

	return Optimize(query)
}
