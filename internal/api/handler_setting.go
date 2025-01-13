package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MizuchiLabs/mantrae/internal/config"
	"github.com/MizuchiLabs/mantrae/internal/db"
)

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
