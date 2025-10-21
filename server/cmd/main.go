package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/mizuchilabs/mantrae/server/internal/api/server"
	"github.com/mizuchilabs/mantrae/server/internal/config"
	"github.com/mizuchilabs/mantrae/server/internal/tasks"
)

func main() {
	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app, err := config.Setup(ctx)
	if err != nil {
		slog.Error("Setup failed", "error", err)
		os.Exit(1)
	}

	// Start background jobs
	scheduler := tasks.NewScheduler(ctx, app)
	scheduler.Start()

	srv := server.NewServer(app)
	if err := srv.Start(ctx); err != nil {
		slog.Error("Server error", "error", err)
		os.Exit(1)
	}
}
