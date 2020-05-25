package ast

const lexerDefinition = `
		Use = use
		Method = ^from\s+|^into\s+|^update\s+|^to\s+|^delete\s+
		As = as
		In = in\s+

		Headers = headers
		With = with
		Only = only
		Timeout = timeout
		Hidden = hidden
		IgnoreErrors = ignore-errors
		MaxAge = max-age
		CacheControl = cache-control
		SMaxAge = s-max-age

		Arrow = ->
		Flatten = flatten
		Base64 = base64
		Json = json
		Matches = matches 

		Colon = :
		LeftBrace = {
		RightBrace = }
		LeftBracket = \[
		RightBracket = \]
		LeftParentheses = \( 
		RightParentheses = \)
		Dollar = \$
		Dot = \.

		Float = [+-]?[0-9]+[.]{1}[0-9]+
		Int = [-+]?[0-9]+
		String = ".*?"
		Ident = [A-Za-z0-9_-]+

		Equal = =
		Comment = \/\/(.*)?
		whitespace = [,\s]+
`

const (
	FromMethod   = "from "
	IntoMethod   = "into "
	UpdateMethod = "update "
	ToMethod     = "to "
	DeleteMethod = "delete "
)

const (
	WithKeyword         = "with"
	OnlyKeyword         = "only"
	TimeoutKeyword      = "timeout"
	HeadersKeyword      = "headers"
	HiddenKeyword       = "hidden"
	MaxAgeKeyword       = "max-age"
	SmaxAgeKeyword      = "s-max-age"
	IgnoreErrorsKeyword = "ignore-errors"
)

type Query struct {
	Use    []Use   `(Use @@)*`
	Blocks []Block `@@*`
}

type Use struct {
	Key   string   `(@MaxAge | @CacheControl | @SMaxAge | @Timeout) Equal`
	Value UseValue `@@`
}

type UseValue struct {
	Int    *int    `@Int`
	String *string `| @String`
}

type Block struct {
	Method     string      `@Method`
	Resource   string      `@Ident`
	Alias      string      `(As @Ident)?`
	In         []string    `(In @Ident (Dot @Ident)*)?`
	Qualifiers []Qualifier `@@*`
}

type Qualifier struct {
	With         []WithItem    `(With @@+)`
	Only         []Filter      `| (Only @@+)`
	Hidden       bool          `| (@Hidden)`
	Timeout      *TimeoutValue `| (Timeout @@)`
	Headers      []HeaderItem  `| (Headers @@+)`
	MaxAge       *MaxAgeValue  `| (MaxAge Equal @@)`
	SMaxAge      *SMaxAgeValue `| (SMaxAge Equal @@)`
	IgnoreErrors bool          `| (@IgnoreErrors)`
}

type Filter struct {
	Field []string `@Ident (Dot @Ident)*`
	Match string   `(Arrow Matches LeftParentheses @String RightParentheses)?`
}

type WithItem struct {
	Key     string `@(Ident (Dot Ident)*) Equal`
	Value   Value  `@@ (Arrow (`
	Flatten bool   `	@Flatten |`
	Base64  bool   `	@Base64  |`
	Json    bool   `	@Json))?`
}

type Value struct {
	List      []*Value      `LeftBracket (@@)* RightBracket`
	Object    []ObjectEntry `| LeftBrace (@@)* RightBrace`
	Variable  *string       `| Dollar @Ident`
	Primitive *Primitive    `| @@`
}

type ObjectEntry struct {
	Key   string      `(@String | @Ident) Colon`
	Value ObjectValue `@@`
}

type ObjectValue struct {
	Nested    []ObjectEntry  `LeftBrace (@@)* RightBrace`
	List      []*ObjectValue `| LeftBracket (@@)* RightBracket`
	Variable  *string        `| Dollar @Ident`
	Primitive *Primitive     `| @@`
}

type Primitive struct {
	String *string   `@String`
	Int    *int      `| @Int`
	Float  *float64  `| @Float`
	Chain  []Chained `| @@ (Dot @@)*`
}

type Chained struct {
	PathVariable string `Dollar @Ident`
	PathItem     string `| @Ident`
}

type HeaderItem struct {
	Key   string      `@Ident Equal`
	Value HeaderValue `@@`
}

type HeaderValue struct {
	Variable *string `Dollar @Ident`
	String   *string `| @String`
}

type variableOrInt struct {
	Variable *string `Dollar @Ident`
	Int      *int    `| @Int`
}

type TimeoutValue variableOrInt
type MaxAgeValue variableOrInt
type SMaxAgeValue variableOrInt
