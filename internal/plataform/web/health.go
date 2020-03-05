package web

import "github.com/valyala/fasthttp"

type Check struct{}

func NewCheck() Check {
	return Check{}
}

func (c Check) health(ctx *fasthttp.RequestCtx) error {
	ctx.Response.SetBodyString("I'm healthy! :)")
	return nil
}

func (c Check) resourceStatus(ctx *fasthttp.RequestCtx) error {
	ctx.Response.SetBodyString("Up and running! :)")
	return nil
}
