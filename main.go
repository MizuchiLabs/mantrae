package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/MizuchiLabs/mantrae/internal/api"
	"github.com/MizuchiLabs/mantrae/internal/config"
	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/pkg/dns"
	"github.com/MizuchiLabs/mantrae/pkg/traefik"
	"github.com/lmittmann/tint"
)

// Set up global logger with specified configuration
func init() {
	logger := slog.New(tint.NewHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

func main() {
	if err := db.InitDB(); err != nil {
		slog.Error("Failed to initialize database", "error", err)
		return
	}
	defer db.DB.Close() // Close the database connection when the program exits

	// Parse command-line flags and set default settings
	var flags config.Flags
	if err := flags.Parse(); err != nil {
		slog.Error("Failed to parse flags", "error", err)
		return
	}

	// Schedule backups
	if err := config.ScheduleBackups(); err != nil {
		slog.Error("Failed to schedule backups", "error", err)
	}

	// Create a context that will be used to signal background processes to stop
	ctx, cancel := context.WithCancel(context.Background())

	// Start the background sync processes
	go traefik.Sync(ctx)
	go dns.Sync(ctx)

	srv := &http.Server{
		Addr:              ":" + strconv.Itoa(flags.Port),
		Handler:           api.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	slog.Info("Listening on port", "port", flags.Port)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("ListenAndServe", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")
	cancel()

	ctxShutdown, cancelShutdown := context.WithTimeout(ctx, 5*time.Second)
	defer cancelShutdown()
	if err := srv.Shutdown(ctxShutdown); err != nil {
		slog.Error("Server forced to shutdown:", "error", err)
	}
}
