package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/b2wdigital/restQL-golang/internal/parser"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
)

func New() fasthttp.RequestHandler {
	r := fasthttprouter.New()

	r.POST("/validate-query", validateQuery)

	r.GET("/health", func(ctx *fasthttp.RequestCtx) { ctx.Response.SetBodyString("I'm healthy! :)") })
	r.GET("/resource-status", func(ctx *fasthttp.RequestCtx) { ctx.Response.SetBodyString("Up and running! :)") })
	r.NotFound = func(ctx *fasthttp.RequestCtx) { ctx.Response.SetBodyString("There is nothing here. =/") }

	return r.Handler
}

func validateQuery(ctx *fasthttp.RequestCtx) {
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
