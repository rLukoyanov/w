package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Setup(level string) error {
	if level == "" {
		level = "info"
	}

	parsedLevel, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		return fmt.Errorf("invalid log level %q: %w", level, err)
	}

	zerolog.SetGlobalLevel(parsedLevel)
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	// Configure goose to use zerolog
	goose.SetLogger(&GooseLogger{})

	return nil
}

// GooseLogger adapts zerolog to goose's Logger interface
type GooseLogger struct{}

func (l *GooseLogger) Fatalf(format string, v ...interface{}) {
	log.Fatal().Msgf(format, v...)
}

func (l *GooseLogger) Printf(format string, v ...interface{}) {
	log.Info().Msgf(format, v...)
}
