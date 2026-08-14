// Package env provides typed environment variable retrieval with generic support.
package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// SupportedTypes is the constraint interface for types that GetEnv can return.
type SupportedTypes interface {
	~string | ~int | ~float64 | ~bool | ~[]string | time.Duration
}

// GetEnv retrieves an environment variable by key.
// if the variable is not set, it returns the provided default value.
//
// Panics if the conversion fails or the type is unsupported.
func GetEnv[T SupportedTypes](key string, def T) T {
	val := os.Getenv(key)
	if val == "" {
		return def
	}

	var result T
	switch any(result).(type) {
	case int:
		v, err := strconv.Atoi(val)
		if err != nil {
			panic(fmt.Errorf("failed to convert %q to int: %w", val, err))
		}

		return any(v).(T)

	case bool:
		v, err := strconv.ParseBool(val)
		if err != nil {
			panic(fmt.Errorf("failed to convert %q to bool: %w", val, err))
		}

		return any(v).(T)

	case float64:
		v, err := strconv.ParseFloat(val, 64)
		if err != nil {
			panic(fmt.Errorf("failed to convert %q to float64: %w", val, err))
		}

		return any(v).(T)

	case string:
		return any(val).(T)

	case []string:
		if val == "" {
			return any([]string{}).(T)
		}
		parts := strings.Split(val, ",")

		return any(parts).(T)

	case time.Duration:
		dur, err := time.ParseDuration(val)
		if err != nil {
			panic(fmt.Errorf("failed to parse duration %q: %w", val, err))
		}

		return any(dur).(T)

	default:
		panic(fmt.Sprintf("unsupported type: %T", result))
	}
}
