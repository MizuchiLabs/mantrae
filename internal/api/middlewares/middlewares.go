package middlewares

import (
	"database/sql"
	"net/http"

	"github.com/MizuchiLabs/mantrae/internal/app"
)

type Middleware func(http.Handler) http.Handler

type MiddlewareHandler struct {
	db     *sql.DB
	config app.Config
}

// NewMiddleware creates a new middleware set with configuration
func NewMiddlewareHandler(db *sql.DB, config app.Config) *MiddlewareHandler {
	return &MiddlewareHandler{
		db:     db,
		config: config,
	}
}

// Chain combines multiple middlewares into a single middleware
func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}
