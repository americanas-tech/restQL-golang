package parser

import (
	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
	"github.com/b2wdigital/restQL-golang/v4/internal/parser/ast"
)

type Parser interface {
	Parse(queryStr string) (domain.Query, error)
}

type parser struct {
	astGenerator ast.Generator
}

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
