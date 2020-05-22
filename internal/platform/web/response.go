package web

import (
	"bytes"
	"encoding/json"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/platform/logger"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
)

func Respond(ctx *fasthttp.RequestCtx, data interface{}, statusCode int, headers map[string]string) error {
	ctx.Response.Header.SetContentType("application/json; charset=utf-8")
	ctx.Response.SetStatusCode(statusCode)
	for k, v := range headers {
		ctx.Response.Header.Set(k, v)
	}

	if data != nil {
		encoder := json.NewEncoder(ctx.Response.BodyWriter())
		if err := encoder.Encode(&data); err != nil {
			return err
		}
	}

	return nil
}

func RespondError(ctx *fasthttp.RequestCtx, err error, log *logger.Logger) error {

	// If the error was of the type *Error, the handler has
	// a specific status code and error to return.
	if webErr, ok := err.(*Error); ok {
		er := ErrorResponse{
			Error: webErr.Err.Error(),
		}
		if err := Respond(ctx, er, webErr.Status, nil); err != nil {
			return err
		}
		return nil
	}

	log.Error("returning 500 response", err)

	er := ErrorResponse{
		Error: err.Error(),
	}
	if err := Respond(ctx, er, http.StatusInternalServerError, nil); err != nil {
		return err
	}
	return nil
}

type StatementDebugging struct {
	Method          string                 `json:"method,omitempty"`
	Url             string                 `json:"url,omitempty"`
	RequestHeaders  map[string]string      `json:"request-headers,omitempty"`
	ResponseHeaders map[string]string      `json:"response-headers,omitempty"`
	Params          map[string]interface{} `json:"params,omitempty"`
	RequestBody     interface{}            `json:"request-body,omitempty"`
	ResponseTime    int64                  `json:"response-time,omitempty"`
}

type StatementMetadata struct {
	IgnoreErrors string `json:"ignore-errors,omitempty"`
}

type StatementDetails struct {
	Status   int                 `json:"status"`
	Success  bool                `json:"success"`
	Metadata StatementMetadata   `json:"metadata"`
	Debug    *StatementDebugging `json:"debug,omitempty"`
}

type StatementResult struct {
	Details interface{} `json:"details"`
	Result  interface{} `json:"result,omitempty"`
}

type QueryResponse struct {
	StatusCode int
	Body       map[string]StatementResult
	Headers    map[string]string
}

func MakeQueryResponse(queryResult domain.Resources) QueryResponse {
	m := make(map[string]StatementResult)
	for key, response := range queryResult {
		m[string(key)] = parseResponse(response)
	}

	statusCode := CalculateStatusCode(queryResult)
	headers := makeHeaders(queryResult)
	return QueryResponse{Body: m, StatusCode: statusCode, Headers: headers}
}

func parseResponse(response interface{}) StatementResult {
	switch response := response.(type) {
	case domain.DoneResource:
		return StatementResult{Details: parseDetails(response.Details), Result: response.Result}
	case domain.DoneResources:
		details := make([]interface{}, len(response))
		results := make([]interface{}, len(response))

		hasResult := false

		for i, r := range response {
			result := parseResponse(r)

			d := result.Details
			if d != nil {
				details[i] = d
			}

			r := result.Result
			if r != nil {
				hasResult = true
				results[i] = r
			}
		}

		if !hasResult {
			return StatementResult{Details: details, Result: nil}
		}

		return StatementResult{Details: details, Result: results}
	default:
		return StatementResult{}
	}
}

func parseDetails(details domain.Details) StatementDetails {
	var metadata StatementMetadata
	if details.IgnoreErrors {
		metadata.IgnoreErrors = "ignore"
	}

	return StatementDetails{
		Status:   details.Status,
		Success:  details.Success,
		Metadata: metadata,
		Debug:    parseDebug(details.Debug),
	}
}

func parseDebug(debug *domain.Debugging) *StatementDebugging {
	if debug == nil {
		return nil
	}

	return &StatementDebugging{
		Method:          debug.Method,
		Url:             debug.Url,
		RequestHeaders:  debug.RequestHeaders,
		ResponseHeaders: debug.ResponseHeaders,
		Params:          debug.Params,
		RequestBody:     debug.RequestBody,
		ResponseTime:    debug.ResponseTime,
	}
}

func CalculateStatusCode(queryResult domain.Resources) int {
	results := make([]interface{}, len(queryResult))
	index := 0
	for _, r := range queryResult {
		results[index] = r
		index++
	}

	maxStatusCode := findMaxStatusCode(results)

	return maxStatusCode
}

var statusNormalization = map[int]int{0: 500, 204: 200, 201: 200}

func calculateResultStatusCode(result interface{}) int {
	switch r := result.(type) {
	case domain.DoneResource:
		if r.Details.IgnoreErrors {
			return 200
		}

		status := r.Details.Status
		normalizedStatus, found := statusNormalization[status]
		if found {
			return normalizedStatus
		}

		return status
	case domain.DoneResources:
		return findMaxStatusCode(r)
	default:
		return 500
	}
}

func findMaxStatusCode(results []interface{}) int {
	resourceStatuses := make([]int, len(results))
	for i, result := range results {
		resourceStatuses[i] = calculateResultStatusCode(result)
	}

	maxStatusCode := 200
	for _, status := range resourceStatuses {
		if status > maxStatusCode {
			maxStatusCode = status
		}
	}
	return maxStatusCode
}

func makeHeaders(queryResult domain.Resources) map[string]string {
	cacheControl := calculateCacheControl(queryResult)
	cacheControlString := generateCacheControlString(cacheControl)

	headers := make(map[string]string)
	if cacheControlString != "" {
		headers["Cache-Control"] = cacheControlString
	}

	return headers
}

func calculateCacheControl(queryResult domain.Resources) domain.ResourceCacheControl {
	results := make([]interface{}, len(queryResult))
	index := 0
	for _, r := range queryResult {
		results[index] = r
		index++
	}

	return findMinCacheControl(results)
}

func findMinCacheControl(results []interface{}) domain.ResourceCacheControl {
	resourceCacheControls := make([]domain.ResourceCacheControl, len(results))
	for i, result := range results {
		resourceCacheControls[i] = calculateResultCacheControl(result)
	}

	minCacheControl := domain.ResourceCacheControl{
		MaxAge:  domain.ResourceCacheControlValue{Exist: false},
		SMaxAge: domain.ResourceCacheControlValue{Exist: false},
	}

	for _, cc := range resourceCacheControls {
		switch {
		case minCacheControl.NoCache:
			continue
		case cc.NoCache:
			minCacheControl.NoCache = true
		default:
			if !minCacheControl.MaxAge.Exist || cc.MaxAge.Time < minCacheControl.MaxAge.Time {
				minCacheControl.MaxAge = cc.MaxAge
			}

			if !minCacheControl.SMaxAge.Exist || cc.SMaxAge.Time < minCacheControl.SMaxAge.Time {
				minCacheControl.SMaxAge = cc.SMaxAge
			}

			minCacheControl.NoCache = false
		}
	}

	return minCacheControl
}

func calculateResultCacheControl(result interface{}) domain.ResourceCacheControl {
	switch result := result.(type) {
	case domain.DoneResource:
		return result.Details.CacheControl
	case domain.DoneResources:
		return findMinCacheControl(result)
	default:
		return domain.ResourceCacheControl{}
	}
}

func generateCacheControlString(cacheControl domain.ResourceCacheControl) string {
	var buf bytes.Buffer

	if cacheControl.NoCache {
		return "no-cache"
	}

	if cacheControl.MaxAge.Exist {
		buf.WriteString("max-age=")
		buf.WriteString(strconv.Itoa(cacheControl.MaxAge.Time))
	}

	if cacheControl.SMaxAge.Exist {
		if buf.Len() > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString("s-maxage=")
		buf.WriteString(strconv.Itoa(cacheControl.SMaxAge.Time))
	}

	return buf.String()
}
