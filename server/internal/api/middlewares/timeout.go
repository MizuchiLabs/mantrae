package middlewares

import (
	"context"
	"net/http"
	"time"

	"connectrpc.com/connect"
)

const (
	// Timeout is the default timeout for all RPC calls
	Timeout = 5 * time.Second
)

// TimeoutInterceptor adds a timeout to all RPC calls
func TimeoutInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			ctx, cancel := context.WithTimeout(ctx, Timeout)
			defer cancel()

			return next(ctx, req)
		}
	}
}

// HTTPTimeoutMiddleware adds timeout to HTTP handlers
func HTTPTimeoutMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), Timeout)
			defer cancel()

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
