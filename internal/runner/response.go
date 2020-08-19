package runner

import (
	"bytes"
	"regexp"
	"strconv"

	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
)

type DoneResourceOptions struct {
	Debugging    bool
	IgnoreErrors bool
	MaxAge       interface{}
	SMaxAge      interface{}
}

func NewDoneResource(request domain.HttpRequest, response domain.HttpResponse, options DoneResourceOptions) domain.DoneResource {
	dr := domain.DoneResource{
		Status:          response.StatusCode,
		Success:         response.StatusCode >= 200 && response.StatusCode < 400,
		IgnoreErrors:    options.IgnoreErrors,
		CacheControl:    makeCacheControl(response, options),
		Method:          request.Method,
		Url:             response.Url,
		RequestParams:   request.Query,
		RequestBody:     request.Body,
		RequestHeaders:  request.Headers,
		ResponseHeaders: response.Headers,
		ResponseBody:    response.Body,
		ResponseTime:    response.Duration.Milliseconds(),
	}

	return dr
}

func NewErrorResponse(err error, request domain.HttpRequest, response domain.HttpResponse, options DoneResourceOptions) domain.DoneResource {
	return domain.DoneResource{
		Status:          response.StatusCode,
		Success:         false,
		IgnoreErrors:    options.IgnoreErrors,
		ResponseBody:    err.Error(),
		Method:          request.Method,
		Url:             response.Url,
		RequestParams:   request.Query,
		RequestBody:     request.Body,
		RequestHeaders:  request.Headers,
		ResponseHeaders: response.Headers,
		ResponseTime:    response.Duration.Milliseconds(),
	}
}

func NewEmptyChainedResponse(params []string, options DoneResourceOptions) domain.DoneResource {
	var buf bytes.Buffer

	buf.WriteString("The request was skipped due to missing { ")
	for _, p := range params {
		buf.WriteString(":")
		buf.WriteString(p)
		buf.WriteString(" ")
	}
	buf.WriteString("} param value")

	return domain.DoneResource{
		Status:       400,
		Success:      false,
		IgnoreErrors: options.IgnoreErrors,
		ResponseBody: buf.String(),
	}
}

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

func makeCacheControl(response domain.HttpResponse, options DoneResourceOptions) domain.ResourceCacheControl {
	headerCacheControl, headerFound := getCacheControlOptionsFromHeader(response)
	defaultCacheControl, defaultFound := getDefaultCacheControlOptions(options)

	if !headerFound && !defaultFound {
		return domain.ResourceCacheControl{}
	}

	switch {
	case !headerFound && !defaultFound:
		return domain.ResourceCacheControl{}
	case !headerFound:
		return defaultCacheControl
	case !defaultFound:
		return headerCacheControl
	default:
		return bestCacheControl(headerCacheControl, defaultCacheControl)
	}
}

func bestCacheControl(first domain.ResourceCacheControl, second domain.ResourceCacheControl) domain.ResourceCacheControl {
	result := domain.ResourceCacheControl{}

	if first.NoCache || second.NoCache {
		result.NoCache = true
		return result
	}

	result.MaxAge = bestCacheControlValue(first.MaxAge, second.MaxAge)
	result.SMaxAge = bestCacheControlValue(first.SMaxAge, second.SMaxAge)

	return result
}

func bestCacheControlValue(first domain.ResourceCacheControlValue, second domain.ResourceCacheControlValue) domain.ResourceCacheControlValue {
	switch {
	case !first.Exist && !second.Exist:
		return domain.ResourceCacheControlValue{Exist: false}
	case !first.Exist:
		return second
	case !second.Exist:
		return first
	default:
		time := min(first.Time, second.Time)
		return domain.ResourceCacheControlValue{Exist: true, Time: time}
	}
}

func min(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func getDefaultCacheControlOptions(options DoneResourceOptions) (cc domain.ResourceCacheControl, found bool) {
	maxAge, ok := options.MaxAge.(int)
	if ok {
		found = true
		cc.MaxAge = domain.ResourceCacheControlValue{Exist: true, Time: maxAge}
	}

	smaxAge, ok := options.SMaxAge.(int)
	if ok {
		found = true
		cc.SMaxAge = domain.ResourceCacheControlValue{Exist: true, Time: smaxAge}
	}

	return cc, found
}

var maxAgeHeaderRegex = regexp.MustCompile("max-age=(\\d+)")
var smaxAgeHeaderRegex = regexp.MustCompile("s-maxage=(\\d+)")
var noCacheHeaderRegex = regexp.MustCompile("no-cache")

func getCacheControlOptionsFromHeader(response domain.HttpResponse) (cc domain.ResourceCacheControl, found bool) {
	cacheControl, ok := response.Headers["Cache-Control"]
	if !ok {
		return domain.ResourceCacheControl{}, false
	}

	if noCacheHeaderRegex.MatchString(cacheControl) {
		return domain.ResourceCacheControl{NoCache: true}, true
	}

	maxAgeMatches := maxAgeHeaderRegex.FindAllStringSubmatch(cacheControl, -1)
	maxAgeValue, ok := extractCacheControlValueFromHeader(maxAgeMatches)
	if ok {
		found = true
		cc.MaxAge = domain.ResourceCacheControlValue{Exist: true, Time: maxAgeValue}
	}

	smaxAgeMatches := smaxAgeHeaderRegex.FindAllStringSubmatch(cacheControl, -1)
	smaxAgeValue, ok := extractCacheControlValueFromHeader(smaxAgeMatches)
	if ok {
		found = true
		cc.SMaxAge = domain.ResourceCacheControlValue{Exist: true, Time: smaxAgeValue}
	}

	return cc, found
}

func extractCacheControlValueFromHeader(header [][]string) (int, bool) {
	if len(header) <= 0 || len(header[0]) < 2 {
		return 0, false
	}

	headerValue := header[0][1]
	time, err := strconv.Atoi(headerValue)
	if err != nil {
		return 0, false
	}

	return time, true
}
