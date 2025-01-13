// Package api provides handlers for the API
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/internal/traefik"
)

// Profiles -------------------------------------------------------------------

// GetProfiles retrieves all profiles from the database.
func GetProfiles(w http.ResponseWriter, r *http.Request) {
	profiles, err := db.Query.ListProfiles(r.Context())
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to get profiles: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	writeJSON(w, profiles)
}

// GetProfile fetches a single profile by its ID.
func GetProfile(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}

	profile, err := db.Query.GetProfileByID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Profile not found: %s", err.Error()), http.StatusNotFound)
		return
	}

	writeJSON(w, profile)
}

// UpsertProfile updates an existing profile identified by its ID in the database.
func UpsertProfile(w http.ResponseWriter, r *http.Request) {
	var profile db.Profile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to decode profile: %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}

	if err := profile.Verify(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := db.Query.UpsertProfile(r.Context(), db.UpsertProfileParams{
		Name:     profile.Name,
		Url:      profile.Url,
		Username: profile.Username,
		Password: profile.Password,
		Tls:      profile.Tls,
	})
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to update profile: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	go traefik.GetTraefikConfig()
	writeJSON(w, data)
}

// DeleteProfile removes an existing profile from the database by its ID.
func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}

	if _, err := db.Query.GetProfileByID(r.Context(), id); err != nil {
		http.Error(w, fmt.Sprintf("Profile not found: %s", err.Error()), http.StatusNotFound)
		return
	}

	if err := db.Query.DeleteProfileByID(r.Context(), id); err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to delete profile: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
}
