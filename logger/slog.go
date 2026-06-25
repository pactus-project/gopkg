package logger

import (
	"io"
	"log/slog"
	"os"
)

// Slog is a slog-based implementation of the Slogger interface.
type Slog struct {
	log *slog.Logger
}

// SlogHandler is a factory function that returns a configured slog.Logger.
type SlogHandler func() *slog.Logger

// DefaultSlog is the pre-built default slog logger instance.
var DefaultSlog = NewSlog(nil)

// NewSlog creates a new Slog logger using functional options.
func NewSlog(handler SlogHandler) *Slog {
	if handler == nil {
		handler = WithTextHandler(os.Stdout, slog.LevelInfo)
	}

	return &Slog{
		log: handler(),
	}
}

// WithJSONHandler returns a logger with JSON formatting and custom level.
func WithJSONHandler(w io.Writer, level slog.Level) SlogHandler {
	return func() *slog.Logger {
		if w == nil {
			w = os.Stdout
		}
		handler := slog.NewJSONHandler(w, &slog.HandlerOptions{
			Level: level,
		})

		return slog.New(handler)
	}
}

// WithTextHandler returns a logger with text formatting and custom level.
func WithTextHandler(w io.Writer, level slog.Level) SlogHandler {
	return func() *slog.Logger {
		if w == nil {
			w = os.Stdout
		}
		handler := slog.NewTextHandler(w, &slog.HandlerOptions{
			Level: level,
		})

		return slog.New(handler)
	}
}

// Debug logs a message at debug level.
func (s *Slog) Debug(msg string, args ...any) {
	s.log.Debug(msg, args...)
}

// Info logs a message at info level.
func (s *Slog) Info(msg string, args ...any) {
	s.log.Info(msg, args...)
}

// Warn logs a message at warning level.
func (s *Slog) Warn(msg string, args ...any) {
	s.log.Warn(msg, args...)
}

// Error logs a message at error level.
func (s *Slog) Error(msg string, args ...any) {
	s.log.Error(msg, args...)
}

// Fatal logs a message at error level and then exits with status 1.
func (s *Slog) Fatal(msg string, args ...any) {
	s.log.Error(msg, args...)
	//nolint:revive // calls to os.Exit is acceptable here.
	os.Exit(1)
}

// With returns a child logger with the given key-value pairs attached.
func (s *Slog) With(args ...any) *Slog {
	return &Slog{
		log: s.log.With(args...),
	}
}
