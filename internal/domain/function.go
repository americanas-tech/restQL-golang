package domain

// Function is the interface implemented by types that
// provide encoding, filters and special behaviour through
// the apply operator.
//
// Target returns the value upon which the function will be applied.
// Arguments return the arguments provided to the Function
// Argument fetches a Function argument by name
// SetArgument immutably updates the value of an argument by name
// Map uses a function to transform the target value preserving the current Function type wrapping it.
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

// Argument fetches a NoMultiplex argument by name
func (nm NoMultiplex) Argument(name string) Arg {
	return Arg{}
}

// SetArgument immutably updates the value of an argument by name
func (nm NoMultiplex) SetArgument(name string, value interface{}) Function {
	return nm
}

// Target return the value upon which NoMultiplex will be applied.
func (nm NoMultiplex) Target() interface{} {
	return nm.Value
}

// Arguments return the arguments provided to NoMultiplex function
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

// Argument fetches a JSON argument by name
func (j JSON) Argument(name string) Arg {
	return Arg{}
}

// SetArgument immutably updates the value of an argument by name
func (j JSON) SetArgument(name string, value interface{}) Function {
	return j
}

// Target return the value upon which JSON will be applied.
func (j JSON) Target() interface{} {
	return j.Value
}

// Argumetns return the arguments provided to JSON function
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

// Argument fetches a Base64 argument by name
func (b Base64) Argument(name string) Arg {
	return Arg{}
}

// SetArgument immutably updates the value of an argument by name
func (b Base64) SetArgument(name string, value interface{}) Function {
	return b
}

// Target return the value upon which Base64 will be applied.
func (b Base64) Target() interface{} {
	return b.Value
}

// Arguments return the arguments provided to Base64 function
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
	Args  []Arg
}

const MatchArgRegex = "regex"

// Argument fetches a Match argument by name
func (m Match) Argument(name string) Arg {
	if name == MatchArgRegex {
		return m.Args[0]
	}

	return Arg{}
}

// SetArgument immutably updates the value of an argument by name
func (m Match) SetArgument(name string, argValue interface{}) Function {
	if name == MatchArgRegex {
		return Match{Value: m.Value, Args: []Arg{{Name: MatchArgRegex, Value: argValue}}}
	}

	return m
}

// Target return the value upon which Match will be applied.
func (m Match) Target() interface{} {
	return m.Value
}

// Arguments return the arguments provided to Match function
func (m Match) Arguments() []Arg {
	return m.Args
}

// Map apply the given function to the Target value
// preserving the Match as a wrapper.
func (m Match) Map(fn func(target interface{}) interface{}) Function {
	return Match{Value: fn(m.Value), Args: m.Args}
}

// AsBody is a Function that define a `with`
// parameter as the request body for statements
// using to, into or patch methods.
type AsBody struct {
	Value interface{}
}

// Argument fetches a AsBody argument by name
func (ab AsBody) Argument(name string) Arg {
	return Arg{}
}

// SetArgument immutably updates the value of an argument by name
func (ab AsBody) SetArgument(name string, value interface{}) Function {
	return ab
}

// Target return the value upon which AsBody will be applied.
func (ab AsBody) Target() interface{} {
	return ab.Value
}

// Arguments return the arguments provided to AsBody function
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

// Argument fetches a Flatten argument by name
func (f Flatten) Argument(name string) Arg {
	return Arg{}
}

// SetArgument immutably updates the value of an argument by name
func (f Flatten) SetArgument(name string, value interface{}) Function {
	return f
}

// Target return the value upon which Flatten will be applied.
func (f Flatten) Target() interface{} {
	return f.Value
}

// Arguments return the arguments provided to Flatten function
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

// Argument fetches a NoExplode argument by name
func (f NoExplode) Argument(name string) Arg {
	return Arg{}
}

// SetArgument immutably updates the value of an argument by name
func (f NoExplode) SetArgument(name string, value interface{}) Function {
	return f
}

// Target return the value upon which NoExplode will be applied.
func (f NoExplode) Target() interface{} {
	return f.Value
}

// Arguments return the arguments provided to NoExplode function
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

// Arguments return the arguments provided to FilterByRegex function
func (f FilterByRegex) Arguments() []Arg {
	return f.Args
}

// Argument fetches a FilterByRegex argument by name
func (f FilterByRegex) Argument(name string) Arg {
	for _, arg := range f.Args {
		if arg.Name == name {
			return arg
		}
	}

	return Arg{}
}

// SetArgument immutably updates the value of an argument by name
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

// NoExplode is a Function that disable object explosion
// when resolving with values.
type AsQuery struct {
	Value interface{}
}

// Argument fetches a AsQuery argument by name
func (f AsQuery) Argument(name string) Arg {
	return Arg{}
}

// SetArgument immutably updates the value of an argument by name
func (f AsQuery) SetArgument(name string, value interface{}) Function {
	return f
}

// Target return the value upon which AsQuery will be applied.
func (f AsQuery) Target() interface{} {
	return f.Value
}

// Arguments return the arguments provided to AsQuery function
func (f AsQuery) Arguments() []Arg {
	return nil
}

// Map apply the given function to the Target value
// preserving the AsQuery as a wrapper.
func (f AsQuery) Map(fn func(target interface{}) interface{}) Function {
	return AsQuery{Value: fn(f.Value)}
}

// NoDuplicate is a Function that removes duplicate elements from lists.
type NoDuplicate struct {
	Value interface{}
}

// Argument fetches a NoDuplicate argument by name
func (b NoDuplicate) Argument(name string) Arg {
	return Arg{}
}

// SetArgument immutably updates the value of an argument by name
func (b NoDuplicate) SetArgument(name string, value interface{}) Function {
	return b
}

// Target return the value upon which NoDuplicate will be applied.
func (b NoDuplicate) Target() interface{} {
	return b.Value
}

// Arguments return the arguments provided to NoDuplicate function
func (b NoDuplicate) Arguments() []Arg {
	return nil
}

// Map apply the given function to the Target value
// preserving the NoDuplicate as a wrapper.
func (b NoDuplicate) Map(fn func(target interface{}) interface{}) Function {
	return NoDuplicate{Value: fn(b.Value)}
}
