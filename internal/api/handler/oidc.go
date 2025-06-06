package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/MizuchiLabs/mantrae/internal/config"
	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/internal/settings"
	"github.com/MizuchiLabs/mantrae/internal/util"
	"github.com/coreos/go-oidc/v3/oidc"
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
		oidcConfig, oauth2Config, _, err := setupOIDCConfig(r.Context(), a)
		if err != nil {
			http.Error(w, "OIDC not configured: "+err.Error(), http.StatusServiceUnavailable)
			return
		}

		if !oidcConfig.Enabled {
			http.Error(w, "OIDC disabled", http.StatusServiceUnavailable)
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
		if oidcConfig.UsePKCE {
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
		oidcConfig, oauth2Config, verifier, err := setupOIDCConfig(r.Context(), a)
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
		if oidcConfig.UsePKCE {
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

		// Generate JWT and set cookie
		expirationTime := time.Now().Add(24 * time.Hour)
		jwtToken, err := util.EncodeUserJWT(user.Username, a.Config.Secret, expirationTime)
		if err != nil {
			http.Error(w, "Failed to generate JWT", http.StatusInternalServerError)
			return
		}

		if err := q.UpdateUserLastLogin(r.Context(), user.ID); err != nil {
			fmt.Printf("Failed to update last login for user %s: %v\n", user.Username, err)
		}

		http.SetCookie(w, &http.Cookie{
			Name:     util.CookieName,
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

func OIDCStatus(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oidcConfig, _, _, err := setupOIDCConfig(r.Context(), a)
		if err != nil {
			slog.Error("Failed to get OIDC config", "error", err)
		}

		response := map[string]interface{}{
			"enabled":       false,
			"provider":      "",
			"loginDisabled": false,
		}

		if err == nil && oidcConfig != nil {
			providerName, _ := a.SM.Get(r.Context(), settings.KeyOIDCProviderName)
			pwLogin, _ := a.SM.Get(r.Context(), settings.KeyPasswordLoginDisabled)
			response["enabled"] = oidcConfig.Enabled
			response["provider"] = providerName.String("")
			response["loginDisabled"] = pwLogin.Bool(false)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// Helper function that handles both config and OIDC setup
func setupOIDCConfig(
	ctx context.Context,
	a *config.App,
) (*OIDCConfig, *oauth2.Config, *oidc.IDTokenVerifier, error) {
	config, err := getOIDCConfig(ctx, a)
	if err != nil {
		return nil, nil, nil, err
	}

	// Create OIDC provider
	provider, err := oidc.NewProvider(ctx, strings.TrimSuffix(config.IssuerURL, "/"))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	// Create OAuth2 config
	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       config.Scopes,
	}

	// For PKCE, don't include client secret
	if config.UsePKCE {
		oauth2Config.ClientSecret = ""
	}

	// Create ID token verifier
	verifier := provider.Verifier(&oidc.Config{
		ClientID: config.ClientID,
	})

	return config, oauth2Config, verifier, nil
}

func getOIDCConfig(ctx context.Context, a *config.App) (*OIDCConfig, error) {
	config := &OIDCConfig{Scopes: []string{"openid", "email", "profile"}}

	sets, err := a.SM.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}

	// Parse settings (same as before but simplified validation)
	if enabled, exists := sets[settings.KeyOIDCEnabled]; exists {
		config.Enabled = enabled.Bool(false)
	}
	if !config.Enabled {
		return config, nil // Return early if disabled
	}

	if pkce, exists := sets[settings.KeyOIDCPKCE]; exists {
		config.UsePKCE = pkce.Bool(false)
	}

	if clientID, exists := sets[settings.KeyOIDCClientID]; exists {
		config.ClientID = clientID.String("")
	}

	if !config.UsePKCE {
		if clientSecret, exists := sets[settings.KeyOIDCClientSecret]; exists {
			config.ClientSecret = clientSecret.String("")
		}
	}

	if serverURL, exists := sets[settings.KeyServerURL]; exists {
		if parsed, err := url.Parse(serverURL.String("")); err == nil {
			config.RedirectURL = strings.TrimSuffix(parsed.String(), "/") + "/api/oidc/callback"
		}
	}

	if issuerURL, exists := sets[settings.KeyOIDCIssuerURL]; exists {
		config.IssuerURL = issuerURL.String("")
	}

	if scopes, exists := sets["oauth_scopes"]; exists && scopes.String("") != "" {
		config.Scopes = strings.Split(scopes.String(""), ",")
		for i := range config.Scopes {
			config.Scopes[i] = strings.TrimSpace(config.Scopes[i])
		}
	}

	return config, nil
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

		newUserID, err := q.CreateUser(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to create OIDC user: %w", err)
		}
		user = &db.User{
			ID:       newUserID,
			Username: username,
			Email:    &userInfo.Email,
			IsAdmin:  false,
		}
	} else {
		// Update existing user's email if verified
		if userInfo.Email != "" && userInfo.EmailVerified {
			if err := q.UpdateUser(ctx, db.UpdateUserParams{
				Username: user.Username,
				Email:    &userInfo.Email,
				IsAdmin:  user.IsAdmin,
				ID:       user.ID,
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
