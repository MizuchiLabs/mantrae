// Package main is the entrypoint for the mantrae agent.
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/mizuchilabs/mantrae/agent/internal/client"
	"github.com/mizuchilabs/mantrae/pkg/logger"
	"github.com/mizuchilabs/mantrae/pkg/meta"
)

func main() {
	version := flag.Bool("version", false, "Show version")
	flag.Parse()

	if *version {
		fmt.Println(meta.Version)
		return
	}

	logger.Setup()

	cfg, err := client.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	slog.Info("Starting agent",
		"server", cfg.ServerURL,
		"profile_id", cfg.ProfileID,
		"agent_id", cfg.AgentID)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	agent := client.NewAgent(cfg)

	// Handle graceful shutdown
	go func() {
		<-quit
		slog.Info("Shutting down agent...")
		cancel()
	}()

	agent.Run(ctx)
}
