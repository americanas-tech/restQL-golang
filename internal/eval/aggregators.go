package eval

import (
	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
)

// ApplyAggregators resolves the `in` keyword in the query,
// taking values from one statement result than setting it
// on target statement result.
func ApplyAggregators(query domain.Query, resources domain.Resources) domain.Resources {
	for _, stmt := range query.Statements {
		if len(stmt.In) == 0 {
			continue
		}

		target := stmt.In[0]
		path := stmt.In[1:]

		originResourceID := domain.NewResourceID(stmt)
		originResource := resources[originResourceID]

		targetResourceID := domain.ResourceID(target)
		targetResource := resources[targetResourceID]

		aggregateOriginOnTarget(path, originResource, targetResource)
		resources[originResourceID] = cleanOriginResult(originResource)
	}

	return resources
}

func aggregateOriginOnTarget(path []string, origin interface{}, target interface{}) {
	switch target := target.(type) {
	case restql.DoneResource:
		body := target.ResponseBody.Unmarshal()
		aggregateOriginOnTarget(path, origin, body)
	case restql.DoneResources:
		aggregateOriginOnListTarget(path, origin, target)
	case []interface{}:
		aggregateOriginOnListTarget(path, origin, target)
	case map[string]interface{}:
		field := path[0]

		nextTarget, found := target[field]
		if !found {
			nextTarget = make(map[string]interface{})
			target[field] = nextTarget
		}

		if len(path) == 1 {
			//setOriginOnTarget(field, origin, target)
			target[field] = parseOrigin(origin)
		} else {
			aggregateOriginOnTarget(path[1:], origin, nextTarget)
		}

	}
}

func aggregateOriginOnListTarget(path []string, origin interface{}, target []interface{}) {
	switch origin := origin.(type) {
	case restql.DoneResource:
		body := origin.ResponseBody.Unmarshal()
		aggregateOriginOnTarget(path, body, target)
	case restql.DoneResources:
		for i, t := range target {
			aggregateOriginOnTarget(path, origin[i], t)
		}
	case []interface{}:
		for i, t := range target {
			aggregateOriginOnTarget(path, origin[i], t)
		}
	default:
		for _, t := range target {
			aggregateOriginOnTarget(path, origin, t)
		}
	}
}

func setOriginOnTarget(field string, origin interface{}, target interface{}) {
	switch target := target.(type) {
	case restql.DoneResource:
		body := target.ResponseBody.Unmarshal()
		setOriginOnTarget(field, origin, body)
	case map[string]interface{}:
		target[field] = parseOrigin(origin)

	}
}

func parseOrigin(origin interface{}) interface{} {
	switch origin := origin.(type) {
	case restql.DoneResource:
		body := origin.ResponseBody.Unmarshal()
		return body
	case restql.DoneResources:
		result := make([]interface{}, len(origin))
		for i, o := range origin {
			result[i] = parseOrigin(o)
		}
		return result
	default:
		return origin
	}
}

func cleanOriginResult(origin interface{}) interface{} {
	switch origin := origin.(type) {
	case restql.DoneResource:
		origin.ResponseBody.SetValue(nil)
		origin.ResponseBody.SetBytes(nil)
		return origin
	case restql.DoneResources:
		result := make(restql.DoneResources, len(origin))
		for i, o := range origin {
			result[i] = cleanOriginResult(o)
		}

		return result
	default:
		return nil
	}
}
