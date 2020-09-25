package restql

import (
	"time"
)

// Body represents a HTTP body in a request or response.
type Body interface{}

// Headers represents all HTTP header in a request or response.
type Headers map[string]string

// HttpRequest represents a HTTP call to be
// made to an upstream dependency defined by the mappings.
type HTTPRequest struct {
	Method  string
	Schema  string
	Host    string
	Path    string
	Query   map[string]interface{}
	Body    Body
	Headers Headers
	Timeout time.Duration
}

// HttpResponse represents a HTTP call result
// from an upstream dependency defined by the mappings.
type HTTPResponse struct {
	URL        string
	StatusCode int
	Body       *ResponseBody
	Headers    Headers
	Duration   time.Duration
}

