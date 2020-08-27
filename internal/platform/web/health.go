package web

import (
	"fmt"
	"github.com/valyala/fasthttp"
)

type check struct {
	build string
}

func newCheck(build string) check {
	return check{build: build}
}

func (c check) Health(ctx *fasthttp.RequestCtx) error {
	ctx.Response.SetBodyString("I'm healthy! :)")
	return nil
}

func (c check) ResourceStatus(ctx *fasthttp.RequestCtx) error {
	ctx.Response.SetBodyString(fmt.Sprintf("RestQL is running with build %s", c.build))
	return nil
}
