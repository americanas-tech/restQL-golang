package runner

import (
	"encoding/json"
	"strconv"

	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
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
		resolvedValue, ok := resolveWithParamValue(value, input)
		if !ok {
			continue
		}

		result[key] = resolvedValue
	}

	return domain.Params{Body: body, Values: result}
}

func resolveWithParamValue(value interface{}, input domain.QueryInput) (interface{}, bool) {
	switch value := value.(type) {
	case domain.Variable:
		return getUniqueParamValue(value.Target, input)
	case domain.Chain:
		return resolveChain(value, input)
	case domain.Function:
		v, ok := resolveWithParamValue(value.Target(), input)
		fnValue := value.Map(func(target interface{}) interface{} { return v })
		return fnValue, ok
	case map[string]interface{}:
		return resolveComplexWithParam(value, input), true
	case []interface{}:
		return resolveListWithParam(value, input), true
	default:
		return value, true
	}
}

func resolveWithBody(body interface{}, input domain.QueryInput) interface{} {
	switch body := body.(type) {
	case domain.Variable:
		p, found := getUniqueParamValue(body.Target, input)
		if !found {
			return nil
		}

		b, err := unmarshalValue(p)
		if err != nil {
			return p
		}

		return b
	case domain.Function:
		return body.Map(func(target interface{}) interface{} {
			return resolveWithBody(target, input)
		})
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
	l := make([]interface{}, 0, len(list))
	for _, val := range list {
		value, ok := resolveWithParamValue(val, input)
		if !ok {
			continue
		}
		l = append(l, value)
	}

	return l
}

func resolveComplexWithParam(object map[string]interface{}, input domain.QueryInput) interface{} {
	m := make(map[string]interface{})
	for key, val := range object {
		value, ok := resolveWithParamValue(val, input)
		if !ok {
			continue
		}
		m[key] = value
	}

	return m
}

func resolveCacheControl(cacheControl domain.CacheControl, input domain.QueryInput) domain.CacheControl {
	var result domain.CacheControl

	switch value := cacheControl.MaxAge.(type) {
	case domain.Variable:
		paramValue, found := getUniqueParamValue(value.Target, input)
		if !found {
			result.MaxAge = nil
		}

		maxAge, ok := castToInt(paramValue)
		if !ok {
			return domain.CacheControl{}
		}

		result.MaxAge = maxAge
	case int:
		result.MaxAge = value
	}

	switch value := cacheControl.SMaxAge.(type) {
	case domain.Variable:
		paramValue, found := getUniqueParamValue(value.Target, input)
		if !found {
			result.SMaxAge = nil
		}

		smaxAge, ok := castToInt(paramValue)
		if !ok {
			return domain.CacheControl{}
		}

		result.SMaxAge = smaxAge
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
		case domain.Chain:
			rc, ok := resolveChain(value, input)
			if !ok {
				continue
			}

			result[key] = rc
		case string:
			result[key] = value
		}
	}

	return result
}

func resolveTimeout(timeout interface{}, input domain.QueryInput) interface{} {
	switch timeout := timeout.(type) {
	case domain.Variable:
		paramValue, found := getUniqueParamValue(timeout.Target, input)
		if !found {
			return nil
		}

		result, ok := castToInt(paramValue)
		if !ok {
			return nil
		}

		return result
	case int:
		return timeout
	default:
		return nil
	}
}

func resolveChain(chain domain.Chain, input domain.QueryInput) (domain.Chain, bool) {
	result := make(domain.Chain, len(chain))
	for i, pathItem := range chain {
		switch pathItem := pathItem.(type) {
		case domain.Variable:
			paramValue, ok := getUniqueParamValue(pathItem.Target, input)
			if !ok {
				return nil, false
			}

			result[i] = paramValue
		default:
			result[i] = pathItem
		}
	}

	return result, true
}

func castToInt(value interface{}) (int, bool) {
	switch value := value.(type) {
	case string:
		result, err := strconv.Atoi(value)
		if err != nil {
			return 0, false
		}
		return result, true
	case int:
		return value, true
	default:
		return 0, false
	}
}

func getUniqueParamValue(name string, input domain.QueryInput) (interface{}, bool) {
	bodyValue, ok := getUniqueParamValueFromBody(name, input.Body)
	if ok {
		return bodyValue, true
	}

	value, ok := input.Params[name]
	if ok {
		return value, true
	}

	headerValue, found := input.Headers[name]
	return headerValue, found
}

func getUniqueParamValueFromBody(name string, body interface{}) (interface{}, bool) {
	b, ok := body.(map[string]interface{})
	if !ok {
		return nil, false
	}

	value, found := b[name]
	return value, found
}
