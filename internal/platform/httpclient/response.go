package httpclient

import (
	"encoding/json"
	"fmt"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"time"
)

func makeResponse(req *fasthttp.Request, res *fasthttp.Response, responseTime time.Duration) (domain.HttpResponse, error) {
	response := domain.HttpResponse{
		Url:        req.URI().String(),
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
		fmt.Printf("responde body: %s\n", string(res.Body()))
		return nil, errors.Wrapf(err, "failed to unmarshal response. error: %v. content: %s", err, string(res.Body()))
	}
	return responseBody, nil
}
