package server

import (
	"context"
	"net/http"
	"time"

	connectcors "connectrpc.com/cors"
	"github.com/mizuchilabs/mantrae/internal/settings"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// WithCORS adds CORS support to a Connect HTTP handler.
func (s *Server) WithCORS(h http.Handler) http.Handler {
	// Always include safe localhost dev URLs
	allowedOrigins := []string{
		"http://localhost:5173",
		"http://127.0.0.1:5173",
		"http://localhost:" + s.Port,
		"http://127.0.0.1:" + s.Port,
	}

	serverURL, ok := s.app.SM.Get(context.Background(), settings.KeyServerURL)
	if ok {
		allowedOrigins = append(allowedOrigins, serverURL)
	}

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: connectcors.AllowedMethods(),
		ExposedHeaders: connectcors.ExposedHeaders(),
		AllowedHeaders: append(
			connectcors.AllowedHeaders(),
			"Authorization",
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Headers",
		),
		AllowCredentials: true,
		MaxAge:           int(2 * time.Hour / time.Second),
	}).Handler(h)

	return h2c.NewHandler(corsHandler, &http2.Server{})
}
