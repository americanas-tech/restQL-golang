package runner

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/b2wdigital/restQL-golang/v5/internal/domain"
	"github.com/b2wdigital/restQL-golang/v5/pkg/restql"
)

// ApplyEncoders transform parameter values with encoder functions applied
// into a Resource collection with the values processed.
func ApplyEncoders(resources domain.Resources, log restql.Logger) domain.Resources {
	for resourceID, statement := range resources {
		if statement, ok := statement.(domain.Statement); ok {
			resources[resourceID] = applyEncoderToStatement(log, statement)
		}
	}

	return resources
}

func applyEncoderToStatement(log restql.Logger, statement domain.Statement) domain.Statement {
	values := statement.With.Values
	for key, value := range values {
		result := applyEncoderToValue(log, value)

		values[key] = result
	}

	body := applyEncoderToBody(log, statement.With.Body)

	statement.With.Body = body
	statement.With.Values = values

	return statement
}

func applyEncoderToBody(log restql.Logger, body interface{}) interface{} {
	switch body := body.(type) {
	case domain.Base64:
		return applyBase64encoder(applyEncoderToBody(log, body.Target()))
	case domain.JSON:
		return applyEncoderToBody(log, body.Target())
	case domain.Flatten:
		return applyFlattenEncoder(log, applyEncoderToBody(log, body.Target()))
	case domain.Function:
		return body.Map(func(target interface{}) interface{} {
			return applyEncoderToBody(log, target)
		})
	default:
		return body
	}
}

func applyEncoderToValue(log restql.Logger, value interface{}) interface{} {
	switch value := value.(type) {
	case domain.Base64:
		target := value.Target()
		if _, ok := target.(domain.Chain); ok {
			return value
		}

		return applyBase64encoder(applyEncoderToValue(log, value.Target()))
	case domain.JSON:
		target := value.Target()
		if _, ok := target.(domain.Chain); ok {
			return value
		}

		return applyJSONEncoder(log, applyEncoderToValue(log, value.Target()))
	case domain.Flatten:
		target := value.Target()
		if _, ok := target.(domain.Chain); ok {
			return value
		}

		return applyFlattenEncoder(log, applyEncoderToValue(log, value.Target()))
	case domain.Function:
		return value.Map(func(target interface{}) interface{} {
			return applyEncoderToValue(log, target)
		})
	case map[string]interface{}:
		m := make(map[string]interface{})
		for k, v := range value {
			m[k] = applyEncoderToValue(log, v)
		}
		return m
	case []interface{}:
		l := make([]interface{}, len(value))
		for i, v := range value {
			l[i] = applyEncoderToValue(log, v)
		}
		return l
	default:
		return value
	}
}

func applyJSONEncoder(log restql.Logger, value interface{}) interface{} {
	data, err := json.Marshal(value)
	if err != nil {
		log.Debug("failed to apply json encoder", "target", value)
	}

	return string(data)
}

func applyBase64encoder(value interface{}) interface{} {
	data := []byte(fmt.Sprintf("%v", value))
	return base64.StdEncoding.EncodeToString(data)
}

func applyFlattenEncoder(log restql.Logger, value interface{}) interface{} {
	if value, ok := value.([]interface{}); ok {
		return flatten(value)
	}

	log.Warn("flatten encoder used on non list value", "value", value)
	return value
}

func flatten(ii []interface{}) []interface{} {
	var res []interface{}
	for _, i := range ii {
		if i == nil {
			continue
		}
		switch t := i.(type) {
		case []interface{}:
			res = append(res, flatten(t)...)
		default:
			res = append(res, i)
		}
	}

	return res
}
