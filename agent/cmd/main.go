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
	// update := flag.Bool("update", false, "Update to latest version")
	version := flag.Bool("version", false, "Show version")

	flag.Parse()

	if *version {
		fmt.Println(meta.Version)
		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	logger.Setup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start agent loop
	go client.Client(ctx)

	<-quit
	slog.Info("Shutting down agent...")
	cancel()
	<-ctx.Done()
}
