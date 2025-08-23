package envv

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type envv[T any] struct {
	name         string
	required     bool
	defaultValue *T
}

type basicEnvv struct {
	name string
}

type typedEnvv[T any] struct {
	name string
}

// Get creates a new basic env var instance for the given env var name.
// This is the entry point for the fluent API to read env vars.
func Get(name string) basicEnvv {
	return basicEnvv{
		name: name,
	}
}

// String creates a string-typed env var.
func (e basicEnvv) String() typedEnvv[string] {
	return typedEnvv[string](e)
}

// Int creates a int-typed env var.

func (e basicEnvv) Int() typedEnvv[int] {
	return typedEnvv[int](e)
}

// Bool creates a bool-typed env var.

func (e basicEnvv) Bool() typedEnvv[bool] {
	return typedEnvv[bool](e)
}

// Float64 creates a float64-typed env var.
func (e basicEnvv) Float64() typedEnvv[float64] {
	return typedEnvv[float64](e)
}

// Duration creates a duration-typed env var.
func (e basicEnvv) Duration() typedEnvv[time.Duration] {
	return typedEnvv[time.Duration](e)
}

// Default sets a default value for the environment variable.
// it can only be used on a typed env var.
func (e typedEnvv[T]) Default(defaultValue T) envv[T] {
	return envv[T]{
		name:         e.name,
		required:     false,
		defaultValue: &defaultValue,
	}
}

// Required marks the env var as required, causing Parse to panic if the variable is not set.
func (e typedEnvv[T]) Required() envv[T] {
	return envv[T]{
		name:         e.name,
		required:     true,
		defaultValue: nil,
	}
}

// Optional marks the env var as optional, allowing it to be unset without causing a panic.
func (e typedEnvv[T]) Optional() envv[T] {
	return envv[T]{
		name:         e.name,
		required:     false,
		defaultValue: nil,
	}
}

// Parse reads and parses the env var according to its type.
// It returns the typed value if successful, or panics if:
// - The variable is required but not set
// - The value cannot be parsed into the expected type
func (e envv[T]) Parse() T {
	value := os.Getenv(e.name)
	if value == "" && e.required {
		panic("required env var is missing: " + e.name)
	}
	if value == "" && e.defaultValue != nil {
		return *e.defaultValue
	}
	var t T
	switch any(t).(type) {
	case string:
		return any(value).(T)
	case int:
		v, err := strconv.Atoi(value)
		if err != nil {
			panic(fmt.Errorf("failed to parse int value %q: %w", value, err))
		}
		return any(v).(T)
	case bool:
		v, err := strconv.ParseBool(value)
		if err != nil {
			panic(fmt.Errorf("failed to parse bool value %q: %w", value, err))
		}
		return any(v).(T)
	case float64:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			panic(fmt.Errorf("failed to parse float64 value %q: %w", value, err))
		}
		return any(v).(T)
	case time.Duration:
		v, err := time.ParseDuration(value)
		if err != nil {
			panic(fmt.Errorf("failed to parse duration value %q: %w", value, err))
		}
		return any(v).(T)
	default:
		panic(fmt.Sprintf("expected type for env var %q", e.name))
	}
}
