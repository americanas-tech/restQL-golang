package web

import (
	"fmt"
	"github.com/valyala/fasthttp"
)

type Check struct {
	build string
}

func NewCheck(build string) Check {
	return Check{build: build}
}

func (c Check) Health(ctx *fasthttp.RequestCtx) error {
	ctx.Response.SetBodyString("I'm healthy! :)")
	return nil
}

func (c Check) ResourceStatus(ctx *fasthttp.RequestCtx) error {
	ctx.Response.SetBodyString(fmt.Sprintf("RestQL is running with build %s", c.build))
	return nil
}
