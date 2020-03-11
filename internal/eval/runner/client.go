package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"time"
)

type Request struct {
	Host    string
	Query   map[string]string
	Body    interface{}
	Headers map[string]string
}

type HttpClient struct {
	client *fasthttp.Client
}

func NewHttpClient() HttpClient {
	c := &fasthttp.Client{
		NoDefaultUserAgentHeader: false,
		ReadTimeout:              3 * time.Second,
		WriteTimeout:             1 * time.Second,
	}

	return HttpClient{client: c}
}

func (hc HttpClient) Do(ctx context.Context, request Request) (interface{}, error) {
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
		return nil, errors.Wrap(err, "request execution failed")
	}

	buf := bytes.NewBuffer(res.Body())
	encoder := json.NewDecoder(buf)

	responseBody := make(map[string]interface{})
	err = encoder.Decode(&responseBody)
	if err != nil {
		return nil, errors.Wrap(err, "response body decode failed")
	}

	return responseBody, nil
}

func makeQueryArgs(request Request) []byte {
	var buf bytes.Buffer
	for key, value := range request.Query {
		buf.WriteString("&")
		buf.WriteString(key)
		buf.WriteString("=")
		buf.WriteString(value)
	}

	return buf.Bytes()
}
