package eval

type ValidationError struct {
	Err error
}

func (ve ValidationError) Error() string {
	return ve.Err.Error()
}

type NotFoundError struct {
	Err error
}

func (ne NotFoundError) Error() string {
	return ne.Err.Error()
}

type ParserError struct {
	Err error
}

func (pe ParserError) Error() string {
	return pe.Err.Error()
}

type QueryOptions struct {
	Namespace string
	Id        string
	Revision  int
}

type QueryInput struct {
	Params  map[string]interface{}
	Headers map[string]string
}

type EnvSource interface {
	GetString(key string) string
}

type FileSource interface {
	Unmarshal(target interface{}) error
}

type Configuration interface {
	Env() EnvSource
	File() FileSource
	Build() string
}

type Logger interface {
	Panic(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
	Error(msg string, err error, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
}
