package parser

import (
	"regexp"
	"strings"

	"github.com/b2wdigital/restQL-golang/v6/internal/domain"
	"github.com/b2wdigital/restQL-golang/v6/internal/parser/ast"
	"github.com/pkg/errors"
)

// Optimize transforms a restQL AST into the internal representation.
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

		if qualifier.DependsOn != "" {
			s.DependsOn = domain.DependsOn{Target: qualifier.DependsOn}
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
		case ast.JSON:
			v = domain.JSON{Value: v}
		case ast.Flatten:
			v = domain.Flatten{Value: v}
		case ast.NoExplode:
			v = domain.NoExplode{Value: v}
		case ast.AsQuery:
			v = domain.AsQuery{Value: v}
		case ast.NoDuplicate:
			v = domain.NoDuplicate{Value: v}
		}
	}

	return v
}

func makeOnlyFilter(onlyQualifier ast.Qualifier) ([]interface{}, error) {
	filters := onlyQualifier.Only

	result := make([]interface{}, len(filters))
	for i, f := range filters {
		var filter interface{} = f.Field
		for _, fn := range f.Functions {
			filterWithFunc, err := applyFunctionToFilter(filter, fn)
			if err != nil {
				return nil, err
			}

			filter = filterWithFunc
		}

		result[i] = filter
	}

	return result, nil
}

func applyFunctionToFilter(field, fn interface{}) (interface{}, error) {
	switch fn := fn.(type) {
	case ast.Match:
		return makeMatchFunction(field, fn)
	case ast.FilterByRegex:
		return makeFilterByRegexFunction(field, fn)
	default:
		return field, nil
	}
}

func makeFilterByRegexFunction(target interface{}, filterByRegexFn ast.FilterByRegex) (domain.Function, error) {
	var fr domain.Function = domain.FilterByRegex{Value: target}

	switch {
	case filterByRegexFn.PathVariable != nil:
		path := domain.Variable{Target: *filterByRegexFn.PathVariable}
		fr = fr.SetArgument(domain.FilterByRegexArgPath, path)
	case filterByRegexFn.PathString != nil:
		fr = fr.SetArgument(domain.FilterByRegexArgPath, *filterByRegexFn.PathString)
	}

	switch {
	case filterByRegexFn.RegexVariable != nil:
		regexVariable := domain.Variable{Target: *filterByRegexFn.RegexVariable}
		fr = fr.SetArgument(domain.FilterByRegexArgRegex, regexVariable)
	case filterByRegexFn.RegexString != nil:
		reg, err := regexp.Compile(*filterByRegexFn.RegexString)
		if err != nil {
			return domain.FilterByRegex{}, err
		}

		fr = fr.SetArgument(domain.FilterByRegexArgRegex, reg)
	}

	return fr, nil
}

func makeMatchFunction(target interface{}, matchFn ast.Match) (domain.Function, error) {
	var match domain.Function = domain.Match{Value: target}

	if matchFn.String != nil {
		arg := *matchFn.String
		regex, err := regexp.Compile(arg)
		if err != nil {
			return domain.Match{}, errors.Wrap(err, "matches function regex argument is invalid")
		}

		match = match.SetArgument(domain.MatchArgRegex, regex)

		return match, nil
	}

	if matchFn.Variable != nil {
		match = match.SetArgument(domain.MatchArgRegex, domain.Variable{Target: *matchFn.Variable})

		return match, nil
	}

	return domain.Match{}, errors.New("no argument provided to matches functions")
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
