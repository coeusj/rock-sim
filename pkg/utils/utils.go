package utils

import (
	"flag"
	"log/slog"
	"os"
	"strings"
	"time"
)

func CreateLogger(logLevel string) *slog.Logger {
	var level slog.Level
	switch strings.ToLower(logLevel) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
}

func DateTimeWithNanoseconds(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.9999999")
}

func GetCliFlags() (string, string) {
	envFile := flag.String("env-file", "../../dev.env", "Environment file path")
	logLevel := flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	flag.Parse()

	return *envFile, *logLevel
}
