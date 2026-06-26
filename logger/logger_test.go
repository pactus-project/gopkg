package logger

import (
	"bytes"
	"errors"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

type Foo struct{}

func (Foo) LogString() string {
	return "foo"
}

type Bar struct{}

func (Bar) String() string {
	return "bar"
}

func setupTestLogger(t *testing.T, conf *Config) {
	t.Helper()

	InitGlobalLogger(t.Context(), conf)
}

func TestNilObjLogger(t *testing.T) {
	conf := DefaultConfig()
	setupTestLogger(t, conf)

	l := NewSubLogger(t.Context(), "test", nil)
	var buf bytes.Buffer
	l.logger = l.logger.Output(&buf)

	l.Info("hello", "error", errors.New("error"))
	assert.Contains(t, buf.String(), "hello")
	assert.Contains(t, buf.String(), "error")
}

func TestSubLogger(t *testing.T) {
	var buf bytes.Buffer
	subLogger := NewSubLogger(t.Context(), "test", Foo{})
	subLogger.logger = subLogger.logger.Output(&buf)

	subLogger.Info("msg")

	out := buf.String()

	assert.Contains(t, out, "info")
	assert.Contains(t, out, "msg")
}

func TestSubLoggerLevel(t *testing.T) {
	conf := DefaultConfig()
	conf.Colorful = false
	setupTestLogger(t, conf)

	globalInst.config.Levels["test"] = "warn"
	subLogger := NewSubLogger(t.Context(), "test", Foo{})
	var buf bytes.Buffer
	subLogger.logger = subLogger.logger.Output(&buf)

	subLogger.Trace("msg")
	subLogger.Debug("msg")
	subLogger.Info("msg")
	subLogger.Warn("msg")
	subLogger.Error("msg")

	out := buf.String()

	assert.Contains(t, out, "foo")
	assert.NotContains(t, out, "trace")
	assert.NotContains(t, out, "debug")
	assert.NotContains(t, out, "info")
	assert.Contains(t, out, "warn")
	assert.Contains(t, out, "error")
}

func TestLogger(t *testing.T) {
	conf := DefaultConfig()
	conf.Colorful = false
	setupTestLogger(t, conf)

	var buf bytes.Buffer
	log.Logger = log.Output(&buf)

	Trace("msg", "trace", "trace")
	Debug("msg", "Debug", "Debug")
	Info("msg", nil)
	Info("msg", "a", nil)
	Info("msg", "b", []byte{1, 2, 3})
	Warn("msg", "x")
	Error("msg", "y", Foo{})

	out := buf.String()

	t.Log(out)

	assert.NotContains(t, out, "trace")
	assert.Contains(t, out, "debug")
	assert.Contains(t, out, "foo")
	assert.Contains(t, out, "010203")
	assert.Contains(t, out, "!INVALID-KEY!")
	assert.Contains(t, out, "!MISSING-VALUE!")
	assert.Contains(t, out, "null")
	assert.Contains(t, out, "info")
	assert.Contains(t, out, "warn")
	assert.Contains(t, out, "error")
}

func TestNilValue(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = log.Output(&buf)

	var foo *Foo

	Info("msg", "null", nil)
	Info("msg", "error", error(nil))
	Info("msg", "stringer", foo)

	out := buf.String()

	t.Log(out)

	assert.Contains(t, out, "null")
	assert.Contains(t, out, "error")
	assert.Contains(t, out, "stringer")
}

func TestInvalidLevel(t *testing.T) {
	conf := DefaultConfig()
	conf.Colorful = false
	setupTestLogger(t, conf)

	var buf1 bytes.Buffer
	log.Logger = log.Logger.Output(&buf1)

	globalInst.config.Levels["test"] = "invalid"
	l := NewSubLogger(t.Context(), "test", Foo{})

	var buf2 bytes.Buffer
	l.logger = l.logger.Output(&buf2)

	l.Error("message", "key", "val")

	out1 := buf1.String()
	out2 := buf2.String()

	t.Log(out1)

	assert.Contains(t, out1, "Unknown Level String")
	assert.NotContains(t, out2, "error")
	assert.NotContains(t, out2, "message")
}

func TestLogStringer(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = log.Output(&buf)

	Info("msg", "obj", Foo{})
	assert.Contains(t, buf.String(), "foo")
}

func TestStringer(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = log.Output(&buf)

	Info("msg", "obj", Bar{})
	assert.Contains(t, buf.String(), "bar")
}

func TestFatal(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = log.Output(&buf)

	// Replace FatalExitFunc to prevent os.Exit during test.
	saved := zerolog.FatalExitFunc
	zerolog.FatalExitFunc = func() {}
	defer func() { zerolog.FatalExitFunc = saved }()

	Fatal("fatal message", "key", "val")

	out := buf.String()
	assert.Contains(t, out, "fatal message")
	assert.Contains(t, out, "key")
	assert.Contains(t, out, "val")
}

func TestPanic(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = log.Output(&buf)

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic, got nil")
		}
		out := buf.String()
		assert.Contains(t, out, "panic message")
	}()

	Panic("panic message")
}
