package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/internal/util"
	"golang.org/x/crypto/bcrypt"
)

type ctxKey string

const (
	AuthUserKey ctxKey = "user"
)

// BasicAuth middleware for simple authentication
func (h *MiddlewareHandler) BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		q := h.app.Conn.GetQuery()
		user, err := q.GetUserByUsername(r.Context(), username)
		if err != nil {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userPassword, err := q.GetUserPassword(r.Context(), user.ID)
		if err != nil {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password)); err != nil {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), AuthUserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// JWT authentication middleware
func (h *MiddlewareHandler) JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		claims, err := util.DecodeUserJWT(tokenString, h.app.Config.Secret)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Verify user exists in database
		q := h.app.Conn.GetQuery()
		user, err := q.GetUserByUsername(r.Context(), claims.Username)
		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), AuthUserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *MiddlewareHandler) AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the username from the request context
		user, ok := r.Context().Value(AuthUserKey).(db.GetUserByUsernameRow)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !user.IsAdmin {
			http.Error(w, "Admin privileges required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
