package web

import (
	"fmt"
	"github.com/b2wdigital/restQL-golang/internal/eval"
	"github.com/b2wdigital/restQL-golang/internal/parser"
	"github.com/b2wdigital/restQL-golang/internal/plataform/conf"
	"github.com/b2wdigital/restQL-golang/internal/plataform/logger"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
)

var ErrInvalidNamespace = errors.New("invalid namespace")
var ErrInvalidQueryId = errors.New("invalid query id")
var ErrInvalidRevision = errors.New("invalid revision")
var ErrInvalidRevisionType = errors.New("invalid revision : must be an integer")

type RestQl struct {
	config    conf.Config
	log       *logger.Logger
	evaluator eval.Evaluator
}

func NewRestQl(c conf.Config, l *logger.Logger, e eval.Evaluator) RestQl {
	return RestQl{config: c, log: l, evaluator: e}
}

func (r RestQl) ValidateQuery(ctx *fasthttp.RequestCtx) error {
	queryTxt := string(ctx.PostBody())
	_, err := parser.Parse(queryTxt)
	if err != nil {
		r.log.Error("an error occurred when parsing query", err)

		e := &Error{
			Err:    errors.Wrap(err, "invalid query"),
			Status: http.StatusUnprocessableEntity,
		}

		return RespondError(ctx, e)
	}

	return Respond(ctx, nil, http.StatusOK)
}

func (r RestQl) RunSavedQuery(ctx *fasthttp.RequestCtx) error {
	options, err := r.makeQueryOptions(ctx)
	if err != nil {
		return RespondError(ctx, err)
	}

	input := r.makeQueryInput(ctx)

	query, err := r.evaluator.SavedQuery(options, input)
	if err != nil {
		return RespondError(ctx, err)
	}

	return Respond(ctx, query, http.StatusOK)
}

func (r RestQl) makeQueryOptions(ctx *fasthttp.RequestCtx) (eval.QueryOptions, error) {
	namespace, err := pathParamString(ctx, "namespace")
	if err != nil {
		r.log.Error("failed to load namespace path param", err)
		return eval.QueryOptions{}, err
	}

	queryId, err := pathParamString(ctx, "queryId")
	if err != nil {
		r.log.Error("failed to load query id path param", err)
		return eval.QueryOptions{}, err
	}

	revisionStr, err := pathParamString(ctx, "revision")
	if err != nil {
		r.log.Error("failed to load revision path param", err)
		return eval.QueryOptions{}, err
	}

	revision, err := strconv.Atoi(revisionStr)
	if err != nil {
		r.log.Debug("failed to convert revision to integer")
		e := &Error{Err: ErrInvalidRevisionType, Status: http.StatusBadRequest}

		return eval.QueryOptions{}, e
	}

	qo := eval.QueryOptions{
		Namespace: namespace,
		Id:        queryId,
		Revision:  revision,
	}

	return qo, nil
}

func (r RestQl) makeQueryInput(ctx *fasthttp.RequestCtx) eval.QueryInput {
	params := make(map[string]interface{})
	ctx.Request.URI().QueryArgs().VisitAll(func(keyByte, valueByte []byte) {
		key := string(keyByte)
		value := string(valueByte)

		if currentValue, ok := params[key]; ok {
			var newValue []string

			switch currentValue := currentValue.(type) {
			case []string:
				newValue = append(currentValue, value)
			default:
				newValue = []string{fmt.Sprintf("%v", currentValue), value}
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

	return eval.QueryInput{
		Params:  params,
		Headers: headers,
	}
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