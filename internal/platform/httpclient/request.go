package httpclient

import (
	"bytes"
	"encoding/json"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
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

func setupRequest(request restql.HTTPRequest, req *fasthttp.Request) error {
	uri := fasthttp.AcquireURI()
	defer func() {
		fasthttp.ReleaseURI(uri)
	}()
	uri.DisablePathNormalizing = true
	uri.SetScheme(request.Schema)
	uri.SetHost(request.Host)
	uri.SetPath(request.Path)
	uri.SetQueryStringBytes(makeQueryArgs(uri.QueryString(), request))

	req.SetRequestURIBytes(uri.FullURI())

	if request.Method == http.MethodPost || request.Method == http.MethodPut || request.Method == http.MethodPatch {
		var data []byte

		body := request.Body
		if strBody, ok := body.(string); ok {
			data = []byte(strBody)
		} else {
			jsonData, err := json.Marshal(body)
			if err != nil {
				return errors.Wrap(err, "failed to marshal request body")
			}

			data = jsonData
		}

		req.SetBody(data)
	}

	for key, value := range request.Headers {
		req.Header.Set(key, value)
	}

	req.Header.SetMethod(request.Method)
	return nil
}

func makeQueryArgs(queryArgs []byte, request restql.HTTPRequest) []byte {
	buf := bytes.NewBuffer(queryArgs)

	for key, value := range request.Query {
		appendQueryArg(buf, key, value)
	}

	return buf.Bytes()
}

func appendQueryArg(buf *bytes.Buffer, key string, value interface{}) {
	if value == nil {
		return
	}

	switch value := value.(type) {
	case string:
		appendStringParam(buf, key, value)
	case bool:
		appendStringParam(buf, key, strconv.FormatBool(value))
	case int:
		appendStringParam(buf, key, strconv.Itoa(value))
	case float64:
		appendStringParam(buf, key, strconv.FormatFloat(value, 'f', -1, 64))
	case map[string]interface{}:
		appendMapParam(buf, key, value)
	case []interface{}:
		for _, v := range value {
			appendQueryArg(buf, key, v)
		}
	}
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

func appendStringParam(buf *bytes.Buffer, key string, value string) {
	if value == "" {
		return
	}

	buf.Write(ampersand)
	buf.WriteString(key)
	buf.Write(equal)
	buf.WriteString(url.QueryEscape(value))
}
