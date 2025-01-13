package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/pkg/util"
	"golang.org/x/crypto/bcrypt"
)

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
