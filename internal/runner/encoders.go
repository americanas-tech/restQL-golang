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
		return applyBase64encoder(body)
	case domain.Json:
		return body.Target
	default:
		return body
	}
}

func applyEncoderToValue(log restql.Logger, value interface{}) interface{} {
	var result interface{}
	switch value := value.(type) {
	case domain.Base64:
		result = applyBase64encoder(value)
	case domain.Json:
		result = applyJsonEncoder(log, value)
	case map[string]interface{}:
		m := make(map[string]interface{})
		for k, v := range value {
			m[k] = applyEncoderToValue(log, v)
		}
		result = m
	case []interface{}:
		l := make([]interface{}, len(value))
		for i, v := range value {
			l[i] = applyEncoderToValue(log, v)
		}
		result = l
	default:
		result = value
	}
	return result
}

func applyJsonEncoder(log restql.Logger, value domain.Json) interface{} {
	if _, ok := value.Target.(domain.Chain); ok {
		return value
	}

	data, err := json.Marshal(value.Target)
	if err != nil {
		log.Debug("failed to apply json encoder", "target", value.Target)
	}

	return string(data)
}

func applyBase64encoder(value domain.Base64) interface{} {
	if _, ok := value.Target.(domain.Chain); ok {
		return value
	}

	data := []byte(fmt.Sprintf("%v", value.Target))
	return base64.StdEncoding.EncodeToString(data)
}
