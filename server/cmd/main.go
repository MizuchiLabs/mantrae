package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/mizuchilabs/mantrae/server/internal/api/server"
	"github.com/mizuchilabs/mantrae/server/internal/config"
)

func main() {
	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app, err := config.Setup(ctx)
	if err != nil {
		slog.Error("Setup failed", "error", err)
		return
	}

	srv := server.NewServer(app)
	if err := srv.Start(ctx); err != nil {
		slog.Error("Server error", "error", err)
		os.Exit(1)
	}
}
