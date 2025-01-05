package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/MizuchiLabs/mantrae/agent/client"
	"github.com/lmittmann/tint"
)

// Set up global logger with specified configuration
func init() {
	levelMap := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}

	logLevel := slog.LevelInfo
	if level, ok := os.LookupEnv("LOG_LEVEL"); ok {
		if l, ok := levelMap[level]; ok {
			logLevel = l
		}
	}
	slog.SetDefault(slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: logLevel})))
}

func main() {
	// update := flag.Bool("update", false, "Update to latest version")
	version := flag.Bool("version", false, "Show version")

	flag.Parse()

	if *version {
		fmt.Println(client.Version)
		return
	}

	// Check if token is set
	if token, ok := os.LookupEnv("TOKEN"); !ok || token == "" {
		slog.Error("missing token")
		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	slog.Info("Starting agent...")
	client.Client(quit)
}
