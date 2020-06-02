package runner

import (
	"encoding/json"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"strconv"
)

func ResolveVariables(resources domain.Resources, input domain.QueryInput) domain.Resources {
	for key, statement := range resources {
		if statement, ok := statement.(domain.Statement); ok {
			statement.With = resolveWith(statement.With, input)
			statement.Timeout = resolveTimeout(statement.Timeout, input)
			statement.Headers = resolveHeaders(statement.Headers, input)
			statement.CacheControl = resolveCacheControl(statement.CacheControl, input)

			resources[key] = statement
		}
	}

	return resources
}

func resolveWith(with domain.Params, input domain.QueryInput) domain.Params {
	if with.Values == nil && with.Body == nil {
		return with
	}

	body := resolveWithBody(with.Body, input)

	result := make(map[string]interface{})

	for key, value := range with.Values {
		switch value := value.(type) {
		case domain.Variable:
			resolvedValue, ok := getUniqueParamValue(value.Target, input)
			if !ok {
				continue
			}

			result[key] = resolvedValue
		case domain.Chain:
			result[key] = resolveChain(value, input)
		case map[string]interface{}:
			result[key] = resolveComplexWithParam(value, input)
		case []interface{}:
			result[key] = resolveListWithParam(value, input)
		default:
			result[key] = value
		}
	}

	return domain.Params{Body: body, Values: result}
}

func resolveWithBody(body interface{}, input domain.QueryInput) interface{} {
	switch body := body.(type) {
	case domain.Variable:
		p, found := input.Params[body.Target]
		if !found {
			return nil
		}

		b, err := unmarshalValue(p)
		if err != nil {
			return p
		}

		return b
	case domain.Flatten:
		return domain.Flatten{Target: resolveWithBody(body.Target, input)}
	case domain.Json:
		return domain.Json{Target: resolveWithBody(body.Target, input)}
	case domain.Base64:
		return domain.Base64{Target: resolveWithBody(body.Target, input)}
	default:
		return nil
	}
}

func unmarshalValue(value interface{}) (interface{}, error) {
	switch value := value.(type) {
	case string:
		var result interface{}
		err := json.Unmarshal([]byte(value), &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	case []interface{}:
		result := make([]interface{}, len(value))
		for i, v := range value {
			u, err := unmarshalValue(v)
			if err != nil {
				return nil, err
			}

			result[i] = u
		}
		return result, nil
	default:
		return value, nil
	}
}

func resolveListWithParam(list []interface{}, input domain.QueryInput) interface{} {
	l := make([]interface{}, len(list))
	for i, val := range list {
		switch val := val.(type) {
		case domain.Variable:
			value, ok := getUniqueParamValue(val.Target, input)
			if !ok {
				continue
			}

			l[i] = value
		case []interface{}:
			l[i] = resolveListWithParam(val, input)
		default:
			l[i] = val
		}
	}

	return l
}

func resolveComplexWithParam(object map[string]interface{}, input domain.QueryInput) interface{} {
	m := make(map[string]interface{})
	for key, val := range object {
		switch val := val.(type) {
		case domain.Variable:
			value, ok := getUniqueParamValue(val.Target, input)
			if !ok {
				continue
			}

			m[key] = value
		case map[string]interface{}:
			m[key] = resolveComplexWithParam(val, input)
		default:
			m[key] = val
		}
	}

	return m
}

func resolveChain(chain domain.Chain, input domain.QueryInput) domain.Chain {
	result := make(domain.Chain, len(chain))
	for i, pathItem := range chain {
		switch pathItem := pathItem.(type) {
		case domain.Variable:
			paramValue, ok := getUniqueParamValue(pathItem.Target, input)
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
		paramValue, ok := getUniqueParamValue(value.Target, input)
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
		paramValue, ok := getUniqueParamValue(value.Target, input)
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
			paramValue, ok := getUniqueParamValue(value.Target, input)
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
		paramValue, ok := getUniqueParamValue(timeout.Target, input)
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

func getUniqueParamValue(name string, input domain.QueryInput) (string, bool) {
	value, ok := input.Params[name]
	if !ok {
		headerValue, found := input.Headers[name]
		return headerValue, found
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
