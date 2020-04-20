package runner

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
)

func UnwrapFlatten(resources domain.Resources) domain.Resources {
	for resourceId, resource := range resources {
		resources[resourceId] = unwrapResource(resource)
	}

	return resources
}

func unwrapResource(resource interface{}) interface{} {
	switch resource := resource.(type) {
	case domain.Statement:
		return unwrapStatement(resource)
	case []interface{}:
		multiplexedResource := make([]interface{}, len(resource))
		for i, r := range resource {
			multiplexedResource[i] = unwrapResource(r)
		}
		return multiplexedResource
	default:
		return resource
	}
}

func unwrapStatement(statement domain.Statement) domain.Statement {
	params := statement.With
	for key, value := range params {
		result := unwrapValue(value)

		params[key] = result
	}

	statement.With = params

	return statement
}

func unwrapValue(value interface{}) interface{} {
	switch value := value.(type) {
	case domain.Flatten:
		return value.Target
	case map[string]interface{}:
		m := make(map[string]interface{})
		for k, v := range value {
			m[k] = unwrapValue(v)
		}
		return m
	case []interface{}:
		l := make([]interface{}, len(value))
		for i, v := range value {
			l[i] = unwrapValue(v)
		}
		return l
	default:
		return value
	}
}
