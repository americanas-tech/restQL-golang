package domain

// Function is the interface implemented by types that
// provide encoding, filters and special behaviour through
// the apply operator.
//
// Target returns the value upon which the function will be applied.
// Map uses a function to transform the target value preserving the
// current Function type wrapping it.
type Function interface {
	Target() interface{}
	Map(fn func(target interface{}) interface{}) Function
}

// NoMultiplex is a Function that disable request
// multiplexing of statements with list parameters.
type NoMultiplex struct {
	Value interface{}
}

// Target return the value upon which NoMultiplex will be applied.
func (nm NoMultiplex) Target() interface{} {
	return nm.Value
}

// Map apply the given function to the Target value
// preserving the NoMultiplex as a wrapper.
func (nm NoMultiplex) Map(fn func(target interface{}) interface{}) Function {
	return NoMultiplex{Value: fn(nm.Value)}
}

// JSON is a Function that encode the target value as json.
type JSON struct {
	Value interface{}
}

// Target return the value upon which JSON will be applied.
func (j JSON) Target() interface{} {
	return j.Value
}

// Map apply the given function to the Target value
// preserving the JSON as a wrapper.
func (j JSON) Map(fn func(target interface{}) interface{}) Function {
	return JSON{Value: fn(j.Value)}
}

// Base64 is a Function that encode the target value as base64.
type Base64 struct {
	Value interface{}
}

// Target return the value upon which Base64 will be applied.
func (b Base64) Target() interface{} {
	return b.Value
}

// Map apply the given function to the Target value
// preserving the Base64 as a wrapper.
func (b Base64) Map(fn func(target interface{}) interface{}) Function {
	return Base64{Value: fn(b.Value)}
}

// Match is a Function that select values from the
// statement result based on the given Arg.
type Match struct {
	Value interface{}
	Arg   interface{}
}

// Target return the value upon which Match will be applied.
func (m Match) Target() interface{} {
	return m.Value
}

// Map apply the given function to the Target value
// preserving the Match as a wrapper.
func (m Match) Map(fn func(target interface{}) interface{}) Function {
	return Match{Value: fn(m.Value), Arg: m.Arg}
}

// AsBody is a Function that define a `with`
// parameter as the request body for statements
// using to, into or patch methods.
type AsBody struct {
	Value interface{}
}

// Target return the value upon which AsBody will be applied.
func (ab AsBody) Target() interface{} {
	return ab.Value
}

// Map apply the given function to the Target value
// preserving the AsBody as a wrapper.
func (ab AsBody) Map(fn func(target interface{}) interface{}) Function {
	return AsBody{Value: fn(ab.Value)}
}

// Flatten is a Function that encode the target value
// as a plain list of value.
type Flatten struct {
	Value interface{}
}

// Target return the value upon which Flatten will be applied.
func (f Flatten) Target() interface{} {
	return f.Value
}

// Map apply the given function to the Target value
// preserving the Flatten as a wrapper.
func (f Flatten) Map(fn func(target interface{}) interface{}) Function {
	return Flatten{Value: fn(f.Value)}
}

// NoExplode is a Function that disable object explosion
// when resolving with values.
type NoExplode struct {
	Value interface{}
}

// Target return the value upon which NoExplode will be applied.
func (f NoExplode) Target() interface{} {
	return f.Value
}

// Map apply the given function to the Target value
// preserving the NoExplode as a wrapper.
func (f NoExplode) Map(fn func(target interface{}) interface{}) Function {
	return NoExplode{Value: fn(f.Value)}
}
