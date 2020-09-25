package restql

import (
	"encoding/json"
)

// ResponseBody is a wrapper that allows restQL to defer JSON parsing
// the HTTP body of an upstream response.
//
// Internally it stores two values: a byte slice and a interface. The most
// common scenario will be of creating a ResponseBody from the byte slice red
// from the HTTP response and returning it to downstream.
//
// When features need to access the content of the JSON response, they can
// use the Unmarshal method to get it.
//
// If the byte slice is unmarshalled or a new value is set on the response body
// with SetValue method, then the Marshal and Unmarshal function will operate
// using this value rather then the byte slice.
type ResponseBody struct {
	log Logger
	jsonBytes []byte
	jsonValue interface{}
}

// NewResponseBodyFromBytes creates a ResponseBody wrapper from
// an HTTP response data.
func NewResponseBodyFromBytes(log Logger, b []byte) *ResponseBody {
	r := &ResponseBody{log: log, jsonBytes: b}

	return r
}

// NewResponseBodyFromValue creates a ResponseBody wrapper from
// a generic value, usually from an error.
func NewResponseBodyFromValue(log Logger, v interface{}) *ResponseBody {
	r := &ResponseBody{log: log}

	r.SetValue(v)
	return r
}

// Bytes return the bytes data wrapped.
func (r *ResponseBody) Bytes() []byte {
	return r.jsonBytes
}

// Value return the generic data wrapped.
func (r *ResponseBody) Value() interface{} {
	return r.jsonValue
}

// SetValue defines a generic value to replace
// the byte slice data.
func (r *ResponseBody) SetValue(v interface{}) {
	r.jsonValue = v
}

// Marshal returns the content of ResponseBody ready to
// be sent to downstream.
//
// This method can process the content in 4 ways:
// - If there is a generic data, marshal it using a JSON parser
//   and return the result as a json.RawMessage.
// - Else, if the byte slice is empty, return nil.
// - Else, if the byte slice is not an valid JSON, stringify it.
// - Finally, if the byte slice is not empty and is a valid json,
//   return it as a json.RawMessage.
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

// Unmarshal returns the content of ResponseBody ready to
// be manipulated internally by restQL.
//
// This method can process the content in 4 ways:
// - If there is a generic data, return it.
// - Else, if the byte slice is empty or is not a valid json,
//   return it as a string.
// - Finally, if it is valid to be manipulated, then unmarshal it
//   and return.
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

// Valid return true if the ResponseBody content
// can be manipulated by restQL.
func (r *ResponseBody) Valid() bool {
	if r.jsonValue != nil {
		return true
	}

	return len(r.jsonBytes) > 0 && json.Valid(r.jsonBytes)
}

// Clear removes all internal content.
func (r *ResponseBody) Clear() {
	r.jsonBytes = nil
	r.jsonValue = nil
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
