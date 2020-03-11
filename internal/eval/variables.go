package eval

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"strconv"
)

func ResolveVariables(query domain.Query, input domain.QueryInput) domain.Query {
	result := domain.Query{
		Use:        query.Use,
		Statements: make([]domain.Statement, len(query.Statements)),
	}

	for i, statement := range query.Statements {
		statement.With = resolveWith(statement.With, input)
		statement.Timeout = resolveTimeout(statement.Timeout, input)
		statement.Headers = resolveHeaders(statement.Headers, input)
		statement.CacheControl = resolveCacheControl(statement.CacheControl, input)

		result.Statements[i] = statement
	}

	return result
}

func resolveWith(with domain.Params, input domain.QueryInput) domain.Params {
	if with == nil {
		return nil
	}

	result := make(domain.Params)

	for key, value := range with {
		switch value := value.(type) {
		case domain.Variable:
			paramValue, ok := input.Params[value.Target]
			if !ok {
				continue
			}

			result[key] = paramValue
		case domain.Chain:
			result[key] = resolveChain(value, input)
		default:
			result[key] = value
		}
	}

	return result
}

func resolveChain(chain domain.Chain, input domain.QueryInput) domain.Chain {
	result := make(domain.Chain, len(chain))
	for i, pathItem := range chain {
		switch pathItem := pathItem.(type) {
		case domain.Variable:
			paramValue, ok := getUniqueParamValue(pathItem.Target, input.Params)
			if !ok {
				continue
			}

			result[i] = paramValue
		default:
			result[i] = pathItem
		}
	}

	return result
}

func resolveCacheControl(cacheControl domain.CacheControl, input domain.QueryInput) domain.CacheControl {
	var result domain.CacheControl

	switch value := cacheControl.MaxAge.(type) {
	case domain.Variable:
		paramValue, ok := getUniqueParamValue(value.Target, input.Params)
		if !ok {
			result.MaxAge = nil
		}

		maxAge, err := strconv.Atoi(paramValue)
		if err != nil {
			result.MaxAge = nil
		}

		result.MaxAge = maxAge
	case int:
		result.MaxAge = value
	}

	switch value := cacheControl.SMaxAge.(type) {
	case domain.Variable:
		paramValue, ok := getUniqueParamValue(value.Target, input.Params)
		if !ok {
			result.SMaxAge = nil
		}

		maxAge, err := strconv.Atoi(paramValue)
		if err != nil {
			result.SMaxAge = nil
		}

		result.SMaxAge = maxAge
	case int:
		result.SMaxAge = value
	}

	return result
}

func resolveHeaders(headers map[string]interface{}, input domain.QueryInput) map[string]interface{} {
	if headers == nil {
		return nil
	}

	result := make(map[string]interface{})

	for key, value := range headers {
		switch value := value.(type) {
		case domain.Variable:
			paramValue, ok := getUniqueParamValue(value.Target, input.Params)
			if !ok {
				continue
			}

			result[key] = paramValue
		case string:
			result[key] = value
		}
	}

	return result
}

func resolveTimeout(timeout interface{}, input domain.QueryInput) interface{} {
	switch timeout := timeout.(type) {
	case domain.Variable:
		paramValue, ok := getUniqueParamValue(timeout.Target, input.Params)
		if !ok {
			return nil
		}

		result, err := strconv.Atoi(paramValue)
		if err != nil {
			return nil
		}

		return result
	case int:
		return timeout
	default:
		return nil
	}
}

func getUniqueParamValue(name string, params map[string]interface{}) (string, bool) {
	value, ok := params[name]
	if !ok {
		return "", false
	}

	switch value := value.(type) {
	case []string:
		return value[0], true
	case string:
		return value, true
	default:
		return "", false
	}
}
