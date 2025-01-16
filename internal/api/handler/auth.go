package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MizuchiLabs/mantrae/internal/config"
	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/internal/mail"
	"github.com/MizuchiLabs/mantrae/internal/util"
	"golang.org/x/crypto/bcrypt"
)

// Login verifies user credentials using a normal password and returns a JWT token if successful.
func Login(q *db.Queries, secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user db.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Failed to decode credentials", http.StatusBadRequest)
			return
		}

		if user.Username == "" || user.Password == "" {
			http.Error(w, "Username or password cannot be empty", http.StatusBadRequest)
			return
		}

		dbUser, err := q.GetUserByUsername(r.Context(), user.Username)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
			http.Error(w, "Invalid username or password "+err.Error(), http.StatusUnauthorized)
			return
		}

		expirationTime := time.Now().Add(24 * time.Hour)
		if r.URL.Query().Get("remember") == "true" {
			expirationTime = time.Now().Add(7 * 24 * time.Hour)
		}

		token, err := util.EncodeUserJWT(user.Username, secret, expirationTime)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := q.UpdateUserLastLogin(r.Context(), dbUser.ID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}

// VerifyToken checks the validity of a JWT token provided in cookies or Authorization header.
func VerifyToken(q *db.Queries, secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var token string

		// Check for token in cookies and Authorization header
		if cookie, err := r.Cookie("token"); err == nil {
			token = cookie.Value
		} else {
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				token = strings.TrimPrefix(authHeader, "Bearer ")
			} else {
				http.Error(w, "Token cannot be empty", http.StatusBadRequest)
				return
			}
		}

		data, err := util.DecodeUserJWT(token, secret)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid token: %s", err.Error()), http.StatusUnauthorized)
			return
		}

		user, err := q.GetUserByUsername(r.Context(), data.Username)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

// ResetPassword allows users to reset their password using a valid JWT token.
func ResetPassword(q *db.Queries, secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Failed to decode credentials", http.StatusBadRequest)
			return
		}

		tokenString := r.Header.Get("Authorization")
		if len(tokenString) < 7 {
			http.Error(w, "Token cannot be empty", http.StatusBadRequest)
			return
		}

		user, err := util.DecodeUserJWT(tokenString[7:], secret)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if user.Username == "" || data.Password == "" {
			http.Error(w, "Username or password cannot be empty", http.StatusBadRequest)
			return
		}

		dbUser, err := q.GetUserByUsername(r.Context(), user.Username)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = q.UpdateUser(r.Context(), db.UpdateUserParams{
			ID:       dbUser.ID,
			Username: dbUser.Username,
			Email:    dbUser.Email,
			Password: string(hash),
			IsAdmin:  dbUser.IsAdmin,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// SendResetEmail sends an email with a password reset link to the user's registered email address.
func SendResetEmail(q *db.Queries, secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := q.GetUserByUsername(r.Context(), r.PathValue("name"))
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if user.Email == nil {
			http.Error(w, fmt.Sprintf("%s has no email address", user.Username), http.StatusBadRequest)
			return
		}

		token, err := util.EncodeUserJWT(user.Username, secret, time.Now().Add(10*time.Minute))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		settings, err := q.ListSettings(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var config config.EmailConfig
		var resetLink string
		for _, setting := range settings {
			switch setting.Key {
			case "email-host":
				config.Host = setting.Value
			case "email-port":
				config.Port = setting.Value
			case "email-username":
				config.Username = setting.Value
			case "email-password":
				config.Password = setting.Value
			case "email-from":
				config.From = setting.Value
			case "server-url":
				resetLink = fmt.Sprintf("%s/login/reset?token=%s", setting.Value, token)
			}
		}
		data := map[string]interface{}{
			"ResetLink": resetLink,
			"Minutes":   10,
		}
		if err := mail.Send(*user.Email, "reset-password", config, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
