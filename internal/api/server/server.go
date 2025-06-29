// Package server initializes and runs the Mantrae server.
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
	"connectrpc.com/validate"
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
	defer func() {
		if err := s.app.Conn.Close(); err != nil {
			slog.Error("failed to close database connection", "error", err)
		}
	}()

	server := &http.Server{
		Addr:              s.Host + ":" + s.Port,
		Handler:           s.WithCORS(s.mux),
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
		MaxHeaderBytes:    8 * 1024, // 8KiB
	}

	// Channel to catch server errors
	serverErr := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		serverURL, ok := s.app.SM.Get("server_url")
		if ok && serverURL == "" {
			serverURL = s.Host + ":" + s.Port
		}
		slog.Info("Server listening on", "address", "127.0.0.1:"+s.Port)
		slog.Info("Agents can connect to", "address", serverURL)
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
	// Protovalidator
	validator, err := validate.NewInterceptor()
	if err != nil {
		log.Fatal(err)
	}

	// Common interceptors
	opts := []connect.HandlerOption{
		connect.WithCompressMinBytes(1024),
		connect.WithInterceptors(
			middlewares.Logging(),
			middlewares.NewAuthInterceptor(s.app),
			middlewares.NewAuditInterceptor(s.app),
			validator,
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
				header.Set("Content-Type", "application/json")
				return connect.NewError(connect.CodeInternal, fmt.Errorf("internal server error"))
			},
		),
	}

	serviceNames := []string{
		mantraev1connect.ProfileServiceName,
		mantraev1connect.UserServiceName,
		mantraev1connect.EntryPointServiceName,
		mantraev1connect.SettingServiceName,
		mantraev1connect.DnsProviderServiceName,
		mantraev1connect.AgentServiceName,
		mantraev1connect.RouterServiceName,
		mantraev1connect.ServiceServiceName,
		mantraev1connect.MiddlewareServiceName,
		mantraev1connect.BackupServiceName,
		mantraev1connect.UtilServiceName,
		mantraev1connect.AuditLogServiceName,
		// mantraev1connect.EventServiceName,
	}

	checker := grpchealth.NewStaticChecker(serviceNames...)
	reflector := grpcreflect.NewStaticReflector(serviceNames...)

	s.mux.Handle(grpchealth.NewHandler(checker))
	s.mux.Handle(grpcreflect.NewHandlerV1(reflector))
	s.mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	// Static files
	staticContent, err := fs.Sub(web.StaticFS, "build")
	if err != nil {
		log.Fatal(err)
	}
	uploadsContent := http.FileServer(http.Dir("./data/uploads"))
	s.mux.Handle("/", http.FileServer(http.FS(staticContent)))
	s.mux.Handle("/uploads/", http.StripPrefix("/uploads/", uploadsContent))

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
	s.mux.Handle(mantraev1connect.NewDnsProviderServiceHandler(
		service.NewDnsProviderService(s.app),
		opts...,
	))
	s.mux.Handle(mantraev1connect.NewAgentServiceHandler(
		service.NewAgentService(s.app),
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
	s.mux.Handle(mantraev1connect.NewBackupServiceHandler(
		service.NewBackupService(s.app),
		opts...,
	))
	s.mux.Handle(mantraev1connect.NewUtilServiceHandler(
		service.NewUtilService(s.app),
		opts...,
	))
	s.mux.Handle(mantraev1connect.NewAuditLogServiceHandler(
		service.NewAuditLogService(s.app),
		opts...,
	))
	// s.mux.Handle(mantraev1connect.NewEventServiceHandler(
	// 	service.NewEventService(s.app),
	// 	opts...,
	// ))

	// Traefik endpoint (HTTP) ------------------------------------------------
	mw := middlewares.NewMiddlewareHandler(s.app)
	logChain := middlewares.Chain(mw.Logger)
	basicChain := middlewares.Chain(mw.Logger, mw.BasicAuth)
	jwtChain := middlewares.Chain(mw.Logger, mw.JWTAuth)

	if s.SecureTraefik {
		s.mux.Handle("GET /api/{name}", basicChain(handler.PublishTraefikConfig(s.app)))
	} else {
		s.mux.Handle("GET /api/{name}", logChain(handler.PublishTraefikConfig(s.app)))
	}

	// Upload handler (HTTP) --------------------------------------------------
	s.mux.Handle("POST /upload/avatar", jwtChain(handler.UploadAvatar(s.app)))
	s.mux.Handle("POST /upload/backup", jwtChain(handler.UploadBackup(s.app)))

	// OIDC handlers (HTTP) ---------------------------------------------------
	s.mux.Handle("GET /oidc/login", logChain(handler.OIDCLogin(s.app)))
	s.mux.Handle("GET /oidc/callback", logChain(handler.OIDCCallback(s.app)))
}
