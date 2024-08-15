// Package api provides handlers for the API
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Helper function to write JSON response
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CreateProfile creates a new profile
func CreateProfile(w http.ResponseWriter, r *http.Request) {
	var profile Profile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := VerifyProfile(profile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	profiles, err := LoadProfiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	profiles = append(profiles, profile)
	if err := SaveProfiles(profiles); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, profile)
}

// GetProfiles returns all profiles
func GetProfiles(w http.ResponseWriter, r *http.Request) {
	profiles, err := LoadProfiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, profiles)
	// syncTraefik()
}

// UpdateProfile updates a single profile
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var updatedProfile Profile
	if err := json.NewDecoder(r.Body).Decode(&updatedProfile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := VerifyProfile(updatedProfile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if r.PathValue("name") == "" {
		http.Error(w, "profile name cannot be empty", http.StatusBadRequest)
		return
	}

	profiles, err := LoadProfiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, profile := range profiles {
		if profile.Name == r.PathValue("name") {
			profiles[i] = updatedProfile
			if err := SaveProfiles(profiles); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			writeJSON(w, updatedProfile)
			return
		}
	}
	http.Error(w, "profile not found", http.StatusNotFound)
}

// DeleteProfile deletes a single profile
func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	profiles, err := LoadProfiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, profile := range profiles {
		if profile.Name == r.PathValue("name") {
			profiles = append(profiles[:i], profiles[i+1:]...)
			if err := SaveProfiles(profiles); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "profile not found", http.StatusNotFound)
}

// GetConfig returns the traefik config for a single profile
func GetConfig(w http.ResponseWriter, r *http.Request) {
	profiles, err := LoadProfiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, profile := range profiles {
		if strings.EqualFold(profile.Name, r.PathValue("name")) {
			w.Header().Set("Content-Type", "application/yaml")
			w.Header().
				Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.yaml", profile.Name))

			yamlConfig, err := ParseConfig(profile.Instance.Dynamic)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if _, err := w.Write(yamlConfig); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
	}
	http.NotFound(w, r)
}
