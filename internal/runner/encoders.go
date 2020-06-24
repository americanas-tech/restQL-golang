package runner

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/pkg/restql"
)

func ApplyEncoders(resources domain.Resources, log restql.Logger) domain.Resources {
	for resourceId, statement := range resources {
		if statement, ok := statement.(domain.Statement); ok {
			resources[resourceId] = applyEncoderToStatement(log, statement)
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

	body := applyEncoderToBody(statement.With.Body)

	statement.With.Body = body
	statement.With.Values = values

	return statement
}

func applyEncoderToBody(body interface{}) interface{} {
	switch body := body.(type) {
	case domain.Base64:
		return applyBase64encoder(applyEncoderToBody(body.Target()))
	case domain.Json:
		return applyEncoderToBody(body.Target())
	case domain.Function:
		return body.Map(func(target interface{}) interface{} {
			return applyEncoderToBody(target)
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
	case domain.Json:
		target := value.Target()
		if _, ok := target.(domain.Chain); ok {
			return value
		}

		return applyJsonEncoder(log, applyEncoderToValue(log, value.Target()))
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

func applyJsonEncoder(log restql.Logger, value interface{}) interface{} {
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
