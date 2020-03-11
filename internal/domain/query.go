package domain

type Query struct {
	Use        map[string]interface{}
	Statements []Statement
}

type Statement struct {
	Method       string
	Resource     string
	Alias        string
	Headers      map[string]interface{}
	Timeout      interface{}
	With         Params
	Only         []interface{}
	Hidden       bool
	CacheControl CacheControl
	IgnoreErrors bool
}

type Params map[string]interface{}

type CacheControl struct {
	MaxAge  interface{}
	SMaxAge interface{}
}

type Variable struct {
	Target string
}

type Chain []interface{}

type Flatten struct {
	Target interface{}
}

type Json struct {
	Target interface{}
}

type Base64 struct {
	Target interface{}
}

type Match struct {
	Target interface{}
	Arg    string
}

type QueryOptions struct {
	Namespace string
	Id        string
	Revision  int
}

type QueryInput struct {
	Params  map[string]interface{}
	Headers map[string]string
}
