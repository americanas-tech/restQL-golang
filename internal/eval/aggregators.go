package eval

import (
	"fmt"
	"github.com/b2wdigital/restQL-golang/v5/internal/domain"
	"github.com/b2wdigital/restQL-golang/v5/pkg/restql"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
)

// ApplyAggregators resolves the `in` keyword in the query,
// taking values from one statement result than setting it
// on target statement result.
func ApplyAggregators(log restql.Logger, query domain.Query, resources domain.Resources) domain.Resources {
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

		err := aggregateOriginOnTarget(path, originResource, targetResource)
		if err != nil {
			log.Error("an error occurred when aggregating the resources", err)
			continue
		}
		resources[originResourceID] = cleanOriginResult(originResource)
	}

	return resources
}

func aggregateOriginOnTarget(path []string, origin interface{}, target interface{}) error {
	switch target := target.(type) {
	case restql.DoneResource:
		body := target.ResponseBody.Unmarshal()
		return aggregateOriginOnTarget(path, origin, body)
	case restql.DoneResources:
		return aggregateOriginOnListTarget(path, origin, target)
	case []interface{}:
		return aggregateOriginOnListTarget(path, origin, target)
	case map[string]interface{}:
		field := path[0]

		nextTarget, targetFieldExist := target[field]
		if !targetFieldExist {
			nextTarget = make(map[string]interface{})
			target[field] = nextTarget
		}

		if len(path) > 1 {
			return aggregateOriginOnTarget(path[1:], origin, nextTarget)
		}

		originValue := parseOrigin(origin)
		if !targetFieldExist {
			target[field] = originValue
			return nil
		}

		targetValue := target[field]
		merged, err := merge(targetValue, originValue)
		if err != nil {
			return err
		}
		target[field] = merged

		return nil
	default:
		return errors.Errorf("unknown target type: %T", target)
	}
}

func merge(target interface{}, origin interface{}) (interface{}, error) {
	switch target := target.(type) {
	case map[string]interface{}:
		return mergeInObject(target, origin)
	case []interface{}:
		return mergeInList(target, origin)
	default:
		return target, errors.Errorf("invalid target type for merge: %T", target)
	}
}

func mergeInObject(target map[string]interface{}, origin interface{}) (interface{}, error) {
	switch origin := origin.(type) {
	case map[string]interface{}:
		err := mergo.Merge(&target, origin, mergo.WithOverride, mergo.WithSliceDeepCopy)
		return target, err
	default:
		return target, errors.Errorf("invalid origin type for merge into object: %T", origin)
	}
}

func mergeInList(target []interface{}, origin interface{}) (interface{}, error) {
	switch origin := origin.(type) {
	case []interface{}:
		size := max(len(target), len(origin))
		l := make([]interface{}, size)

		for i := 0; i < size; i++ {
			if i > len(target)-1 {
				l[i] = origin[i]
				continue
			}

			if i > len(origin)-1 {
				l[i] = target[i]
				continue
			}

			merged, err := merge(target[i], origin[i])
			if err != nil {
				return target, err
			}

			l[i] = merged
		}

		return l, nil
	default:
		return target, errors.Errorf("invalid origin type for merge into list: %T", origin)
	}
}

func max(a int, b int) int {
	if a > b {
		return a
	}

	return b
}

func aggregateOriginOnListTarget(path []string, origin interface{}, target []interface{}) error {
	switch origin := origin.(type) {
	case restql.DoneResource:
		body := origin.ResponseBody.Unmarshal()
		return aggregateOriginOnTarget(path, body, target)
	case restql.DoneResources:
		var err error
		for i, t := range target {
			aggregateErr := aggregateOriginOnTarget(path, origin[i], t)
			if aggregateErr != nil {
				err = fmt.Errorf("failed to aggregate multiplexed resource into target: %v\n%w", aggregateErr, err)
			}
		}
		return err
	case []interface{}:
		var err error
		for i, t := range target {
			aggregateErr := aggregateOriginOnTarget(path, origin[i], t)
			if aggregateErr != nil {
				err = fmt.Errorf("failed to aggregate multiplexed resource into target: %v\n%w", aggregateErr, err)
			}
		}
		return err
	default:
		var err error
		for _, t := range target {
			aggregateErr := aggregateOriginOnTarget(path, origin, t)
			if aggregateErr != nil {
				err = fmt.Errorf("failed to aggregate multiplexed resource into target: %v\n%w", aggregateErr, err)
			}
		}
		return err
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
		origin.ResponseBody.Clear()
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
