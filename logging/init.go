package logging

import (
	"log/slog"
	"os"

	"github.com/ThreeDotsLabs/watermill"
)

// Initializes a new josn logger and sets it as the default logger
func InitSlogLogger() *slog.Logger {
	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(l)

	return l
}

// InitWatermillLogger returns LoggerAdapter (of type SlogLoggerAdapter ) that corresponds to slog default logger.
func InitWatermillLogger() watermill.LoggerAdapter {
	slogLogger := InitSlogLogger()

	return watermill.NewSlogLogger(slogLogger)
}
