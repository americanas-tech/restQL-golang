package runner

import (
	"bytes"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"strconv"
)

const debugParamName = "_debug"

type DoneResourceOptions struct {
	Debugging    bool
	IgnoreErrors bool
}

func NewDoneResource(request domain.HttpRequest, response domain.HttpResponse, options DoneResourceOptions) domain.DoneResource {
	dr := domain.DoneResource{
		Details: domain.Details{
			Status:       response.StatusCode,
			Success:      response.StatusCode >= 200 && response.StatusCode < 400,
			IgnoreErrors: options.IgnoreErrors,
		},
		Result: response.Body,
	}

	if options.Debugging {
		dr.Details.Debug = newDebugging(request, response)
	}

	return dr
}

func newDebugging(request domain.HttpRequest, response domain.HttpResponse) *domain.Debugging {
	return &domain.Debugging{
		Url:             response.Url,
		Params:          request.Query,
		RequestHeaders:  request.Headers,
		ResponseHeaders: response.Headers,
		ResponseTime:    response.Duration.Milliseconds(),
	}
}

func IsDebugEnabled(queryCtx domain.QueryContext) bool {
	param, found := queryCtx.Input.Params[debugParamName]
	if !found {
		return false
	}

	debug, ok := param.(string)
	if !ok {
		return false
	}

	d, err := strconv.ParseBool(debug)
	if err != nil {
		return false
	}

	return d
}

func NewTimeoutResponse(err error, request domain.HttpRequest, response domain.HttpResponse, options DoneResourceOptions) domain.DoneResource {
	dr := domain.DoneResource{
		Details: domain.Details{
			Status:       408,
			Success:      false,
			IgnoreErrors: options.IgnoreErrors,
		},
		Result: err.Error(),
	}

	if options.Debugging {
		dr.Details.Debug = newDebugging(request, response)
	}

	return dr
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
		Details: domain.Details{Status: 400, Success: false, IgnoreErrors: options.IgnoreErrors},
		Result:  buf.String(),
	}
}

func GetEmptyChainedParams(statement domain.Statement) []string {
	var r []string
	for key, value := range statement.With {
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
