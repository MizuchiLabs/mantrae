package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/MizuchiLabs/mantrae/internal/config"
)

// DownloadBackup creates a backup of the database and returns it as a JSON response.
func DownloadBackup(bm *config.BackupManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename, err := bm.GetLatestBackup()
		if err != nil {
			bm.CreateBackup(r.Context())
			filename, err = bm.GetLatestBackup()
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to create backup: %v", err), http.StatusInternalServerError)
				return
			}
		}
		backupPath := filepath.Join(bm.Config.Dir, filename)

		// Validate filename to prevent directory traversal
		if !bm.IsValidBackupFile(filename) {
			http.Error(w, "Invalid backup file", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		w.Header().Set("Content-Type", "application/octet-stream")
		http.ServeFile(w, r, backupPath)
	}
}

// RestoreBackup restores a backup from a provided file.
func RestoreBackup(bm *config.BackupManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Limit request size to prevent memory issues
		r.Body = http.MaxBytesReader(w, r.Body, 100<<20) // 100MB limit

		if err := r.ParseMultipartForm(100 << 20); err != nil {
			http.Error(w, "File too large or invalid form data", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Failed to get uploaded file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Generate unique filename
		filename := fmt.Sprintf("backup_%s%s",
			time.Now().UTC().Format("20060102_150405"),
			filepath.Ext(header.Filename))

		backupPath := filepath.Join(bm.Config.Dir, filename)

		// Create destination file
		dst, err := os.Create(backupPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create backup file: %v", err), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Copy uploaded file to destination
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Failed to save backup file", http.StatusInternalServerError)
			return
		}

		// Attempt to restore the backup
		if err := bm.RestoreBackup(r.Context(), backupPath); err != nil {
			// Clean up the uploaded file if restore fails.
			os.Remove(backupPath)
			http.Error(w, fmt.Sprintf("Failed to restore backup: %v", err), http.StatusInternalServerError)
			return
		}

		response := map[string]string{"message": "Backup restored successfully"}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

func ListBackups(bm *config.BackupManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		backups, err := bm.ListBackups()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to list backups: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(backups); err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
			return
		}
	}
}
