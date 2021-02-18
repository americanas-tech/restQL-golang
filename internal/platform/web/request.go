package web

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"net/url"
)

var errPathParamNotFound = errors.New("invalid path param, not found")

func pathParamString(ctx *fasthttp.RequestCtx, name string) (string, error) {
	param, found := ctx.UserValue(name).(string)
	if !found {
		return "", fmt.Errorf("%w: %s", errPathParamNotFound, name)
	}

	unescaped, err := url.QueryUnescape(param)
	if err != nil {
		return param, nil
	}

	return unescaped, nil
}
