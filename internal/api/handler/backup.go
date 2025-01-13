package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MizuchiLabs/mantrae/internal/config"
)

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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
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

	w.WriteHeader(http.StatusNoContent)
}
