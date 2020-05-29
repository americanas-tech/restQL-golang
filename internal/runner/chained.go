package runner

import "github.com/b2wdigital/restQL-golang/internal/domain"

const (
	EmptyChained = "__EMPTY_CHAINED__"
)

func ResolveChainedValues(resources domain.Resources, doneResources domain.Resources) domain.Resources {
	for resourceId, stmt := range resources {
		resources[resourceId] = resolveStatement(stmt, doneResources)
	}

	return resources
}

func resolveStatement(stmt interface{}, doneResources domain.Resources) interface{} {
	switch stmt := stmt.(type) {
	case domain.Statement:
		params := stmt.With
		for paramName, value := range params {
			params[paramName] = resolveParam(value, doneResources)
		}
	}

	return stmt
}

func resolveParam(value interface{}, doneResources domain.Resources) interface{} {
	switch param := value.(type) {
	case domain.Chain:
		return resolveChainParam(param, doneResources)
	case domain.Flatten:
		return domain.Flatten{Target: resolveParam(param.Target, doneResources)}
	case domain.Json:
		return domain.Json{Target: resolveParam(param.Target, doneResources)}
	case domain.Base64:
		return domain.Base64{Target: resolveParam(param.Target, doneResources)}
	case []interface{}:
		return resolveListParam(param, doneResources)
	case map[string]interface{}:
		return resolveObjectParam(param, doneResources)
	default:
		return value
	}
}

func resolveObjectParam(objectParam map[string]interface{}, doneResources domain.Resources) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range objectParam {
		result[key] = resolveParam(value, doneResources)
	}

	return result
}

func resolveListParam(listParam []interface{}, doneResources domain.Resources) []interface{} {
	result := make([]interface{}, len(listParam))
	copy(result, listParam)

	for i, value := range result {
		result[i] = resolveParam(value, doneResources)
	}

	return result
}

func resolveChainParam(chain domain.Chain, doneResources domain.Resources) interface{} {
	path := toPath(chain)
	resourceId := domain.ResourceId(path[0])

	switch done := doneResources[resourceId].(type) {
	case domain.DoneResources:
		return resolveWithMultiplexedRequests(path[1:], done)
	case domain.DoneResource:
		return resolveWithSingleRequest(path[1:], done)
	default:
		return nil
	}
}

func resolveWithMultiplexedRequests(path []string, doneRequests domain.DoneResources) []interface{} {
	var result []interface{}

	for _, request := range doneRequests {
		switch request := request.(type) {
		case domain.DoneResource:
			result = append(result, resolveWithSingleRequest(path, request))
		case domain.DoneResources:
			result = append(result, resolveWithMultiplexedRequests(path, request))
		}
	}

	return result
}

func resolveWithSingleRequest(path []string, done domain.DoneResource) interface{} {
	if done.Status < 200 || done.Status >= 400 {
		return EmptyChained
	}

	valueFromBody, found := getValueFromBody(path, done.ResponseBody)
	if found {
		return valueFromBody
	}

	valueFromHeader, found := done.ResponseHeaders[path[0]]
	if found {
		return valueFromHeader
	}

	return nil
}

func toPath(chain domain.Chain) []string {
	r := make([]string, len(chain))
	for i, c := range chain {
		r[i] = c.(string)
	}
	return r
}

func getValueFromBody(pathToValue []string, b domain.Body) (interface{}, bool) {
	if len(pathToValue) == 0 {
		return b, true
	}

	switch body := b.(type) {
	case map[string]interface{}:
		v, found := body[pathToValue[0]]
		if !found {
			return nil, false
		}

		return getValueFromBody(pathToValue[1:], v)
	case []interface{}:
		result := make([]interface{}, len(body))
		for i, v := range body {
			result[i], _ = getValueFromBody(pathToValue, v)
		}
		return result, true
	default:
		return body, true
	}
}
