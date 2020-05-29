package eval

import "github.com/b2wdigital/restQL-golang/internal/domain"

func ApplyAggregators(query domain.Query, resources domain.Resources) domain.Resources {
	for _, stmt := range query.Statements {
		if len(stmt.In) == 0 {
			continue
		}

		target := stmt.In[0]
		path := stmt.In[1:]

		originResourceId := domain.NewResourceId(stmt)
		originResource := resources[originResourceId]

		targetResourceId := domain.ResourceId(target)
		targetResource := resources[targetResourceId]

		aggregateOriginOnTarget(path, originResource, targetResource)
		resources[originResourceId] = cleanOriginResult(originResource)
	}

	return resources
}

func aggregateOriginOnTarget(path []string, origin interface{}, target interface{}) {
	switch target := target.(type) {
	case domain.DoneResource:
		aggregateOriginOnTarget(path, origin, target.ResponseBody)
	case domain.DoneResources:
		aggregateOriginOnListTarget(path, origin, target)
	case []interface{}:
		aggregateOriginOnListTarget(path, origin, target)
	case map[string]interface{}:
		field := path[0]
		if len(path) == 1 {
			setOriginOnTarget(field, origin, target)
		} else {
			aggregateOriginOnTarget(path[1:], origin, target[field])
		}

	}
}

func aggregateOriginOnListTarget(path []string, origin interface{}, target []interface{}) {
	switch origin := origin.(type) {
	case domain.DoneResource:
		aggregateOriginOnTarget(path, origin.ResponseBody, target)
	case domain.DoneResources:
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
	case domain.DoneResource:
		setOriginOnTarget(field, origin, target.ResponseBody)
	case map[string]interface{}:
		target[field] = parseOrigin(origin)
	}
}

func parseOrigin(origin interface{}) interface{} {
	switch origin := origin.(type) {
	case domain.DoneResource:
		return origin.ResponseBody
	case domain.DoneResources:
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
	case domain.DoneResource:
		origin.ResponseBody = nil
		return origin
	case domain.DoneResources:
		result := make(domain.DoneResources, len(origin))
		for i, o := range origin {
			result[i] = cleanOriginResult(o)
		}

		return result
	default:
		return nil
	}
}
