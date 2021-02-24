package runner

import (
	"encoding/json"
	"fmt"
	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
	"github.com/pkg/errors"
	"math"
	"strings"
)

// EmptyChained in a token used to represent a chained parameter value
// that could not be resolved due to a failed response from the upstream dependency.
const EmptyChained = "__EMPTY_CHAINED__"

// ErrInvalidChainedParameter represents an error when a chain parameter value
// references an unknown statement.
var ErrInvalidChainedParameter = errors.New("chained parameter targeting unknown statement")

// ResolveChainedValues takes an unresolved Resource collection and replace
// chain parameter values by data present in the done Resource collection.
func ResolveChainedValues(resources domain.Resources, doneResources domain.Resources) domain.Resources {
	for resourceID, stmt := range resources {
		resources[resourceID] = resolveStatement(stmt, doneResources)
	}

	return resources
}

func resolveStatement(stmt interface{}, doneResources domain.Resources) interface{} {
	switch stmt := stmt.(type) {
	case domain.Statement:
		params := stmt.With.Values
		for paramName, value := range params {
			v := resolveValue(value, doneResources, resolverOptions{explode: true})
			if v == nil {
				continue
			}
			params[paramName] = v
		}

		headers := stmt.Headers
		for name, value := range headers {
			resolved := resolveValue(value, doneResources, resolverOptions{explode: true})
			if resolved == nil {
				continue
			}

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

type resolverOptions struct {
	explode bool
}

func resolveValue(value interface{}, doneResources domain.Resources, options resolverOptions) interface{} {
	switch param := value.(type) {
	case domain.Chain:
		return resolveChainParam(param, doneResources)
	case domain.NoExplode:
		return resolveValue(param.Target(), doneResources, resolverOptions{explode: false})
	case domain.Function:
		return param.Map(func(target interface{}) interface{} {
			return resolveValue(target, doneResources, options)
		})
	case []interface{}:
		return resolveListParam(param, doneResources, options)
	case map[string]interface{}:
		return resolveObjectParam(param, doneResources, options)
	default:
		return value
	}
}

func resolveObjectParam(objectParam map[string]interface{}, doneResources domain.Resources, options resolverOptions) interface{} {
	result := make(map[string]interface{})

	for key, value := range objectParam {
		v := resolveValue(value, doneResources, options)
		if v == nil {
			continue
		}

		result[key] = v
	}

	if options.explode {
		return explodeListValuesInNewMaps(result)
	}

	return result
}

func explodeListValuesInNewMaps(m map[string]interface{}) interface{} {
	n := minimumListValueLength(m)

	if n == 0 {
		return m
	}

	result := make([]interface{}, n)
	for i := 0; i < n; i++ {
		newMap := make(map[string]interface{})

		for k, v := range m {
			var newValue interface{}
			switch v := v.(type) {
			case []interface{}:
				newValue = v[i]
			default:
				newValue = v
			}

			if newValue == nil {
				continue
			}

			newMap[k] = newValue
		}

		result[i] = newMap
	}

	return result
}

func minimumListValueLength(m map[string]interface{}) int {
	var result uint64 = math.MaxInt64

	for _, v := range m {
		if v, ok := v.([]interface{}); ok {
			length := uint64(len(v))
			if length < result {
				result = length
			}

		}
	}

	if result == math.MaxInt64 {
		return 0
	}

	return int(result)
}

func resolveListParam(listParam []interface{}, doneResources domain.Resources, options resolverOptions) []interface{} {
	result := make([]interface{}, len(listParam))
	copy(result, listParam)

	for i, value := range result {
		result[i] = resolveValue(value, doneResources, options)
	}

	return result
}

func resolveChainParam(chain domain.Chain, doneResources domain.Resources) interface{} {
	path := toPath(chain)
	resourceID := domain.ResourceID(path[0])

	switch done := doneResources[resourceID].(type) {
	case restql.DoneResources:
		return resolveWithMultiplexedRequests(path[1:], done)
	case restql.DoneResource:
		return resolveWithSingleRequest(path[1:], done)
	default:
		return nil
	}
}

func resolveWithMultiplexedRequests(path []string, doneRequests restql.DoneResources) []interface{} {
	var result []interface{}

	for _, request := range doneRequests {
		switch request := request.(type) {
		case restql.DoneResource:
			result = append(result, resolveWithSingleRequest(path, request))
		case restql.DoneResources:
			result = append(result, resolveWithMultiplexedRequests(path, request))
		}
	}

	return result
}

func resolveWithSingleRequest(path []string, done restql.DoneResource) interface{} {
	if done.Status < 200 || done.Status >= 400 {
		return EmptyChained
	}

	body := done.ResponseBody.Unmarshal()
	valueFromBody, found := getValueFromBody(path, body)
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

func getValueFromBody(pathToValue []string, b restql.Body) (interface{}, bool) {
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
		return nil, false
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

// ValidateChainedValues returns an error if a chain
// parameter value references an unknown statement.
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
	resourceID := domain.ResourceID(path[0])

	_, found := resources[resourceID]
	if !found {
		return fmt.Errorf("%w : %s", ErrInvalidChainedParameter, strings.Join(path, "."))
	}

	return nil
}
