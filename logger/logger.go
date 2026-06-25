// Package logger provides a structured logging abstraction using zerolog
// with global convenience functions and pluggable log targets.
package logger

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"reflect"
	"slices"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogStringer is an interface for types that want to provide a compact
// string representation for logging (e.g. hashes, IDs).
type LogStringer interface {
	LogString() string
}

var globalInst *logger

type logger struct {
	config *Config
	writer io.Writer
}

// InitGlobalLogger initializes the global logger with the given config.
// The context is used to attach initial fields to the root logger.
// Subsequent calls are no-ops.
func InitGlobalLogger(ctx context.Context, conf *Config) {
	if globalInst != nil {
		return
	}

	writers := []io.Writer{}

	if slices.Contains(conf.Targets, "file") {
		fileWriter := &lumberjack.Logger{
			Filename:   conf.Filename,
			MaxSize:    conf.MaxSize,
			MaxBackups: conf.MaxBackups,
			Compress:   conf.Compress,
			MaxAge:     conf.RotateLogAfterDays,
		}
		writers = append(writers, fileWriter)
	}

	if slices.Contains(conf.Targets, "console") {
		if conf.Colorful {
			consoleWriter := &zerolog.ConsoleWriter{
				Out:        os.Stderr,
				TimeFormat: "15:04:05",
			}
			writers = append(writers, consoleWriter)
		} else {
			writers = append(writers, os.Stderr)
		}
	}

	globalInst = &logger{
		config: conf,
		writer: io.MultiWriter(writers...),
	}
	log.Logger = zerolog.New(globalInst.writer).With().Ctx(ctx).Timestamp().Logger()

	lvl, err := zerolog.ParseLevel(conf.Levels["default"])
	if err != nil {
		Warn("invalid default log level", "error", err)
	}
	log.Logger = log.Logger.Level(lvl)
}

func addFields(event *zerolog.Event, keyvals ...any) *zerolog.Event {
	if event == nil {
		return nil
	}

	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "!MISSING-VALUE!")
	}
	for index := 0; index < len(keyvals); index += 2 {
		key, ok := keyvals[index].(string)
		if !ok {
			key = "!INVALID-KEY!"
		}

		value := keyvals[index+1]
		switch typ := value.(type) {
		case LogStringer:
			if isNil(typ) {
				event.Any(key, typ)
			} else {
				event.Str(key, typ.LogString())
			}
		case fmt.Stringer:
			if isNil(typ) {
				event.Any(key, typ)
			} else {
				event.Stringer(key, typ)
			}
		case error:
			event.AnErr(key, typ)
		case []byte:
			event.Str(key, hex.EncodeToString(typ))
		default:
			event.Any(key, typ)
		}
	}

	return event
}

// Trace logs a message at trace level.
func Trace(msg string, keyvals ...any) {
	addFields(log.Trace(), keyvals...).Msg(msg)
}

// Debug logs a message at debug level.
func Debug(msg string, keyvals ...any) {
	addFields(log.Debug(), keyvals...).Msg(msg)
}

// Info logs a message at info level.
func Info(msg string, keyvals ...any) {
	addFields(log.Info(), keyvals...).Msg(msg)
}

// Warn logs a message at warning level.
func Warn(msg string, keyvals ...any) {
	addFields(log.Warn(), keyvals...).Msg(msg)
}

// Error logs a message at error level.
func Error(msg string, keyvals ...any) {
	addFields(log.Error(), keyvals...).Msg(msg)
}

// Fatal logs a message at fatal level, then calls os.Exit(1).
// Use FatalExitFunc to override the exit behavior.
func Fatal(msg string, keyvals ...any) {
	addFields(log.Fatal(), keyvals...).Msg(msg)
}

// Panic logs a message at panic level, then panics.
func Panic(msg string, keyvals ...any) {
	addFields(log.Panic(), keyvals...).Msg(msg)
}

func isNil(val any) bool {
	if val == nil {
		return true
	}

	defer func() { _ = recover() }()

	return reflect.ValueOf(val).IsNil()
}
