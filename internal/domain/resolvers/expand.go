package resolvers

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"math"
)

type parameter struct {
	key   string
	value interface{}
}

type listParameters struct {
	key   string
	value []interface{}
}

func MultiplexStatements(query domain.Query) domain.Query {
	var multiplexed []domain.Statement
	for _, stmt := range query.Statements {
		s := multiplex(stmt)
		multiplexed = append(multiplexed, s...)
	}

	q := query
	q.Statements = multiplexed

	return q
}

func multiplex(statement domain.Statement) []domain.Statement {
	params := statement.With
	if params == nil {
		return []domain.Statement{statement}
	}

	listParams := getListParams(params)
	if len(listParams) == 0 {
		return []domain.Statement{statement}
	}

	statementsParameters := zipListParams(listParams)

	result := make([]domain.Statement, len(statementsParameters))
	for i, parameters := range statementsParameters {
		newStmt := copyStatement(statement)
		for _, p := range parameters {
			newStmt.With[p.key] = p.value
		}

		result[i] = newStmt
	}

	return result
}

func getListParams(params map[string]interface{}) []listParameters {
	var result []listParameters
	for key, val := range params {
		switch val := val.(type) {
		case domain.Flatten:
			continue
		case []interface{}:
			result = append(result, listParameters{key: key, value: val})
		}
	}

	return result
}

func zipListParams(listParams []listParameters) [][]parameter {
	statementCount := minimumListParamLength(listParams)
	result := make([][]parameter, statementCount)

	var i uint64
	for i = 0; i < statementCount; i++ {
		statementParameters := make([]parameter, len(listParams))
		for j, lp := range listParams {
			statementParameters[j] = parameter{key: lp.key, value: lp.value[i]}
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
			newStatement.With[k] = v
		}
	}

	if statement.Only != nil {
		newStatement.Only = make([]interface{}, len(statement.Only))
		copy(newStatement.Only, statement.Only)
	}

	return newStatement
}
