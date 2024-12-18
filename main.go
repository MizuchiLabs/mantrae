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
	"github.com/MizuchiLabs/mantrae/internal/tasks"
	"github.com/MizuchiLabs/mantrae/pkg/util"
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

	logLevel, exists := levelMap[strings.ToLower(util.App.LogLevel)]
	if !exists {
		logLevel = slog.LevelInfo
	}

	opts := &tint.Options{Level: logLevel}
	logger := slog.New(tint.NewHandler(os.Stdout, opts))
	slog.SetDefault(logger)
}

func main() {
	if err := db.InitDB(); err != nil {
		slog.Error("Failed to initialize database", "error", err)
		return
	}
	defer db.DB.Close() // Close the database connection when the program exits

	// Parse command-line flags and set default settings
	if err := config.Parse(); err != nil {
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
	tasks.StartSync(ctx)

	// Start the WebUI server
	srv := &http.Server{
		Addr:              ":" + util.App.Port,
		Handler:           api.Server(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	slog.Info("Server running on", "port", util.App.Port)
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
