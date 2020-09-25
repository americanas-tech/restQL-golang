package restql

import (
	"encoding/json"
)

type ResponseBody struct {
	log Logger
	jsonBytes []byte
	jsonValue interface{}
}

func NewResponseBodyFromBytes(log Logger, b []byte) *ResponseBody {
	r := &ResponseBody{log: log}

	r.SetBytes(b)
	return r
}

func NewResponseBodyFromValue(log Logger, v interface{}) *ResponseBody {
	r := &ResponseBody{log: log}

	r.SetValue(v)
	return r
}

func (r *ResponseBody) SetBytes(b []byte) {
	r.jsonBytes = b
}

func (r *ResponseBody) SetValue(v interface{}) {
	r.jsonValue = v
}

func (r *ResponseBody) Marshal() ([]byte, error) {
	if r.jsonValue != nil {
		return json.Marshal(r.jsonValue)
	}

	return r.jsonBytes, nil
}

func (r *ResponseBody) Unmarshal() interface{} {
	if r.jsonValue != nil {
		return r.jsonValue
	}

	bodyByte := r.jsonBytes
	if !r.Valid() {
		return string(bodyByte)
	}

	var responseBody interface{}
	err := json.Unmarshal(bodyByte, &responseBody)
	if err != nil {
		body := string(bodyByte)
		r.log.Error("failed to unmarshal response body", err, "body", body)

		return body
	}

	r.jsonValue = responseBody
	return responseBody
}

func (r *ResponseBody) Valid() bool {
	if r.jsonValue != nil {
		return true
	}

	return len(r.jsonBytes) > 0 && json.Valid(r.jsonBytes)
}