package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/settings"
	"github.com/mizuchilabs/mantrae/internal/store/db"
	"github.com/mizuchilabs/mantrae/pkg/meta"
	"golang.org/x/oauth2"
)

type OIDCConfig struct {
	Enabled      bool     `json:"enabled"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	RedirectURL  string   `json:"redirect_url"`
	IssuerURL    string   `json:"issuer_url"`
	Scopes       []string `json:"scopes"`
	Provider     string   `json:"provider"`
	UsePKCE      bool     `json:"use_pkce"`
}

type OIDCUserInfo struct {
	Sub               string `json:"sub"`
	Email             string `json:"email"`
	EmailVerified     bool   `json:"email_verified"`
	Name              string `json:"name"`
	PreferredUsername string `json:"preferred_username"`
}

func OIDCLogin(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oauth2Config, _, err := getOIDCConfig(r.Context(), r, a)
		if err != nil {
			http.Error(w, "OIDC not configured: "+err.Error(), http.StatusServiceUnavailable)
			return
		}

		// Generate state for CSRF protection
		state, err := generateRandomState()
		if err != nil {
			http.Error(w, "Failed to generate state", http.StatusInternalServerError)
			return
		}

		// Store state in cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "oauth_state",
			Value:    state,
			Path:     "/",
			MaxAge:   600,
			HttpOnly: true,
			Secure:   r.TLS != nil,
			SameSite: http.SameSiteLaxMode,
		})

		opts := []oauth2.AuthCodeOption{oauth2.AccessTypeOffline}

		// Add PKCE if enabled
		if oauth2Config.ClientSecret == "" {
			verifier := oauth2.GenerateVerifier()
			http.SetCookie(w, &http.Cookie{
				Name:     "pkce_verifier",
				Value:    verifier,
				Path:     "/",
				MaxAge:   600,
				HttpOnly: true,
				Secure:   r.TLS != nil,
				SameSite: http.SameSiteLaxMode,
			})
			opts = append(opts, oauth2.S256ChallengeOption(verifier))
		}

		authURL := oauth2Config.AuthCodeURL(state, opts...)
		http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
	}
}

func OIDCCallback(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oauth2Config, verifier, err := getOIDCConfig(r.Context(), r, a)
		if err != nil {
			http.Error(w, "OIDC not configured: "+err.Error(), http.StatusServiceUnavailable)
			return
		}

		// Verify state
		stateCookie, err := r.Cookie("oauth_state")
		if err != nil || stateCookie.Value != r.URL.Query().Get("state") {
			http.Error(w, "Invalid state parameter", http.StatusBadRequest)
			return
		}

		// Clear state cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "oauth_state",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})

		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "No authorization code received", http.StatusBadRequest)
			return
		}

		opts := []oauth2.AuthCodeOption{}

		// Handle PKCE
		if oauth2Config.ClientSecret == "" {
			verifierCookie, err := r.Cookie("pkce_verifier")
			if err != nil {
				http.Error(w, "PKCE verifier not found", http.StatusBadRequest)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "pkce_verifier",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
			})

			opts = append(opts, oauth2.VerifierOption(verifierCookie.Value))
		}

		// Exchange code for token
		token, err := oauth2Config.Exchange(r.Context(), code, opts...)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("Token exchange failed: %v", err),
				http.StatusInternalServerError,
			)
			return
		}

		// Verify ID token
		idToken, ok := token.Extra("id_token").(string)
		if !ok {
			http.Error(w, "No id_token in response", http.StatusInternalServerError)
			return
		}

		verifiedToken, err := verifier.Verify(r.Context(), idToken)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("Token verification failed: %v", err),
				http.StatusInternalServerError,
			)
			return
		}

		// Extract user info from verified token
		var userInfo OIDCUserInfo
		if err := verifiedToken.Claims(&userInfo); err != nil {
			http.Error(
				w,
				fmt.Sprintf("Failed to parse claims: %v", err),
				http.StatusInternalServerError,
			)
			return
		}

		// Find or create user
		q := a.Conn.GetQuery()
		user, err := findOrCreateOIDCUser(r.Context(), q, &userInfo)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("Failed to process user: %v", err),
				http.StatusInternalServerError,
			)
			return
		}

		// Generate JWT
		expirationTime := time.Now().Add(24 * time.Hour)
		jwtToken, err := meta.EncodeUserToken(user.ID, a.Secret, expirationTime)
		if err != nil {
			http.Error(w, "Failed to generate JWT", http.StatusInternalServerError)
			return
		}

		if err := q.UpdateUserLastLogin(r.Context(), user.ID); err != nil {
			fmt.Printf("Failed to update last login for user %s: %v\n", user.Username, err)
		}

		http.SetCookie(w, &http.Cookie{
			Name:     meta.CookieName,
			Value:    jwtToken,
			Path:     "/",
			MaxAge:   int(expirationTime.Unix() - time.Now().Unix()),
			HttpOnly: true,
			Secure:   r.TLS != nil,
			SameSite: http.SameSiteLaxMode,
		})
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}

func getOIDCConfig(
	ctx context.Context,
	r *http.Request,
	a *config.App,
) (*oauth2.Config, *oidc.IDTokenVerifier, error) {
	sets := a.SM.GetAll()

	// Parse settings (same as before but simplified validation)
	if enabled, ok := sets[settings.KeyOIDCEnabled]; ok {
		if !settings.AsBool(enabled) {
			return nil, nil, errors.New("oidc disabled")
		}
	}

	config := &oauth2.Config{}
	if clientID, ok := sets[settings.KeyOIDCClientID]; ok {
		config.ClientID = clientID
	}
	if clientSecret, ok := sets[settings.KeyOIDCClientSecret]; ok {
		config.ClientSecret = clientSecret
	}
	if pkce, ok := sets[settings.KeyOIDCPKCE]; ok {
		if settings.AsBool(pkce) {
			config.ClientSecret = ""
		}
	}

	config.RedirectURL = getRedirectURL(r)
	issuerURL, ok := sets[settings.KeyOIDCIssuerURL]
	if !ok {
		return nil, nil, errors.New("oidc issuer url not set")
	}
	provider, err := oidc.NewProvider(ctx, strings.TrimSuffix(issuerURL, "/"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}
	config.Endpoint = provider.Endpoint()

	config.Scopes = []string{"openid", "email", "profile"}
	if scopes, exists := sets["oauth_scopes"]; exists && scopes != "" {
		config.Scopes = strings.Split(scopes, ",")
		for i := range config.Scopes {
			config.Scopes[i] = strings.TrimSpace(config.Scopes[i])
		}
	}

	// Create ID token verifier
	verifier := provider.Verifier(&oidc.Config{
		ClientID: config.ClientID,
	})

	return config, verifier, nil
}

func getRedirectURL(r *http.Request) string {
	proto := "http"
	if r.TLS != nil {
		proto = "https"
	} else if forwarded := r.Header.Get("X-Forwarded-Proto"); forwarded != "" {
		proto = forwarded
	}

	// check url for redirect url
	if redirectURL := r.URL.Query().Get("redirect"); redirectURL != "" {
		return redirectURL
	}

	host := r.Host
	if forwardedHost := r.Header.Get("X-Forwarded-Host"); forwardedHost != "" {
		host = forwardedHost
	}
	return fmt.Sprintf("%s://%s/oidc/callback", proto, host)
}

func findOrCreateOIDCUser(
	ctx context.Context,
	q *db.Queries,
	userInfo *OIDCUserInfo,
) (*db.User, error) {
	var user *db.User

	// Try to find existing user by email or username
	if userInfo.Email != "" {
		// First try to find by email
		if existingUser, emailErr := q.GetUserByEmail(ctx, &userInfo.Email); emailErr == nil {
			user = &db.User{
				ID:       existingUser.ID,
				Username: existingUser.Username,
				Email:    existingUser.Email,
				IsAdmin:  existingUser.IsAdmin,
			}
		}
	}

	if user == nil && userInfo.PreferredUsername != "" {
		// Try to find by username
		if existingUser, usernameErr := q.GetUserByUsername(ctx, userInfo.PreferredUsername); usernameErr == nil {
			user = &db.User{
				ID:       existingUser.ID,
				Username: existingUser.Username,
				Email:    existingUser.Email,
				IsAdmin:  existingUser.IsAdmin,
			}
		}
	}

	if user == nil {
		// Create new user
		username := userInfo.PreferredUsername
		if username == "" {
			// Fallback to email prefix if no preferred username
			if userInfo.Email != "" {
				username = strings.Split(userInfo.Email, "@")[0]
			} else {
				username = fmt.Sprintf("oidc_user_%s", userInfo.Sub)
			}
		}

		// Ensure username is unique
		originalUsername := username
		counter := 1
		for {
			if _, err := q.GetUserByUsername(ctx, username); err != nil {
				break // Username is available
			}
			username = fmt.Sprintf("%s_%d", originalUsername, counter)
			counter++
		}

		params := db.CreateUserParams{
			Username: username,
			Email:    &userInfo.Email,
			IsAdmin:  false,
		}

		newUser, err := q.CreateUser(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to create OIDC user: %w", err)
		}
		user = &db.User{
			ID:       newUser.ID,
			Username: username,
			Email:    &userInfo.Email,
			IsAdmin:  false,
		}
	} else {
		// Update existing user's email if verified
		if userInfo.Email != "" && userInfo.EmailVerified {
			if _, err := q.UpdateUser(ctx, db.UpdateUserParams{
				ID:       user.ID,
				Username: user.Username,
				Email:    &userInfo.Email,
				IsAdmin:  user.IsAdmin,
			}); err != nil {
				return nil, fmt.Errorf("failed to update user email: %w", err)
			}
			user.Email = &userInfo.Email
		}
	}

	return user, nil
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
