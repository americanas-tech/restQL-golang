package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/b2wdigital/restQL-golang/internal/eval"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"time"
)

type HttpClient struct {
	client *fasthttp.Client
}

func New() HttpClient {
	c := &fasthttp.Client{
		NoDefaultUserAgentHeader: false,
		ReadTimeout:              3 * time.Second,
		WriteTimeout:             1 * time.Second,
	}

	return HttpClient{client: c}
}

func (hc HttpClient) Do(ctx context.Context, request eval.Request) (eval.Response, error) {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
	}()

	uri := fasthttp.URI{}
	uri.SetHost(request.Host)
	uri.SetQueryStringBytes(makeQueryArgs(request))

	req.SetRequestURI(uri.String())

	err := hc.client.Do(req, res)
	if err != nil {
		return eval.Response{}, errors.Wrap(err, "request execution failed")
	}

	buf := bytes.NewBuffer(res.Body())
	encoder := json.NewDecoder(buf)

	responseBody := make(map[string]interface{})
	err = encoder.Decode(&responseBody)
	if err != nil {
		return eval.Response{}, errors.Wrap(err, "response body decode failed")
	}

	response := eval.Response{
		StatusCode: res.StatusCode(),
		Body:       responseBody,
		Headers:    nil,
	}

	return response, nil
}

func makeQueryArgs(request eval.Request) []byte {
	var buf bytes.Buffer
	for key, value := range request.Query {
		buf.WriteString("&")
		buf.WriteString(key)
		buf.WriteString("=")
		buf.WriteString(value)
	}

	return buf.Bytes()
}
