package httpclient

import (
	"encoding/json"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"time"
)

func makeResponse(req *fasthttp.Request, res *fasthttp.Response, responseTime time.Duration) (domain.HttpResponse, error) {
	responseBody, err := unmarshalBody(res)
	if err != nil {
		return domain.HttpResponse{}, err
	}

	headers := readHeaders(res)

	response := domain.HttpResponse{
		Url:        req.URI().String(),
		StatusCode: res.StatusCode(),
		Body:       responseBody,
		Headers:    headers,
		Duration:   responseTime,
	}
	return response, nil
}

func unmarshalBody(res *fasthttp.Response) (interface{}, error) {
	var responseBody interface{}
	err := json.Unmarshal(res.Body(), &responseBody)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal response. error: %v. content: %s", err, string(res.Body()))
	}
	return responseBody, nil
}
