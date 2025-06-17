package server

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"github.com/caarlos0/env/v11"
	"github.com/mizuchilabs/mantrae/internal/api/handler"
	"github.com/mizuchilabs/mantrae/internal/api/middlewares"
	"github.com/mizuchilabs/mantrae/internal/api/service"
	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1/mantraev1connect"
	"github.com/mizuchilabs/mantrae/web"
)

const elementsHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>API Documentation</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
	<script src="https://unpkg.com/@stoplight/elements/web-components.min.js"></script>
    <link rel="stylesheet" href="https://unpkg.com/@stoplight/elements/styles.min.css">
</head>
<body>
    <elements-api
        apiDescriptionUrl="/openapi.yaml"
        router="hash"
        layout="sidebar"
    />
</body>
</html>
`

type Server struct {
	Host          string `env:"HOST"           envDefault:"0.0.0.0"`
	Port          string `env:"PORT"           envDefault:"3000"`
	SecureTraefik bool   `env:"SECURE_TRAEFIK" envDefault:"false"`
	mux           *http.ServeMux
	app           *config.App
}

func NewServer(app *config.App) *Server {
	cfg, err := env.ParseAs[Server]()
	if err != nil {
		log.Fatal(err)
	}

	return &Server{
		Host: cfg.Host,
		Port: cfg.Port,
		mux:  http.NewServeMux(),
		app:  app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.registerServices()
	defer s.app.Conn.Close()

	server := &http.Server{
		Addr:              s.Host + ":" + s.Port,
		Handler:           middlewares.WithCORS(s.mux, s.app),
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
		MaxHeaderBytes:    8 * 1024, // 8KiB
	}

	// Channel to catch server errors
	serverErr := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		slog.Info("Server starting", "host", s.Host, "port", s.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// Wait for context cancellation or server error
	select {
	case <-ctx.Done():
		slog.Info("Shutting down server gracefully...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server shutdown error: %w", err)
		}
		return nil

	case err := <-serverErr:
		return fmt.Errorf("server error: %w", err)
	}
}

func (s *Server) registerServices() {
	// Common interceptors
	opts := []connect.HandlerOption{
		connect.WithCompressMinBytes(1024),
		connect.WithInterceptors(
			middlewares.Authentication(s.app),
			middlewares.Logging(),
		),
		connect.WithRecover(
			func(ctx context.Context, spec connect.Spec, header http.Header, panic any) error {
				// Log the panic with context
				slog.Error("panic recovered in RPC call",
					"method", spec.Procedure,
					"panic", panic,
					"trace", string(debug.Stack()),
				)
				header.Set("X-Error-Type", "panic")
				return connect.NewError(connect.CodeInternal, fmt.Errorf("internal server error"))
			},
		),
	}

	// Static files
	staticContent, err := fs.Sub(web.StaticFS, "build")
	if err != nil {
		log.Fatal(err)
	}
	s.mux.Handle("/", http.FileServer(http.FS(staticContent)))

	serviceNames := []string{
		mantraev1connect.ProfileServiceName,
		mantraev1connect.UserServiceName,
		mantraev1connect.EntryPointServiceName,
		mantraev1connect.SettingServiceName,
		mantraev1connect.AgentServiceName,
		mantraev1connect.AgentManagementServiceName,
		mantraev1connect.RouterServiceName,
		mantraev1connect.ServiceServiceName,
		mantraev1connect.MiddlewareServiceName,
	}

	checker := grpchealth.NewStaticChecker(serviceNames...)
	reflector := grpcreflect.NewStaticReflector(serviceNames...)

	s.mux.Handle(grpchealth.NewHandler(checker))
	s.mux.Handle(grpcreflect.NewHandlerV1(reflector))
	s.mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	// Serve OpenAPI specs file
	s.mux.HandleFunc("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "proto/gen/openapi/openapi.yaml")
	})

	// Serve Elements UI
	s.mux.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		if _, err := w.Write([]byte(elementsHTML)); err != nil {
			slog.Error("failed to write elements HTML", "error", err)
		}
	})

	// Service implementations
	s.mux.Handle(mantraev1connect.NewProfileServiceHandler(
		service.NewProfileService(s.app),
		opts...,
	))
	s.mux.Handle(mantraev1connect.NewUserServiceHandler(
		service.NewUserService(s.app),
		opts...,
	))
	s.mux.Handle(mantraev1connect.NewEntryPointServiceHandler(
		service.NewEntryPointService(s.app),
		opts...,
	))
	s.mux.Handle(mantraev1connect.NewSettingServiceHandler(
		service.NewSettingService(s.app),
		opts...,
	))
	s.mux.Handle(mantraev1connect.NewAgentServiceHandler(
		service.NewAgentService(s.app),
		opts...,
	))
	s.mux.Handle(mantraev1connect.NewAgentManagementServiceHandler(
		service.NewAgentManagementService(s.app),
		opts...,
	))
	s.mux.Handle(mantraev1connect.NewRouterServiceHandler(
		service.NewRouterService(s.app),
		opts...,
	))
	s.mux.Handle(mantraev1connect.NewServiceServiceHandler(
		service.NewServiceService(s.app),
		opts...,
	))
	s.mux.Handle(mantraev1connect.NewMiddlewareServiceHandler(
		service.NewMiddlewareService(s.app),
		opts...,
	))

	// Traefik endpoint (HTTP) ------------------------------------------------
	mw := middlewares.NewMiddlewareHandler(s.app)
	logChain := middlewares.Chain(mw.Logger)
	basicChain := middlewares.Chain(mw.Logger, mw.BasicAuth)

	if s.SecureTraefik {
		s.mux.Handle("GET /{name}", basicChain(handler.PublishTraefikConfig(s.app)))
	} else {
		s.mux.Handle("GET /{name}", logChain(handler.PublishTraefikConfig(s.app)))
	}

	// TODO: OIDC
	// TODO: Public IP
}
