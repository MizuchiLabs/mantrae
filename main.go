package main

import (
	"embed"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/MizuchiLabs/mantrae/api"
	"github.com/lmittmann/tint"
)

//go:embed all:web/build
var webFS embed.FS

// Set up global logger with specified configuration
func init() {
	logger := slog.New(tint.NewHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	go api.FetchTraefikConfig()
}

// statusRecorder is a wrapper around http.ResponseWriter to capture the status code
type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

// Logging middleware to log HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
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
		slog.Info("Request",
			"method", r.Method,
			"url", r.URL.Path,
			"protocol", r.Proto,
			"status", recorder.statusCode,
			"duration", duration,
		)
	})
}

func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/profiles", api.CreateProfile)
	mux.HandleFunc("GET /api/profiles", api.GetProfiles)
	mux.HandleFunc("PUT /api/profiles/{name}", api.UpdateProfile)
	mux.HandleFunc("DELETE /api/profiles/{name}", api.DeleteProfile)
	mux.HandleFunc("PUT /api/routers/{profile}/{router}", api.UpdateRouter)
	mux.HandleFunc("DELETE /api/routers/{profile}/{router}", api.DeleteRouter)
	mux.HandleFunc("PUT /api/services/{profile}/{service}", api.UpdateService)
	mux.HandleFunc("DELETE /api/services/{profile}/{service}", api.DeleteService)
	mux.HandleFunc("PUT /api/middlewares/{profile}/{middleware}", api.UpdateMiddleware)
	mux.HandleFunc("DELETE /api/middlewares/{profile}/{middleware}", api.DeleteMiddleware)
	mux.HandleFunc("GET /api/{name}", api.GetConfig)

	staticContent, err := fs.Sub(webFS, "web/build")
	if err != nil {
		slog.Error("Sub", "error", err)
		return
	}

	mux.Handle("/", http.FileServer(http.FS(staticContent)))
	cors := enableCors(loggingMiddleware(mux))

	// Start the background sync process
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go api.Sync(wg)

	log.Println("Listening on port 3000")
	if err := http.ListenAndServe(":3000", cors); err != nil {
		slog.Error("ListenAndServe", "error", err)
		return
	}
	wg.Wait()
}
