package parser

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/parser/ast"
)

func Parse(queryStr string) (domain.Query, error) {
	queryAst, err := ast.Parse(queryStr)

	err = ValidateQuery(queryAst)
	if err != nil {
		return domain.Query{}, err
	}

	statements := mapToStatements(queryAst.Blocks)
	query := domain.Query{Statements: statements}

	if queryAst.Use != nil {
		query.Use = makeUse(queryAst)
	}

	return query, nil
}

func makeUse(queryAst *ast.Query) map[string]interface{} {
	result := map[string]interface{}{}
	for _, use := range queryAst.Use {
		if use.Value.String != nil {
			result[use.Key] = *use.Value.String
		} else {
			result[use.Key] = *use.Value.Int
		}
	}
	return result
}

func mapToStatements(fromBlocks []ast.Block) []domain.Statement {
	result := make([]domain.Statement, len(fromBlocks))

	for i, block := range fromBlocks {
		result[i] = makeStatement(block)
	}

	return result
}

func makeStatement(block ast.Block) domain.Statement {
	s := domain.Statement{
		Method:   block.Method,
		Resource: block.Resource,
		Alias:    block.Alias,
	}
	for _, qualifier := range block.Qualifiers {
		if qualifier.With != nil {
			s.With = makeParams(qualifier)
		}

		if qualifier.Only != nil {
			s.Only = makeOnlyFilter(qualifier)
		}

		if qualifier.Timeout != nil {
			s.Timeout = makeTimeout(qualifier)
		}

		if qualifier.Headers != nil {
			s.Headers = makeHeaders(qualifier)
		}

		if qualifier.MaxAge != nil {
			value := makeMaxAge(qualifier)
			s.CacheControl.MaxAge = value
		}

		if qualifier.SMaxAge != nil {
			value := makeSMaxAge(qualifier)
			s.CacheControl.SMaxAge = value
		}

		s.Hidden = qualifier.Hidden || s.Hidden
		s.IgnoreErrors = qualifier.IgnoreErrors || s.IgnoreErrors
	}

	return s
}

func makeParams(wq ast.Qualifier) domain.Params {
	p := domain.Params{}
	for _, item := range wq.With {
		v := getValue(item.Value)

		if item.Flatten {
			v = domain.Flatten{v}
		}

		if item.Json {
			v = domain.Json{v}
		}

		if item.Base64 {
			v = domain.Base64{v}
		}

		p[item.Key] = v
	}

	return p
}

func makeOnlyFilter(onlyQualifier ast.Qualifier) []interface{} {
	filters := onlyQualifier.Only

	result := make([]interface{}, len(filters))
	for i, f := range filters {
		if f.Match != "" {
			result[i] = domain.Match{Target: f.Field, Arg: f.Match}
		} else {
			result[i] = f.Field
		}
	}

	return result
}

func makeHeaders(qualifier ast.Qualifier) map[string]interface{} {
	result := map[string]interface{}{}

	for _, header := range qualifier.Headers {
		k := header.Key
		v := header.Value

		if v.String != nil {
			result[k] = *v.String
		}

		if v.Variable != nil {
			result[k] = v.Variable
		}
	}

	return result
}

func makeTimeout(qualifier ast.Qualifier) interface{} {
	v := qualifier.Timeout
	if v.Int != nil {
		return *v.Int
	}

	if v.Variable != nil {
		return *v.Variable
	}

	return nil
}

func makeSMaxAge(qualifier ast.Qualifier) interface{} {
	v := qualifier.SMaxAge
	if v.Int != nil {
		return *v.Int
	}

	if v.Variable != nil {
		return v.Variable
	}

	return nil
}

func makeMaxAge(qualifier ast.Qualifier) interface{} {
	v := qualifier.MaxAge
	if v.Int != nil {
		return *v.Int
	}

	if v.Variable != nil {
		return v.Variable
	}

	return nil
}

func getValue(value ast.Value) interface{} {
	if value.Variable != nil {
		return value.Variable
	}

	if value.Primitive != nil {
		return getPrimitive(value.Primitive)
	}

	if value.List != nil {
		result := make([]interface{}, len(value.List))
		for i, v := range value.List {
			result[i] = getValue(*v)
		}

		return result
	}

	if value.Object != nil {
		return getMap(value.Object)
	}

	return nil
}

func getMap(entries []ast.ObjectEntry) map[string]interface{} {
	result := map[string]interface{}{}

	for _, entry := range entries {
		if entry.Value.Primitive != nil {
			result[entry.Key] = getPrimitive(entry.Value.Primitive)
		}

		if entry.Value.Nested != nil {
			result[entry.Key] = getMap(entry.Value.Nested)
		}
	}

	return result
}

func getPrimitive(primitive *ast.Primitive) interface{} {
	if primitive.Int != nil {
		return *primitive.Int
	}

	if primitive.String != nil {
		return *primitive.String
	}

	if primitive.Float != nil {
		return *primitive.Float
	}

	if primitive.Chain != nil {
		return makeChain(primitive)
	}

	return nil
}

func makeChain(primitive *ast.Primitive) domain.Chain {
	result := make(domain.Chain, len(primitive.Chain))
	for i, chained := range primitive.Chain {
		if chained.PathVariable != "" {
			result[i] = chained.PathVariable
		}

		if chained.PathItem != "" {
			result[i] = chained.PathItem
		}
	}
	return result
}
