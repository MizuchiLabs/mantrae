package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/MizuchiLabs/mantrae/internal/api"
	"github.com/MizuchiLabs/mantrae/pkg/traefik"
	"github.com/MizuchiLabs/mantrae/pkg/util"
	"github.com/lmittmann/tint"
)

//go:embed all:web/build
var webFS embed.FS

// Set up global logger with specified configuration
func init() {
	logger := slog.New(tint.NewHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	if err := util.GenerateCreds(); err != nil {
		slog.Error("Failed to generate creds", "error", err)
	}
	var profiles traefik.Profiles
	if err := profiles.Load(); err != nil {
		slog.Error("Failed to get traefik config", "error", err)
	}
	go traefik.GetTraefikConfig()
}

func main() {
	version := flag.Bool("version", false, "Print version and exit")
	port := flag.Int("port", 3000, "Port to listen on")
	url := flag.String(
		"url",
		"",
		"Specify the URL of the Traefik instance (e.g. http://localhost:8080)",
	)
	username := flag.String(
		"username",
		"",
		"Specify the username for the Traefik instance",
	)
	password := flag.String(
		"password",
		"",
		"Specify the password for the Traefik instance",
	)
	flag.Parse()

	if *version {
		fmt.Println(util.Version)
		os.Exit(0)
	}

	if *url != "" {
		var profiles traefik.Profiles
		if err := profiles.SetDefaultProfile(*url, *username, *password); err != nil {
			slog.Error("Failed to add default profile", "error", err)
			return
		}
	}

	mux := api.Routes()
	middle := api.Chain(api.Log, api.Cors)

	staticContent, err := fs.Sub(webFS, "web/build")
	if err != nil {
		slog.Error("Sub", "error", err)
		return
	}

	mux.Handle("/", http.FileServer(http.FS(staticContent)))

	// Start the background sync process
	go traefik.Sync()

	srv := &http.Server{
		Addr:              ":" + strconv.Itoa(*port),
		Handler:           middle(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	slog.Info("Listening on port", "port", *port)
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
