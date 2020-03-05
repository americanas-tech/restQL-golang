package web

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"net/http"
)

func Respond(ctx *fasthttp.RequestCtx, data interface{}, statusCode int) error {
	ctx.Response.Header.SetContentType("application/json; charset=utf-8")
	ctx.Response.SetStatusCode(statusCode)

	if data != nil {
		encoder := json.NewEncoder(ctx.Response.BodyWriter())
		if err := encoder.Encode(&data); err != nil {
			return err
		}
	}

	return nil
}

func RespondError(ctx *fasthttp.RequestCtx, err error) error {

	// If the error was of the type *Error, the handler has
	// a specific status code and error to return.
	if webErr, ok := err.(*Error); ok {
		er := ErrorResponse{
			Error: webErr.Err.Error(),
		}
		if err := Respond(ctx, er, webErr.Status); err != nil {
			return err
		}
		return nil
	}

	er := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}
	if err := Respond(ctx, er, http.StatusInternalServerError); err != nil {
		return err
	}
	return nil
}
