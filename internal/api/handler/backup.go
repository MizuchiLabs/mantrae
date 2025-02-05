package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/MizuchiLabs/mantrae/internal/backup"
	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/internal/source"
	"github.com/MizuchiLabs/mantrae/internal/traefik"
)

// DownloadBackup creates a new backup of the database and returns it
func DownloadBackup(bm *backup.BackupManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a new backup
		if err := bm.Create(r.Context()); err != nil {
			http.Error(
				w,
				fmt.Sprintf("Failed to create backup: %v", err),
				http.StatusInternalServerError,
			)
			return
		}
		files, err := bm.Backend.List(r.Context())
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("Failed to list backups: %v", err),
				http.StatusInternalServerError,
			)
			return
		}

		filename := files[0].Name
		reader, err := bm.Backend.Retrieve(r.Context(), filename)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("Failed to retrieve backup: %v", err),
				http.StatusInternalServerError,
			)
			return
		}
		defer reader.Close()

		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		w.Header().Set("Content-Type", "application/octet-stream")
		if _, err := io.Copy(w, reader); err != nil {
			http.Error(
				w,
				fmt.Sprintf("Failed to write backup data: %v", err),
				http.StatusInternalServerError,
			)
			return
		}
	}
}

func DownloadBackupByName(bm *backup.BackupManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := r.PathValue("filename")

		if !bm.IsValidBackupFile(filename) {
			http.Error(w, "Invalid backup filename", http.StatusBadRequest)
			return
		}

		reader, err := bm.Backend.Retrieve(r.Context(), filename)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("Failed to retrieve backup: %v", err),
				http.StatusInternalServerError,
			)
			return
		}
		defer reader.Close()

		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		w.Header().Set("Content-Type", "application/octet-stream")
		if _, err := io.Copy(w, reader); err != nil {
			http.Error(
				w,
				fmt.Sprintf("Failed to write backup data: %v", err),
				http.StatusInternalServerError,
			)
			return
		}
	}
}

func CreateBackup(bm *backup.BackupManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := bm.Create(r.Context())
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("Failed to create backup: %v", err),
				http.StatusInternalServerError,
			)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

// RestoreBackup restores a backup from a provided file.
func RestoreBackup(bm *backup.BackupManager) http.HandlerFunc {
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

		extension := filepath.Ext(header.Filename)
		if extension != ".db" {
			http.Error(w, "Invalid file type", http.StatusBadRequest)
			return
		}

		// Generate unique filename
		filename := fmt.Sprintf("backup_%s%s",
			time.Now().UTC().Format("20060102_150405"),
			filepath.Ext(header.Filename))

		// Store the uploaded backup using the backend
		if err = bm.Backend.Store(r.Context(), filename, file); err != nil {
			http.Error(
				w,
				fmt.Sprintf("Failed to store backup file: %v", err),
				http.StatusInternalServerError,
			)
			return
		}

		// Attempt to restore the backup
		if err := bm.Restore(r.Context(), filename); err != nil {
			// Clean up the uploaded file if restore fails
			if err = bm.Backend.Delete(r.Context(), filename); err != nil {
				http.Error(
					w,
					fmt.Sprintf("Failed to delete backup file: %v", err),
					http.StatusInternalServerError,
				)
				return
			}
			http.Error(
				w,
				fmt.Sprintf("Failed to restore backup: %v", err),
				http.StatusInternalServerError,
			)
			return
		}

		time.Sleep(100 * time.Millisecond)

		// Test the connection
		if err := bm.Conn.Ping(); err != nil {
			http.Error(
				w,
				"Database connection failed after restore",
				http.StatusInternalServerError,
			)
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

// RestoreBackup restores a backup from a provided file.
func RestoreDynamicConfig(bm *backup.BackupManager) http.HandlerFunc {
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

		// Check if yaml
		extension := filepath.Ext(header.Filename)
		if extension != ".yaml" && extension != ".yml" && extension != ".json" {
			http.Error(w, "Invalid file type", http.StatusBadRequest)
			return
		}

		profileID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dc, err := traefik.ConvertFileToDynamicConfig(file, extension)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		uploadedConfig := traefik.ConvertDynamicToRuntime(dc)

		// Attempt to restore the backup
		q := bm.Conn.GetQuery()
		currentConfig, err := q.GetLocalTraefikConfig(r.Context(), profileID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		mergedConfig := traefik.MergeConfigs(currentConfig.Config, uploadedConfig)
		if err := q.UpsertTraefikConfig(r.Context(), db.UpsertTraefikConfigParams{
			ProfileID: profileID,
			Source:    source.Local,
			Config:    mergedConfig,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func ListBackups(bm *backup.BackupManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		backups, err := bm.Backend.List(r.Context())
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("Failed to list backups: %v", err),
				http.StatusInternalServerError,
			)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(backups); err != nil {
			http.Error(
				w,
				fmt.Sprintf("Failed to encode response: %v", err),
				http.StatusInternalServerError,
			)
			return
		}
	}
}

func DeleteBackup(bm *backup.BackupManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := r.PathValue("filename")
		if err := bm.Backend.Delete(r.Context(), filename); err != nil {
			http.Error(
				w,
				fmt.Sprintf("Failed to delete backup: %v", err),
				http.StatusInternalServerError,
			)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
