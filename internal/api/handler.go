// Package api provides handlers for the API
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

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

// UpdateConfig updates the config for a single profile
func UpdateConfig(w http.ResponseWriter, r *http.Request) {
	var config traefik.Dynamic
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "Failed to decode config", http.StatusBadRequest)
		return
	}

	if err := traefik.UpdateConfig(config.ProfileID, &config); err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to update config: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	// Update the DNS records immediately
	go dns.UpdateDNS()

	writeJSON(w, config)
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
	cmd := exec.Command("sqlite3", util.Path(util.DBName), ".dump")
	sqlFile, err := os.CreateTemp("", "backup-*.sql")
	if err != nil {
		http.Error(w, "Failed to create SQL file", http.StatusInternalServerError)
	}
	defer os.Remove(sqlFile.Name())

	cmd.Stdout = sqlFile
	if err := cmd.Run(); err != nil {
		http.Error(w, "Failed to run sqlite3 command", http.StatusInternalServerError)
		return
	}

	// Reopen the temp file to start reading from the beginning
	if _, err := sqlFile.Seek(0, 0); err != nil {
		http.Error(w, "Failed to seek SQL file", http.StatusInternalServerError)
		return
	}

	content := fmt.Sprintf(
		"attachment; filename=backup-%s.sql",
		time.Now().Format("2006-01-02"),
	)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", content)

	if _, err := io.Copy(w, sqlFile); err != nil {
		http.Error(w, "Failed to write file to response", http.StatusInternalServerError)
	}
}

func UploadBackup(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if err = db.DB.Close(); err != nil {
		http.Error(
			w,
			"Failed to close existing database connection",
			http.StatusInternalServerError,
		)
		return
	}

	restoreFile, err := os.CreateTemp("", "restore-*.sql")
	if err != nil {
		http.Error(w, "Failed to create SQL file", http.StatusInternalServerError)
		return
	}
	defer restoreFile.Close()
	defer os.Remove(restoreFile.Name())

	if _, err = io.Copy(restoreFile, file); err != nil {
		http.Error(w, "Failed to copy file to content", http.StatusInternalServerError)
		return
	}

	// Reopen the temp file to start reading from the beginning
	if _, err = restoreFile.Seek(0, 0); err != nil {
		http.Error(w, "Failed to seek SQL file", http.StatusInternalServerError)
		return
	}

	// Delete Database
	if err = os.Remove(util.Path(util.DBName)); err != nil && !os.IsNotExist(err) {
		http.Error(w, "Failed to delete existing database", http.StatusInternalServerError)
		return
	}

	cmd := exec.Command("sqlite3", util.Path(util.DBName), ".read "+restoreFile.Name())
	if err = cmd.Run(); err != nil {
		http.Error(w, "Failed to run sqlite3 command", http.StatusInternalServerError)
		return
	}

	// Reopen the database connection
	if err = db.InitDB(); err != nil {
		http.Error(w, "Failed to initialize database", http.StatusInternalServerError)
		return
	}

	// Broadcast the update to all clients
	util.Broadcast <- "profiles"

	// Respond to the client that the restore was successful.
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write([]byte("Database restored successfully")); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

// Extra ---------------------------------------------------------------------

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
