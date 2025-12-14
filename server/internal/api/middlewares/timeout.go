package middlewares

import (
	"context"
	"time"

	"connectrpc.com/connect"
)

const (
	// Timeout is the default timeout for all RPC calls
	Timeout = 5 * time.Second
)

// NewTimeoutInterceptor adds a timeout to all RPC calls
func NewTimeoutInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// Skip if client already provided a deadline
			if _, hasDeadline := ctx.Deadline(); hasDeadline {
				return next(ctx, req)
			}

			// Apply timeout
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			return next(ctx, req)
		}
	}
}
