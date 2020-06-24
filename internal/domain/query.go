package domain

import (
	"fmt"
	"regexp"
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

type Function interface {
	Target() interface{}
	Map(fn func(target interface{}) interface{}) Function
}

type NoMultiplex struct {
	Value interface{}
}

func (nm NoMultiplex) Target() interface{} {
	return nm.Value
}

func (nm NoMultiplex) Map(fn func(target interface{}) interface{}) Function {
	return NoMultiplex{Value: fn(nm.Value)}
}

type Json struct {
	Value interface{}
}

func (j Json) Target() interface{} {
	return j.Value
}

func (j Json) Map(fn func(target interface{}) interface{}) Function {
	return Json{Value: fn(j.Value)}
}

type Base64 struct {
	Value interface{}
}

func (b Base64) Target() interface{} {
	return b.Value
}

func (b Base64) Map(fn func(target interface{}) interface{}) Function {
	return Base64{Value: fn(b.Value)}
}

type Match struct {
	Value interface{}
	Arg   *regexp.Regexp
}

func (m Match) Target() interface{} {
	return m.Value
}

func (m Match) Map(fn func(target interface{}) interface{}) Function {
	return Match{Value: fn(m.Value), Arg: m.Arg}
}

type QueryOptions struct {
	Namespace string
	Id        string
	Revision  int
	Tenant    string
}

type QueryInput struct {
	Params  map[string]interface{}
	Body    interface{}
	Headers map[string]string
}

type QueryContext struct {
	Mappings map[string]Mapping
	Options  QueryOptions
	Input    QueryInput
}

type SavedQuery struct {
	Text       string
	Deprecated bool
}

type ErrQueryRevisionDeprecated struct {
	Revision int
}

func (e ErrQueryRevisionDeprecated) Error() string {
	return fmt.Sprintf("the revision %d is deprecated", e.Revision)
}
