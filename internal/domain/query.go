package domain

import (
	"fmt"
)

const (
	FromMethod   string = "from"
	ToMethod            = "to"
	IntoMethod          = "into"
	UpdateMethod        = "update"
	DeleteMethod        = "delete"
)

type Query struct {
	Use        Modifiers
	Statements []Statement
}

type Modifiers map[string]interface{}

type Statement struct {
	Method       string
	Resource     string
	Alias        string
	In           []string
	Headers      map[string]interface{}
	Timeout      interface{}
	With         Params
	Only         []interface{}
	Hidden       bool
	CacheControl CacheControl
	IgnoreErrors bool
}

type Params struct {
	Body   interface{}
	Values map[string]interface{}
}

type CacheControl struct {
	MaxAge  interface{}
	SMaxAge interface{}
}

type Variable struct {
	Target string
}

type Chain []interface{}

type ErrQueryRevisionDeprecated struct {
	Revision int
}

func (e ErrQueryRevisionDeprecated) Error() string {
	return fmt.Sprintf("the revision %d is deprecated", e.Revision)
}
