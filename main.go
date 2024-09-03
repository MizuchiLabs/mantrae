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
	"github.com/MizuchiLabs/mantrae/pkg/traefik"
	"github.com/MizuchiLabs/mantrae/pkg/util"
	"github.com/lmittmann/tint"
)

// Set up global logger with specified configuration
func init() {
	logger := slog.New(tint.NewHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

func main() {
	db, err := db.InitDB()
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		return
	}
	defer db.Close() // Close the database connection when the program exits

	flags := config.ParseFlags() // Parse command-line flags
	util.SetDefaultAdminUser()   // Set default admin user

	// Start the background sync processes
	go traefik.Sync()
	// go dns.Sync()

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown:", "error", err)
	}
}
