package ast

// restQL language keywords.
const (
	FromMethod          = "from"
	IntoMethod          = "into"
	UpdateMethod        = "update"
	ToMethod            = "to"
	DeleteMethod        = "delete"
	WithKeyword         = "with"
	OnlyKeyword         = "only"
	HeadersKeyword      = "headers"
	HiddenKeyword       = "hidden"
	TimeoutKeyword      = "timeout"
	MaxAgeKeyword       = "max-age"
	SmaxAgeKeyword      = "s-max-age"
	IgnoreErrorsKeyword = "ignore-errors"
	NoMultiplex         = "no-multiplex"
	Base64              = "base64"
	JSON                = "json"
	AsBody              = "as-body"
	Flatten             = "flatten"
)

// Query is the root of the restQL AST.
type Query struct {
	Use    []Use
	Blocks []Block
}

// Use is the syntax node representing the `use` clause.
type Use struct {
	Key   string
	Value UseValue
}

// UseValue is the syntax node representing
// the `use` clause possible values.
type UseValue struct {
	Int    *int
	String *string
}

// Block is the syntax node representing a statement.
type Block struct {
	Method     string
	Resource   string
	Alias      string
	In         []string
	Qualifiers []Qualifier
}

// Qualifier is the syntax node representing statement
// clauses: `with`, `only`, `hidden`, `headers`, `timeout`
// `max-age`, `s-max-age` and `ignore-errors`.
type Qualifier struct {
	With         *Parameters
	Only         []Filter
	Headers      []HeaderItem
	Hidden       bool
	Timeout      *TimeoutValue
	MaxAge       *MaxAgeValue
	SMaxAge      *SMaxAgeValue
	IgnoreErrors bool
}

// Filter is the syntax node representing entries
// in the `only` clause.
type Filter struct {
	Field []string
	Match *Match
}

// Match is the syntax node representing the
// `matches` function.
type Match struct {
	String   *string
	Variable *string
}

// Parameters is the syntax node representing
// the `with` clause.
type Parameters struct {
	Body      *ParameterBody
	KeyValues []KeyValue
}

// ParameterBody is the syntax node representing
// the dynamic body feature of the `with` clause.
type ParameterBody struct {
	Target    string
	Functions []string
}

// KeyValue is the syntax node representing
// parameters in the `with` clause.
type KeyValue struct {
	Key       string
	Value     Value
	Functions []string
}

// Value is the syntax node representing
// possible types used in the `with` clause
// parameters.
type Value struct {
	List      []Value
	Object    []ObjectEntry
	Variable  *string
	Primitive *Primitive
}

// ObjectEntry is the syntax node representing
// an object value.
type ObjectEntry struct {
	Key   string
	Value Value
}

// Primitive is the syntax node representing
// the basic restQL value types.
type Primitive struct {
	String  *string
	Int     *int
	Float   *float64
	Boolean *bool
	Chain   []Chained
	Null    bool
}

// Chained is the syntax node representing
// a chain value.
type Chained struct {
	PathVariable string
	PathItem     string
}

// HeaderItem is the syntax node representing
// an entry of the `headers` clause.
type HeaderItem struct {
	Key   string
	Value HeaderValue
}

// HeaderValue is the syntax node representing
// a `headers` clause entry value.
type HeaderValue struct {
	Variable *string
	String   *string
	Chain    []Chained
}

type variableOrInt struct {
	Variable *string
	Int      *int
}

// TimeoutValue is the syntax node representing
// the value in the `timeout` clause.
type TimeoutValue variableOrInt

// MaxAgeValue is the syntax node representing
// the value in the `max-age` clause.
type MaxAgeValue variableOrInt

// SMaxAgeValue is the syntax node representing
// the value in the `s-max-age` clause.
type SMaxAgeValue variableOrInt

// Generator encapsulate the parsing implementation
// used to transform a query string into an AST.
type Generator struct{}

// New constructs an AST generator.
func New() (Generator, error) {
	return Generator{}, nil
}

const noFilename = ""

// Parse transform a query string into an AST.
func (g Generator) Parse(query string) (*Query, error) {
	parse, err := Parse(noFilename, []byte(query))
	if err != nil {
		return nil, err
	}

	q := parse.(Query)
	return &q, nil
}
