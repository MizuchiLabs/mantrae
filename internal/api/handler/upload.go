package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"slices"
	"time"

	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/storage"
)

func UploadAvatar(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Limit request size to prevent memory issues
		r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10MB limit

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "File too large or invalid form data", http.StatusBadRequest)
			return
		}
		defer r.MultipartForm.RemoveAll()

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Failed to get uploaded file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		extension := filepath.Ext(header.Filename)
		allowedExtensions := []string{".png", ".jpg", ".jpeg"}
		if !slices.Contains(allowedExtensions, extension) {
			http.Error(w, "Invalid file type", http.StatusBadRequest)
			return
		}

		username := r.URL.Query().Get("username")
		if username == "" {
			http.Error(w, "Username not provided", http.StatusBadRequest)
			return
		}
		// Generate unique filename
		filename := fmt.Sprintf(
			"avatar_%s_%s%s",
			username,
			time.Now().UTC().Format("20060102_150405"),
			filepath.Ext(header.Filename),
		)

		storePath, err := storage.GetBackend(r.Context(), a.SM, "uploads")
		if err != nil {
			http.Error(w, "Failed to get storage backend", http.StatusInternalServerError)
			return
		}
		storePath.Store(r.Context(), filename, file)

		response := map[string]string{"message": "Avatar updated successfully"}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

func UploadBackup(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Limit request size to prevent memory issues
		r.Body = http.MaxBytesReader(w, r.Body, 100<<20) // 100MB limit

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "File too large or invalid form data", http.StatusBadRequest)
			return
		}
		defer r.MultipartForm.RemoveAll()

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
		filename := fmt.Sprintf("upload_%s%s",
			time.Now().UTC().Format("20060102_150405"),
			filepath.Ext(header.Filename))

		// Store the uploaded backup using the backend
		if err = a.BM.Storage.Store(r.Context(), filename, file); err != nil {
			http.Error(
				w,
				fmt.Sprintf("Failed to store backup file: %v", err),
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
