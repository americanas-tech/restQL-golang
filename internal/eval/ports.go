package eval

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
