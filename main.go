package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
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
	opts := &tint.Options{}
	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		opts.Level = slog.LevelDebug
	case "info":
		opts.Level = slog.LevelInfo
	case "warn":
		opts.Level = slog.LevelWarn
	case "error":
		opts.Level = slog.LevelError
	default:
		opts.Level = slog.LevelInfo
	}
	logger := slog.New(tint.NewHandler(os.Stdout, opts))
	slog.SetDefault(logger)
}

func main() {
	if _, ok := os.LookupEnv("SECRET"); !ok {
		slog.Error("SECRET environment variable not set")
		return
	}
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
	defer cancel()

	// Start the background sync processes
	go traefik.Sync(ctx)
	go dns.Sync(ctx)

	// TODO: Later
	// Start the grpc server
	// go api.Server(flags.Agent.Port)

	// Start the WebUI server
	srv := &http.Server{
		Addr:              ":" + flags.Port,
		Handler:           api.Routes(flags.UseAuth),
		ReadHeaderTimeout: 5 * time.Second,
	}

	slog.Info("Server running on", "port", flags.Port)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("ListenAndServe", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")
	cancel()

	ctxShutdown, cancelShutdown := context.WithTimeout(ctx, 3*time.Second)
	defer cancelShutdown()
	if err := srv.Shutdown(ctxShutdown); err != nil {
		slog.Error("Server forced to shutdown:", "error", err)
	}
}
