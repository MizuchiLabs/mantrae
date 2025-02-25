package handler

import (
	"encoding/json"
	"net/http"

	"github.com/MizuchiLabs/mantrae/internal/util"
	"github.com/MizuchiLabs/mantrae/pkg/build"
)

// GetVersion returns the current version of Mantrae as a plain text response
func GetVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := map[string]string{"version": build.Version}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetPublicIP attempts to resolve the public IP address of the current machine
func GetPublicIP(w http.ResponseWriter, r *http.Request) {
	machineIPs, err := util.GetPublicIPsCached()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(machineIPs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
