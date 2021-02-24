package runner

import (
	"bytes"
	"github.com/b2wdigital/restQL-golang/v5/pkg/restql"
	"strconv"
	"strings"

	"github.com/b2wdigital/restQL-golang/v5/internal/domain"
)

// DoneResourceOptions represents information
// from the statement that should be passed
// to the result.
type DoneResourceOptions struct {
	Debugging    bool
	IgnoreErrors bool
	MaxAge       interface{}
	SMaxAge      interface{}
}

// NewDoneResource constructs a DoneResourceOptions value.
func NewDoneResource(request restql.HTTPRequest, response restql.HTTPResponse, options DoneResourceOptions) restql.DoneResource {
	dr := restql.DoneResource{
		Status:          response.StatusCode,
		Success:         response.StatusCode >= 200 && response.StatusCode < 400,
		IgnoreErrors:    options.IgnoreErrors,
		CacheControl:    makeCacheControl(response, options),
		Method:          request.Method,
		URL:             response.URL,
		RequestParams:   request.Query,
		RequestBody:     request.Body,
		RequestHeaders:  request.Headers,
		ResponseHeaders: response.Headers,
		ResponseBody:    response.Body,
		ResponseTime:    response.Duration.Milliseconds(),
	}

	return dr
}

// NewErrorResponse builds a DoneResource value for a failed HTTP call.
func NewErrorResponse(log restql.Logger, err error, request restql.HTTPRequest, response restql.HTTPResponse, options DoneResourceOptions) restql.DoneResource {
	rb := restql.NewResponseBodyFromValue(log, err.Error())

	return restql.DoneResource{
		Status:          response.StatusCode,
		Success:         false,
		IgnoreErrors:    options.IgnoreErrors,
		ResponseBody:    rb,
		Method:          request.Method,
		URL:             response.URL,
		RequestParams:   request.Query,
		RequestBody:     request.Body,
		RequestHeaders:  request.Headers,
		ResponseHeaders: response.Headers,
		ResponseTime:    response.Duration.Milliseconds(),
	}
}

// NewEmptyChainedResponse builds a DoneResource for a statement
// with unresolved chain parameters.
func NewEmptyChainedResponse(log restql.Logger, params []string, options DoneResourceOptions) restql.DoneResource {
	var buf bytes.Buffer

	buf.WriteString("The request was skipped due to missing { ")
	for _, p := range params {
		buf.WriteString(":")
		buf.WriteString(p)
		buf.WriteString(" ")
	}
	buf.WriteString("} param value")

	rb := restql.NewResponseBodyFromValue(log, buf.String())
	return restql.DoneResource{
		Status:       400,
		Success:      false,
		IgnoreErrors: options.IgnoreErrors,
		ResponseBody: rb,
	}
}

// GetEmptyChainedParams returns the chain parameters that
// could not be resolved.
func GetEmptyChainedParams(statement domain.Statement) []string {
	var r []string
	for key, value := range statement.With.Values {
		if isEmptyChained(value) {
			r = append(r, key)
		}
	}

	return r
}

func isEmptyChained(value interface{}) bool {
	switch value := value.(type) {
	case map[string]interface{}:
		for _, v := range value {
			if isEmptyChained(v) {
				return true
			}
		}

		return false
	case []interface{}:
		for _, v := range value {
			if isEmptyChained(v) {
				return true
			}
		}

		return false
	default:
		return value == EmptyChained
	}
}

func makeCacheControl(response restql.HTTPResponse, options DoneResourceOptions) restql.ResourceCacheControl {
	headerCacheControl, headerFound := getCacheControlOptionsFromHeader(response)
	defaultCacheControl, defaultFound := getDefaultCacheControlOptions(options)

	if !headerFound && !defaultFound {
		return restql.ResourceCacheControl{}
	}

	switch {
	case !headerFound && !defaultFound:
		return restql.ResourceCacheControl{}
	case !headerFound:
		return defaultCacheControl
	case !defaultFound:
		return headerCacheControl
	default:
		return bestCacheControl(headerCacheControl, defaultCacheControl)
	}
}

func bestCacheControl(first restql.ResourceCacheControl, second restql.ResourceCacheControl) restql.ResourceCacheControl {
	result := restql.ResourceCacheControl{}

	if first.NoCache || second.NoCache {
		result.NoCache = true
		return result
	}

	result.MaxAge = bestCacheControlValue(first.MaxAge, second.MaxAge)
	result.SMaxAge = bestCacheControlValue(first.SMaxAge, second.SMaxAge)

	return result
}

func bestCacheControlValue(first restql.ResourceCacheControlValue, second restql.ResourceCacheControlValue) restql.ResourceCacheControlValue {
	switch {
	case !first.Exist && !second.Exist:
		return restql.ResourceCacheControlValue{Exist: false}
	case !first.Exist:
		return second
	case !second.Exist:
		return first
	default:
		time := min(first.Time, second.Time)
		return restql.ResourceCacheControlValue{Exist: true, Time: time}
	}
}

func min(a int, b int) int {
	if a < b {
		return a
	}

	return b
}

func getDefaultCacheControlOptions(options DoneResourceOptions) (cc restql.ResourceCacheControl, found bool) {
	maxAge, ok := options.MaxAge.(int)
	if ok {
		found = true
		cc.MaxAge = restql.ResourceCacheControlValue{Exist: true, Time: maxAge}
	}

	smaxAge, ok := options.SMaxAge.(int)
	if ok {
		found = true
		cc.SMaxAge = restql.ResourceCacheControlValue{Exist: true, Time: smaxAge}
	}

	return cc, found
}

func isComma(r rune) bool {
	return r == ','
}

func findCacheControlHeader(response restql.HTTPResponse) (string, bool) {
	for k, v := range response.Headers {
		if strings.EqualFold(k, "Cache-Control") {
			return v, true
		}
	}

	return "", false
}

func getCacheControlOptionsFromHeader(response restql.HTTPResponse) (cc restql.ResourceCacheControl, found bool) {
	cacheControl, ok := findCacheControlHeader(response)

	if !ok {
		return restql.ResourceCacheControl{}, false
	}

	cacheControlFields := strings.FieldsFunc(cacheControl, isComma)

	for _, ccField := range cacheControlFields {
		ccField = strings.TrimSpace(ccField)

		if strings.EqualFold(ccField, "no-cache") {
			return restql.ResourceCacheControl{NoCache: true}, true
		}

		keyValue := strings.Split(ccField, "=")
		if len(keyValue) < 2 {
			continue
		}

		key, value := keyValue[0], keyValue[1]

		if strings.EqualFold(key, "max-age") {
			timeValue, err := strconv.Atoi(value)
			if err != nil {
				continue
			}

			found = true
			cc.MaxAge = restql.ResourceCacheControlValue{Exist: true, Time: timeValue}
		}

		if strings.EqualFold(key, "s-maxage") {
			timeValue, err := strconv.Atoi(value)
			if err != nil {
				continue
			}

			found = true
			cc.SMaxAge = restql.ResourceCacheControlValue{Exist: true, Time: timeValue}
		}
	}

	return cc, found
}
