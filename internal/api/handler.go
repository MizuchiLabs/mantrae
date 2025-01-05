// Package api provides handlers for the API
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/MizuchiLabs/mantrae/internal/config"
	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/internal/dns"
	"github.com/MizuchiLabs/mantrae/internal/tasks"
	"github.com/MizuchiLabs/mantrae/internal/traefik"
	"github.com/MizuchiLabs/mantrae/pkg/util"
	"golang.org/x/crypto/bcrypt"
	"sigs.k8s.io/yaml"
)

// Helper function to write JSON response to the HTTP response writer.
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Authentication -------------------------------------------------------------

// Login verifies user credentials using a normal password and returns a JWT token if successful.
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
	if err != nil ||
		bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)) != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	if r.URL.Query().Get("remember") == "true" {
		expirationTime = time.Now().Add(7 * 24 * time.Hour)
	}

	token, err := util.EncodeUserJWT(user.Username, expirationTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"token": token})
}

// VerifyToken checks the validity of a JWT token provided in cookies or Authorization header.
func VerifyToken(w http.ResponseWriter, r *http.Request) {
	var token string

	// Check for token in cookies and Authorization header
	if cookie, err := r.Cookie("token"); err == nil {
		token = cookie.Value
	} else {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			http.Error(w, "Token cannot be empty", http.StatusBadRequest)
			return
		}
	}

	data, err := util.DecodeUserJWT(token)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid token: %s", err.Error()), http.StatusUnauthorized)
		return
	}

	writeJSON(w, data)
}

// ResetPassword allows users to reset their password using a valid JWT token.
func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Failed to decode credentials", http.StatusBadRequest)
		return
	}

	tokenString := r.Header.Get("Authorization")
	if len(tokenString) < 7 {
		http.Error(w, "Token cannot be empty", http.StatusBadRequest)
		return
	}

	user, err := util.DecodeUserJWT(tokenString[7:])
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	if user.Username == "" || data.Password == "" {
		http.Error(w, "Username or password cannot be empty", http.StatusBadRequest)
		return
	}

	dbUser, err := db.Query.GetUserByUsername(context.Background(), user.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err = db.Query.UpdateUser(context.Background(), db.UpdateUserParams{
		ID:       dbUser.ID,
		Username: dbUser.Username,
		Email:    dbUser.Email,
		Password: string(hash),
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// SendResetEmail sends an email with a password reset link to the user's registered email address.
func SendResetEmail(w http.ResponseWriter, r *http.Request) {
	user, err := db.Query.GetUserByUsername(context.Background(), r.PathValue("name"))
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if user.Email == nil {
		http.Error(w, fmt.Sprintf("%s has no email address", user.Username), http.StatusBadRequest)
		return
	}

	token, err := util.EncodeUserJWT(user.Username, time.Now().Add(10*time.Minute))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	settings, err := db.Query.ListSettings(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var config util.EmailConfig
	var resetLink string
	for _, setting := range settings {
		switch setting.Key {
		case "email-host":
			config.EmailHost = setting.Value
		case "email-port":
			config.EmailPort = setting.Value
		case "email-username":
			config.EmailUsername = setting.Value
		case "email-password":
			config.EmailPassword = setting.Value
		case "email-from":
			config.EmailFrom = setting.Value
		case "server-url":
			resetLink = fmt.Sprintf("%s/login/reset?token=%s", setting.Value, token)
		}
	}
	data := map[string]interface{}{
		"ResetLink": resetLink,
		"Minutes":   10,
	}
	if err := util.SendMail(*user.Email, "reset-password", config, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Version --------------------------------------------------------------------

// GetVersion returns the current version of Mantrae as a plain text response.
func GetVersion(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(util.Version)); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

// Events ---------------------------------------------------------------------

// GetEvents streams server-sent events (SSE) for real-time updates.
func GetEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Register the client to receive updates
	util.ClientsMutex.Lock()
	util.Clients[w] = true
	util.ClientsMutex.Unlock()

	defer func() {
		// Unregister the client when the connection is closed
		util.ClientsMutex.Lock()
		delete(util.Clients, w)
		util.ClientsMutex.Unlock()
	}()

	for {
		select {
		case message := <-util.Broadcast:
			// Serialize the EventMessage to JSON
			data, err := json.Marshal(message)
			if err != nil {
				fmt.Printf("Error marshalling message: %v\n", err)
				continue
			}
			// Send the data to the client
			fmt.Fprintf(w, "data: %s\n\n", data)
			w.(http.Flusher).Flush()
		case <-r.Context().Done():
			return
		}
	}
}

// Profiles -------------------------------------------------------------------

// GetProfiles retrieves all profiles from the database.
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

// GetProfile fetches a single profile by its ID.
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

// CreateProfile inserts a new profile into the database.
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

// UpdateProfile updates an existing profile identified by its ID in the database.
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

// DeleteProfile removes an existing profile from the database by its ID.
func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}

	if _, err := db.Query.GetProfileByID(context.Background(), id); err != nil {
		http.Error(w, fmt.Sprintf("Profile not found: %s", err.Error()), http.StatusNotFound)
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

// UpsertUser updates a single user
func UpsertUser(w http.ResponseWriter, r *http.Request) {
	var user db.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode user: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if err := user.Verify(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.ID == 0 {
		data, err := db.Query.CreateUser(context.Background(), db.CreateUserParams{
			Username: user.Username,
			Password: user.Password,
			Email:    user.Email,
			IsAdmin:  user.IsAdmin,
		})
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("Failed to update user: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}
		writeJSON(w, data)
	} else {
		params := db.UpdateUserParams{
			ID:       user.ID,
			Username: user.Username,
			Password: user.Password,
			Email:    user.Email,
			IsAdmin:  user.IsAdmin,
		}
		data, err := db.Query.UpdateUser(context.Background(), params)
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

// UpsertProvider inserts or updates a provider in the database based on the provided data.
func UpsertProvider(w http.ResponseWriter, r *http.Request) {
	var provider db.Provider
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

	data, err := db.Query.UpsertProvider(context.Background(), db.UpsertProviderParams(provider))
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to upsert provider: %s", err.Error()),
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

// Entrypoints ---------------------------------------------------------------------

// GetEntryPoints returns the config for a single profile
func GetEntryPoints(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}

	entrypoints, err := db.Query.ListEntryPointsByProfileID(context.Background(), id)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to get entry points: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	for i := range entrypoints {
		if err := entrypoints[i].DecodeFields(); err != nil {
			slog.Error("Failed to decode entry point", "name", entrypoints[i].Name, "error", err)
		}
	}
	writeJSON(w, entrypoints)
}

// Router ---------------------------------------------------------------------

// GetRouters retrieves all routers for a given profile by its ID.
func GetRouters(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}

	routers, err := db.Query.ListRoutersByProfileID(context.Background(), id)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to get routers: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	for i := range routers {
		if err := routers[i].DecodeFields(); err != nil {
			slog.Error("Failed to decode router", "name", routers[i].Name, "error", err)
		}
	}
	writeJSON(w, routers)
}

// UpsertRouter inserts or updates a router in the database based on the provided data.
func UpsertRouter(w http.ResponseWriter, r *http.Request) {
	var router db.Router
	if err := json.NewDecoder(r.Body).Decode(&router); err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to decode router: %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}

	// Use default provider if not set
	provider, err := db.Query.GetDefaultProvider(context.Background())
	if err == nil && provider.ID != 0 && router.DnsProvider == nil && router.ID == "" {
		router.DnsProvider = &provider.ID
	}

	// Update DNS
	updateDNS := false
	oldRouter, err := db.Query.GetRouterByID(context.Background(), router.ID)
	if err == nil && oldRouter.Rule != router.Rule && oldRouter.DnsProvider != nil {
		updateDNS = true
	}

	if err = router.Verify(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := db.Query.UpsertRouter(context.Background(), db.UpsertRouterParams(router))
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to upsert router: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	if err := data.DecodeFields(); err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to decode router: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	if updateDNS {
		slog.Debug("Refreshing DNS record", "router", data.Name)
		go dns.UpdateDNS()
	}

	go tasks.Refresh()
	writeJSON(w, data)
}

// DeleteRouter removes a specific router from the database using its ID.
func DeleteRouter(w http.ResponseWriter, r *http.Request) {
	router, err := db.Query.GetRouterByID(context.Background(), r.PathValue("id"))
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Router not found: %s", err.Error()),
			http.StatusNotFound,
		)
		return
	}

	// Delete service too
	service, err := db.Query.GetServiceByName(context.Background(), db.GetServiceByNameParams{
		Name:      router.Service,
		ProfileID: router.ProfileID,
	})
	if err != nil {
		slog.Error("Failed to get service", "name", router.Service)
	} else {
		if err := db.Query.DeleteServiceByID(context.Background(), service.ID); err != nil {
			http.Error(
				w,
				fmt.Sprintf("Failed to delete service: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}
	}

	if err := db.Query.DeleteRouterByID(context.Background(), r.PathValue("id")); err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to delete router: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Services -------------------------------------------------------------------

// GetServices fetches all services for a given profile by its ID.
func GetServices(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}

	services, err := db.Query.ListServicesByProfileID(context.Background(), id)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to get services: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	for i := range services {
		if err := services[i].DecodeFields(); err != nil {
			slog.Error("Failed to decode router", "name", services[i].Name, "error", err)
		}
	}
	writeJSON(w, services)
}

// UpsertService inserts or updates a service in the database based on the provided data.
func UpsertService(w http.ResponseWriter, r *http.Request) {
	var service db.Service
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

	_, err := db.Query.GetServiceByID(context.Background(), service.ID)
	if err == nil {
		if err = db.Query.DeleteServiceByID(context.Background(), service.ID); err != nil {
			http.Error(
				w,
				fmt.Sprintf("Failed to delete service: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}
	}

	data, err := db.Query.UpsertService(context.Background(), db.UpsertServiceParams(service))
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to upsert service: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	if err := data.DecodeFields(); err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to decode service: %s", err.Error()),
			http.StatusInternalServerError,
		)
	}

	writeJSON(w, data)
}

// DeleteService removes a specific service from the database using its ID.
func DeleteService(w http.ResponseWriter, r *http.Request) {
	err := db.Query.DeleteServiceByID(context.Background(), r.PathValue("id"))
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to delete service: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Middlewares ----------------------------------------------------------------

// GetMiddlewares fetches all middlewares for a given profile by its ID.
func GetMiddlewares(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}

	middlewares, err := db.Query.ListMiddlewaresByProfileID(context.Background(), id)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to get middlewares: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	for i := range middlewares {
		if err := middlewares[i].DecodeFields(); err != nil {
			http.Error(
				w,
				fmt.Sprintf("Failed to decode middleware: %s", err.Error()),
				http.StatusInternalServerError,
			)
		}
	}

	writeJSON(w, middlewares)
}

// UpsertMiddleware inserts or updates a middleware in the database based on the provided data.
func UpsertMiddleware(w http.ResponseWriter, r *http.Request) {
	var middleware db.Middleware
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

	_, err := db.Query.GetMiddlewareByID(context.Background(), middleware.ID)
	if err == nil {
		if err = db.Query.DeleteMiddlewareByID(context.Background(), middleware.ID); err != nil {
			http.Error(
				w,
				fmt.Sprintf("Failed to delete middleware: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}
	}

	data, err := db.Query.UpsertMiddleware(
		context.Background(),
		db.UpsertMiddlewareParams(middleware),
	)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to upsert middleware: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	if err := data.DecodeFields(); err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to decode service: %s", err.Error()),
			http.StatusInternalServerError,
		)
	}

	writeJSON(w, data)
}

// DeleteMiddleware removes a specific middleware from the database using its ID.
func DeleteMiddleware(w http.ResponseWriter, r *http.Request) {
	err := db.Query.DeleteMiddlewareByID(context.Background(), r.PathValue("id"))
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to delete middleware: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Settings -------------------------------------------------------------------

// GetSettings retrieves all settings from the database.
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

// GetSetting fetches a single setting by its key.
func GetSetting(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	setting, err := db.Query.GetSettingByKey(context.Background(), key)
	if err != nil {
		http.Error(w, fmt.Sprintf("Setting not found: %s", err.Error()), http.StatusNotFound)
		return
	}
	writeJSON(w, setting)
}

// UpdateSetting updates an existing setting in the database based on the provided data.
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

// Agents ---------------------------------------------------------------------

// GetAgents retrieves all agents for a given profile by its ID.
func GetAgents(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}

	agents, err := db.Query.ListAgentsByProfileID(context.Background(), id)
	if err != nil {
		http.Error(w, "Failed to get agents", http.StatusInternalServerError)
		return
	}

	for i := range agents {
		if err := agents[i].DecodeFields(); err != nil {
			slog.Error("Failed to decode agent", "name", agents[i].ID, "error", err)
		}
	}
	writeJSON(w, agents)
}

// UpsertAgent inserts or updates an agent in the database based on the provided data.
func UpsertAgent(w http.ResponseWriter, r *http.Request) {
	var agent db.Agent
	if err := json.NewDecoder(r.Body).Decode(&agent); err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to decode agent: %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}

	if err := agent.Verify(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := db.Query.UpsertAgent(context.Background(), db.UpsertAgentParams(agent))
	if err != nil {
		http.Error(w, "Failed to upsert agent", http.StatusInternalServerError)
		return
	}

	if err := data.DecodeFields(); err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to decode service: %s", err.Error()),
			http.StatusInternalServerError,
		)
	}

	writeJSON(w, data)
}

// DeleteAgent removes a specific agent from the database using its ID and type (hard/soft).
func DeleteAgent(w http.ResponseWriter, r *http.Request) {
	agentID := r.PathValue("id")

	if err := db.Query.DeleteAgentByID(context.Background(), agentID); err != nil {
		http.Error(w, "Failed to delete agent", http.StatusInternalServerError)
		return
	}

	// Delete all connected routers
	routers, err := db.Query.ListRoutersByAgentID(context.Background(), &agentID)
	if err != nil {
		http.Error(w, "Failed to get routers", http.StatusInternalServerError)
		return
	}
	for _, router := range routers {
		if err := db.Query.DeleteRouterByID(context.Background(), router.ID); err != nil {
			http.Error(w, "Failed to delete router", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

// RegenerateAgentToken generates a new JWT token for an agent
func RegenerateAgentToken(w http.ResponseWriter, r *http.Request) {
	agent, err := db.Query.GetAgentByID(context.Background(), r.PathValue("id"))
	if err != nil {
		http.Error(w, "Failed to get agent", http.StatusInternalServerError)
		return
	}

	setting, err := db.Query.GetSettingByKey(context.Background(), "server-url")
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to get settings: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	token, err := util.EncodeAgentJWT(agent.ProfileID, setting.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"token": token})
}

// Backup ---------------------------------------------------------------------

// DownloadBackup creates a backup of the database and returns it as a JSON response.
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

// UploadBackup restores a backup from a provided file.
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

// DeleteRouterDNS deletes DNS records associated with a given router.
func DeleteRouterDNS(w http.ResponseWriter, r *http.Request) {
	var router db.Router
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

// GetMiddlewarePlugins retrieves middleware plugins available for Traefik.
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

// GetTraefikOverview fetches an overview of the connected Traefik instance for a given profile.
func GetTraefikOverview(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	profile, err := db.Query.GetProfileByID(context.Background(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Profile not found: %s", err.Error()), http.StatusNotFound)
	}

	// Fetch overview
	data, err := traefik.GetOverview(profile)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to fetch overview: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, data)
}

// GetTraefikConfig generates and returns the dynamic configuration for Traefik in JSON or YAML format.
func GetTraefikConfig(w http.ResponseWriter, r *http.Request) {
	profile, err := db.Query.GetProfileByName(context.Background(), r.PathValue("name"))
	if err != nil {
		http.Error(w, fmt.Sprintf("Profile not found: %s", err.Error()), http.StatusNotFound)
		return
	}

	dynamicConfig, err := traefik.GenerateConfig(profile.ID)
	if err != nil {
		util.Broadcast <- util.EventMessage{
			Type:    "config_error",
			Message: fmt.Sprintf("Failed to generate config: %s", err.Error()),
		}
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
	util.Broadcast <- util.EventMessage{
		Type:    "config_error",
		Message: "",
	}
}

// GetPublicIP attempts to resolve the public IP address of a Traefik instance by its profile ID.
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
