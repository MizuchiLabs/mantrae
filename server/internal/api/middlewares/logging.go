package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"connectrpc.com/connect"
)

// statusRecorder is a wrapper around http.ResponseWriter to capture the status code
type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code and writes the header
func (rec *statusRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

// Implement the http.Flusher interface to forward Flush calls to the underlying ResponseWriter
func (rec *statusRecorder) Flush() {
	// Check if the underlying ResponseWriter supports flushing
	if flusher, ok := rec.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Logger middleware to log HTTP requests
func (h *MiddlewareHandler) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/_app/") || r.URL.Path == "/favicon.ico" {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rec, r)
		duration := time.Since(start)

		level := slog.LevelDebug
		switch {
		case rec.statusCode >= 500:
			level = slog.LevelError
		case rec.statusCode >= 400:
			level = slog.LevelWarn
		}

		slog.Log(r.Context(), level, "http_request",
			slog.String("method", r.Method),
			slog.String("url", r.URL.Path),
			slog.Int("status", rec.statusCode),
			slog.String("protocol", r.Proto),
			slog.Int64("duration_ms", duration.Milliseconds()),
		)
	})
}

func Logging() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if req.Spec().Procedure == "HealthCheck" {
				return next(ctx, req)
			}
			start := time.Now()
			resp, err := next(ctx, req)
			duration := time.Since(start)

			logger := slog.With(
				slog.String("method", req.Spec().Procedure),
				slog.String("peer", req.Peer().Addr),
				slog.String("protocol", req.Peer().Protocol),
				slog.Int64("duration_ms", duration.Milliseconds()),
			)

			if err != nil {
				logger.With(slog.String("error", err.Error())).Error("rpc_call")
			} else {
				logger.Debug("rpc_call")
			}

			return resp, err
		}
	}
}
