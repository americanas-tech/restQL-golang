package logger

import (
	"fmt"
	"github.com/b2wdigital/restQL-golang/internal/plataform/conf"
	"github.com/rs/zerolog"
	"io"
)

type logConf struct {
	Enable    bool   `yaml:"enable"`
	Timestamp bool   `yaml:"timestamp"`
	Level     string `yaml:"level"`
	Format    string `yaml:"format"`
}

type configFile struct {
	Logging logConf `yaml:"logging"`
}

var defaultConfig = configFile{Logging: logConf{Enable: true, Timestamp: true, Level: "info", Format: "json"}}

func New(w io.Writer, config conf.Config) *Logger {
	lc := defaultConfig
	config.File().Unmarshal(&lc)

	output := w
	if lc.Logging.Format == "pretty" {
		output = zerolog.ConsoleWriter{Out: w}
	}

	logger := zerolog.New(output)

	if lc.Logging.Timestamp {
		logger = logger.With().Timestamp().Logger()
		zerolog.TimestampFieldName = "timestamp"
	}

	level, err := zerolog.ParseLevel(lc.Logging.Level)
	if err != nil {
		logger = logger.Level(level)
	}

	if !lc.Logging.Enable {
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
	for i := 0; i < len(fields); i += 2 {
		key := fmt.Sprintf("%v", fields[i])
		value := fields[i+1]

		fieldMap[key] = value
	}
	return fieldMap
}
