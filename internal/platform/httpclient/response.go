package httpclient

import (
	"github.com/b2wdigital/restQL-golang/v6/pkg/restql"
	"github.com/valyala/fasthttp"
	"time"
)

func unmarshalBody(log restql.Logger, response *fasthttp.Response) (*restql.ResponseBody, error) {
	bodyByte := response.Body()
	bb := make([]byte, len(bodyByte))
	copy(bb, bodyByte)

	rb := restql.NewResponseBodyFromBytes(log, bb)
	if !rb.Valid() {
		return rb, errInvalidJSON
	}

	return rb, nil
}

func readHeaders(res *fasthttp.Response) restql.Headers {
	h := make(restql.Headers)
	res.Header.VisitAll(func(key, value []byte) {
		h[string(key)] = string(value)
	})

	return h
}

func makeErrorResponse(requestURL string, responseTime time.Duration, statusCode int) restql.HTTPResponse {
	return restql.HTTPResponse{
		URL:        requestURL,
		StatusCode: statusCode,
		Duration:   responseTime,
	}
}
