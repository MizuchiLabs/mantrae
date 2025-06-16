package server

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"github.com/mizuchilabs/mantrae/internal/api/middlewares"
	"github.com/mizuchilabs/mantrae/internal/api/service"
	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/util"
	"github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1/mantraev1connect"
	"github.com/mizuchilabs/mantrae/web"
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
	// Start the event processor before registering services
	util.StartEventProcessor(ctx)

	s.registerServices()
	defer s.app.Conn.Close()
	host := s.app.Config.Server.Host
	port := s.app.Config.Server.Port
	allowedOrigins := s.getAllowedOrigins(ctx)

	server := &http.Server{
		Addr:              host + ":" + port,
		Handler:           middlewares.CORS(allowedOrigins...)(s.mux),
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
		MaxHeaderBytes:    8 * 1024, // 8KiB
	}

	// Channel to catch server errors
	serverErr := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		slog.Info("Server starting", "host", host, "port", port)
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

func (s *Server) getAllowedOrigins(ctx context.Context) []string {
	var origins []string

	// Always allow development frontend
	devOrigin := "http://127.0.0.1:5173"
	origins = append(origins, devOrigin)

	// Get server URL from settings for production
	if serverURL, err := s.app.SM.Get(ctx, "server_url"); err == nil {
		if url := serverURL.String(""); url != "" && url != devOrigin {
			origins = append(origins, strings.TrimSuffix(url, "/"))
		}
	}

	// Remove duplicates
	seen := make(map[string]bool)
	var uniqueOrigins []string
	for _, origin := range origins {
		if !seen[origin] {
			uniqueOrigins = append(uniqueOrigins, origin)
			seen[origin] = true
		}
	}

	slog.Debug("CORS allowed origins", "origins", uniqueOrigins)
	return uniqueOrigins
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

	// Routes
	s.routes()

	serviceNames := []string{
		mantraev1connect.ProfileServiceName,
		mantraev1connect.UserServiceName,
		mantraev1connect.EntryPointServiceName,
		mantraev1connect.SettingServiceName,
		mantraev1connect.AgentServiceName,
	}

	checker := grpchealth.NewStaticChecker(serviceNames...)
	reflector := grpcreflect.NewStaticReflector(serviceNames...)

	s.mux.Handle(grpchealth.NewHandler(checker))
	s.mux.Handle(grpcreflect.NewHandlerV1(reflector))
	s.mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

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
}
