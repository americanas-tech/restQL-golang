package logger

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/b2wdigital/restQL-golang/pkg/restql"
	"github.com/rs/zerolog"
)

var noOpLogger = New(ioutil.Discard, LogOptions{})

type LogOptions struct {
	Enable               bool
	TimestampFieldName   string
	TimestampFieldFormat string
	Level                string
	Format               string
}

func New(w io.Writer, options LogOptions) *Logger {
	output := w
	if options.Format == "pretty" {
		output = zerolog.ConsoleWriter{Out: w}
	}

	logger := zerolog.New(output)

	if len(options.TimestampFieldName) > 0 {
		logger = logger.With().Timestamp().Logger()
		zerolog.TimestampFieldName = options.TimestampFieldName
	}

	if len(options.TimestampFieldFormat) > 0 {
		zerolog.TimeFieldFormat = options.TimestampFieldFormat
	}

	level, err := zerolog.ParseLevel(options.Level)
	if err == nil {
		logger = logger.Level(level)
	}

	if !options.Enable {
		logger.Level(zerolog.Disabled)
	}

	return &Logger{zLogger: logger}
}

type Logger struct {
	zLogger zerolog.Logger
}

func (l *Logger) Panic(msg string, fields ...interface{}) {
	entry := l.zLogger.Panic()
	fieldMap := makeFieldMap(fields)

	entry.Fields(fieldMap).Msg(msg)
}

func (l *Logger) Fatal(msg string, fields ...interface{}) {
	fieldMap := makeFieldMap(fields)

	l.zLogger.Fatal().Fields(fieldMap).Msg(msg)
}

func (l *Logger) Error(msg string, err error, fields ...interface{}) {
	fieldMap := makeFieldMap(fields)

	l.zLogger.Error().Err(err).Fields(fieldMap).Msg(msg)
}

func (l *Logger) Warn(msg string, fields ...interface{}) {
	fieldMap := makeFieldMap(fields)

	l.zLogger.Warn().Fields(fieldMap).Msg(msg)
}

func (l *Logger) Info(msg string, fields ...interface{}) {
	fieldMap := makeFieldMap(fields)

	l.zLogger.Info().Fields(fieldMap).Msg(msg)
}

func (l *Logger) Debug(msg string, fields ...interface{}) {
	fieldMap := makeFieldMap(fields)

	l.zLogger.Debug().Fields(fieldMap).Msg(msg)
}

func (l *Logger) With(key string, value interface{}) restql.Logger {
	cl := l.zLogger.With().Str(key, fmt.Sprintf("%v", value)).Logger()
	return &Logger{zLogger: cl}
}

func makeFieldMap(fields []interface{}) map[string]interface{} {
	fieldMap := make(map[string]interface{})
	for i := 0; i <= len(fields)-2; i += 2 {
		key := fmt.Sprintf("%v", fields[i])
		value := fields[i+1]

		fieldMap[key] = value
	}
	return fieldMap
}
