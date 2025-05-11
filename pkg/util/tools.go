package util

import (
	"os"
	"strconv"
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
