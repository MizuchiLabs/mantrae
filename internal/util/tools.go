// Package util provides utility functions.
package util

import (
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
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

func ResolvePath(path string) string {
	basePath := GetEnv("BASE_PATH", "data")

	// If the provided path is absolute, return it as-is
	if filepath.IsAbs(path) {
		return path
	}

	// Create the base directory if it doesn't exist
	if err := os.MkdirAll(basePath, 0o750); err != nil {
		log.Printf("Warning: failed to create base directory: %v", err)
	}

	return filepath.Join(basePath, path)
}

// CopyFile copies a file from src to dst safely
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(filepath.Clean(src))
	if err != nil {
		return err
	}
	defer func() {
		if err = srcFile.Close(); err != nil {
			slog.Error("failed to close source file", "error", err)
		}
	}()

	dstFile, err := os.Create(filepath.Clean(dst))
	if err != nil {
		return err
	}
	defer func() {
		if err = dstFile.Close(); err != nil {
			slog.Error("failed to close destination file", "error", err)
		}
	}()

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return dstFile.Sync()
}

func CleanSliceStr(in []string) []string {
	out := in[:0]
	for _, s := range in {
		if s != "" {
			out = append(out, s)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
