package httpclient

import (
	"bytes"
	"encoding/json"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"net/http"
	"net/url"
	"strconv"
)

var (
	ampersand = []byte("&")
	equal     = []byte("=")
)

func setupRequest(request domain.HttpRequest, req *fasthttp.Request) error {
	uri := fasthttp.URI{DisablePathNormalizing: true}
	uri.SetScheme(request.Schema)
	uri.SetHost(request.Uri)
	uri.SetQueryStringBytes(makeQueryArgs(uri.QueryString(), request))

	uriStr := uri.String()
	req.SetRequestURI(uriStr)

	if request.Method == http.MethodPost || request.Method == http.MethodPut || request.Method == http.MethodPatch {
		data, err := json.Marshal(request.Body)
		if err != nil {
			return errors.Wrap(err, "failed to marshal request body")
		}

		req.SetBody(data)
	}

	for key, value := range request.Headers {
		req.Header.Set(key, value)
	}

	req.Header.SetMethod(request.Method)
	return nil
}

func readHeaders(res *fasthttp.Response) domain.Headers {
	h := make(domain.Headers)
	res.Header.VisitAll(func(key, value []byte) {
		h[string(key)] = string(value)
	})

	return h
}

func makeQueryArgs(queryArgs []byte, request domain.HttpRequest) []byte {
	buf := bytes.NewBuffer(queryArgs)

	for key, value := range request.Query {
		switch value := value.(type) {
		case string:
			appendStringParam(buf, key, value)
		case int:
			appendStringParam(buf, key, strconv.Itoa(value))
		case float64:
			appendStringParam(buf, key, strconv.FormatFloat(value, 'f', -1, 64))
		case map[string]interface{}:
			appendMapParam(buf, key, value)
		case []interface{}:
			appendListParam(buf, key, value)
		}
	}

	return buf.Bytes()
}

func appendMapParam(buf *bytes.Buffer, key string, value map[string]interface{}) {
	data, err := json.Marshal(value)
	if err != nil {
		return
	}

	buf.Write(ampersand)
	buf.WriteString(key)
	buf.Write(equal)
	buf.WriteString(url.QueryEscape(string(data)))
}

func appendListParam(buf *bytes.Buffer, key string, value []interface{}) {
	for _, v := range value {
		s, ok := v.(string)
		if !ok {
			continue
		}

		buf.Write(ampersand)
		buf.WriteString(key)
		buf.Write(equal)
		buf.WriteString(url.QueryEscape(s))
	}
}

func appendStringParam(buf *bytes.Buffer, key string, value string) {
	buf.Write(ampersand)
	buf.WriteString(key)
	buf.Write(equal)
	buf.WriteString(url.QueryEscape(value))
}
