// Package api provides handlers for the API
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MizuchiLabs/mantrae/internal/config"
	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/pkg/dns"
	"github.com/MizuchiLabs/mantrae/pkg/traefik"
	"github.com/MizuchiLabs/mantrae/pkg/util"
	"golang.org/x/crypto/bcrypt"
)

// Helper function to write JSON response
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Authentication -------------------------------------------------------------

// Login verifies the user credentials
func Login(w http.ResponseWriter, r *http.Request) {
	var user db.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Failed to decode credentials", http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username or password cannot be empty", http.StatusBadRequest)
		return
	}

	dbUser, err := db.Query.GetUserByUsername(context.Background(), user.Username)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Check if the password is correct
	if err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := GenerateJWT(user.Username)
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

// Version --------------------------------------------------------------------

// GetVersion returns the current version of Mantrae
func GetVersion(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(util.Version)); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

// Events ---------------------------------------------------------------------

func GetEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Register the client to receive updates
	util.ClientsMutex.Lock()
	util.Clients[util.Broadcast] = true
	util.ClientsMutex.Unlock()

	defer func() {
		// Unregister the client when the connection is closed
		util.ClientsMutex.Lock()
		delete(util.Clients, util.Broadcast)
		util.ClientsMutex.Unlock()
	}()

	for {
		select {
		case message := <-util.Broadcast:
			fmt.Fprintf(w, "data: %s\n\n", message)
			w.(http.Flusher).Flush() // Force the data to be sent to the client
		case <-r.Context().Done():
			return
		}
	}
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

	if err := profile.Verify(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := db.Query.CreateProfile(context.Background(), profile)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to create profile: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	go traefik.GetTraefikConfig()
	writeJSON(w, data)
}

// UpdateProfile updates a single profile
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var profile db.UpdateProfileParams
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, "Failed to decode profile", http.StatusBadRequest)
		return
	}

	if err := profile.Verify(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := db.Query.UpdateProfile(context.Background(), profile)
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

// DeleteProfile deletes a single profile
func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	if err := db.Query.DeleteProfileByID(context.Background(), id); err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to delete profile: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Users ----------------------------------------------------------------------

// GetUsers returns all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := db.Query.ListUsers(context.Background())
	if err != nil {
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}
	writeJSON(w, users)
}

// GetUser returns a single user
func GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	user, err := db.Query.GetUserByID(context.Background(), id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	writeJSON(w, user)
}

// CreateUser creates a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user db.CreateUserParams
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Failed to decode user", http.StatusBadRequest)
		return
	}

	if err := user.Verify(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := db.Query.CreateUser(context.Background(), user)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to create user: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, data)
}

// UpdateUser updates a single user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user db.UpdateUserParams
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Failed to decode user", http.StatusBadRequest)
		return
	}

	if err := user.Verify(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := db.Query.UpdateUser(context.Background(), user)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to update user: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, data)
}

// DeleteUser deletes a single user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if err := db.Query.DeleteUserByID(context.Background(), id); err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to delete user: %s", err.Error()),
			http.StatusInternalServerError,
		)
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

	if err := provider.Verify(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := db.Query.CreateProvider(context.Background(), provider)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to create provider: %s", err.Error()),
			http.StatusInternalServerError,
		)
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

	if err := provider.Verify(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := db.Query.UpdateProvider(context.Background(), provider)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to update profile: %s", err.Error()),
			http.StatusInternalServerError,
		)
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
		http.Error(
			w,
			fmt.Sprintf("Failed to delete provider: %s", err.Error()),
			http.StatusInternalServerError,
		)
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
		http.Error(
			w,
			fmt.Sprintf("Failed to decode config: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, data)
}

func UpdateRouter(w http.ResponseWriter, r *http.Request) {
	var router traefik.Router
	if err := json.NewDecoder(r.Body).Decode(&router); err != nil {
		http.Error(w, "Failed to decode config", http.StatusBadRequest)
		return
	}

	if err := router.Verify(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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
		http.Error(
			w,
			fmt.Sprintf("Failed to decode config: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	data.Routers[router.Name] = router
	newConfig, err := traefik.UpdateConfig(config.ProfileID, data)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to update config: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	// Update the DNS records immediately
	go dns.UpdateDNS()

	writeJSON(w, newConfig)
}

func UpdateService(w http.ResponseWriter, r *http.Request) {
	var service traefik.Service
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		http.Error(w, "Failed to decode config", http.StatusBadRequest)
		return
	}

	if err := service.Verify(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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
		http.Error(
			w,
			fmt.Sprintf("Failed to decode config: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	data.Services[service.Name] = service
	newConfig, err := traefik.UpdateConfig(config.ProfileID, data)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to update config: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, newConfig)
}

func UpdateMiddleware(w http.ResponseWriter, r *http.Request) {
	var middleware traefik.Middleware
	if err := json.NewDecoder(r.Body).Decode(&middleware); err != nil {
		http.Error(w, "Failed to decode config", http.StatusBadRequest)
		return
	}

	if err := middleware.Verify(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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
		http.Error(
			w,
			fmt.Sprintf("Failed to decode config: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	data.Middlewares[middleware.Name] = middleware
	newConfig, err := traefik.UpdateConfig(config.ProfileID, data)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to update config: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, newConfig)
}

func DeleteRouter(w http.ResponseWriter, r *http.Request) {
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
		http.Error(
			w,
			fmt.Sprintf("Failed to decode config: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	delete(data.Routers, r.PathValue("name"))
	delete(data.Services, r.PathValue("name"))
	newConfig, err := traefik.UpdateConfig(config.ProfileID, data)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to update config: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, newConfig)
}

func DeleteMiddleware(w http.ResponseWriter, r *http.Request) {
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
		http.Error(
			w,
			fmt.Sprintf("Failed to decode config: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	delete(data.Middlewares, r.PathValue("name"))
	newConfig, err := traefik.UpdateConfig(config.ProfileID, data)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to update config: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, newConfig)
}

// Settings -------------------------------------------------------------------

// GetSettings returns all settings
func GetSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := db.Query.ListSettings(context.Background())
	if err != nil {
		http.Error(w, "Failed to get settings", http.StatusInternalServerError)
		return
	}
	writeJSON(w, settings)
}

// GetSetting returns a single setting
func GetSetting(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	setting, err := db.Query.GetSettingByKey(context.Background(), key)
	if err != nil {
		http.Error(w, "Setting not found", http.StatusNotFound)
		return
	}
	writeJSON(w, setting)
}

// UpdateSetting updates a single setting
func UpdateSetting(w http.ResponseWriter, r *http.Request) {
	var setting db.UpdateSettingParams
	if err := json.NewDecoder(r.Body).Decode(&setting); err != nil {
		http.Error(w, "Failed to decode setting", http.StatusBadRequest)
		return
	}

	if err := setting.Verify(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := db.Query.UpdateSetting(context.Background(), setting)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to update setting: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	// Check if the updated setting affects the backup configuration
	if setting.Key == "backup-enabled" || setting.Key == "backup-schedule" {
		if err := config.ScheduleBackups(); err != nil {
			http.Error(w, "Failed to reconfigure backup", http.StatusInternalServerError)
			return
		}
	}

	writeJSON(w, data)
}

// Backup ---------------------------------------------------------------------

// DownloadBackup returns a backup of the database
func DownloadBackup(w http.ResponseWriter, r *http.Request) {
	data, err := config.DumpBackup(r.Context())
	if err != nil {
		http.Error(w, "Failed to get backup", http.StatusInternalServerError)
		return
	}
	writeJSON(w, data)
}

func UploadBackup(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	var data config.BackupData
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		http.Error(w, "Failed to decode backup data", http.StatusBadRequest)
		return
	}
	if err := config.RestoreBackup(r.Context(), &data); err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to restore backup: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Utility --------------------------------------------------------------------

// DeleteRouterDNS deletes the DNS records for a router
func DeleteRouterDNS(w http.ResponseWriter, r *http.Request) {
	var router traefik.Router
	if err := json.NewDecoder(r.Body).Decode(&router); err != nil {
		http.Error(w, "Failed to decode config", http.StatusBadRequest)
		return
	}

	go dns.DeleteDNS(router)

	w.WriteHeader(http.StatusOK)
}

// GetTraefikConfig returns the traefik config
func GetTraefikConfig(w http.ResponseWriter, r *http.Request) {
	profile, err := db.Query.GetProfileByName(context.Background(), r.PathValue("name"))
	if err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	config, err := db.Query.GetConfigByProfileID(context.Background(), profile.ID)
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

	if _, err := w.Write(yamlConfig); err != nil {
		http.Error(w, "Failed to write traefik config", http.StatusInternalServerError)
		return
	}
}

func GetPublicIP(w http.ResponseWriter, r *http.Request) {
	ip, err := util.GetPublicIP()
	if err != nil {
		http.Error(w, "Failed to get public IP", http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"ip": ip})
}
