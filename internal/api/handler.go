// Package api provides handlers for the API
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/MizuchiLabs/mantrae/internal/config"
	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/pkg/dns"
	"github.com/MizuchiLabs/mantrae/pkg/traefik"
	"github.com/MizuchiLabs/mantrae/pkg/util"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
	"golang.org/x/crypto/bcrypt"
	"sigs.k8s.io/yaml"
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

	token, err := util.EncodeUserJWT(user.Username)
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

	_, err := util.DecodeJWT(tokenString[7:])
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

// Agents ---------------------------------------------------------------------

// GetAgents returns all agents
func GetAgents(w http.ResponseWriter, r *http.Request) {
	agents, err := db.Query.ListAgents(context.Background())
	if err != nil {
		http.Error(w, "Failed to get agents", http.StatusInternalServerError)
		return
	}
	writeJSON(w, agents)
}

// GetAgent returns an agent
func GetAgent(w http.ResponseWriter, r *http.Request) {
	agent, err := db.Query.GetAgentByID(context.Background(), r.PathValue("id"))
	if err != nil {
		http.Error(w, "Failed to get agent", http.StatusInternalServerError)
		return
	}
	writeJSON(w, agent)
}

// DeleteAgent deletes an agent
func DeleteAgent(w http.ResponseWriter, r *http.Request) {
	if err := db.Query.DeleteAgentByID(context.Background(), r.PathValue("id")); err != nil {
		http.Error(w, "Failed to delete agent", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetAgentToken returns an agent token
func GetAgentToken(w http.ResponseWriter, r *http.Request) {
	var data struct {
		ServerURL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if data.ServerURL == "" {
		http.Error(w, "Server URL cannot be empty", http.StatusBadRequest)
		return
	}

	token, err := util.EncodeAgentJWT(data.ServerURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"token": token})
}

// Profiles -------------------------------------------------------------------

// GetProfiles returns all profiles but without the dynamic data
func GetProfiles(w http.ResponseWriter, r *http.Request) {
	profiles, err := db.Query.ListProfiles(context.Background())
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

// GetProfile returns a single profile
func GetProfile(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}

	profile, err := db.Query.GetProfileByID(context.Background(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Profile not found: %s", err.Error()), http.StatusNotFound)
		return
	}

	writeJSON(w, profile)
}

// CreateProfile creates a new profile
func CreateProfile(w http.ResponseWriter, r *http.Request) {
	var profile db.CreateProfileParams
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
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
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
		http.Error(
			w,
			fmt.Sprintf("Failed to get users: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	writeJSON(w, users)
}

// GetUser returns a single user
func GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	user, err := db.Query.GetUserByID(context.Background(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("User not found: %s", err.Error()), http.StatusNotFound)
		return
	}
	writeJSON(w, user)
}

// CreateUser creates a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user db.CreateUserParams
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode user: %s", err.Error()), http.StatusBadRequest)
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
		http.Error(w, fmt.Sprintf("Failed to decode user: %s", err.Error()), http.StatusBadRequest)
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
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
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
		http.Error(
			w,
			fmt.Sprintf("Failed to get providers: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, providers)
}

// GetProvider returns a single provider
func GetProvider(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	provider, err := db.Query.GetProviderByID(context.Background(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Provider not found: %s", err.Error()), http.StatusNotFound)
		return
	}

	writeJSON(w, provider)
}

// CreateProvider creates a new provider
func CreateProvider(w http.ResponseWriter, r *http.Request) {
	var provider db.CreateProviderParams
	if err := json.NewDecoder(r.Body).Decode(&provider); err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to decode provider: %s", err.Error()),
			http.StatusBadRequest,
		)
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
		http.Error(
			w,
			fmt.Sprintf("Failed to decode provider: %s", err.Error()),
			http.StatusBadRequest,
		)
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
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
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
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}

	data, err := traefik.DecodeFromDB(id)
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
		http.Error(
			w,
			fmt.Sprintf("Failed to decode router: %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}

	if err := router.Verify(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}

	data, err := traefik.DecodeFromDB(id)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to decode config: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	data.Routers[router.Name] = router
	newConfig, err := traefik.EncodeToDB(data)
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
		http.Error(
			w,
			fmt.Sprintf("Failed to decode service: %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}

	if err := service.Verify(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	data, err := traefik.DecodeFromDB(id)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to decode config: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	data.Services[service.Name] = service
	newConfig, err := traefik.EncodeToDB(data)
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
		http.Error(
			w,
			fmt.Sprintf("Failed to decode middleware: %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}

	if err := middleware.Verify(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	data, err := traefik.DecodeFromDB(id)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to decode config: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	data.Middlewares[middleware.Name] = middleware
	newConfig, err := traefik.EncodeToDB(data)
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

func AddPluginMiddleware(w http.ResponseWriter, r *http.Request) {
	var plugin traefik.Plugin
	if err := json.NewDecoder(r.Body).Decode(&plugin); err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to decode plugin: %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}

	var pluginRaw map[string]any
	if err := yaml.Unmarshal([]byte(plugin.Snippet.Yaml), &pluginRaw); err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to decode plugin: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	pluginName := ""
	importParts := strings.Split(plugin.Name, "/")
	if len(importParts) > 0 {
		pluginName = importParts[len(importParts)-1]
	}
	myPluginName := fmt.Sprintf("my-%s", pluginName)

	pluginData := pluginRaw["http"].(map[string]any)["middlewares"].(map[string]any)[myPluginName].(map[string]any)["plugin"].(map[string]any)[pluginName].(map[string]any)
	newMiddleware := traefik.Middleware{
		Name:           pluginName + "@http",
		Provider:       "http",
		Type:           strings.ToLower(pluginName),
		MiddlewareType: "http",
		Plugin:         make(map[string]dynamic.PluginConf),
	}
	newMiddleware.Plugin[pluginName] = pluginData

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}

	data, err := traefik.DecodeFromDB(id)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to decode config: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	data.Middlewares[newMiddleware.Name] = newMiddleware
	newConfig, err := traefik.EncodeToDB(data)
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
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}

	data, err := traefik.DecodeFromDB(id)
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
	newConfig, err := traefik.EncodeToDB(data)
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
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}

	data, err := traefik.DecodeFromDB(id)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to decode config: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	delete(data.Middlewares, r.PathValue("name"))
	newConfig, err := traefik.EncodeToDB(data)
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
		http.Error(
			w,
			fmt.Sprintf("Failed to get settings: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	writeJSON(w, settings)
}

// GetSetting returns a single setting
func GetSetting(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	setting, err := db.Query.GetSettingByKey(context.Background(), key)
	if err != nil {
		http.Error(w, fmt.Sprintf("Setting not found: %s", err.Error()), http.StatusNotFound)
		return
	}
	writeJSON(w, setting)
}

// UpdateSetting updates a single setting
func UpdateSetting(w http.ResponseWriter, r *http.Request) {
	var setting db.UpdateSettingParams
	if err := json.NewDecoder(r.Body).Decode(&setting); err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to decode setting: %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}

	if err := setting.Verify(); err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to verify setting: %s", err.Error()),
			http.StatusBadRequest,
		)
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
			http.Error(
				w,
				fmt.Sprintf("Failed to schedule backups: %s", err.Error()),
				http.StatusInternalServerError,
			)
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
		http.Error(
			w,
			fmt.Sprintf("Failed to create backup: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	writeJSON(w, data)
}

func UploadBackup(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to upload backup: %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}
	defer file.Close()

	var data config.BackupData
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to decode backup: %s", err.Error()),
			http.StatusBadRequest,
		)
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
		http.Error(
			w,
			fmt.Sprintf("Failed to decode router: %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}

	go dns.DeleteDNS(router)

	w.WriteHeader(http.StatusOK)
}

func GetMiddlewarePlugins(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://plugins.traefik.io/api/services/plugins")
	if err != nil {
		http.Error(w, "Failed to fetch plugins", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Parse the response body into the Go Plugin struct
	var plugins []traefik.Plugin
	if err := json.NewDecoder(resp.Body).Decode(&plugins); err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to parse plugins: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	// Filter out non-middleware plugins
	var middlewarePlugins []traefik.Plugin
	for _, plugin := range plugins {
		if plugin.Type == "middleware" {
			middlewarePlugins = append(middlewarePlugins, plugin)
		}
	}

	writeJSON(w, middlewarePlugins)
}

// GetTraefikConfig returns the traefik config
func GetTraefikConfig(w http.ResponseWriter, r *http.Request) {
	profile, err := db.Query.GetProfileByName(context.Background(), r.PathValue("name"))
	if err != nil {
		http.Error(w, fmt.Sprintf("Profile not found: %s", err.Error()), http.StatusNotFound)
		return
	}

	data, err := traefik.DecodeFromDB(profile.ID)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to decode config: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	dynamicConfig, err := traefik.GenerateConfig(*data)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to generate config: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	var outConfig []byte
	acceptYAML := r.URL.Query().Get("yaml") == "true"
	contentType := "application/json"

	// Marshal the config as YAML if requested, otherwise default to JSON
	if acceptYAML {
		outConfig, err = yaml.Marshal(dynamicConfig)
		contentType = "application/yaml"
	} else {
		outConfig, err = json.Marshal(dynamicConfig)
	}
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to marshal config: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	// Set the appropriate content type and write the response
	w.Header().Set("Content-Type", contentType)
	if _, err := w.Write(outConfig); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

// GetPublicIP tries to guess the public ip of the Traefik instance
func GetPublicIP(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}

	profile, err := db.Query.GetProfileByID(context.Background(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Profile not found: %s", err.Error()), http.StatusNotFound)
		return
	}

	// Parse the URL
	u, err := url.Parse(profile.Url)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid URL: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// Check if it's an IP address
	if net.ParseIP(u.Hostname()) != nil {
		if !net.ParseIP(u.Hostname()).IsLoopback() {
			writeJSON(w, map[string]string{"ip": u.Hostname()})
			return
		}
	}

	// If it's a valid hostname, resolve to IP
	ips, err := net.LookupHost(u.Hostname())
	if err == nil && len(ips) > 0 {
		if !net.ParseIP(ips[0]).IsLoopback() {
			writeJSON(w, map[string]string{"ip": ips[0]})
			return
		}
	}

	ip, err := util.GetPublicIP()
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to get public IP: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	writeJSON(w, map[string]string{"ip": ip})
}
