package env_test

import (
	"testing"
	"time"

	"github.com/pactus-project/gopkg/env"
	"github.com/stretchr/testify/assert"
)

func TestGetEnvEmpty(t *testing.T) {
	assert.Empty(t, env.GetEnv("MY_STRING", ""))
	assert.Equal(t, []string{}, env.GetEnv("MY_STRING_LIST", []string{}))
}

// TestGetEnv verifies that environment variables are correctly parsed into supported types.
func TestGetEnv(t *testing.T) {
	t.Setenv("MY_INT", "1")
	t.Setenv("MY_BOOL", "true")
	t.Setenv("MY_FLOAT", "3.14")
	t.Setenv("MY_STRING", "str")
	t.Setenv("MY_STRING_LIST", "str1,str2")
	t.Setenv("MY_DURATION", "5m")

	assert.Equal(t, 1, env.GetEnv("MY_INT", 0))
	assert.True(t, env.GetEnv("MY_BOOL", false))
	assert.InEpsilon(t, 3.14, env.GetEnv("MY_FLOAT", 0.0), 0.0001)
	assert.Equal(t, "str", env.GetEnv("MY_STRING", ""))
	assert.Equal(t, []string{"str1", "str2"}, env.GetEnv("MY_STRING_LIST", []string{}))
	assert.Equal(t, time.Minute*5, env.GetEnv("MY_DURATION", time.Duration(0)))
}

// TestGetEnvWithDefault verifies that default values are used when environment variables are not set.
func TestGetEnvWithDefault(t *testing.T) {
	assert.Equal(t, 1, env.GetEnv("MY_INT", 1))
	assert.False(t, env.GetEnv("MY_BOOL", false))
	assert.True(t, env.GetEnv("MY_BOOL", true))
	assert.InEpsilon(t, 3.14, env.GetEnv("MY_FLOAT", 3.14), 0.0001)
	assert.Equal(t, "str", env.GetEnv("MY_STRING", "str"))
	assert.Equal(t, []string{"str1", "str2"}, env.GetEnv("MY_STRING_LIST", []string{"str1", "str2"}))
	assert.Equal(t, time.Second*5, env.GetEnv("MY_DURATION", time.Second*5))
}

// TestGetEnvWrongType checks that GetEnv panics when default values cannot be parsed into the desired type.
func TestGetEnvWrongType(t *testing.T) {
	assert.Panics(t, func() {
		t.Setenv("MY_INT", "one")
		env.GetEnv("MY_INT", 0)
	})
	assert.Panics(t, func() {
		t.Setenv("MY_BOOL", "ok")
		env.GetEnv("MY_BOOL", false)
	})
	assert.Panics(t, func() {
		t.Setenv("MY_FLOAT", "pi")
		env.GetEnv("MY_FLOAT", 0)
	})
	assert.Panics(t, func() {
		t.Setenv("MY_DURATION", "two seconds")
		env.GetEnv("MY_DURATION", time.Duration(0))
	})
}
