package runner

import (
	"encoding/json"
	"fmt"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/pkg/errors"
	"strings"
)

const (
	EmptyChained = "__EMPTY_CHAINED__"
)

var ErrInvalidChainedParameter = errors.New("chained parameter targeting unknown statement")

func ResolveChainedValues(resources domain.Resources, doneResources domain.Resources) domain.Resources {
	for resourceId, stmt := range resources {
		resources[resourceId] = resolveStatement(stmt, doneResources)
	}

	return resources
}

func resolveStatement(stmt interface{}, doneResources domain.Resources) interface{} {
	switch stmt := stmt.(type) {
	case domain.Statement:
		params := stmt.With.Values
		for paramName, value := range params {
			params[paramName] = resolveValue(value, doneResources)
		}

		headers := stmt.Headers
		for name, value := range headers {
			resolved := resolveValue(value, doneResources)
			headerValue, err := stringify(resolved)
			if err != nil {
				headers[name] = EmptyChained
				continue
			}
			headers[name] = headerValue
		}

	case []interface{}:
		result := make([]interface{}, len(stmt))
		for i, s := range stmt {
			result[i] = resolveStatement(s, doneResources)
		}
		return result
	}

	return stmt
}

func stringify(value interface{}) (string, error) {
	switch value := value.(type) {
	case string:
		return value, nil
	case map[string]interface{}:
		b, err := json.Marshal(value)
		return string(b), err
	case []interface{}:
		b, err := json.Marshal(value)
		return string(b), err
	default:
		return fmt.Sprintf("%v", value), nil
	}
}

func resolveValue(value interface{}, doneResources domain.Resources) interface{} {
	switch param := value.(type) {
	case domain.Chain:
		return resolveChainParam(param, doneResources)
	case domain.Function:
		return param.Map(func(target interface{}) interface{} {
			return resolveValue(target, doneResources)
		})
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
		result[key] = resolveValue(value, doneResources)
	}

	return result
}

func resolveListParam(listParam []interface{}, doneResources domain.Resources) []interface{} {
	result := make([]interface{}, len(listParam))
	copy(result, listParam)

	for i, value := range result {
		result[i] = resolveValue(value, doneResources)
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

	valueFromHeader, found := getValueFromHeader(path[0], done.ResponseHeaders)
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
	if b == nil {
		return nil, false
	}

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

func getValueFromHeader(name string, headers map[string]string) (string, bool) {
	name = strings.ToLower(name)
	for k, v := range headers {
		if strings.ToLower(k) == name {
			return v, true
		}
	}

	return "", false
}

func ValidateChainedValues(resources domain.Resources) error {
	for _, stmt := range resources {
		err := validateStatement(stmt, resources)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateStatement(stmt interface{}, resources domain.Resources) error {
	switch stmt := stmt.(type) {
	case domain.Statement:
		params := stmt.With.Values
		for _, value := range params {
			err := validateParam(value, resources)
			if err != nil {
				return err
			}
		}
		return nil
	case []interface{}:
		for _, s := range stmt {
			err := validateStatement(s, resources)
			if err != nil {
				return err
			}
		}
		return nil
	}

	return nil
}

func validateParam(value interface{}, resources domain.Resources) error {
	switch param := value.(type) {
	case domain.Chain:
		return validateChainParam(param, resources)
	case domain.Function:
		return validateParam(param.Target(), resources)
	case []interface{}:
		return validateListParam(param, resources)
	case map[string]interface{}:
		return validateObjectParam(param, resources)
	default:
		return nil
	}
}

func validateObjectParam(objectParam map[string]interface{}, resources domain.Resources) error {
	for _, value := range objectParam {
		err := validateParam(value, resources)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateListParam(listParam []interface{}, resources domain.Resources) error {
	for _, value := range listParam {
		err := validateParam(value, resources)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateChainParam(chain domain.Chain, resources domain.Resources) error {
	path := toPath(chain)
	resourceId := domain.ResourceId(path[0])

	_, found := resources[resourceId]
	if !found {
		return fmt.Errorf("%w : %s", ErrInvalidChainedParameter, strings.Join(path, "."))
	}

	return nil
}
