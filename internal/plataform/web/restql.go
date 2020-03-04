package web

import (
	"bytes"
	"encoding/json"
	"github.com/b2wdigital/restQL-golang/internal/parser"
	"github.com/b2wdigital/restQL-golang/internal/plataform/conf"
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
	ctx.SetContentType("application/json")

	encoder := json.NewEncoder(ctx.Response.BodyWriter())

	queryTxt := bytes.NewBuffer(ctx.PostBody()).String()
	_, err := parser.Parse(queryTxt)
	if err != nil {
		log.Printf("[ERROR] an error ocurrend when parsing query : %v", err)

		ctx.SetStatusCode(http.StatusUnprocessableEntity)

		encoder.Encode(struct {
			Msg string
		}{"invalid query"})
		return
	}

	ctx.SetStatusCode(http.StatusOK)
	encoder.Encode(struct {
		Msg string
	}{"valid query"})
}
