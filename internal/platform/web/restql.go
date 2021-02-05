package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/b2wdigital/restQL-golang/v4/internal/eval"
	"github.com/b2wdigital/restQL-golang/v4/internal/parser"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/conf"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/web/middleware"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

var jsonContentType = "application/json"

var (
	errInvalidRevisionType     = errors.New("invalid revision : must be an integer")
	errInvalidTenant           = errors.New("invalid tenant : no value provided")
	errFailedToReadRequestBody = errors.New("failed to read and unmarshal request body")
)

type restQl struct {
	config    *conf.Config
	log       restql.Logger
	evaluator eval.Evaluator
	parser    parser.Parser
}

func newRestQl(l restql.Logger, cfg *conf.Config, e eval.Evaluator, p parser.Parser) restQl {
	return restQl{config: cfg, log: l, evaluator: e, parser: p}
}

func (r restQl) ValidateQuery(ctx *fasthttp.RequestCtx) error {
	queryTxt := string(ctx.PostBody())
	_, err := r.parser.Parse(queryTxt)
	if err != nil {
		r.log.Error("an error occurred when parsing queryRevision", err)
		e := fmt.Errorf("%w: %s", parser.ErrInvalidQuery, err)

		return RespondError(ctx, e, errToStatusCode)
	}

	return Respond(ctx, nil, http.StatusOK, nil)
}

func (r restQl) RunAdHocQuery(reqCtx *fasthttp.RequestCtx) error {
	ctx := middleware.GetNativeContext(reqCtx)
	ctx = restql.WithLogger(reqCtx, r.log)

	tenant, err := makeTenant(reqCtx, r.config.Tenant)
	if err != nil {
		r.log.Error("failed to build query options", err)
		return RespondError(reqCtx, err, errToStatusCode)
	}
	options := restql.QueryOptions{Tenant: tenant}

	input, err := makeQueryInput(reqCtx, r.log)
	if err != nil {
		r.log.Error("failed to build query input", err)
		return RespondError(reqCtx, err, errToStatusCode)
	}

	queryTxt := string(reqCtx.PostBody())

	result, err := r.evaluator.AdHocQuery(ctx, queryTxt, options, input)
	if err != nil {
		r.log.Error("failed to evaluated adhoc query", err)

		adhocErrToStatusCode := make(map[error]int)
		for err, status := range errToStatusCode {
			adhocErrToStatusCode[err] = status
		}
		adhocErrToStatusCode[eval.ErrParser] = http.StatusBadRequest

		return RespondError(reqCtx, err, adhocErrToStatusCode)
	}

	debugEnabled := isDebugEnabled(input)
	response, err := MakeQueryResponse(result, debugEnabled)
	if err != nil {
		return RespondError(reqCtx, err, errToStatusCode)
	}

	return Respond(reqCtx, response.Body, response.StatusCode, response.Headers)
}

func (r restQl) RunSavedQuery(reqCtx *fasthttp.RequestCtx) error {
	log := r.log.With("restql-endpoint", string(reqCtx.Request.URI().Path()))
	log = log.With("request-id", string(reqCtx.Request.Header.Peek("X-TID")))

	ctx := middleware.GetNativeContext(reqCtx)
	ctx = restql.WithLogger(ctx, log)

	options, err := makeQueryOptions(reqCtx, log, r.config.Tenant)
	if err != nil {
		log.Error("failed to build query options", err)
		return RespondError(reqCtx, err, errToStatusCode)
	}

	input, err := makeQueryInput(reqCtx, log)
	if err != nil {
		log.Error("failed to build query input", err)
		return RespondError(reqCtx, err, errToStatusCode)
	}

	result, err := r.evaluator.SavedQuery(ctx, options, input)
	if err != nil {
		log.Error("failed to evaluated saved queryRevision", err)

		return RespondError(reqCtx, err, errToStatusCode)
	}

	debugEnabled := isDebugEnabled(input)
	response, err := MakeQueryResponse(result, debugEnabled)
	if err != nil {
		return RespondError(reqCtx, err, errToStatusCode)
	}

	return Respond(reqCtx, response.Body, response.StatusCode, response.Headers)
}

func makeQueryOptions(ctx *fasthttp.RequestCtx, log restql.Logger, envTenant string) (restql.QueryOptions, error) {
	namespace, err := pathParamString(ctx, "namespace")
	if err != nil {
		log.Error("failed to load namespace path param", err)
		return restql.QueryOptions{}, err
	}

	queryID, err := pathParamString(ctx, "queryId")
	if err != nil {
		log.Error("failed to load queryRevision id path param", err)
		return restql.QueryOptions{}, err
	}

	revisionStr, err := pathParamString(ctx, "revision")
	if err != nil {
		log.Error("failed to load revision path param", err)
		return restql.QueryOptions{}, err
	}

	revision, err := strconv.Atoi(revisionStr)
	if err != nil {
		log.Debug("failed to convert revision to integer")
		return restql.QueryOptions{}, errInvalidRevisionType
	}

	tenant, err := makeTenant(ctx, envTenant)
	if err != nil {
		return restql.QueryOptions{}, err
	}

	qo := restql.QueryOptions{
		Namespace: namespace,
		Id:        queryID,
		Revision:  revision,
		Tenant:    tenant,
	}

	return qo, nil
}

func makeTenant(ctx *fasthttp.RequestCtx, envTenant string) (string, error) {
	var tenant string

	if envTenant != "" {
		tenant = envTenant
	} else {
		tenant = string(ctx.QueryArgs().Peek("tenant"))
	}

	if tenant == "" {
		return "", errInvalidTenant
	}
	return tenant, nil
}

func makeQueryInput(ctx *fasthttp.RequestCtx, log restql.Logger) (restql.QueryInput, error) {
	params := make(map[string]interface{})
	ctx.Request.URI().QueryArgs().VisitAll(func(keyByte, valueByte []byte) {
		key := string(keyByte)
		value := string(valueByte)

		if currentValue, ok := params[key]; ok {
			var newValue []interface{}

			switch currentValue := currentValue.(type) {
			case []interface{}:
				newValue = append(currentValue, value)
			default:
				newValue = []interface{}{currentValue, value}
			}

			params[key] = newValue
		} else {
			params[key] = value
		}

	})

	headers := make(map[string]string)
	ctx.Request.Header.VisitAll(func(key, value []byte) {
		headers[string(key)] = string(value)
	})

	input := restql.QueryInput{
		Params:  params,
		Headers: headers,
	}

	contentType := string(ctx.Request.Header.ContentType())
	if contentType == jsonContentType {
		requestBody := ctx.Request.Body()
		if len(requestBody) > 0 {
			var b interface{}
			err := json.Unmarshal(requestBody, &b)
			if err != nil {
				log.Error("failed to unmarshal request body", err)
				return restql.QueryInput{}, fmt.Errorf("%w: %s", errFailedToReadRequestBody, err)
			}

			input.Body = b
		}
	}

	return input, nil
}

const debugParamName = "_debug"

func isDebugEnabled(queryInput restql.QueryInput) bool {
	param, found := queryInput.Params[debugParamName]
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
