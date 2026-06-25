package logger

import (
	"log/slog"
	"os"
	"sync"
)

var (
	globLogger Logger
	once       sync.Once
)

// InitGlobalLogger init global logger.
func InitGlobalLogger() {
	once.Do(func() {
		globLogger = NewSlog(WithTextHandler(os.Stdout, slog.LevelDebug))
	})
}

// Debug logs a debug-level message via the global logger.
func Debug(msg string, args ...any) {
	log(msg, slog.LevelDebug, args...)
}

// Info logs an info-level message via the global logger.
func Info(msg string, args ...any) {
	log(msg, slog.LevelInfo, args...)
}

// Warn logs a warning-level message via the global logger.
func Warn(msg string, args ...any) {
	log(msg, slog.LevelWarn, args...)
}

// Error logs an error-level message via the global logger.
func Error(msg string, args ...any) {
	log(msg, slog.LevelError, args...)
}

// Fatal logs an error-level message via the global logger and then exits.
func Fatal(msg string, args ...any) {
	log(msg, slog.LevelError, args...)
	//nolint:revive // calls to os.Exit is acceptable here.
	os.Exit(1)
}

func log(msg string, level slog.Level, args ...any) {
	if globLogger == nil {
		InitGlobalLogger()
	}

	switch level {
	case slog.LevelDebug:
		globLogger.Debug(msg, args...)
	case slog.LevelInfo:
		globLogger.Info(msg, args...)
	case slog.LevelWarn:
		globLogger.Warn(msg, args...)
	case slog.LevelError:
		globLogger.Error(msg, args...)
	}
}
