package middlewares

import (
	"log/slog"
	"net/http"
	"strings"
	"time"
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

		status := recorder.statusCode

		// Log the request details
		if status >= 500 {
			slog.Error("Request",
				"method", r.Method,
				"url", r.URL.Path,
				"protocol", r.Proto,
				"status", recorder.statusCode,
				"duration", duration,
			)
		}
		if status >= 400 && status < 500 {
			slog.Warn("Request",
				"method", r.Method,
				"url", r.URL.Path,
				"protocol", r.Proto,
				"status", recorder.statusCode,
				"duration", duration,
			)
			return
		}
		if status >= 200 && status < 400 {
			slog.Info("Request",
				"method", r.Method,
				"url", r.URL.Path,
				"protocol", r.Proto,
				"status", recorder.statusCode,
				"duration", duration,
			)
			return
		}
	})
}
