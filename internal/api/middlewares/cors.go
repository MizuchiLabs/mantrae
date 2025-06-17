package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	connectcors "connectrpc.com/cors"
	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/settings"
	"github.com/rs/cors"
)

const (
	defaultDevOrigin = "http://127.0.0.1:5173"
)

func CORS(allowedOrigins ...string) func(http.Handler) http.Handler {
	// Default to dev origin if none provided
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{defaultDevOrigin}
	}

	// Create a map for faster origin lookup
	originMap := make(map[string]bool)
	for _, origin := range allowedOrigins {
		originMap[strings.TrimSpace(origin)] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Set CORS headers
			if origin != "" && originMap[origin] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().
				Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// WithCORS adds CORS support to a Connect HTTP handler.
func WithCORS(h http.Handler, app *config.App) http.Handler {
	serverURL, err := app.Conn.GetQuery().GetSetting(context.Background(), settings.KeyServerURL)
	if err != nil {
		slog.Error("Failed to get server URL", "error", err)
		return h
	}
	if serverURL.Value == "" {
		slog.Error("Server URL not set")
		return h
	}
	middleware := cors.New(cors.Options{
		AllowedOrigins: []string{serverURL.Value},
		AllowedMethods: connectcors.AllowedMethods(),
		AllowedHeaders: connectcors.AllowedHeaders(),
		ExposedHeaders: connectcors.ExposedHeaders(),
	})
	return middleware.Handler(h)
}
