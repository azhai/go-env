// Package env provides a simple environment variable manager for Go.
// It supports loading environment variables from a .env file and the system environment.
package env

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// Version is the current version of the env package.
const Version = "0.0.4"

// Env represents an environment variable manager.
// It stores environment variables in a map and supports loading from a file and the system environment.
type Env struct {
	storage map[string]Entry
}

// New creates a new Env instance and loads environment variables from the default .env file.
func New() *Env {
	return NewWithFile(".env")
}

// NewWithFile creates a new Env instance and loads environment variables from the specified file.
// It initializes the storage map and attempts to Load the file.
// If the file cannot be opened or read, the error is ignored.
func NewWithFile(filename string) *Env {
	env := &Env{storage: make(map[string]Entry)}
	_ = env.Load(os.Open(filename))
	return env
}

// Load reads environment variables from a reader and stores them in the internal storage.
// It skips comments (lines starting with #) and empty lines.
// It splits each line into key-value pairs using the first '=' character.
// It trims whitespace from keys and values and removes surrounding quotes if present.
// If there is an error reading the file, it returns the error.
func (e *Env) Load(reader io.ReadCloser, err error) error {
	if err != nil {
		return err
	}
	// Close the file and handle any potential error
	defer func() {
		if closeErr := reader.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.TrimLeftFunc(scanner.Text(), unicode.IsSpace)

		// Skip comments or empty lines
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		// Split key and value by the first occurrence of '='
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimRightFunc(parts[0], unicode.IsSpace)
		value := strings.TrimSpace(parts[1])

		// Remove surrounding quotes if any
		if len(value) >= 2 {
			firstChar, lastChar := value[0], value[len(value)-1]
			if (firstChar == '"' && lastChar == '"') || (firstChar == '\'' && lastChar == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		// Store the key-value pair in the storage
		e.storage[key] = Entry{Key: key, Value: value}
	}

	return scanner.Err()
}

// Lookup searches for an environment variable by key.
// It first checks the internal storage. If not found, it checks the system environment variables.
// If the variable is found, it returns the Entry and true; otherwise, it returns an empty Entry and false.
// If the variable is found in the system environment, it is added to the internal storage.
func (e *Env) Lookup(key string) (Entry, bool) {
	if entry, ok := e.storage[key]; ok {
		return entry, true
	}

	ee := os.Getenv(key)
	if ee != "" {
		entry := Entry{Key: key, Value: ee}
		e.storage[key] = entry
		return entry, true
	}

	return Entry{}, false
}

// Get retrieves and expands the value of an environment variable by key.
func (e *Env) Get(key string) string {
	if s := e.GetStr(key); s != "" {
		return os.Expand(s, func(k string) string {
			return e.GetStr(k)
		})
	}
	return ""
}

// GetStr retrieves the string value of an environment variable by key.
// If the variable is not found, it returns the fallback value.
func (e *Env) GetStr(key string, fallback ...string) string {
	if entry, ok := e.Lookup(key); ok || len(fallback) == 0 {
		return entry.Str()
	}
	return fallback[0]
}

// GetInt retrieves the integer value of an environment variable by key.
// If the variable is not found or the value cannot be converted to an integer, it returns the fallback value.
func (e *Env) GetInt(key string, fallback ...int) int {
	if entry, ok := e.Lookup(key); ok || len(fallback) == 0 {
		return entry.Int()
	}
	return fallback[0]
}

// GetInt64 retrieves the 64-bit integer value of an environment variable by key.
// If the variable is not found or the value cannot be converted to a 64-bit integer, it returns the fallback value.
func (e *Env) GetInt64(key string, fallback ...int64) int64 {
	if entry, ok := e.Lookup(key); ok || len(fallback) == 0 {
		return entry.Int64()
	}
	return fallback[0]
}

// GetBool retrieves the boolean value of an environment variable by key.
// It supports "yes", "no", "true", and "false" as valid boolean values.
// If the variable is not found or the value cannot be converted to a boolean, it returns the fallback value.
func (e *Env) GetBool(key string, fallback ...bool) bool {
	if entry, ok := e.Lookup(key); ok || len(fallback) == 0 {
		return entry.Bool()
	}
	return fallback[0]
}

// Entry represents an environment variable entry with a key and a value.
type Entry struct {
	Key   string
	Value string
}

// Str returns the string value of the entry.
func (e *Entry) Str() string {
	return e.Value
}

// Int converts the entry's value to an integer.
// If the conversion fails, it returns 0.
func (e *Entry) Int() int {
	if val, err := strconv.Atoi(e.Value); err == nil {
		return val
	}
	return 0
}

// Int64 converts the entry's value to a 64-bit integer.
// If the conversion fails, it returns 0.
func (e *Entry) Int64() int64 {
	if val, err := strconv.ParseInt(e.Value, 10, 64); err == nil {
		return val
	}
	return 0
}

// Bool converts the entry's value to a boolean.
// It supports "yes", "no", "true", and "false" as valid boolean values.
// If the conversion fails, it returns false.
func (e *Entry) Bool() bool {
	e.Value = strings.ToLower(e.Value)
	if e.Value == "yes" {
		return true
	} else if e.Value == "no" {
		return false
	}
	if val, err := strconv.ParseBool(e.Value); err == nil {
		return val
	}
	return false
}
