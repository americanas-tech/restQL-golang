package ast

import (
	"fmt"
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer/regex"
	"github.com/pkg/errors"
)

type Generator struct {
	participleParser *participle.Parser
}

func New() (Generator, error) {
	newParser, err := makeParser()
	if err != nil {
		return Generator{}, err
	}

	return Generator{participleParser: newParser}, nil
}

func (g Generator) Parse(query string) (q *Query, err error) {
	q = &Query{}
	if err = g.participleParser.ParseString(query, q); err != nil {
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

	parser, err := participle.Build(
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
