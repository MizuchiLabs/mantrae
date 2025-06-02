package middlewares

import (
	"net/http"
	"strings"
)

const (
	defaultDevOrigin = "http://127.0.0.1:5173"
)

func CORS(allowedOrigins ...string) func(http.Handler) http.Handler {
	// Default to dev origin if none provided
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{defaultDevOrigin}
	}

	// Create a map for faster origin lookup
	originMap := make(map[string]bool)
	for _, origin := range allowedOrigins {
		originMap[strings.TrimSpace(origin)] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Set CORS headers
			if origin != "" && originMap[origin] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().
				Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
