package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/MizuchiLabs/mantrae/agent/client"
	"github.com/MizuchiLabs/mantrae/pkg/build"
	"github.com/MizuchiLabs/mantrae/pkg/logger"
)

func main() {
	// update := flag.Bool("update", false, "Update to latest version")
	version := flag.Bool("version", false, "Show version")

	flag.Parse()

	if *version {
		fmt.Println(build.Version)
		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	logger.Setup()
	slog.Info("Starting agent...")
	client.Client(quit)
}
