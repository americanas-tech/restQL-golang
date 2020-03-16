package runner

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"strconv"
)

func ResolveVariables(query domain.Query, params map[string]interface{}) domain.Query {
	result := domain.Query{
		Use:        query.Use,
		Statements: make([]domain.Statement, len(query.Statements)),
	}

	for i, statement := range query.Statements {
		statement.With = resolveWith(statement.With, params)
		statement.Timeout = resolveTimeout(statement.Timeout, params)
		statement.Headers = resolveHeaders(statement.Headers, params)
		statement.CacheControl = resolveCacheControl(statement.CacheControl, params)

		result.Statements[i] = statement
	}

	return result
}

func resolveWith(with domain.Params, params map[string]interface{}) domain.Params {
	if with == nil {
		return nil
	}

	result := make(domain.Params)

	for key, value := range with {
		switch value := value.(type) {
		case domain.Variable:
			paramValue, ok := params[value.Target]
			if !ok {
				continue
			}

			result[key] = paramValue
		case domain.Chain:
			result[key] = resolveChain(value, params)
		case map[string]interface{}:
			result[key] = resolveComplexWithParam(value, params)
		case []interface{}:
			result[key] = resolveListWithParam(value, params)
		default:
			result[key] = value
		}
	}

	return result
}

func resolveListWithParam(list []interface{}, parameters map[string]interface{}) interface{} {
	l := make([]interface{}, len(list))
	for i, val := range list {
		switch val := val.(type) {
		case domain.Variable:
			paramValue, ok := parameters[val.Target]
			if !ok {
				continue
			}

			l[i] = paramValue
		default:
			l[i] = val
		}
	}

	return l
}

func resolveComplexWithParam(object map[string]interface{}, parameters map[string]interface{}) interface{} {
	m := make(map[string]interface{})
	for key, val := range object {
		switch val := val.(type) {
		case domain.Variable:
			paramValue, ok := parameters[val.Target]
			if !ok {
				continue
			}

			m[key] = paramValue
		default:
			m[key] = val
		}
	}

	return m
}

func resolveChain(chain domain.Chain, params map[string]interface{}) domain.Chain {
	result := make(domain.Chain, len(chain))
	for i, pathItem := range chain {
		switch pathItem := pathItem.(type) {
		case domain.Variable:
			paramValue, ok := getUniqueParamValue(pathItem.Target, params)
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

func resolveCacheControl(cacheControl domain.CacheControl, params map[string]interface{}) domain.CacheControl {
	var result domain.CacheControl

	switch value := cacheControl.MaxAge.(type) {
	case domain.Variable:
		paramValue, ok := getUniqueParamValue(value.Target, params)
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
		paramValue, ok := getUniqueParamValue(value.Target, params)
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

func resolveHeaders(headers map[string]interface{}, params map[string]interface{}) map[string]interface{} {
	if headers == nil {
		return nil
	}

	result := make(map[string]interface{})

	for key, value := range headers {
		switch value := value.(type) {
		case domain.Variable:
			paramValue, ok := getUniqueParamValue(value.Target, params)
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

func resolveTimeout(timeout interface{}, params map[string]interface{}) interface{} {
	switch timeout := timeout.(type) {
	case domain.Variable:
		paramValue, ok := getUniqueParamValue(timeout.Target, params)
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
