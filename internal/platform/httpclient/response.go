package httpclient

import (
	"encoding/json"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"time"
)

func makeResponse(requestUrl string, res *fasthttp.Response, responseTime time.Duration) (domain.HttpResponse, error) {
	response := domain.HttpResponse{
		Url:        requestUrl,
		StatusCode: res.StatusCode(),
		Headers:    readHeaders(res),
		Duration:   responseTime,
	}

	if len(res.Body()) > 0 {
		responseBody, err := unmarshalBody(res)
		if err != nil {
			return domain.HttpResponse{}, err
		}

		response.Body = responseBody
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

func makeErrorResponse(requestUrl string, responseTime time.Duration, err error) domain.HttpResponse {
	statusCode := 0
	if err == fasthttp.ErrTimeout {
		statusCode = 408
	}

	return domain.HttpResponse{
		Url:        requestUrl,
		StatusCode: statusCode,
		Duration:   responseTime,
	}
}
