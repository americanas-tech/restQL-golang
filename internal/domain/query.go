package domain

import (
	"fmt"
)

// Methods available to be used in query statements.
const (
	FromMethod   string = "from"
	ToMethod            = "to"
	IntoMethod          = "into"
	UpdateMethod        = "update"
	DeleteMethod        = "delete"
)

// Query is the internal representation of the restQL language.
type Query struct {
	Use        Modifiers
	Statements []Statement
}

// Modifiers is the internal representation of the `use` clause.
type Modifiers map[string]interface{}

// Statement is the internal representation of a query statement.
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

// Params is the internal representation of the `with` clause.
type Params struct {
	Body   interface{}
	Values map[string]interface{}
}

// CacheControl is the internal representation of the `max-age` and `s-max-age` clauses.
type CacheControl struct {
	MaxAge  interface{}
	SMaxAge interface{}
}

// Variable is the internal representation of a variable parameter value.
type Variable struct {
	Target string
}

// Chain is the internal representation of a chain parameter value.
type Chain []interface{}

// ErrQueryRevisionDeprecated represents an error from fetching
// a query marked as deprecated.
type ErrQueryRevisionDeprecated struct {
	Revision int
}

func (e ErrQueryRevisionDeprecated) Error() string {
	return fmt.Sprintf("the revision %d is deprecated", e.Revision)
}
