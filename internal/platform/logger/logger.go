package logger

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/b2wdigital/restQL-golang/v6/pkg/restql"
	"github.com/rs/zerolog"
)

var noOpLogger = New(ioutil.Discard, LogOptions{})

// LogOptions represents available logging configurations
type LogOptions struct {
	Enable               bool
	TimestampFieldName   string
	TimestampFieldFormat string
	Level                string
	Format               string
}

type zeroLogger struct {
	zLogger zerolog.Logger
}

// New constructs a zeroLogger instance.
func New(w io.Writer, options LogOptions) restql.Logger {
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

	return &zeroLogger{zLogger: logger}
}

func (zl *zeroLogger) Panic(msg string, fields ...interface{}) {
	entry := zl.zLogger.Panic()
	fieldMap := makeFieldMap(fields)

	entry.Fields(fieldMap).Msg(msg)
}

func (zl *zeroLogger) Fatal(msg string, fields ...interface{}) {
	fieldMap := makeFieldMap(fields)

	zl.zLogger.Fatal().Fields(fieldMap).Msg(msg)
}

func (zl *zeroLogger) Error(msg string, err error, fields ...interface{}) {
	fieldMap := makeFieldMap(fields)

	zl.zLogger.Error().Err(err).Fields(fieldMap).Msg(msg)
}

func (zl *zeroLogger) Warn(msg string, fields ...interface{}) {
	fieldMap := makeFieldMap(fields)

	zl.zLogger.Warn().Fields(fieldMap).Msg(msg)
}

func (zl *zeroLogger) Info(msg string, fields ...interface{}) {
	fieldMap := makeFieldMap(fields)

	zl.zLogger.Info().Fields(fieldMap).Msg(msg)
}

func (zl *zeroLogger) Debug(msg string, fields ...interface{}) {
	fieldMap := makeFieldMap(fields)

	zl.zLogger.Debug().Fields(fieldMap).Msg(msg)
}

func (zl *zeroLogger) With(key string, value interface{}) restql.Logger {
	cl := zl.zLogger.With().Str(key, fmt.Sprintf("%v", value)).Logger()
	return &zeroLogger{zLogger: cl}
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
