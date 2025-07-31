package utils

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func SetupZerologLogger() zerolog.Logger {
	// Configure zerolog JSON formatting
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.TimestampFieldName = "timestamp"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "message"
	zerolog.CallerFieldName = "caller"

	logger := zerolog.New(os.Stdout).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	return logger
}
