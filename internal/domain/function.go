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
	Arguments() []Arg
	Argument(string) Arg
	SetArgument(string, interface{}) Function
	Map(fn func(target interface{}) interface{}) Function
}

type Arg struct {
	Name  string
	Value interface{}
}

// NoMultiplex is a Function that disable request
// multiplexing of statements with list parameters.
type NoMultiplex struct {
	Value interface{}
}

func (nm NoMultiplex) Argument(name string) Arg {
	return Arg{}
}

func (nm NoMultiplex) SetArgument(name string, value interface{}) Function {
	return nm
}

// Target return the value upon which NoMultiplex will be applied.
func (nm NoMultiplex) Target() interface{} {
	return nm.Value
}

// Args return the arguments provided to NoMultiplex function
func (nm NoMultiplex) Arguments() []Arg {
	return nil
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

func (j JSON) Argument(name string) Arg {
	return Arg{}
}

func (j JSON) SetArgument(name string, value interface{}) Function {
	return j
}

// Target return the value upon which JSON will be applied.
func (j JSON) Target() interface{} {
	return j.Value
}

// Args return the arguments provided to JSON function
func (j JSON) Arguments() []Arg {
	return nil
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

func (b Base64) Argument(name string) Arg {
	return Arg{}
}

func (b Base64) SetArgument(name string, value interface{}) Function {
	return b
}

// Target return the value upon which Base64 will be applied.
func (b Base64) Target() interface{} {
	return b.Value
}

// Args return the arguments provided to Base64 function
func (b Base64) Arguments() []Arg {
	return nil
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

func (m Match) Argument(name string) Arg {
	if name == "regex" {
		return Arg{Name: "regex", Value: m.Arg}
	}

	return Arg{}
}

func (m Match) SetArgument(name string, argValue interface{}) Function {
	if name == "regex" {
		return Match{Value: m.Value, Arg: argValue}
	}

	return m
}

// Target return the value upon which Match will be applied.
func (m Match) Target() interface{} {
	return m.Value
}

// Args return the arguments provided to Match function
func (m Match) Arguments() []Arg {
	return []Arg{{Name: "regex", Value: m.Arg}}
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

func (ab AsBody) Argument(name string) Arg {
	return Arg{}
}

func (ab AsBody) SetArgument(name string, value interface{}) Function {
	return ab
}

// Target return the value upon which AsBody will be applied.
func (ab AsBody) Target() interface{} {
	return ab.Value
}

// Args return the arguments provided to AsBody function
func (ab AsBody) Arguments() []Arg {
	return nil
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

func (f Flatten) Argument(name string) Arg {
	return Arg{}
}

func (f Flatten) SetArgument(name string, value interface{}) Function {
	return f
}

// Target return the value upon which Flatten will be applied.
func (f Flatten) Target() interface{} {
	return f.Value
}

// Args return the arguments provided to Flatten function
func (f Flatten) Arguments() []Arg {
	return nil
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

func (f NoExplode) Argument(name string) Arg {
	return Arg{}
}

func (f NoExplode) SetArgument(name string, value interface{}) Function {
	return f
}

// Target return the value upon which NoExplode will be applied.
func (f NoExplode) Target() interface{} {
	return f.Value
}

// Args return the arguments provided to NoExplode function
func (f NoExplode) Arguments() []Arg {
	return nil
}

// Map apply the given function to the Target value
// preserving the NoExplode as a wrapper.
func (f NoExplode) Map(fn func(target interface{}) interface{}) Function {
	return NoExplode{Value: fn(f.Value)}
}

// NoExplode is a Function that disable object explosion
// when resolving with values.
type FilterByRegex struct {
	Value interface{}
	Args  []Arg
}

const (
	FilterByRegexArgRegex = "regex"
	FilterByRegexArgPath  = "path"
)

func NewFilterByRegex(target, path, regex interface{}) FilterByRegex {
	return FilterByRegex{
		Value: target,
		Args: []Arg{
			{Name: "path", Value: path},
			{Name: "regex", Value: regex},
		},
	}
}

// Target return the value upon which FilterByRegex will be applied.
func (f FilterByRegex) Target() interface{} {
	return f.Value
}

// Args return the arguments provided to FilterByRegex function
func (f FilterByRegex) Arguments() []Arg {
	return f.Args
}

func (f FilterByRegex) Argument(name string) Arg {
	for _, arg := range f.Args {
		if arg.Name == name {
			return arg
		}
	}

	return Arg{}
}

func (f FilterByRegex) SetArgument(name string, value interface{}) Function {
	newArg := Arg{Name: name, Value: value}
	args := f.Args

	for i, arg := range args {
		if arg.Name == name {
			args[i] = newArg
			return FilterByRegex{Value: f.Value, Args: args}
		}
	}

	args = append(args, newArg)
	return FilterByRegex{Value: f.Value, Args: args}
}

// Map apply the given function to the Target value
// preserving the FilterByRegex as a wrapper.
func (f FilterByRegex) Map(fn func(target interface{}) interface{}) Function {
	return FilterByRegex{Value: fn(f.Value), Args: f.Args}
}
