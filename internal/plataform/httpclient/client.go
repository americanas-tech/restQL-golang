package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/plataform/logger"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"time"
)

type HttpClient struct {
	client *fasthttp.Client
	log    *logger.Logger
}

func New(log *logger.Logger) HttpClient {
	c := &fasthttp.Client{
		NoDefaultUserAgentHeader: false,
		ReadTimeout:              3 * time.Second,
		WriteTimeout:             1 * time.Second,
	}

	return HttpClient{client: c, log: log}
}

func (hc HttpClient) Do(ctx context.Context, request domain.Request) (domain.Response, error) {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
	}()

	uri := fasthttp.URI{DisablePathNormalizing: true}
	uri.SetScheme(request.Schema)
	uri.SetHost(request.Uri)
	uri.SetQueryStringBytes(makeQueryArgs(request))

	uriStr := uri.String()
	hc.log.Debug("request uri build", "uri", uriStr)
	req.SetRequestURI(uriStr)

	for key, value := range request.Headers {
		req.Header.Set(key, value)
	}

	err := hc.client.Do(req, res)
	if err != nil {
		return domain.Response{}, errors.Wrap(err, "request execution failed")
	}

	responseBody, err := hc.unmarshalBody(res)
	if err != nil {
		return domain.Response{}, err
	}

	headers := readHeaders(res)

	response := domain.Response{
		StatusCode: res.StatusCode(),
		Body:       responseBody,
		Headers:    headers,
	}

	return response, nil
}

func readHeaders(res *fasthttp.Response) domain.Headers {
	h := make(domain.Headers)
	res.Header.VisitAll(func(key, value []byte) {
		h[string(key)] = string(value)
	})

	return h
}

func (hc HttpClient) unmarshalBody(res *fasthttp.Response) (interface{}, error) {
	var responseBody interface{}
	err := json.Unmarshal(res.Body(), &responseBody)
	if err != nil {
		hc.log.Debug("failed to unmarshal response", "error", err, "response", string(res.Body()))
		return nil, errors.Wrap(err, "response body decode failed")
	}
	return responseBody, nil
}

var (
	ampersand = []byte("&")
	equal     = []byte("=")
)

func makeQueryArgs(request domain.Request) []byte {
	var buf bytes.Buffer
	for key, value := range request.Query {
		buf.Write(ampersand)
		buf.WriteString(key)
		buf.Write(equal)
		buf.WriteString(value)
	}

	return buf.Bytes()
}
