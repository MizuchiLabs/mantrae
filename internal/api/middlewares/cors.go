package middlewares

import (
	"net/http"
	"time"

	connectcors "connectrpc.com/cors"
	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/settings"
	"github.com/rs/cors"
)

// WithCORS adds CORS support to a Connect HTTP handler.
func WithCORS(h http.Handler, app *config.App, port string) http.Handler {
	// Always include safe localhost dev URLs
	allowedOrigins := []string{
		"http://localhost:5173",
		"http://127.0.0.1:5173",
		"http://localhost:" + port,
		"http://127.0.0.1:" + port,
	}

	serverURL, ok := app.SM.Get(settings.KeyServerURL)
	if ok {
		allowedOrigins = append(allowedOrigins, serverURL)
	}

	return cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: connectcors.AllowedMethods(),
		// AllowedHeaders: connectcors.AllowedHeaders(),
		AllowedHeaders: []string{"*"},
		ExposedHeaders: connectcors.ExposedHeaders(),
		MaxAge:         int(2 * time.Hour / time.Second),
	}).Handler(h)
}
