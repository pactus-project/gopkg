package logger

import (
	"context"

	"github.com/rs/zerolog"
)

// SubLogger is a named child logger that can attach an object
// (implementing LogStringer) to every message and has its own log level.
type SubLogger struct {
	logger zerolog.Logger
	name   string
	obj    LogStringer
}

// NewSubLogger creates a new SubLogger with the given name and optional object.
// It derives its log level from the global config's Levels map (falling back to "default").
func NewSubLogger(ctx context.Context, name string, obj LogStringer) *SubLogger {
	sub := &SubLogger{
		logger: zerolog.New(getGlobalInst(ctx).writer).With().Ctx(ctx).Timestamp().Logger(),
		name:   name,
		obj:    obj,
	}

	lvlStr := getGlobalInst(ctx).config.Levels[name]
	if lvlStr == "" {
		lvlStr = getGlobalInst(ctx).config.Levels["default"]
	}

	lvl, err := zerolog.ParseLevel(lvlStr)
	if err != nil {
		Warn("invalid log level", "error", err, "name", name)
	}
	sub.logger = sub.logger.Level(lvl)

	return sub
}

func (sl *SubLogger) logObj(event *zerolog.Event, msg string, keyvals ...any) {
	if event == nil {
		return
	}

	if sl.obj != nil {
		event = event.Str(sl.name, sl.obj.LogString())
	}

	addFields(event, keyvals...).Msg(msg)
}

// Trace logs a message at trace level.
func (sl *SubLogger) Trace(msg string, keyvals ...any) {
	sl.logObj(sl.logger.Trace(), msg, keyvals...)
}

// Debug logs a message at debug level.
func (sl *SubLogger) Debug(msg string, keyvals ...any) {
	sl.logObj(sl.logger.Debug(), msg, keyvals...)
}

// Info logs a message at info level.
func (sl *SubLogger) Info(msg string, keyvals ...any) {
	sl.logObj(sl.logger.Info(), msg, keyvals...)
}

// Warn logs a message at warning level.
func (sl *SubLogger) Warn(msg string, keyvals ...any) {
	sl.logObj(sl.logger.Warn(), msg, keyvals...)
}

// Error logs a message at error level.
func (sl *SubLogger) Error(msg string, keyvals ...any) {
	sl.logObj(sl.logger.Error(), msg, keyvals...)
}

// Fatal logs a message at fatal level, then calls os.Exit(1).
// Use FatalExitFunc to override the exit behavior.
func (sl *SubLogger) Fatal(msg string, keyvals ...any) {
	sl.logObj(sl.logger.Fatal(), msg, keyvals...)
}

// Panic logs a message at panic level, then panics.
func (sl *SubLogger) Panic(msg string, keyvals ...any) {
	sl.logObj(sl.logger.Panic(), msg, keyvals...)
}

// SetObj replaces the attached object for this sub-logger.
func (sl *SubLogger) SetObj(obj LogStringer) {
	sl.obj = obj
}
