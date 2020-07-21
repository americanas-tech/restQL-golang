package domain

import "regexp"

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

type AsBody struct {
	Value interface{}
}

func (ab AsBody) Target() interface{} {
	return ab.Value
}

func (ab AsBody) Map(fn func(target interface{}) interface{}) Function {
	return AsBody{Value: fn(ab.Value)}
}
