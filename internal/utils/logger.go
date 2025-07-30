package utils

import (
	"github.com/rs/zerolog"
	"os"
	"time"
)

func SetupZerologLogger() zerolog.Logger {
	// Configure zerolog
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.TimestampFieldName = "timestamp"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "message"

	// Create logger with JSON output to stdout
	logger := zerolog.New(os.Stdout).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Logger()
	return logger
}
