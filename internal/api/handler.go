// Package api provides handlers for the API
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/pkg/traefik"
)

// Helper function to write JSON response
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Helper function to update maps in case of name changes
func updateName[K comparable, V any](m map[K]V, oldName, newName K) {
	if oldName == newName {
		return
	}

	if _, ok := m[oldName]; ok {
		m[newName] = m[oldName]
		delete(m, oldName)
	}
}

// Authentication -------------------------------------------------------------

// Login verifies the user credentials
func Login(w http.ResponseWriter, r *http.Request) {
	var creds db.Credential
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Failed to decode credentials", http.StatusBadRequest)
		return
	}

	if creds.Username == "" || creds.Password == "" {
		http.Error(w, "Username or password cannot be empty", http.StatusBadRequest)
		return
	}

	if _, err := db.Query.ValidateAuth(context.Background(), db.ValidateAuthParams{
		Username: creds.Username,
		Password: creds.Password,
	}); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := GenerateJWT(creds.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"token": token})
}

func VerifyToken(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if len(tokenString) < 7 {
		http.Error(w, "Token cannot be empty", http.StatusBadRequest)
		return
	}

	_, err := ValidateJWT(tokenString[7:])
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Profiles -------------------------------------------------------------------

// GetProfiles returns all profiles but without the dynamic data
func GetProfiles(w http.ResponseWriter, r *http.Request) {
	profiles, err := db.Query.ListProfiles(context.Background())
	if err != nil {
		http.Error(w, "Failed to get profiles", http.StatusInternalServerError)
		return
	}
	writeJSON(w, profiles)
}

// GetProfile returns a single profile
func GetProfile(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	profile, err := db.Query.GetProfileByID(context.Background(), id)
	if err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}
	writeJSON(w, profile)
}

// CreateProfile creates a new profile
func CreateProfile(w http.ResponseWriter, r *http.Request) {
	var profile db.CreateProfileParams
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, "Failed to decode profile", http.StatusBadRequest)
		return
	}

	data, err := db.Query.CreateProfile(context.Background(), profile)
	if err != nil {
		http.Error(w, "Failed to create profile", http.StatusInternalServerError)
		return
	}

	writeJSON(w, data)
}

// UpdateProfile updates a single profile
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var profile db.UpdateProfileParams
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, "Failed to decode profile", http.StatusBadRequest)
		return
	}

	data, err := db.Query.UpdateProfile(context.Background(), profile)
	if err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	go traefik.GetTraefikConfig()
	writeJSON(w, data)
}

// DeleteProfile deletes a single profile
func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	if err := db.Query.DeleteProfileByID(context.Background(), id); err != nil {
		http.Error(w, "Failed to delete profile", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Providers ------------------------------------------------------------------

// GetProviders returns all providers
func GetProviders(w http.ResponseWriter, r *http.Request) {
	providers, err := db.Query.ListProviders(context.Background())
	if err != nil {
		http.Error(w, "Failed to get providers", http.StatusInternalServerError)
		return
	}

	writeJSON(w, providers)
}

// GetProvider returns a single provider
func GetProvider(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Provider not found", http.StatusNotFound)
		return
	}
	provider, err := db.Query.GetProviderByID(context.Background(), id)
	if err != nil {
		http.Error(w, "Provider not found", http.StatusNotFound)
		return
	}

	writeJSON(w, provider)
}

// CreateProvider creates a new provider
func CreateProvider(w http.ResponseWriter, r *http.Request) {
	var provider db.CreateProviderParams
	if err := json.NewDecoder(r.Body).Decode(&provider); err != nil {
		http.Error(w, "Failed to decode provider", http.StatusBadRequest)
		return
	}

	data, err := db.Query.CreateProvider(context.Background(), provider)
	if err != nil {
		http.Error(w, "Failed to create provider", http.StatusInternalServerError)
		return
	}

	writeJSON(w, data)
}

// UpdateProvider updates a single provider
func UpdateProvider(w http.ResponseWriter, r *http.Request) {
	var provider db.UpdateProviderParams
	if err := json.NewDecoder(r.Body).Decode(&provider); err != nil {
		http.Error(w, "Failed to decode provider", http.StatusBadRequest)
		return
	}

	data, err := db.Query.UpdateProvider(context.Background(), provider)
	if err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	writeJSON(w, data)
}

// DeleteProvider deletes a single provider
func DeleteProvider(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}
	if err := db.Query.DeleteProviderByID(context.Background(), id); err != nil {
		http.Error(w, "Failed to delete provider", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Config ---------------------------------------------------------------------

// GetConfig returns the config for a single profile
func GetConfig(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}
	config, err := db.Query.GetConfigByProfileID(context.Background(), id)
	if err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}
	data, err := traefik.DecodeConfig(config)
	if err != nil {
		http.Error(w, "Failed to decode config", http.StatusInternalServerError)
		return
	}

	writeJSON(w, data)
}

func UpdateConfig(w http.ResponseWriter, r *http.Request) {
	var config traefik.Dynamic
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "Failed to decode config", http.StatusBadRequest)
		return
	}

	if err := traefik.UpdateConfig(config.ProfileID, &config); err != nil {
		http.Error(w, "Failed to update config", http.StatusInternalServerError)
		return
	}

	writeJSON(w, config)
}

// GetTraefikConfig returns the traefik config
func GetTraefikConfig(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}
	config, err := db.Query.GetConfigByProfileID(context.Background(), id)
	if err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}
	data, err := traefik.DecodeConfig(config)
	if err != nil {
		http.Error(w, "Failed to decode config", http.StatusInternalServerError)
		return
	}

	yamlConfig, err := traefik.GenerateConfig(*data)
	if err != nil {
		http.Error(w, "Failed to generate traefik config", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/yaml")
	w.Header().
		Set("Content-Disposition", fmt.Sprintf("attachment; filename=dynamic.yaml"))

	if _, err := w.Write(yamlConfig); err != nil {
		http.Error(w, "Failed to write traefik config", http.StatusInternalServerError)
		return
	}
}
