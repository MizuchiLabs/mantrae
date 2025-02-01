package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MizuchiLabs/mantrae/internal/app"
	"github.com/MizuchiLabs/mantrae/internal/config"
	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/internal/mail"
	"github.com/MizuchiLabs/mantrae/internal/util"
	"golang.org/x/crypto/bcrypt"
)

// Login verifies user credentials using a normal password and returns a JWT token if successful.
func Login(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := a.Conn.GetQuery()
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
		if err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
			http.Error(w, "Invalid username or password "+err.Error(), http.StatusUnauthorized)
			return
		}

		expirationTime := time.Now().Add(24 * time.Hour)
		if r.URL.Query().Get("remember") == "true" {
			expirationTime = time.Now().Add(7 * 24 * time.Hour)
		}

		token, err := util.EncodeUserJWT(user.Username, a.Config.Secret, expirationTime)
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

// VerifyJWT checks the validity of a JWT token provided in cookies or Authorization header.
func VerifyJWT(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := a.Conn.GetQuery()
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

		data, err := util.DecodeUserJWT(token, a.Config.Secret)
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

// VerifyOTP allows users to login using an OTP token, to reset their password
func VerifyOTP(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := a.Conn.GetQuery()
		var data struct {
			Username string `json:"username"`
			Token    string `json:"token"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Failed to decode credentials", http.StatusBadRequest)
			return
		}

		user, err := q.GetUserByUsername(r.Context(), data.Username)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// Verify token
		if user.Otp == nil || user.OtpExpiry == nil {
			http.Error(w, "No reset token found", http.StatusUnauthorized)
			return
		}

		if *user.Otp != data.Token {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if time.Now().After(*user.OtpExpiry) {
			http.Error(w, "Token expired", http.StatusUnauthorized)
			return
		}

		expirationTime := time.Now().Add(1 * time.Hour)
		token, err := util.EncodeUserJWT(user.Username, a.Config.Secret, expirationTime)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := q.UpdateUserLastLogin(r.Context(), user.ID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}

// SendResetEmail sends an email with a password reset link to the user's registered email address.
func SendResetEmail(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := a.Conn.GetQuery()
		user, err := q.GetUserByUsername(r.Context(), r.PathValue("name"))
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if user.Email == nil {
			http.Error(
				w,
				fmt.Sprintf("%s has no email address", user.Username),
				http.StatusBadRequest,
			)
			return
		}

		// Generate OTP
		expiresAt := time.Now().Add(10 * time.Minute)
		token, err := util.GenerateOTP()
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		settings, err := q.ListSettings(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = q.UpdateUserResetToken(r.Context(), db.UpdateUserResetTokenParams{
			ID:        user.ID,
			Otp:       &token,
			OtpExpiry: &expiresAt,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var config app.EmailConfig
		for _, setting := range settings {
			switch setting.Key {
			case "email_host":
				config.Host = setting.Value
			case "email_port":
				config.Port = setting.Value
			case "email_username":
				config.Username = setting.Value
			case "email_password":
				config.Password = setting.Value
			case "email_from":
				config.From = setting.Value
			}
		}
		data := map[string]interface{}{
			"Token": token,
			"Date":  expiresAt.Format(time.RFC3339),
		}
		if err := mail.Send(*user.Email, "reset-password", config, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
