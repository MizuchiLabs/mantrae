package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MizuchiLabs/mantrae/internal/db"
)

// DNS Providers ---------------------------------------------------------------

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

	data, err := db.Query.UpsertProvider(context.Background(), db.UpsertProviderParams(db.UpsertProviderParams{
		Name:       provider.Name,
		Type:       provider.Type,
		ExternalIp: provider.ExternalIp,
		ApiKey:     provider.ApiKey,
		ApiUrl:     provider.ApiUrl,
		ZoneType:   provider.ZoneType,
		Proxied:    provider.Proxied,
		IsActive:   provider.IsActive,
	}))
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
