// Package api provides handlers for the API
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/MizuchiLabs/mantrae/pkg/traefik"
	"github.com/MizuchiLabs/mantrae/pkg/util"
)

// Helper function to write JSON response
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds util.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if creds.Username == "" || creds.Password == "" {
		http.Error(w, "username and password cannot be empty", http.StatusBadRequest)
		return
	}

	var valid util.Credentials
	if err := valid.GetCreds(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if creds.Username != valid.Username || creds.Password != valid.Password {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
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
	tokenString := r.Header.Get("Authorization")[7:]
	if tokenString == "" {
		http.Error(w, "token cannot be empty", http.StatusBadRequest)
		return
	}

	_, err := ValidateJWT(tokenString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// CreateProfile creates a new profile
func CreateProfile(w http.ResponseWriter, r *http.Request) {
	var profile traefik.Profile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := profile.Verify(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	profiles, err := traefik.LoadProfiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	profiles = append(profiles, profile)
	if err := traefik.SaveProfiles(profiles); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, profile)
}

// GetProfiles returns all profiles
func GetProfiles(w http.ResponseWriter, r *http.Request) {
	profiles, err := traefik.LoadProfiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, profiles)
}

// UpdateProfile updates a single profile
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var updatedProfile traefik.Profile
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

	profiles, err := traefik.LoadProfiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, profile := range profiles {
		if profile.Name == r.PathValue("name") {
			profiles[i] = updatedProfile
			if err := traefik.SaveProfiles(profiles); err != nil {
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
	profiles, err := traefik.LoadProfiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, profile := range profiles {
		if profile.Name == r.PathValue("name") {
			profiles = append(profiles[:i], profiles[i+1:]...)
			if err := traefik.SaveProfiles(profiles); err != nil {
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
	var router traefik.Router
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

	profiles, err := traefik.LoadProfiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, profile := range profiles {
		if strings.EqualFold(profile.Name, profileName) {
			if routerName != router.Name {
				delete(profiles[i].Client.Dynamic.Routers, routerName)
			}
			profiles[i].Client.Dynamic.Routers[router.Name] = router
			if err := traefik.SaveProfiles(profiles); err != nil {
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
	profiles, err := traefik.LoadProfiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, profile := range profiles {
		if profile.Name == r.PathValue("profile") {
			delete(profiles[i].Client.Dynamic.Routers, r.PathValue("router"))
			delete(profiles[i].Client.Dynamic.Services, r.PathValue("router"))
			if err := traefik.SaveProfiles(profiles); err != nil {
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
	var service traefik.Service
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

	profiles, err := traefik.LoadProfiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, profile := range profiles {
		if strings.EqualFold(profile.Name, r.PathValue("profile")) {
			if serviceName != service.Name {
				delete(profiles[i].Client.Dynamic.Services, serviceName)
			}
			profiles[i].Client.Dynamic.Services[service.Name] = service
			if err := traefik.SaveProfiles(profiles); err != nil {
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
	profiles, err := traefik.LoadProfiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, profile := range profiles {
		if profile.Name == r.PathValue("profile") {
			delete(profiles[i].Client.Dynamic.Services, r.PathValue("service"))
			delete(profiles[i].Client.Dynamic.Routers, r.PathValue("service"))
			if err := traefik.SaveProfiles(profiles); err != nil {
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
	var middleware traefik.Middleware
	if err := json.NewDecoder(r.Body).Decode(&middleware); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	profileName := r.PathValue("profile")
	middlewareName := r.PathValue("middleware")
	if profileName == "" || middlewareName == "" {
		http.Error(w, "profile or middleware name cannot be empty", http.StatusBadRequest)
		return
	}

	profiles, err := traefik.LoadProfiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, profile := range profiles {
		if strings.EqualFold(profile.Name, profileName) {
			if middlewareName != middleware.Name {
				delete(profiles[i].Client.Dynamic.Middlewares, middlewareName)
			}
			profiles[i].Client.Dynamic.Middlewares[middleware.Name] = middleware
			if err := traefik.SaveProfiles(profiles); err != nil {
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
	profiles, err := traefik.LoadProfiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, profile := range profiles {
		if profile.Name == r.PathValue("profile") {
			delete(profiles[i].Client.Dynamic.Middlewares, r.PathValue("middleware"))
			if err := traefik.SaveProfiles(profiles); err != nil {
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
	profiles, err := traefik.LoadProfiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, profile := range profiles {
		if strings.EqualFold(profile.Name, r.PathValue("name")) {
			w.Header().Set("Content-Type", "application/yaml")
			w.Header().
				Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.yaml", profile.Name))

			yamlConfig, err := traefik.GenerateConfig(profile.Client.Dynamic)
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
