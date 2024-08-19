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
	traefik.GetTraefikConfig()
	writeJSON(w, profile)
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
	var updatedProfile traefik.Profile
	if err := json.NewDecoder(r.Body).Decode(&updatedProfile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := updatedProfile.Verify(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	profileName := r.PathValue("name")
	if profileName == "" {
		http.Error(w, "profile name cannot be empty", http.StatusBadRequest)
		return
	}

	var p traefik.Profiles
	if err := p.Load(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(p.Profiles) == 0 {
		http.Error(w, "no profiles configured", http.StatusBadRequest)
		return
	}

	if _, ok := p.Profiles[profileName]; !ok {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	if updatedProfile.Name != profileName {
		delete(p.Profiles, profileName)
	}

	p.Profiles[profileName] = updatedProfile
	if err := p.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, p.Profiles[profileName])
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
	w.WriteHeader(http.StatusNoContent)
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

	if _, ok := p.Profiles[profileName]; !ok {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	if routerName != router.Name {
		delete(p.Profiles[profileName].Client.Dynamic.Routers, routerName)
	}

	p.Profiles[profileName].Client.Dynamic.Routers[routerName] = router
	if err := p.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, p.Profiles[profileName])
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

	delete(p.Profiles[profileName].Client.Dynamic.Routers, routerName)
	delete(p.Profiles[profileName].Client.Dynamic.Services, routerName)
	if err := p.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, p.Profiles[r.PathValue("profile")])
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

	if _, ok := p.Profiles[profileName]; !ok {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	if serviceName != service.Name {
		delete(p.Profiles[profileName].Client.Dynamic.Services, serviceName)
	}

	p.Profiles[profileName].Client.Dynamic.Services[serviceName] = service
	if err := p.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, p.Profiles[profileName])
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

	delete(p.Profiles[profileName].Client.Dynamic.Services, serviceName)
	delete(p.Profiles[profileName].Client.Dynamic.Routers, serviceName)
	if err := p.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, p.Profiles[r.PathValue("profile")])
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

	if _, ok := p.Profiles[profileName]; !ok {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	if middlewareName != middleware.Name {
		delete(p.Profiles[profileName].Client.Dynamic.Middlewares, middlewareName)
	}

	p.Profiles[profileName].Client.Dynamic.Middlewares[middleware.Name] = middleware
	if err := p.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, p.Profiles[profileName])
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

	delete(p.Profiles[profileName].Client.Dynamic.Middlewares, middlewareName)
	if err := p.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, p.Profiles[profileName])
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

	if _, ok := p.Profiles[profileName]; !ok {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	yamlConfig, err := traefik.GenerateConfig(p.Profiles[profileName].Client.Dynamic)
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
