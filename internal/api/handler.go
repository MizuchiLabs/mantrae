// Package api provides handlers for the API
package api

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	tokenString := r.Header.Get("Authorization")
	if len(tokenString) < 7 {
		http.Error(w, "token cannot be empty", http.StatusBadRequest)
		return
	}

	_, err := ValidateJWT(tokenString[7:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetProfiles returns all profiles but without the dynamic data
func GetProfiles(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]traefik.Profile, len(traefik.ProfileData.Profiles))
	for name, profile := range traefik.ProfileData.Profiles {
		data[name] = traefik.Profile{
			Name:     profile.Name,
			URL:      profile.URL,
			Username: profile.Username,
			Password: profile.Password,
			TLS:      profile.TLS,
		}
	}
	writeJSON(w, data)
}

// GetProfile returns a single profile
func GetProfile(w http.ResponseWriter, r *http.Request) {
	profileName := r.PathValue("name")
	if profileName == "" {
		http.Error(w, "profile name cannot be empty", http.StatusBadRequest)
		return
	}

	profile, ok := traefik.ProfileData.Profiles[profileName]
	if !ok {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	writeJSON(w, profile)
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

	traefik.ProfileData.Profiles[profile.Name] = profile
	if err := traefik.ProfileData.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	go traefik.GetTraefikConfig()
	writeJSON(w, traefik.ProfileData.Profiles)
}

// UpdateProfile updates a single profile
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	profileName := r.PathValue("name")
	if profileName == "" {
		http.Error(w, "profile name cannot be empty", http.StatusBadRequest)
		return
	}

	var profile traefik.Profile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := profile.Verify(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(traefik.ProfileData.Profiles) == 0 {
		http.Error(w, "no profiles configured", http.StatusBadRequest)
		return
	}

	updateName(traefik.ProfileData.Profiles, profileName, profile.Name)

	traefik.ProfileData.Profiles[profile.Name] = profile
	if err := traefik.ProfileData.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	go traefik.GetTraefikConfig()
	writeJSON(w, profile)
}

// DeleteProfile deletes a single profile
func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	profileName := r.PathValue("name")
	if profileName == "" {
		http.Error(w, "profile name cannot be empty", http.StatusBadRequest)
		return
	}

	if _, ok := traefik.ProfileData.Profiles[profileName]; !ok {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	delete(traefik.ProfileData.Profiles, profileName)
	if err := traefik.ProfileData.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, traefik.ProfileData.Profiles)
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

	profile, ok := traefik.ProfileData.Profiles[profileName]
	if !ok {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	// Initialize the Routers map if it is nil
	if profile.Dynamic.Routers == nil {
		profile.Dynamic.Routers = make(map[string]traefik.Router)
	}

	// If the router name is being changed, delete the old entry
	updateName(profile.Dynamic.Routers, routerName, router.Name)

	profile.Dynamic.Routers[router.Name] = router
	traefik.ProfileData.Profiles[profileName] = profile // Update the profile in the map
	if err := traefik.ProfileData.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, profile)
}

// DeleteRouter deletes a single router and it's services
func DeleteRouter(w http.ResponseWriter, r *http.Request) {
	profileName := r.PathValue("profile")
	routerName := r.PathValue("router")
	if profileName == "" || routerName == "" {
		http.Error(w, "profile or router name cannot be empty", http.StatusBadRequest)
		return
	}

	profile, ok := traefik.ProfileData.Profiles[profileName]
	if !ok {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	delete(profile.Dynamic.Routers, routerName)
	delete(profile.Dynamic.Services, routerName)
	if err := traefik.ProfileData.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, profile)
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

	profile, ok := traefik.ProfileData.Profiles[profileName]
	if !ok {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	if profile.Dynamic.Services == nil {
		profile.Dynamic.Services = make(map[string]traefik.Service)
	}

	updateName(profile.Dynamic.Services, serviceName, service.Name)

	profile.Dynamic.Services[service.Name] = service
	traefik.ProfileData.Profiles[profileName] = profile // Update the profile in the map
	if err := traefik.ProfileData.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, profile)
}

// DeleteService deletes a single service and its router
func DeleteService(w http.ResponseWriter, r *http.Request) {
	profileName := r.PathValue("profile")
	serviceName := r.PathValue("service")
	if profileName == "" || serviceName == "" {
		http.Error(w, "profile or service name cannot be empty", http.StatusBadRequest)
		return
	}

	profile, ok := traefik.ProfileData.Profiles[profileName]
	if !ok {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	delete(profile.Dynamic.Services, serviceName)
	delete(profile.Dynamic.Routers, serviceName)
	if err := traefik.ProfileData.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, profile)
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

	profile, ok := traefik.ProfileData.Profiles[profileName]
	if !ok {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	if profile.Dynamic.Middlewares == nil {
		profile.Dynamic.Middlewares = make(map[string]traefik.Middleware)
	}

	updateName(profile.Dynamic.Middlewares, middlewareName, middleware.Name)
	if err := middleware.Verify(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	profile.Dynamic.Middlewares[middleware.Name] = middleware
	traefik.ProfileData.Profiles[profileName] = profile // Update the profile in the map
	if err := traefik.ProfileData.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, profile)
}

// DeleteMiddleware deletes a single middleware and it's services
func DeleteMiddleware(w http.ResponseWriter, r *http.Request) {
	profileName := r.PathValue("profile")
	middlewareName := r.PathValue("middleware")
	if profileName == "" || middlewareName == "" {
		http.Error(w, "profile or middleware name cannot be empty", http.StatusBadRequest)
		return
	}

	profile, ok := traefik.ProfileData.Profiles[profileName]
	if !ok {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	delete(profile.Dynamic.Middlewares, middlewareName)
	if err := traefik.ProfileData.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, profile)
}

// GetConfig returns the traefik config for a single profile
func GetConfig(w http.ResponseWriter, r *http.Request) {
	profileName := r.PathValue("name")
	if profileName == "" {
		http.Error(w, "profile name cannot be empty", http.StatusBadRequest)
		return
	}

	profile, ok := traefik.ProfileData.Profiles[profileName]
	if !ok {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	yamlConfig, err := traefik.GenerateConfig(profile.Dynamic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/yaml")
	w.Header().
		Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.yaml", profileName))

	if _, err := w.Write(yamlConfig); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
