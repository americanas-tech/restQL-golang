package plugins

import (
	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
)

// DecodeQueryResult transforms a resolved Resources collection
// into a generic map.
func DecodeQueryResult(queryResult domain.Resources) map[string]interface{} {
	m := make(map[string]interface{})
	for key, resource := range queryResult {
		m[string(key)] = parseResource(resource)
	}
	return m
}

func parseResource(resource interface{}) map[string]interface{} {
	m := make(map[string]interface{})

	switch resource := resource.(type) {
	case domain.DoneResource:
		m["details"] = parseDetails(resource)
		m["result"] = resource.ResponseBody

		return m
	case domain.DoneResources:
		details := make([]interface{}, len(resource))
		results := make([]interface{}, len(resource))

		hasResult := false

		for i, r := range resource {
			result := parseResource(r)

			d := result["details"]
			if d != nil {
				details[i] = d
			}

			r := result["result"]
			if r != nil {
				hasResult = true
				results[i] = r
			}
		}

		if !hasResult {
			m["details"] = details
			m["result"] = nil

			return m
		}

		m["details"] = details
		m["result"] = results

		return m
	default:
		return m
	}
}

func parseDetails(resource domain.DoneResource) map[string]interface{} {
	debug := map[string]interface{}{
		"method":          resource.Method,
		"url":             resource.URL,
		"requestHeaders":  resource.RequestHeaders,
		"params":          resource.RequestParams,
		"responseHeaders": resource.ResponseHeaders,
		"requestBody":     resource.RequestBody,
		"responseTime":    resource.ResponseTime,
	}

	metadata := make(map[string]interface{})
	if resource.IgnoreErrors {
		metadata["ignore-errors"] = true
	}

	return map[string]interface{}{
		"status":    resource.Status,
		"success":   resource.Success,
		"metadata":  metadata,
		"debugging": debug,
	}
}
