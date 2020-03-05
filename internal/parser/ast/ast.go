package ast

import (
	"fmt"
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer/regex"
	"github.com/pkg/errors"
)

var parser *participle.Parser

func init() {
	if parser != nil {
		return
	}

	newParser, err := makeParser()
	if err != nil {
		panic(err)
	}

	parser = newParser
}

func Parse(query string) (q *Query, err error) {
	q = &Query{}
	if err = parser.ParseString(query, q); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to parse query: %s", query))
	}

	return q, nil
}

func makeParser() (*participle.Parser, error) {
	QUERY := &Query{}
	definition, err := regex.New(lexerDefinition)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to compile lexer definition"))
	}

	parser, err = participle.Build(
		QUERY,
		participle.Lexer(definition),
		participle.Unquote("String"),
		participle.Elide("Comment"),
	)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to compile parser definition"))
	}

	return parser, nil
}
