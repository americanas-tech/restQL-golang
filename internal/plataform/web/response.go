package web

import (
	"encoding/json"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/valyala/fasthttp"
	"net/http"
)

func Respond(ctx *fasthttp.RequestCtx, data interface{}, statusCode int) error {
	ctx.Response.Header.SetContentType("application/json; charset=utf-8")
	ctx.Response.SetStatusCode(statusCode)

	if data != nil {
		encoder := json.NewEncoder(ctx.Response.BodyWriter())
		if err := encoder.Encode(&data); err != nil {
			return err
		}
	}

	return nil
}

func RespondError(ctx *fasthttp.RequestCtx, err error) error {

	// If the error was of the type *Error, the handler has
	// a specific status code and error to return.
	if webErr, ok := err.(*Error); ok {
		er := ErrorResponse{
			Error: webErr.Err.Error(),
		}
		if err := Respond(ctx, er, webErr.Status); err != nil {
			return err
		}
		return nil
	}

	er := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}
	if err := Respond(ctx, er, http.StatusInternalServerError); err != nil {
		return err
	}
	return nil
}

type StatementDebugging struct {
	Url             string                 `json:"url,omitempty"`
	RequestHeaders  map[string]string      `json:"request-headers,omitempty"`
	ResponseHeaders map[string]string      `json:"response-headers,omitempty"`
	Params          map[string]interface{} `json:"params,omitempty"`
	ResponseTime    int64                  `json:"response-time,omitempty"`
}

type StatementDetails struct {
	Status  int                 `json:"status"`
	Success bool                `json:"success"`
	Debug   *StatementDebugging `json:"debug,omitempty"`
}

type StatementResult struct {
	Details interface{} `json:"details"`
	Result  interface{} `json:"result"`
}

type QueryResponse map[string]StatementResult

func MakeQueryResponse(queryResult domain.Resources) QueryResponse {
	m := make(QueryResponse)
	for key, response := range queryResult {
		m[string(key)] = parseResponse(response)
	}

	return m
}

func parseResponse(response interface{}) StatementResult {
	switch response := response.(type) {
	case domain.DoneResource:
		return StatementResult{Details: parseDetails(response.Details), Result: response.Result}
	case domain.DoneResources:
		details := make([]interface{}, len(response))
		results := make([]interface{}, len(response))

		for i, r := range response {
			result := parseResponse(r)
			details[i] = result.Details
			results[i] = result.Result
		}

		return StatementResult{Details: details, Result: results}
	default:
		return StatementResult{}
	}
}

func parseDetails(details domain.Details) StatementDetails {
	return StatementDetails{
		Status:  details.Status,
		Success: details.Success,
		Debug:   parseDebug(details.Debug),
	}
}

func parseDebug(debug *domain.Debugging) *StatementDebugging {
	if debug == nil {
		return nil
	}

	return &StatementDebugging{
		Url:             debug.Url,
		RequestHeaders:  debug.RequestHeaders,
		ResponseHeaders: debug.ResponseHeaders,
		Params:          debug.Params,
		ResponseTime:    debug.ResponseTime,
	}
}
