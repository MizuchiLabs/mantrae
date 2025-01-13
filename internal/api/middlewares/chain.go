package middlewares

import "net/http"

// Chain middlewares
func Chain(middlewares ...Middleware) Middleware {
	return func(final http.HandlerFunc) http.HandlerFunc {
		for _, middleware := range middlewares {
			final = middleware(final)
		}
		return final
	}
}
