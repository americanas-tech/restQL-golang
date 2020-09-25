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

func (r *ResponseBody) GetBytes() []byte {
	return r.jsonBytes
}

func (r *ResponseBody) SetBytes(b []byte) {
	r.jsonBytes = b
}

func (r *ResponseBody) GetValue() interface{} {
	return r.jsonValue
}

func (r *ResponseBody) SetValue(v interface{}) {
	r.jsonValue = v
}

func (r *ResponseBody) Marshal() (interface{}, error) {
	if r.jsonValue != nil {
		b, err := json.Marshal(r.jsonValue)
		if err != nil {
			return nil, err
		}

		return json.RawMessage(b), nil
	}

	if len(r.jsonBytes) == 0 {
		return nil, nil
	}

	if !json.Valid(r.jsonBytes) {
		return string(r.jsonBytes), nil
	}

	return json.RawMessage(r.jsonBytes), nil
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

// ResourceCacheControlValue represents the values a cache control
// directive is able to have.
// It can either be present and have a integer time value or
// not be present in the upstream response.
type ResourceCacheControlValue struct {
	Exist bool
	Time  int
}

// ResourceCacheControl represent cache control directives
// returned by upstream during statement resolution.
type ResourceCacheControl struct {
	NoCache bool
	MaxAge  ResourceCacheControlValue
	SMaxAge ResourceCacheControlValue
}

// DoneResource represents a statement result.
type DoneResource struct {
	Status          int
	Success         bool
	IgnoreErrors    bool
	CacheControl    ResourceCacheControl
	Method          string
	URL             string
	RequestParams   map[string]interface{}
	RequestHeaders  map[string]string
	RequestBody     interface{}
	ResponseHeaders map[string]string
	ResponseBody    *ResponseBody
	ResponseTime    int64
}

// DoneResources represents a multiplexed statement result.
type DoneResources []interface{}
