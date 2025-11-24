// Package server initializes and runs the Mantrae server.
package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"runtime/debug"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/validate"
	"github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1/mantraev1connect"
	"github.com/mizuchilabs/mantrae/server/internal/api/handler"
	"github.com/mizuchilabs/mantrae/server/internal/api/middlewares"
	"github.com/mizuchilabs/mantrae/server/internal/api/service"
	"github.com/mizuchilabs/mantrae/server/internal/config"
)

type Server struct {
	mux *http.ServeMux
	app *config.App
}

func NewServer(app *config.App) *Server {
	return &Server{
		mux: http.NewServeMux(),
		app: app,
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
		Addr:              "0.0.0.0:" + s.app.BasePort(),
		Handler:           s.WithCORS(s.mux),
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    8 * 1024, // 8KiB
	}

	// Channel to catch server errors
	serverErr := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		serverURL, ok := s.app.SM.Get(ctx, "server_url")
		if ok && serverURL == "" {
			serverURL = "http://127.0.0.1:" + s.app.BasePort()
		}
		slog.Info("Server listening on", "address", "http://127.0.0.1:"+s.app.BasePort())
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
	opts := []connect.HandlerOption{
		connect.WithCompressMinBytes(1024),
		connect.WithInterceptors(
			middlewares.Logging(),
			middlewares.TimeoutInterceptor(),
			middlewares.NewAuthInterceptor(s.app),
			middlewares.NewAuditInterceptor(s.app),
			validate.NewInterceptor(),
		),
		connect.WithRecover(
			func(ctx context.Context, spec connect.Spec, header http.Header, panic any) error {
				if s.app.Debug {
					slog.Error("panic recovered in RPC call",
						"method", spec.Procedure,
						"panic", panic,
						"stack", string(debug.Stack()),
					)
				} else {
					slog.Error("panic recovered in RPC call",
						"method", spec.Procedure,
						"panic", panic,
					)
				}
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
		mantraev1connect.ServersTransportServiceName,
		mantraev1connect.BackupServiceName,
		mantraev1connect.UtilServiceName,
		mantraev1connect.TraefikInstanceServiceName,
		mantraev1connect.AuditLogServiceName,
	}
	s.registerHealthAndReflection(serviceNames)

	// Debug
	if s.app.Debug {
		s.mux.HandleFunc("/debug/pprof/", pprof.Index)
		s.mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		s.mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		s.mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		s.mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	// Static files
	s.WithStatic()

	// Serve OpenAPI spec
	s.OpenAPIHandler()

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
	s.mux.Handle(mantraev1connect.NewServersTransportServiceHandler(
		service.NewServersTransportService(s.app),
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
	s.mux.Handle(mantraev1connect.NewTraefikInstanceServiceHandler(
		service.NewTraefikInstanceService(s.app),
		opts...,
	))
	s.mux.Handle(mantraev1connect.NewAuditLogServiceHandler(
		service.NewAuditLogService(s.app),
		opts...,
	))

	// HTTP middlewares -------------------------------------------------------
	mw := middlewares.NewMiddlewareHandler(s.app)
	logChain := middlewares.Chain(mw.Logger)
	jwtChain := middlewares.Chain(mw.Logger, mw.JWTAuth)

	// Traefik endpoint (HTTP) ------------------------------------------------
	s.mux.Handle("GET /api/{name}", logChain(handler.PublishTraefikConfig(s.app)))

	// Upload handler (HTTP) --------------------------------------------------
	s.mux.Handle("POST /upload/avatar/{id}", jwtChain(handler.UploadAvatar(s.app)))
	s.mux.Handle("POST /upload/backup/{id}", jwtChain(handler.UploadBackup(s.app)))

	// OIDC handlers (HTTP) ---------------------------------------------------
	s.mux.Handle("GET /oidc/login", logChain(handler.OIDCLogin(s.app)))
	s.mux.Handle("GET /oidc/callback", logChain(handler.OIDCCallback(s.app)))
}

func (s *Server) registerHealthAndReflection(serviceNames []string) {
	checker := grpchealth.NewStaticChecker(serviceNames...)
	reflector := grpcreflect.NewStaticReflector(serviceNames...)

	s.mux.Handle(grpchealth.NewHandler(checker))
	s.mux.Handle(grpcreflect.NewHandlerV1(reflector))
	s.mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
	s.mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}
