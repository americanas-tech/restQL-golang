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
			switch param := value.(type) {
			case domain.Chain:
				params[paramName] = resolveChainParam(param, doneResources)
			case []interface{}:
				params[paramName] = resolveListParam(param, doneResources)
			case map[string]interface{}:
				params[paramName] = resolveObjectParam(param, doneResources)
			}
		}
	}

	return stmt
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
	if done.Details.Status > 199 && done.Details.Status < 400 {
		return getValue(path, done.Result)
	}
	return EmptyChained
}

func toPath(chain domain.Chain) []string {
	r := make([]string, len(chain))
	for i, c := range chain {
		r[i] = c.(string)
	}
	return r
}

func resolveObjectParam(objectParam map[string]interface{}, doneResources domain.Resources) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range objectParam {
		switch param := value.(type) {
		case domain.Chain:
			result[key] = resolveChainParam(param, doneResources)
		default:
			result[key] = param
		}
	}

	return result
}

func resolveListParam(listParam []interface{}, doneResources domain.Resources) []interface{} {
	result := make([]interface{}, len(listParam))
	copy(result, listParam)

	for i, value := range result {
		switch param := value.(type) {
		case domain.Chain:
			result[i] = resolveChainParam(param, doneResources)
		}
	}

	return result
}

func getValue(pathToValue []string, b domain.Body) interface{} {
	if len(pathToValue) == 0 {
		return b
	}

	switch body := b.(type) {
	case map[string]interface{}:
		return getValue(pathToValue[1:], body[pathToValue[0]])
	case []interface{}:
		result := make([]interface{}, len(body))
		for i, v := range body {
			result[i] = getValue(pathToValue, v)
		}
		return result
	default:
		return body
	}
}
