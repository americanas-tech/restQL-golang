package runner

import (
	"math"

	"github.com/b2wdigital/restQL-golang/v6/internal/domain"
)

const (
	valuesParamType string = "values"
	bodyParamType          = "body"
)

type parameter struct {
	paramType string
	path      []string
	value     interface{}
}

type listParameters struct {
	paramType string
	path      []string
	value     interface{}
}

// MultiplexStatements creates a statement for each value in a
// list parameter value.
// In case of multiple list parameter values it makes the cartesian
// product of all the lists and makes a statement for each value in a
// result product.
func MultiplexStatements(resources domain.Resources) domain.Resources {
	for resourceID, stmt := range resources {
		switch stmt := stmt.(type) {
		case domain.Statement:
			resources[resourceID] = multiplex(stmt)
		default:
			resources[resourceID] = stmt
		}
	}

	return resources
}

func multiplex(statement domain.Statement) interface{} {
	values := statement.With.Values
	body := statement.With.Body
	if values == nil && body == nil {
		return statement
	}

	listParams := getListParamsFromValues(values)

	bodyParams := getListParamsFromBody(body)

	listParams = append(listParams, bodyParams...)
	if len(listParams) == 0 {
		return statement
	}

	statementsParameters := zipListParams(listParams)

	result := make([]interface{}, len(statementsParameters))
	for i, parameters := range statementsParameters {
		newStmt := copyStatement(statement)
		for _, p := range parameters {
			if p.paramType == valuesParamType {
				setParameterOnStatement(newStmt.With.Values, p.path, p.value)
			}

			if p.paramType == bodyParamType {
				newStmt.With.Body = p.value
			}
		}

		result[i] = multiplex(newStmt)
	}

	return result
}

func getListParamsFromBody(body interface{}) []listParameters {
	var result []listParameters
	if body, ok := body.([]interface{}); ok {
		result = append(result, listParameters{paramType: bodyParamType, value: body})
	}

	return result
}

func setParameterOnStatement(params interface{}, path []string, value interface{}) {
	if params, ok := params.(map[string]interface{}); ok {
		if len(path) == 1 {
			params[path[0]] = value
			return
		}

		field, ok := params[path[0]]
		if !ok {
			return
		}

		setParameterOnStatement(field, path[1:], value)
	}
}

func getListParamsFromValues(values map[string]interface{}) []listParameters {
	var result []listParameters
	for key, val := range values {
		lp := findListParameters([]string{key}, val)
		result = append(result, lp...)
	}

	return result
}

func findListParameters(path []string, val interface{}) []listParameters {
	switch val := val.(type) {
	case domain.NoMultiplex:
		return []listParameters{}
	case domain.AsQuery:
		parameters := findListParameters(path, val.Target())
		if len(parameters) == 0 {
			return []listParameters{}
		}

		result := make([]listParameters, len(parameters))
		for i, lp := range parameters {
			fn := val.Map(func(_ interface{}) interface{} { return lp.value })
			result[i] = listParameters{path: lp.path, paramType: lp.paramType, value: fn}
		}

		return result
	case map[string]interface{}:
		var result []listParameters
		for k, v := range val {
			lp := findListParameters(append(path, k), v)
			result = append(result, lp...)
		}
		return result
	case []interface{}:
		return []listParameters{{path: path, paramType: valuesParamType, value: val}}
	default:
		return []listParameters{}
	}
}

func zipListParams(listParams []listParameters) [][]parameter {
	statementCount := minimumListParamLength(listParams)
	result := make([][]parameter, statementCount)

	var i uint64
	for i = 0; i < statementCount; i++ {
		statementParameters := make([]parameter, len(listParams))
		for j, lp := range listParams {
			v := lp.value

			switch v := v.(type) {
			case domain.AsQuery:
				fn := v.Map(func(target interface{}) interface{} {
					l := target.([]interface{})
					return l[i]
				})

				statementParameters[j] = parameter{path: lp.path, paramType: lp.paramType, value: fn}
			case []interface{}:
				statementParameters[j] = parameter{path: lp.path, paramType: lp.paramType, value: v[i]}
			}
		}

		result[i] = statementParameters
	}

	return result
}

func minimumListParamLength(listParams []listParameters) uint64 {
	var result uint64 = math.MaxUint64
	for _, lp := range listParams {
		length := valueLength(lp)
		if length < result {
			result = length
		}
	}

	return result
}

func valueLength(lp listParameters) uint64 {
	switch v := lp.value.(type) {
	case []interface{}:
		return uint64(len(v))
	case domain.AsQuery:
		target := v.Target().([]interface{})
		return uint64(len(target))
	default:
		return math.MaxUint64
	}
}

func copyStatement(statement domain.Statement) domain.Statement {
	newStatement := statement

	if statement.Headers != nil {
		newStatement.Headers = make(map[string]interface{}, len(statement.Headers))
		for k, v := range statement.Headers {
			newStatement.Headers[k] = v
		}
	}

	if statement.With.Values != nil {
		newStatement.With.Values = make(map[string]interface{}, len(statement.With.Values))
		for k, v := range statement.With.Values {
			newStatement.With.Values[k] = copyValue(v)
		}
	}

	if statement.Only != nil {
		newStatement.Only = make([]interface{}, len(statement.Only))
		copy(newStatement.Only, statement.Only)
	}

	return newStatement
}

func copyValue(value interface{}) interface{} {
	switch value := value.(type) {
	case map[string]interface{}:
		m := make(map[string]interface{})
		for k, v := range value {
			m[k] = copyValue(v)
		}
		return m
	default:
		return value
	}
}
