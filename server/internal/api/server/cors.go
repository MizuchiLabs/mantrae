package server

import (
	"net/http"
	"time"

	connectcors "connectrpc.com/cors"
	"github.com/mizuchilabs/mantrae/pkg/util"
	"github.com/rs/cors"
)

func (s *Server) WithCORS(h http.Handler) http.Handler {
	allowedOrigins := []string{
		util.OriginOnly(s.app.BaseURL),
		util.OriginOnly(s.app.FrontendURL),
	}

	return cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   connectcors.AllowedMethods(),
		ExposedHeaders:   connectcors.ExposedHeaders(),
		AllowedHeaders:   connectcors.AllowedHeaders(),
		AllowCredentials: true,
		MaxAge:           int(2 * time.Hour / time.Second),
	}).Handler(h)
}
