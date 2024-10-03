package main

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/MizuchiLabs/mantrae/agent/client"
	"github.com/lmittmann/tint"
)

// Set up global logger
func init() {
	logger := slog.New(tint.NewHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

func main() {
	token := flag.String("token", "", "Authentication token (required)")
	// update := flag.Bool("update", false, "Update to latest version")
	version := flag.Bool("version", false, "Show version")

	flag.Parse()

	if *token == "" {
		*token = client.LoadToken()
		if len(*token) == 0 {
			slog.Error("missing token")
			return
		}
	} else {
		client.SaveToken(*token)
	}

	if *version {
		log.Println("v0.0.1")
		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	slog.Info("Starting agent...")
	client.Client(quit)
}
