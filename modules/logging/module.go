package logging

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

const timeFormat = "2006-01-02T15:04:05.999Z07:00"

func NewLogger() zerolog.Logger {
	zerolog.TimeFieldFormat = timeFormat

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: timeFormat}

	var writers []io.Writer
	writers = append(writers, consoleWriter)

	log := zerolog.New(nil).With().Timestamp().Logger()
	logger := zerolog.Logger.Output(log, zerolog.MultiLevelWriter(writers...))

	fs, err := os.Create("/tmp/log")
	if err != nil {
		logger.Err(err)
		return logger
	}
	writers = append(writers, fs)

	log = zerolog.New(nil).With().Timestamp().Logger()
	logger = zerolog.Logger.Output(log, zerolog.MultiLevelWriter(writers...))
	return logger
}

// Injected from the result of NewLogger. Ptr so that it can be used in Fx and elsewhere easily
func NewPtrLogger(logger zerolog.Logger) *zerolog.Logger {
	return &logger
}

var Module = fx.Provide(NewLogger, NewPtrLogger)
