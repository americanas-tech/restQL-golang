package runner

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"math"
)

type parameter struct {
	path  []string
	value interface{}
}

type listParameters struct {
	path  []string
	value []interface{}
}

func MultiplexStatements(resources domain.Resources) domain.Resources {
	for resourceId, stmt := range resources {
		switch stmt := stmt.(type) {
		case domain.Statement:
			resources[resourceId] = multiplex(stmt)
		default:
			resources[resourceId] = stmt
		}
	}

	return resources
}

func multiplex(statement domain.Statement) interface{} {
	params := statement.With
	if params == nil {
		return statement
	}

	listParams := getListParams(params)
	if len(listParams) == 0 {
		return statement
	}

	statementsParameters := zipListParams(listParams)

	result := make([]interface{}, len(statementsParameters))
	for i, parameters := range statementsParameters {
		newStmt := copyStatement(statement)
		for _, p := range parameters {
			setParameterOnStatement(newStmt.With, p.path, p.value)
		}

		result[i] = multiplex(newStmt)
	}

	return result
}

func setParameterOnStatement(params interface{}, path []string, value interface{}) {
	switch params := params.(type) {
	case domain.Params:
		if len(path) == 1 {
			params[path[0]] = value
			return
		}

		field, ok := params[path[0]]
		if !ok {
			return
		}

		setParameterOnStatement(field, path[1:], value)
	case map[string]interface{}:
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

func getListParams(params map[string]interface{}) []listParameters {
	var result []listParameters
	for key, val := range params {
		lp := findListParameters([]string{key}, val)
		result = append(result, lp...)
	}

	return result
}

func findListParameters(path []string, val interface{}) []listParameters {
	switch val := val.(type) {
	case domain.Flatten:
		return []listParameters{}
	case map[string]interface{}:
		var result []listParameters
		for k, v := range val {
			lp := findListParameters(append(path, k), v)
			result = append(result, lp...)
		}
		return result
	case []interface{}:
		return []listParameters{{path: path, value: val}}
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
			statementParameters[j] = parameter{path: lp.path, value: lp.value[i]}
		}

		result[i] = statementParameters
	}

	return result
}

func minimumListParamLength(listParams []listParameters) uint64 {
	var result uint64 = math.MaxUint64
	for _, lp := range listParams {
		length := uint64(len(lp.value))
		if length < result {
			result = length
		}
	}

	return result
}

func copyStatement(statement domain.Statement) domain.Statement {
	newStatement := statement

	if statement.Headers != nil {
		newStatement.Headers = make(map[string]interface{}, len(statement.Headers))
		for k, v := range statement.Headers {
			newStatement.Headers[k] = v
		}
	}

	if statement.With != nil {
		newStatement.With = make(map[string]interface{}, len(statement.With))
		for k, v := range statement.With {
			newStatement.With[k] = copyValue(v)
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
