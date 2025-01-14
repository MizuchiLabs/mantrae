package middlewares

import (
	"net/http"

	"github.com/MizuchiLabs/mantrae/internal/config"
	"github.com/MizuchiLabs/mantrae/internal/db"
)

type Middleware func(http.Handler) http.Handler

type MiddlewareHandler struct {
	db     *db.Queries
	config config.Config
}

// NewMiddleware creates a new middleware set with configuration
func NewMiddlewareHandler(db *db.Queries, config config.Config) *MiddlewareHandler {
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
