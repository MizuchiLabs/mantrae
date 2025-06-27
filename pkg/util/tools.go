package util

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// GetEnv returns the value of the environment variable `key` or `fallback`
func GetEnv[T any](key string, fallback T) T {
	if value, ok := os.LookupEnv(key); ok {
		// We need to convert the string to type T
		var result T
		switch any(fallback).(type) {
		case string:
			result = any(value).(T)
		case int:
			if v, err := strconv.Atoi(value); err == nil {
				result = any(v).(T)
			} else {
				return fallback
			}
		case bool:
			if v, err := strconv.ParseBool(value); err == nil {
				result = any(v).(T)
			} else {
				return fallback
			}
		case float64:
			if v, err := strconv.ParseFloat(value, 64); err == nil {
				result = any(v).(T)
			} else {
				return fallback
			}
		default:
			return fallback
		}
		return result
	}
	return fallback
}

// IsTest returns true if the current program is running in a test environment
func IsTest() bool {
	return strings.HasSuffix(os.Args[0], ".test")
}

func ResolvePath(path string) string {
	basePath := GetEnv("BASE_PATH", "data")

	// If the provided path is absolute, return it as-is
	if filepath.IsAbs(path) {
		return path
	}

	// Create the base directory if it doesn't exist
	if err := os.MkdirAll(basePath, 0755); err != nil {
		log.Printf("Warning: failed to create base directory: %v", err)
	}

	return filepath.Join(basePath, path)
}
