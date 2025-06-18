package middlewares

import (
	"net/http"

	"github.com/mizuchilabs/mantrae/internal/config"
)

type Middleware func(http.Handler) http.Handler

type MiddlewareHandler struct {
	app *config.App
}

// NewMiddlewareHandler creates a new middleware set with configuration
func NewMiddlewareHandler(app *config.App) *MiddlewareHandler {
	return &MiddlewareHandler{app: app}
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
