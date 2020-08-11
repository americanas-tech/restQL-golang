package parser

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/parser/ast"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

func Optimize(queryAst *ast.Query) (domain.Query, error) {
	statements, err := mapToStatements(queryAst.Blocks)
	if err != nil {
		return domain.Query{}, err
	}

	query := domain.Query{Statements: statements}

	if queryAst.Use != nil {
		query.Use = makeUse(queryAst)
	}

	return query, nil
}

func makeUse(queryAst *ast.Query) map[string]interface{} {
	result := map[string]interface{}{}
	for _, use := range queryAst.Use {
		key := strings.Trim(use.Key, " ")
		if use.Value.String != nil {
			result[key] = *use.Value.String
		} else {
			result[key] = *use.Value.Int
		}
	}
	return result
}

func mapToStatements(fromBlocks []ast.Block) ([]domain.Statement, error) {
	result := make([]domain.Statement, len(fromBlocks))

	for i, block := range fromBlocks {
		statement, err := makeStatement(block)
		if err != nil {
			return nil, nil
		}

		result[i] = statement
	}

	return result, nil
}

func makeStatement(block ast.Block) (domain.Statement, error) {
	s := domain.Statement{
		Method:   strings.TrimSpace(block.Method),
		Resource: block.Resource,
		Alias:    block.Alias,
		In:       block.In,
	}
	for _, qualifier := range block.Qualifiers {
		if qualifier.With != nil {
			s.With = makeParams(qualifier)
		}

		if qualifier.Only != nil {
			filter, err := makeOnlyFilter(qualifier)
			if err != nil {
				return domain.Statement{}, err
			}

			s.Only = filter
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

	return s, nil
}

func makeParams(wq ast.Qualifier) domain.Params {
	values := make(map[string]interface{})
	for _, item := range wq.With.KeyValues {
		v := getValue(item.Value)

		v = applyFunctions(v, item.Functions)

		values[item.Key] = v
	}

	p := domain.Params{
		Values: values,
	}

	parameterBody := wq.With.Body
	if parameterBody == nil {
		return p
	}

	var body interface{}
	body = domain.Variable{Target: parameterBody.Target}

	body = applyFunctions(body, parameterBody.Functions)

	p.Body = body

	return p
}

func applyFunctions(v interface{}, functions []string) interface{} {
	for _, fn := range functions {
		switch fn {
		case ast.NoMultiplex:
			v = domain.NoMultiplex{Value: v}
		case ast.AsBody:
			v = domain.AsBody{Value: v}
		case ast.Base64:
			v = domain.Base64{Value: v}
		case ast.Json:
			v = domain.Json{Value: v}
		case ast.Flatten:
			v = domain.Flatten{Value: v}
		}
	}

	return v
}

func makeOnlyFilter(onlyQualifier ast.Qualifier) ([]interface{}, error) {
	filters := onlyQualifier.Only

	result := make([]interface{}, len(filters))
	for i, f := range filters {
		if f.Match != "" {
			regex, err := regexp.Compile(f.Match)
			if err != nil {
				return nil, errors.Wrap(err, "matches function regex argument is invalid")
			}
			result[i] = domain.Match{Value: f.Field, Arg: regex}
		} else {
			result[i] = f.Field
		}
	}

	return result, nil
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
			result[k] = domain.Variable{Target: *v.Variable}
		}

		if v.Chain != nil {
			result[k] = makeChain(v.Chain)
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
		return domain.Variable{Target: *v.Variable}
	}

	return nil
}

func makeSMaxAge(qualifier ast.Qualifier) interface{} {
	v := qualifier.SMaxAge
	if v.Int != nil {
		return *v.Int
	}

	if v.Variable != nil {
		return domain.Variable{Target: *v.Variable}
	}

	return nil
}

func makeMaxAge(qualifier ast.Qualifier) interface{} {
	v := qualifier.MaxAge
	if v.Int != nil {
		return *v.Int
	}

	if v.Variable != nil {
		return domain.Variable{Target: *v.Variable}
	}

	return nil
}

func getValue(value ast.Value) interface{} {
	if value.Variable != nil {
		return domain.Variable{Target: *value.Variable}
	}

	if value.Primitive != nil {
		return getPrimitive(value.Primitive)
	}

	if value.List != nil {
		result := make([]interface{}, len(value.List))
		for i, v := range value.List {
			result[i] = getValue(v)
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
		result[entry.Key] = getValue(entry.Value)
	}

	return result
}

func getPrimitive(primitive *ast.Primitive) interface{} {
	if primitive.Null {
		return nil
	}

	if primitive.Int != nil {
		return *primitive.Int
	}

	if primitive.Boolean != nil {
		return *primitive.Boolean
	}

	if primitive.String != nil {
		return *primitive.String
	}

	if primitive.Float != nil {
		return *primitive.Float
	}

	if primitive.Chain != nil {
		return makeChain(primitive.Chain)
	}

	return nil
}

func makeChain(chainedValue []ast.Chained) domain.Chain {
	result := make(domain.Chain, len(chainedValue))
	for i, chained := range chainedValue {
		if chained.PathVariable != "" {
			result[i] = domain.Variable{Target: chained.PathVariable}
		}

		if chained.PathItem != "" {
			result[i] = chained.PathItem
		}
	}
	return result
}
