package runner

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
)

func UnwrapNoMultiplex(resources domain.Resources) domain.Resources {
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
	params := statement.With.Values
	for key, value := range params {
		result := unwrapValue(value)

		params[key] = result
	}

	body := unwrapBody(statement.With.Body)

	statement.With.Body = body
	statement.With.Values = params

	return statement
}

func unwrapBody(body interface{}) interface{} {
	switch body := body.(type) {
	case domain.NoMultiplex:
		return body.Target()
	default:
		return body
	}
}

func unwrapValue(value interface{}) interface{} {
	switch value := value.(type) {
	case domain.NoMultiplex:
		return value.Target()
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
