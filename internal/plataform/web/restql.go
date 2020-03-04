package web

import (
	"bytes"
	"github.com/b2wdigital/restQL-golang/internal/parser"
	"github.com/b2wdigital/restQL-golang/internal/plataform/conf"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
)

type RestQl struct {
	config conf.Config
}

func NewRestQl(config conf.Config) RestQl {
	return RestQl{config: config}
}

func (r RestQl) validateQuery(ctx *fasthttp.RequestCtx) {
	queryTxt := bytes.NewBuffer(ctx.PostBody()).String()
	_, err := parser.Parse(queryTxt)
	if err != nil {
		log.Printf("[ERROR] an error ocurrend when parsing query : %v", err)

		e := &Error{
			Err:    errors.Wrap(err, "invalid query"),
			Status: http.StatusUnprocessableEntity,
		}

		RespondError(ctx, e)

		return
	}

	Respond(ctx, nil, http.StatusOK)
}
