package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
	"github.com/b2wdigital/restQL-golang/v4/internal/eval"
	"github.com/b2wdigital/restQL-golang/v4/internal/parser"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/conf"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/logger"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/web/middleware"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

var jsonContentType = "application/json"

var (
	ErrInvalidNamespace    = errors.New("invalid namespace")
	ErrInvalidQueryId      = errors.New("invalid query id")
	ErrInvalidRevision     = errors.New("invalid revision")
	ErrInvalidRevisionType = errors.New("invalid revision : must be an integer")
	ErrInvalidTenant       = errors.New("invalid tenant : no value provided")
)

type RestQl struct {
	config    *conf.Config
	log       *logger.Logger
	evaluator eval.Evaluator
	parser    parser.Parser
}

func NewRestQl(l *logger.Logger, cfg *conf.Config, e eval.Evaluator, p parser.Parser) RestQl {
	return RestQl{config: cfg, log: l, evaluator: e, parser: p}
}

func (r RestQl) ValidateQuery(ctx *fasthttp.RequestCtx) error {
	queryTxt := string(ctx.PostBody())
	_, err := r.parser.Parse(queryTxt)
	if err != nil {
		r.log.Error("an error occurred when parsing query", err)

		e := &Error{
			Err:    errors.Wrap(err, "invalid query"),
			Status: http.StatusUnprocessableEntity,
		}

		return RespondError(ctx, e)
	}

	return Respond(ctx, nil, http.StatusOK, nil)
}

func (r RestQl) RunAdHocQuery(ctx *fasthttp.RequestCtx) error {
	context := middleware.GetNativeContext(ctx)
	context = restql.WithLogger(ctx, r.log)

	tenant, err := makeTenant(ctx, r.config.Tenant)
	if err != nil {
		r.log.Error("failed to build query options", err)
		return RespondError(ctx, NewRequestError(err, http.StatusBadRequest))
	}
	options := domain.QueryOptions{Tenant: tenant}

	input, err := makeQueryInput(ctx, r.log)
	if err != nil {
		r.log.Error("failed to build query input", err)
		return RespondError(ctx, NewRequestError(err, http.StatusBadRequest))
	}

	queryTxt := string(ctx.PostBody())

	result, err := r.evaluator.AdHocQuery(context, queryTxt, options, input)
	if err != nil {
		r.log.Error("failed to evaluated adhoc query", err)

		switch err := err.(type) {
		case eval.ValidationError:
			return RespondError(ctx, NewRequestError(err, http.StatusUnprocessableEntity))
		case eval.NotFoundError:
			return RespondError(ctx, NewRequestError(err, http.StatusNotFound))
		case eval.ParserError:
			return RespondError(ctx, NewRequestError(err, http.StatusBadRequest))
		case eval.TimeoutError:
			return RespondError(ctx, NewRequestError(err, http.StatusRequestTimeout))
		default:
			return RespondError(ctx, err)
		}
	}

	debugEnabled := isDebugEnabled(input)
	response := MakeQueryResponse(result, debugEnabled)
	return Respond(ctx, response.Body, response.StatusCode, response.Headers)
}

func (r RestQl) RunSavedQuery(ctx *fasthttp.RequestCtx) error {
	log := r.log.With("restql-endpoint", string(ctx.Request.URI().Path()))
	log = log.With("request-id", string(ctx.Request.Header.Peek("X-TID")))

	context := middleware.GetNativeContext(ctx)
	context = restql.WithLogger(context, log)

	options, err := makeQueryOptions(ctx, log, r.config.Tenant)
	if err != nil {
		log.Error("failed to build query options", err, "query", ctx.RequestURI())
		return RespondError(ctx, NewRequestError(err, http.StatusBadRequest))
	}
	queryIdentifier := fmt.Sprintf("%s/%s/%d", options.Namespace, options.Id, options.Revision)

	input, err := makeQueryInput(ctx, log)
	if err != nil {
		log.Error("failed to build query input", err, "query", queryIdentifier)
		return RespondError(ctx, NewRequestError(err, http.StatusBadRequest))
	}

	result, err := r.evaluator.SavedQuery(context, options, input)
	if err != nil {
		log.Error("failed to evaluated saved query", err, "query", queryIdentifier)

		switch err := err.(type) {
		case eval.ValidationError:
			return RespondError(ctx, NewRequestError(err, http.StatusUnprocessableEntity))
		case eval.NotFoundError:
			return RespondError(ctx, NewRequestError(err, http.StatusNotFound))
		case eval.TimeoutError:
			return RespondError(ctx, NewRequestError(err, http.StatusRequestTimeout))
		case eval.ParserError:
			return RespondError(ctx, NewRequestError(err, http.StatusInternalServerError))
		case eval.MappingError:
			return RespondError(ctx, NewRequestError(err, http.StatusInternalServerError))
		case domain.ErrQueryRevisionDeprecated:
			return RespondError(ctx, NewRequestError(err, http.StatusBadRequest))
		default:
			return RespondError(ctx, err)
		}
	}

	debugEnabled := isDebugEnabled(input)
	response := MakeQueryResponse(result, debugEnabled)
	return Respond(ctx, response.Body, response.StatusCode, response.Headers)
}

func makeQueryOptions(ctx *fasthttp.RequestCtx, log restql.Logger, envTenant string) (domain.QueryOptions, error) {
	namespace, err := pathParamString(ctx, "namespace")
	if err != nil {
		log.Error("failed to load namespace path param", err)
		return domain.QueryOptions{}, err
	}

	queryId, err := pathParamString(ctx, "queryId")
	if err != nil {
		log.Error("failed to load query id path param", err)
		return domain.QueryOptions{}, err
	}

	revisionStr, err := pathParamString(ctx, "revision")
	if err != nil {
		log.Error("failed to load revision path param", err)
		return domain.QueryOptions{}, err
	}

	revision, err := strconv.Atoi(revisionStr)
	if err != nil {
		log.Debug("failed to convert revision to integer")
		return domain.QueryOptions{}, ErrInvalidRevisionType
	}

	tenant, err := makeTenant(ctx, envTenant)
	if err != nil {
		return domain.QueryOptions{}, err
	}

	qo := domain.QueryOptions{
		Namespace: namespace,
		Id:        queryId,
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
		return "", ErrInvalidTenant
	}
	return tenant, nil
}

func makeQueryInput(ctx *fasthttp.RequestCtx, log restql.Logger) (domain.QueryInput, error) {
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

	input := domain.QueryInput{
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
				return domain.QueryInput{}, err
			}

			input.Body = b
		}
	}

	return input, nil
}

var paramNameToError = map[string]error{
	"namespace": ErrInvalidNamespace,
	"query":     ErrInvalidQueryId,
	"revision":  ErrInvalidRevision,
}

func pathParamString(ctx *fasthttp.RequestCtx, name string) (string, error) {
	param, found := ctx.UserValue(name).(string)
	if !found {
		e, ok := paramNameToError[name]
		if !ok {
			e = errors.New(fmt.Sprintf("path param not found : %s", name))
		}

		return "", &Error{
			Err:    e,
			Status: http.StatusUnprocessableEntity,
		}
	}

	return param, nil
}

const debugParamName = "_debug"

func isDebugEnabled(queryInput domain.QueryInput) bool {
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
