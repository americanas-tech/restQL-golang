package web

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"net/http"
)

var (
	errInvalidNamespace  = errors.New("invalid namespace")
	errInvalidQueryID    = errors.New("invalid query id")
	errInvalidRevision   = errors.New("invalid revision")
	errInvalidTenantName = errors.New("invalid tenant name")
)

var pathParamNameToError = map[string]error{
	"namespace":  errInvalidNamespace,
	"query":      errInvalidQueryID,
	"revision":   errInvalidRevision,
	"tenantName": errInvalidTenantName,
}

func pathParamString(ctx *fasthttp.RequestCtx, name string) (string, error) {
	param, found := ctx.UserValue(name).(string)
	if !found {
		e, ok := pathParamNameToError[name]
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
