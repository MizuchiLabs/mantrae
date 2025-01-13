package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/internal/traefik"
	"github.com/MizuchiLabs/mantrae/pkg/util"
	"sigs.k8s.io/yaml"
)

func PublishTraefikConfig(w http.ResponseWriter, r *http.Request) {
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
