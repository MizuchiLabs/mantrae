package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/mizuchilabs/mantrae/internal/api/server"
	"github.com/mizuchilabs/mantrae/internal/config"
)

func main() {
	// Graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app, err := config.Setup(ctx)
	if err != nil {
		slog.Error("Setup failed", "error", err)
		return
	}

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	srv := server.NewServer(app)
	if err := srv.Start(ctx); err != nil {
		slog.Error("Server error", "error", err)
		os.Exit(1)
	}
}
