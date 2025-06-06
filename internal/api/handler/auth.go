package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MizuchiLabs/mantrae/internal/api/middlewares"
	"github.com/MizuchiLabs/mantrae/internal/config"
	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/internal/mail"
	"github.com/MizuchiLabs/mantrae/internal/settings"
	"github.com/MizuchiLabs/mantrae/internal/util"
	"golang.org/x/crypto/bcrypt"
)

// Login verifies user credentials using a normal password and returns a JWT token if successful.
func Login(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if password login is enabled
		enabled, err := a.SM.Get(r.Context(), settings.KeyPasswordLoginDisabled)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if enabled.Bool(false) {
			http.Error(w, "Password login is disabled", http.StatusUnauthorized)
			return
		}

		q := a.Conn.GetQuery()
		var request db.User
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Failed to decode credentials", http.StatusBadRequest)
			return
		}

		if request.Username == "" || request.Password == "" {
			http.Error(w, "Username or password cannot be empty", http.StatusBadRequest)
			return
		}

		user, err := q.GetUserByUsername(r.Context(), request.Username)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		userPassword, err := q.GetUserPassword(r.Context(), user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(request.Password)); err != nil {
			http.Error(w, "Invalid username or password.", http.StatusUnauthorized)
			return
		}

		remember := r.URL.Query().Get("remember") == "true"
		expirationTime := time.Now().Add(24 * time.Hour)
		if remember {
			expirationTime = time.Now().Add(30 * 24 * time.Hour)
		}

		jwt, err := util.EncodeUserJWT(request.Username, a.Config.Secret, expirationTime)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := q.UpdateUserLastLogin(r.Context(), user.ID); err != nil {
			fmt.Printf("Failed to update last login for user %s: %v\n", user.Username, err)
		}

		http.SetCookie(w, &http.Cookie{
			Name:     util.CookieName,
			Value:    jwt,
			Path:     "/",
			MaxAge:   int(expirationTime.Unix() - time.Now().Unix()),
			HttpOnly: true,
			Secure:   r.TLS != nil,
			SameSite: http.SameSiteLaxMode,
		})

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     util.CookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})
}

// Verify returns the currently logged in user
func Verify(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middlewares.AuthUserKey).(db.GetUserByUsernameRow)

	response := map[string]any{"user": user}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
		jwt, err := util.EncodeUserJWT(user.Username, a.Config.Secret, expirationTime)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := q.UpdateUserLastLogin(r.Context(), user.ID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     util.CookieName,
			Value:    jwt,
			Path:     "/",
			MaxAge:   int(expirationTime.Unix() - time.Now().Unix()),
			HttpOnly: true,
			Secure:   r.TLS != nil,
			SameSite: http.SameSiteLaxMode,
		})

		response := map[string]any{"user": user}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
		expiresAt := time.Now().Add(15 * time.Minute)
		token, err := util.GenerateOTP()
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		sets, err := a.SM.GetAll(r.Context())
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

		var config mail.EmailConfig
		for _, s := range sets {
			switch s.Value {
			case settings.KeyEmailHost:
				config.Host = s.Value.(string)
			case settings.KeyEmailPort:
				config.Port = s.Value.(string)
			case settings.KeyEmailUser:
				config.Username = s.Value.(string)
			case settings.KeyEmailPassword:
				config.Password = s.Value.(string)
			case settings.KeyEmailFrom:
				config.From = s.Value.(string)
			}
		}
		data := map[string]any{
			"Token": token,
			"Date":  expiresAt.Format("Jan 2, 2006 at 15:04"),
		}
		if err := mail.Send(*user.Email, "reset-password", config, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
