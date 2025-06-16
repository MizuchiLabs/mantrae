package handler

import (
	"encoding/json"
	"net/http"

	"github.com/mizuchilabs/mantrae/internal/db"
	"github.com/mizuchilabs/mantrae/internal/settings"
	"github.com/mizuchilabs/mantrae/internal/util"
)

func ListSettings(sm *settings.SettingsManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		settings, err := sm.GetAll(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(settings); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetSetting(sm *settings.SettingsManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setting, err := sm.Get(r.Context(), r.PathValue("key"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(setting); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func UpsertSetting(sm *settings.SettingsManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var setting db.UpsertSettingParams
		if err := json.NewDecoder(r.Body).Decode(&setting); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := sm.Set(r.Context(), setting.Key, setting.Value, setting.Description); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		util.Broadcast <- util.EventMessage{
			Type:     util.EventTypeUpdate,
			Category: util.EventCategorySetting,
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
