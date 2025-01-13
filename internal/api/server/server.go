package server

import (
	"context"
	"database/sql"
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
	"github.com/MizuchiLabs/mantrae/agent/proto/gen/agent/v1/agentv1connect"
	"github.com/MizuchiLabs/mantrae/internal/api/grpc"
	"github.com/MizuchiLabs/mantrae/internal/api/middlewares"
	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/pkg/util"
	"github.com/MizuchiLabs/mantrae/web"
)

type Server struct {
	mux *http.ServeMux
	db  *db.Queries
}

func NewServer(DB *sql.DB) *Server {
	return &Server{
		mux: http.NewServeMux(),
		db:  db.New(DB),
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.registerServices()

	server := &http.Server{
		Addr:              util.App.Host + ":" + util.App.Port,
		Handler:           middlewares.CORS(s.mux),
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
		MaxHeaderBytes:    8 * 1024, // 8KiB
	}

	// Channel to catch server errors
	serverErr := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		slog.Info("Server starting",
			"host", util.App.Host,
			"port", util.App.Port,
		)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// Wait for context cancellation or server error
	select {
	case <-ctx.Done():
		slog.Info("Shutting down server gracefully...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("Server shutdown error", "error", err)
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
		// connect.WithInterceptors(
		// 	middleware.Authentication(),
		// 	middleware.Logging(),
		// ),
		connect.WithRecover(func(ctx context.Context, spec connect.Spec, header http.Header, panic any) error {
			// Log the panic with context
			slog.Error("panic recovered in RPC call",
				"method", spec.Procedure,
				"panic", panic,
				"trace", string(debug.Stack()),
			)
			header.Set("X-Error-Type", "panic")
			return connect.NewError(connect.CodeInternal, fmt.Errorf("internal server error"))
		}),
	}

	// Static files
	staticContent, err := fs.Sub(web.StaticFS, "build")
	if err != nil {
		log.Fatal(err)
	}
	s.mux.Handle("/", http.FileServer(http.FS(staticContent)))

	// Routes
	s.routes()

	// Reflection service
	reflector := grpcreflect.NewStaticReflector(
		"agent.v1.AgentService",
	)
	// Health check service
	checker := grpchealth.NewStaticChecker(
		"agent.v1.AgentService",
	)

	s.mux.Handle(grpchealth.NewHandler(checker))
	s.mux.Handle(grpcreflect.NewHandlerV1(reflector))
	s.mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	// // Serve OpenAPI specs file
	// fsHandler := http.FileServer(http.Dir("proto/gen/openapi"))
	// s.mux.Handle("/openapi/", http.StripPrefix("/openapi/", fsHandler))

	// // Serve Swagger UI
	// s.mux.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
	// 	httpSwagger.Handler(
	// 		httpSwagger.URL("/openapi/api.swagger.yaml"),
	// 	).ServeHTTP(w, r)
	// })

	// Service implementations
	s.mux.Handle(agentv1connect.NewAgentServiceHandler(&grpc.AgentServer{}, opts...))
}
