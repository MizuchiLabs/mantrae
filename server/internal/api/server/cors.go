package server

import (
	"net/http"
	"time"

	connectcors "connectrpc.com/cors"
	"github.com/mizuchilabs/mantrae/pkg/util"
	"github.com/rs/cors"
)

// WithCORS adds CORS support to a Connect HTTP handler.
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

// WithCORS adds CORS support to a Connect HTTP handler.
// func (s *Server) WithCORS(h http.Handler) http.Handler {
// 	// Always include safe localhost dev URLs
// 	allowedOrigins := []string{
// 		"http://localhost:5173",
// 		"http://127.0.0.1:5173",
// 		"http://localhost:" + s.Port,
// 		"http://127.0.0.1:" + s.Port,
// 	}

// 	serverURL, ok := s.app.SM.Get(context.Background(), settings.KeyServerURL)
// 	if ok {
// 		allowedOrigins = append(allowedOrigins, serverURL)
// 	}

// 	corsHandler := cors.New(cors.Options{
// 		AllowedOrigins: allowedOrigins,
// 		AllowedMethods: connectcors.AllowedMethods(),
// 		ExposedHeaders: connectcors.ExposedHeaders(),
// 		AllowedHeaders: append(
// 			connectcors.AllowedHeaders(),
// 			"Authorization",
// 			"Access-Control-Allow-Origin",
// 			"Access-Control-Allow-Credentials",
// 			"Access-Control-Allow-Headers",
// 		),
// 		AllowCredentials: true,
// 		MaxAge:           int(2 * time.Hour / time.Second),
// 	}).Handler(h)

// 	return h2c.NewHandler(corsHandler, &http2.Server{})
// }
