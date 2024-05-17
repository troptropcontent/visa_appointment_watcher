package logging

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Filename string
}

var Logger *zerolog.Logger

func Init(config Config) {
	logger := New(config)
	Logger = &logger
}

func New(config Config) zerolog.Logger {
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}
	outputs := []io.Writer{
		consoleWriter,
	}

	if config.Filename != "" {
		outputs = append(outputs, &lumberjack.Logger{
			Filename:   config.Filename,
			MaxSize:    5, //
			MaxBackups: 10,
			MaxAge:     14,
		})
	}

	output := zerolog.MultiLevelWriter(outputs...)

	return zerolog.New(output).With().Timestamp().Caller().Logger()
}
