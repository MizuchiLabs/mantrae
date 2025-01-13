package handler

import (
	"net/http"

	"github.com/MizuchiLabs/mantrae/pkg/util"
)

// GetVersion returns the current version of Mantrae as a plain text response.
func GetVersion(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(util.Version)); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

// GetPublicIP attempts to resolve the public IP address of a Traefik instance by its profile ID.
// func GetPublicIP(w http.ResponseWriter, r *http.Request) {
// 	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
// 		return
// 	}

// 	profile, err := db.Query.GetProfileByID(context.Background(), id)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Profile not found: %s", err.Error()), http.StatusNotFound)
// 		return
// 	}

// 	// Parse the URL
// 	u, err := url.Parse(profile.Url)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Invalid URL: %s", err.Error()), http.StatusBadRequest)
// 		return
// 	}

// 	// Check if it's an IP address
// 	if net.ParseIP(u.Hostname()) != nil {
// 		if !net.ParseIP(u.Hostname()).IsLoopback() {
// 			writeJSON(w, map[string]string{"ip": u.Hostname()})
// 			return
// 		}
// 	}

// 	// If it's a valid hostname, resolve to IP
// 	ips, err := net.LookupHost(u.Hostname())
// 	if err == nil && len(ips) > 0 {
// 		if !net.ParseIP(ips[0]).IsLoopback() {
// 			writeJSON(w, map[string]string{"ip": ips[0]})
// 			return
// 		}
// 	}

// 	ip, err := util.GetPublicIP()
// 	if err != nil {
// 		http.Error(
// 			w,
// 			fmt.Sprintf("Failed to get public IP: %s", err.Error()),
// 			http.StatusInternalServerError,
// 		)
// 		return
// 	}
// 	writeJSON(w, map[string]string{"ip": ip})
// }
