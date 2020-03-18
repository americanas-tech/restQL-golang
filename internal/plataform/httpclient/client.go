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

func (hc HttpClient) Do(ctx context.Context, request domain.HttpRequest) (domain.HttpResponse, error) {
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
		return domain.HttpResponse{}, errors.Wrap(err, "request execution failed")
	}

	responseBody, err := hc.unmarshalBody(res)
	if err != nil {
		return domain.HttpResponse{}, err
	}

	headers := readHeaders(res)

	response := domain.HttpResponse{
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

func makeQueryArgs(request domain.HttpRequest) []byte {
	var buf bytes.Buffer
	for key, value := range request.Query {
		switch value := value.(type) {
		case string:
			appendStringParam(buf, key, value)
		case []interface{}:
			appendListParam(buf, key, value)
		}
	}

	return buf.Bytes()
}

func appendListParam(buf bytes.Buffer, key string, value []interface{}) {
	for _, v := range value {
		s, ok := v.(string)
		if !ok {
			continue
		}

		buf.Write(ampersand)
		buf.WriteString(key)
		buf.Write(equal)
		buf.WriteString(s)
	}
}

func appendStringParam(buf bytes.Buffer, key string, value string) {
	buf.Write(ampersand)
	buf.WriteString(key)
	buf.Write(equal)
	buf.WriteString(value)
}
