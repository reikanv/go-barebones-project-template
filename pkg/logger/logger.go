package logger

import (
	"log/slog"
	"os"
)

func InitLogger(isDebug bool) {
	var logLevel slog.Level = slog.LevelError
	if isDebug {
		logLevel = slog.LevelDebug
	}
	handler := slog.Handler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(slog.New(handler))
}
