package ast

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
	Json                = "json"
	AsBody              = "as-body"
)

type Query struct {
	Use    []Use
	Blocks []Block
}

type Use struct {
	Key   string
	Value UseValue
}

type UseValue struct {
	Int    *int
	String *string
}

type Block struct {
	Method     string
	Resource   string
	Alias      string
	In         []string
	Qualifiers []Qualifier
}

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

type Filter struct {
	Field []string
	Match string
}

type Parameters struct {
	Body      *ParameterBody
	KeyValues []KeyValue
}

type ParameterBody struct {
	Target    string
	Functions []string
}

type KeyValue struct {
	Key       string
	Value     Value
	Functions []string
}

type Value struct {
	List      []Value
	Object    []ObjectEntry
	Variable  *string
	Primitive *Primitive
}

type ObjectEntry struct {
	Key   string
	Value Value
}

type Primitive struct {
	String  *string
	Int     *int
	Float   *float64
	Boolean *bool
	Chain   []Chained
	Null    bool
}

type Chained struct {
	PathVariable string
	PathItem     string
}

type HeaderItem struct {
	Key   string
	Value HeaderValue
}

type HeaderValue struct {
	Variable *string
	String   *string
	Chain    []Chained
}

type variableOrInt struct {
	Variable *string
	Int      *int
}

type TimeoutValue variableOrInt
type MaxAgeValue variableOrInt
type SMaxAgeValue variableOrInt

type Generator struct {
}

func New() (Generator, error) {
	return Generator{}, nil
}

const noFilename = ""

func (g Generator) Parse(query string) (*Query, error) {
	parse, err := Parse(noFilename, []byte(query))
	if err != nil {
		return nil, err
	}

	q := parse.(Query)
	return &q, nil
}
