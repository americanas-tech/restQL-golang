package parser

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/parser/ast"
)

type Parser struct {
	astGenerator ast.Generator
}

func New() (Parser, error) {
	generator, err := ast.New()
	if err != nil {
		return Parser{}, err
	}

	return Parser{astGenerator: generator}, nil
}

func (p Parser) Parse(queryStr string) (domain.Query, error) {
	query, err := p.astGenerator.Parse(queryStr)
	if err != nil {
		return domain.Query{}, err
	}

	return Optimize(query)
}
