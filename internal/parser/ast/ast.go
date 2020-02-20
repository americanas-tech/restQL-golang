package ast

import (
	"fmt"
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer/regex"
	"github.com/pkg/errors"
)

func Parse(query string) (*Query, error) {
	definition, err := regex.New(lexerDefinition)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to compile lexer definition"))
	}

	QUERY := &Query{}
	parser := participle.MustBuild(
		QUERY,
		participle.Lexer(definition),
		participle.Unquote("String"),
		participle.Elide("Comment"),
	)

	if err := parser.ParseString(query, QUERY); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to parse query: %s", query))
	}

	return QUERY, nil
}
