package web

import (
	"bytes"
	"github.com/b2wdigital/restQL-golang/internal/parser"
	"github.com/b2wdigital/restQL-golang/internal/plataform/conf"
	"github.com/b2wdigital/restQL-golang/internal/plataform/logger"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"net/http"
)

type RestQl struct {
	config conf.Config
	log    logger.Logger
}

func NewRestQl(config conf.Config, log logger.Logger) RestQl {
	return RestQl{config: config, log: log}
}

func (r RestQl) validateQuery(ctx *fasthttp.RequestCtx) error {
	queryTxt := bytes.NewBuffer(ctx.PostBody()).String()
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
