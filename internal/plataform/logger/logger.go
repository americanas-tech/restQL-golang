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

type Logger interface {
	Panic(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
	Error(msg string, err error, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
}

func New(w io.Writer, config conf.Config) Logger {
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

	return &ZerologLogger{logger: logger}
}

type ZerologLogger struct {
	logger zerolog.Logger
}

func (z *ZerologLogger) Panic(msg string, fields ...interface{}) {
	entry := z.logger.Panic()
	fieldMap := makeFieldMap(fields)

	entry.Fields(fieldMap).Msg(msg)
}

func (z *ZerologLogger) Fatal(msg string, fields ...interface{}) {
	fieldMap := makeFieldMap(fields)

	z.logger.Fatal().Fields(fieldMap).Msg(msg)
}

func (z *ZerologLogger) Error(msg string, err error, fields ...interface{}) {
	fieldMap := makeFieldMap(fields)

	z.logger.Error().Err(err).Fields(fieldMap).Msg(msg)
}

func (z *ZerologLogger) Warn(msg string, fields ...interface{}) {
	fieldMap := makeFieldMap(fields)

	z.logger.Warn().Fields(fieldMap).Msg(msg)
}

func (z *ZerologLogger) Info(msg string, fields ...interface{}) {
	fieldMap := makeFieldMap(fields)

	z.logger.Info().Fields(fieldMap).Msg(msg)
}

func (z *ZerologLogger) Debug(msg string, fields ...interface{}) {
	fieldMap := makeFieldMap(fields)

	z.logger.Debug().Fields(fieldMap).Msg(msg)
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
