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

	if err := profile.Verify(); err != nil {
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
}

// UpdateProfile updates a single profile
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var updatedProfile Profile
	if err := json.NewDecoder(r.Body).Decode(&updatedProfile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := updatedProfile.Verify(); err != nil {
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
			writeJSON(w, profiles[i])
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

// UpdateRouter updates or creates a router
func UpdateRouter(w http.ResponseWriter, r *http.Request) {
	var router Router
	if err := json.NewDecoder(r.Body).Decode(&router); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := router.Verify(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	profileName := r.PathValue("profile")
	routerName := r.PathValue("router")
	if profileName == "" || routerName == "" {
		http.Error(w, "profile or router name cannot be empty", http.StatusBadRequest)
		return
	}

	profiles, err := LoadProfiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, profile := range profiles {
		if strings.EqualFold(profile.Name, profileName) {
			if routerName != router.Name {
				delete(profiles[i].Instance.Dynamic.Routers, routerName)
			}
			profiles[i].Instance.Dynamic.Routers[router.Name] = router
			if err := SaveProfiles(profiles); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			writeJSON(w, profiles[i])
			return
		}
	}
	http.Error(w, "profile not found", http.StatusNotFound)
}

// DeleteRouter deletes a single router and it's services
func DeleteRouter(w http.ResponseWriter, r *http.Request) {
	profiles, err := LoadProfiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, profile := range profiles {
		if profile.Name == r.PathValue("profile") {
			delete(profiles[i].Instance.Dynamic.Routers, r.PathValue("router"))
			delete(profiles[i].Instance.Dynamic.Services, r.PathValue("router"))
			if err := SaveProfiles(profiles); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			writeJSON(w, profiles[i])
			return
		}
	}
	http.Error(w, "profile not found", http.StatusNotFound)
}

// UpdateService updates or creates a service
func UpdateService(w http.ResponseWriter, r *http.Request) {
	var service Service
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := service.Verify(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	profileName := r.PathValue("profile")
	serviceName := r.PathValue("service")
	if profileName == "" || serviceName == "" {
		http.Error(w, "profile or service name cannot be empty", http.StatusBadRequest)
		return
	}

	profiles, err := LoadProfiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, profile := range profiles {
		if strings.EqualFold(profile.Name, r.PathValue("profile")) {
			if serviceName != service.Name {
				delete(profiles[i].Instance.Dynamic.Services, serviceName)
			}
			profiles[i].Instance.Dynamic.Services[service.Name] = service
			if err := SaveProfiles(profiles); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			writeJSON(w, profiles[i])
			return
		}
	}
	http.Error(w, "profile not found", http.StatusNotFound)
}

// DeleteService deletes a single service and its router
func DeleteService(w http.ResponseWriter, r *http.Request) {
	profiles, err := LoadProfiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, profile := range profiles {
		if profile.Name == r.PathValue("profile") {
			delete(profiles[i].Instance.Dynamic.Services, r.PathValue("service"))
			delete(profiles[i].Instance.Dynamic.Routers, r.PathValue("service"))
			if err := SaveProfiles(profiles); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			writeJSON(w, profiles[i])
			return
		}
	}
	http.Error(w, "profile not found", http.StatusNotFound)
}

// UpdateMiddleware updates or creates a middleware
func UpdateMiddleware(w http.ResponseWriter, r *http.Request) {
	var middleware Middleware
	if err := json.NewDecoder(r.Body).Decode(&middleware); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if r.PathValue("profile") == "" {
		http.Error(w, "profile name cannot be empty", http.StatusBadRequest)
		return
	}

	profiles, err := LoadProfiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, profile := range profiles {
		if strings.EqualFold(profile.Name, r.PathValue("profile")) {
			profiles[i].Instance.Dynamic.Middlewares[middleware.Name] = middleware
			if err := SaveProfiles(profiles); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			writeJSON(w, profiles[i])
			return
		}
	}
	http.Error(w, "profile not found", http.StatusNotFound)
}

// DeleteMiddleware deletes a single middleware and it's services
func DeleteMiddleware(w http.ResponseWriter, r *http.Request) {
	profiles, err := LoadProfiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, profile := range profiles {
		if profile.Name == r.PathValue("profile") {
			delete(profiles[i].Instance.Dynamic.Middlewares, r.PathValue("middleware"))
			if err := SaveProfiles(profiles); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			writeJSON(w, profiles[i])
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
