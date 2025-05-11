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

// Log middleware to log HTTP requests
func (h *MiddlewareHandler) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Capture the response status code
		recorder := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		// Serve the request
		next.ServeHTTP(recorder, r)
		duration := time.Since(start)

		if strings.HasPrefix(r.URL.Path, "/_app/") {
			return
		}

		// Log the request details
		status := recorder.statusCode

		msg := "HTTP request"
		fields := []any{
			"method", r.Method,
			"url", r.URL.Path,
			"status", status,
			"protocol", r.Proto,
			"duration_ms", duration.Milliseconds(),
		}

		switch {
		case status >= 500:
			slog.Error(msg, fields...)
		case status >= 400:
			slog.Warn(msg, fields...)
		default:
			slog.Info(msg, fields...)
		}
	})
}

func Logging() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			start := time.Now()
			resp, err := next(ctx, req)
			duration := time.Since(start)

			if req.Spec().Procedure == "HealthCheck" {
				return resp, err
			}

			msg := "RPC call"
			fields := []any{
				"method", req.Spec().Procedure,
				"peer", req.Peer().Addr,
				"protocol", req.Peer().Protocol,
				"duration_ms", duration.Milliseconds(),
			}

			if err != nil {
				slog.Error(msg, append(fields, "error", err)...)
			} else {
				slog.Debug(msg, fields...)
			}

			return resp, err
		}
	}
}
