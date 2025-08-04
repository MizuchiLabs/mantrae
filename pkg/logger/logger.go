// Package logger provides logging setup and configuration for the application.
package logger

import (
	"log/slog"
	"strings"

	"github.com/mizuchilabs/mantrae/pkg/util"
)

func Setup() {
	levelMap := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}
	formatMap := map[string]string{
		"json": FormatJSON,
		"text": FormatText,
	}
	level := util.GetEnv("LOG_LEVEL", "info")
	logLevel, exists := levelMap[strings.ToLower(level)]
	if !exists {
		logLevel = slog.LevelInfo
	}
	format := util.GetEnv("LOG_FORMAT", "text")
	logFormat, exists := formatMap[strings.ToLower(format)]
	if !exists {
		logFormat = FormatText
	}
	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	log := slog.New(NewHandler(opts, logFormat))
	slog.SetDefault(log)
}
