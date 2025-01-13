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

	"github.com/MizuchiLabs/mantrae/internal/config"
	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/internal/traefik"
	"github.com/MizuchiLabs/mantrae/pkg/util"
)

// Helper function to write JSON response to the HTTP response writer.
func writeJSON(w http.ResponseWriter, data any) {
	// Set proper content type
	w.Header().Set("Content-Type", "application/json")

	// If result is already JSON string, write it directly
	if resultJSON, ok := data.(string); ok {
		w.Write([]byte(resultJSON))
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Version ---------------------------------------------------------------------

// GetVersion returns the current version of Mantrae as a plain text response.
func GetVersion(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(util.Version)); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

// Plugins ---------------------------------------------------------------------

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

// Backup ----------------------------------------------------------------------

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

// Utility ---------------------------------------------------------------------

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

// Events ----------------------------------------------------------------------

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
