package api

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/pkg/util"
	"golang.org/x/crypto/bcrypt"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

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
func Log(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Capture the response status code
		recorder := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		// Serve the request
		next(recorder, r)
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
		if status >= 200 && status < 400 && util.App.LogLevel == "debug" {
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

// Cors middleware to enable CORS
func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// BasicAuth middleware to authenticate requests
func BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user, err := db.Query.GetUserByUsername(context.Background(), username)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

// JWT middleware to authenticate requests
func JWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var token string

		// Check for token in cookies
		if cookie, err := r.Cookie("token"); err == nil {
			token = cookie.Value
		} else {
			// If no cookie, check for token in Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				token = strings.TrimPrefix(authHeader, "Bearer ")
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		// Validate JWT token and decode claims
		claims, err := util.DecodeUserJWT(token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add claims to request context for access in subsequent middlewares
		ctx := context.WithValue(r.Context(), "username", claims.Username)
		r = r.WithContext(ctx)

		next(w, r)
	}
}

func AdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the username from the request context
		username, ok := r.Context().Value("username").(string)
		if !ok || username == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := db.Query.GetUserByUsername(context.Background(), username)
		if err != nil {
			http.Error(w, "Admin privileges required", http.StatusForbidden)
			return
		}
		if !user.IsAdmin {
			http.Error(w, "Admin privileges required", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}

// Chain middlewares
func Chain(middlewares ...Middleware) Middleware {
	return func(final http.HandlerFunc) http.HandlerFunc {
		for _, middleware := range middlewares {
			final = middleware(final)
		}
		return final
	}
}
