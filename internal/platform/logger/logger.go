package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"io"
)

type LogOptions struct {
	Enable    bool
	Timestamp bool
	Level     string
	Format    string
}

func New(w io.Writer, options LogOptions) *Logger {
	output := w
	if options.Format == "pretty" {
		output = zerolog.ConsoleWriter{Out: w}
	}

	logger := zerolog.New(output)

	if options.Timestamp {
		logger = logger.With().Timestamp().Logger()
		zerolog.TimestampFieldName = "timestamp"
	}

	level, err := zerolog.ParseLevel(options.Level)
	if err == nil {
		logger = logger.Level(level)
	}

	if !options.Enable {
		zerolog.SetGlobalLevel(zerolog.Disabled)
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

func makeFieldMap(fields []interface{}) map[string]interface{} {
	fieldMap := make(map[string]interface{})
	for i := 0; i <= len(fields)-2; i += 2 {
		key := fmt.Sprintf("%v", fields[i])
		value := fields[i+1]

		fieldMap[key] = value
	}
	return fieldMap
}
