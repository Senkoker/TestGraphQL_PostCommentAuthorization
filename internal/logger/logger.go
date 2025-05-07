package logger

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func LoggerInit(loggerdebug bool) {
	if loggerdebug {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

func GetLogger() *slog.Logger {
	return logger
}
