package web

import "github.com/valyala/fasthttp"

type Check struct{}

func NewCheck() Check {
	return Check{}
}

func (c Check) health(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetBodyString("I'm healthy! :)")
}

func (c Check) resourceStatus(ctx *fasthttp.RequestCtx) {
	ctx.Response.SetBodyString("Up and running! :)")
}
