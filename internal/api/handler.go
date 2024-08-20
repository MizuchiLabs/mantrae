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

	var p traefik.Profiles
	if err := p.Load(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(p.Profiles) == 0 {
		p.Profiles = make(map[string]traefik.Profile)
	}

	p.Profiles[profile.Name] = profile
	if err := p.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	go traefik.GetTraefikConfig()
	writeJSON(w, p.Profiles)
}

// GetProfiles returns all profiles
func GetProfiles(w http.ResponseWriter, r *http.Request) {
	var p traefik.Profiles
	if err := p.Load(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, p.Profiles)
}

// UpdateProfile updates a single profile
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var p traefik.Profile
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := p.Verify(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	profileName := r.PathValue("name")
	if profileName == "" {
		http.Error(w, "profile name cannot be empty", http.StatusBadRequest)
		return
	}

	var profiles traefik.Profiles
	if err := profiles.Load(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(profiles.Profiles) == 0 {
		http.Error(w, "no profiles configured", http.StatusBadRequest)
		return
	}

	if _, ok := profiles.Profiles[profileName]; !ok {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	if p.Name != profileName {
		delete(profiles.Profiles, profileName)
	}

	profiles.Profiles[p.Name] = p
	if err := profiles.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	go traefik.GetTraefikConfig()
	writeJSON(w, profiles.Profiles)
}

// DeleteProfile deletes a single profile
func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	var p traefik.Profiles
	if err := p.Load(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	profileName := r.PathValue("name")
	if profileName == "" {
		http.Error(w, "profile name cannot be empty", http.StatusBadRequest)
		return
	}

	if _, ok := p.Profiles[profileName]; !ok {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	delete(p.Profiles, profileName)
	if err := p.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, p.Profiles)
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

	var p traefik.Profiles
	if err := p.Load(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	profile, ok := p.Profiles[profileName]
	if !ok {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	// Initialize the Routers map if it is nil
	if profile.Dynamic.Routers == nil {
		profile.Dynamic.Routers = make(map[string]traefik.Router)
	}

	// If the router name is being changed, delete the old entry
	if routerName != router.Name {
		delete(profile.Dynamic.Routers, routerName)
	}

	profile.Dynamic.Routers[router.Name] = router
	p.Profiles[profileName] = profile // Update the profile in the map
	if err := p.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, p.Profiles)
}

// DeleteRouter deletes a single router and it's services
func DeleteRouter(w http.ResponseWriter, r *http.Request) {
	var p traefik.Profiles
	if err := p.Load(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	profileName := r.PathValue("profile")
	routerName := r.PathValue("router")
	if profileName == "" || routerName == "" {
		http.Error(w, "profile or router name cannot be empty", http.StatusBadRequest)
		return
	}

	if _, ok := p.Profiles[profileName]; !ok {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	delete(p.Profiles[profileName].Dynamic.Routers, routerName)
	delete(p.Profiles[profileName].Dynamic.Services, routerName)
	if err := p.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, p.Profiles)
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

	var p traefik.Profiles
	if err := p.Load(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	profile, ok := p.Profiles[profileName]
	if !ok {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	if serviceName != service.Name {
		delete(profile.Dynamic.Services, serviceName)
	}

	if profile.Dynamic.Services == nil {
		profile.Dynamic.Services = make(map[string]traefik.Service)
	}

	profile.Dynamic.Services[service.Name] = service
	p.Profiles[profileName] = profile // Update the profile in the map
	if err := p.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, p.Profiles)
}

// DeleteService deletes a single service and its router
func DeleteService(w http.ResponseWriter, r *http.Request) {
	var p traefik.Profiles
	if err := p.Load(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	profileName := r.PathValue("profile")
	serviceName := r.PathValue("service")
	if profileName == "" || serviceName == "" {
		http.Error(w, "profile or service name cannot be empty", http.StatusBadRequest)
		return
	}
	if _, ok := p.Profiles[r.PathValue("profile")]; !ok {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	delete(p.Profiles[profileName].Dynamic.Services, serviceName)
	delete(p.Profiles[profileName].Dynamic.Routers, serviceName)
	if err := p.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, p.Profiles)
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

	var p traefik.Profiles
	if err := p.Load(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	profile, ok := p.Profiles[profileName]
	if !ok {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	if middlewareName != middleware.Name {
		delete(p.Profiles[profileName].Dynamic.Middlewares, middlewareName)
	}

	if profile.Dynamic.Middlewares == nil {
		profile.Dynamic.Middlewares = make(map[string]traefik.Middleware)
	}

	profile.Dynamic.Middlewares[middleware.Name] = middleware
	p.Profiles[profileName] = profile // Update the profile in the map
	if err := p.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, p.Profiles)
}

// DeleteMiddleware deletes a single middleware and it's services
func DeleteMiddleware(w http.ResponseWriter, r *http.Request) {
	var p traefik.Profiles
	if err := p.Load(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	profileName := r.PathValue("profile")
	middlewareName := r.PathValue("middleware")
	if profileName == "" || middlewareName == "" {
		http.Error(w, "profile or middleware name cannot be empty", http.StatusBadRequest)
		return
	}

	if _, ok := p.Profiles[profileName]; !ok {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	delete(p.Profiles[profileName].Dynamic.Middlewares, middlewareName)
	if err := p.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, p.Profiles)
}

// GetConfig returns the traefik config for a single profile
func GetConfig(w http.ResponseWriter, r *http.Request) {
	var p traefik.Profiles
	if err := p.Load(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	profileName := r.PathValue("name")
	if profileName == "" {
		http.Error(w, "profile name cannot be empty", http.StatusBadRequest)
		return
	}

	profile, ok := p.Profiles[profileName]
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
