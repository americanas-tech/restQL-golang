package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
	"net/http"
	"strconv"

	"github.com/valyala/fasthttp"
)

// Respond write the information back to the client.
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

// RespondError translate the error and write it back to the client.
func RespondError(ctx *fasthttp.RequestCtx, err error) error {

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

	er := ErrorResponse{
		Error: err.Error(),
	}
	if err := Respond(ctx, er, http.StatusInternalServerError, nil); err != nil {
		return err
	}
	return nil
}

// StatementDebugging represents the client format of debugging information
type StatementDebugging struct {
	Method          string                 `json:"method,omitempty"`
	URL             string                 `json:"url,omitempty"`
	RequestHeaders  map[string]string      `json:"request-headers,omitempty"`
	ResponseHeaders map[string]string      `json:"response-headers,omitempty"`
	Params          map[string]interface{} `json:"params,omitempty"`
	RequestBody     interface{}            `json:"request-body,omitempty"`
	ResponseTime    int64                  `json:"response-time,omitempty"`
}

// StatementMetadata represents the client format of metadata
type StatementMetadata struct {
	IgnoreErrors string `json:"ignore-errors,omitempty"`
}

// StatementDetails represents the client format of the statement details
type StatementDetails struct {
	Status   int                 `json:"status"`
	Success  bool                `json:"success"`
	Metadata StatementMetadata   `json:"metadata"`
	Debug    *StatementDebugging `json:"debug,omitempty"`
}

// StatementResult represents the client format of the statement result
type StatementResult struct {
	Details interface{} `json:"details"`
	Result  interface{} `json:"result,omitempty"`
}

// QueryResponse represents the client format of the query result
type QueryResponse struct {
	StatusCode int
	Body       map[string]StatementResult
	Headers    map[string]string
}

// MakeQueryResponse create a query execution response for the client.
func MakeQueryResponse(queryResult domain.Resources, debug bool) (QueryResponse, error) {
	m := make(map[string]StatementResult)
	for key, resource := range queryResult {
		r, err := parseResource(resource, debug)
		if err != nil {
			return QueryResponse{}, err
		}

		m[string(key)] = r
	}

	statusCode := CalculateStatusCode(queryResult)
	headers := makeHeaders(queryResult)
	return QueryResponse{Body: m, StatusCode: statusCode, Headers: headers}, nil
}

func parseResource(resource interface{}, debug bool) (StatementResult, error) {
	switch resource := resource.(type) {
	case restql.DoneResource:
		body, err := resource.ResponseBody.Marshal()
		if err != nil {
			return StatementResult{}, err
		}

		return StatementResult{Details: parseDetails(resource, debug), Result: body}, nil
	case restql.DoneResources:
		details := make([]interface{}, len(resource))
		results := make([]interface{}, len(resource))

		hasResult := false

		for i, r := range resource {
			result, err := parseResource(r, debug)
			if err != nil {
				return StatementResult{}, err
			}

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
			return StatementResult{Details: details, Result: nil}, nil
		}

		return StatementResult{Details: details, Result: results}, nil
	default:
		return StatementResult{}, nil
	}
}

func parseDetails(resource restql.DoneResource, debug bool) StatementDetails {
	var metadata StatementMetadata
	if resource.IgnoreErrors {
		metadata.IgnoreErrors = "ignore"
	}

	sd := StatementDetails{
		Status:   resource.Status,
		Success:  resource.Success,
		Metadata: metadata,
	}

	if debug {
		sd.Debug = parseDebug(resource)
	}

	return sd
}

func parseDebug(resource restql.DoneResource) *StatementDebugging {
	return &StatementDebugging{
		Method:          resource.Method,
		URL:             resource.URL,
		RequestHeaders:  resource.RequestHeaders,
		ResponseHeaders: resource.ResponseHeaders,
		Params:          resource.RequestParams,
		RequestBody:     resource.RequestBody,
		ResponseTime:    resource.ResponseTime,
	}
}

// CalculateStatusCode returns the greater status in all
// statement results to be used as the response status code.
// It applies the following normalization to statement result status codes:
//
// 0 => 500
// 204 => 200
// 201 => 200
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
	case restql.DoneResource:
		if r.IgnoreErrors {
			return 200
		}

		status := r.Status
		normalizedStatus, found := statusNormalization[status]
		if found {
			return normalizedStatus
		}

		return status
	case restql.DoneResources:
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
	resourceHeaders := makeResourceHeaders(queryResult)
	ccHeaders := makeCacheControlHeaders(queryResult)

	return appendMap(resourceHeaders, ccHeaders)
}

func makeResourceHeaders(queryResult domain.Resources) map[string]string {
	headers := make(map[string]string)
	for resourceID, response := range queryResult {
		switch response := response.(type) {
		case restql.DoneResource:
			headers = appendMap(headers, getResponseHeader(string(resourceID), response))
		}
	}

	return headers
}

func getResponseHeader(resourceID string, response restql.DoneResource) map[string]string {
	if !response.Success {
		return nil
	}

	headers := make(map[string]string)
	for hn, hv := range response.ResponseHeaders {
		headerName := fmt.Sprintf("%s-%s", resourceID, hn)
		headers[headerName] = hv
	}

	return headers
}

func makeCacheControlHeaders(queryResult domain.Resources) map[string]string {
	cacheControl := calculateCacheControl(queryResult)
	cacheControlString := generateCacheControlString(cacheControl)

	headers := make(map[string]string)
	if cacheControlString != "" {
		headers["Cache-Control"] = cacheControlString
	}
	return headers
}

func calculateCacheControl(queryResult domain.Resources) restql.ResourceCacheControl {
	results := make([]interface{}, len(queryResult))
	index := 0
	for _, r := range queryResult {
		results[index] = r
		index++
	}

	return findMinCacheControl(results)
}

func findMinCacheControl(results []interface{}) restql.ResourceCacheControl {
	resourceCacheControls := make([]restql.ResourceCacheControl, len(results))
	for i, result := range results {
		resourceCacheControls[i] = calculateResultCacheControl(result)
	}

	minCacheControl := restql.ResourceCacheControl{
		MaxAge:  restql.ResourceCacheControlValue{Exist: false},
		SMaxAge: restql.ResourceCacheControlValue{Exist: false},
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

func calculateResultCacheControl(result interface{}) restql.ResourceCacheControl {
	switch result := result.(type) {
	case restql.DoneResource:
		return result.CacheControl
	case restql.DoneResources:
		return findMinCacheControl(result)
	default:
		return restql.ResourceCacheControl{}
	}
}

func generateCacheControlString(cacheControl restql.ResourceCacheControl) string {
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

func appendMap(m1 map[string]string, m2 map[string]string) map[string]string {
	m := make(map[string]string, len(m1)+len(m2))
	for k, v := range m1 {
		m[k] = v
	}

	for k, v := range m2 {
		m[k] = v
	}

	return m
}
