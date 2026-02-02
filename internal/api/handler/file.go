package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/traefik"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
	"gopkg.in/yaml.v3"
)

func sanitizeFilename(name string) string {
	return filepath.Base(strings.ReplaceAll(name, "..", ""))
}

func UploadBackup(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 100MB limit
		if err := r.ParseMultipartForm(100 << 20); err != nil {
			http.Error(w, "File too large or invalid form data", http.StatusBadRequest)
			return
		}
		defer func() {
			if err := r.MultipartForm.RemoveAll(); err != nil {
				slog.Error("failed to close request body", "error", err)
			}
		}()

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Failed to get uploaded file", http.StatusBadRequest)
			return
		}
		defer func() {
			if err = file.Close(); err != nil {
				slog.Error("failed to close uploaded file", "error", err)
			}
		}()

		extension := filepath.Ext(sanitizeFilename(header.Filename))
		allowedExtensions := []string{".db", ".json", ".yaml", ".yml"}
		if !slices.Contains(allowedExtensions, extension) {
			http.Error(w, "Invalid file type", http.StatusBadRequest)
			return
		}

		// Handle sqlite db backups
		if extension == ".db" {
			filename := fmt.Sprintf("upload_%s%s",
				time.Now().UTC().Format("20060102_150405"),
				filepath.Ext(header.Filename))

			if err = a.BM.Storage.Store(r.Context(), filename, file); err != nil {
				http.Error(
					w,
					fmt.Sprintf("Failed to store backup file: %v", err),
					http.StatusInternalServerError,
				)
				return
			}
			if err = a.BM.Restore(r.Context(), filename); err != nil {
				http.Error(
					w,
					fmt.Sprintf("Failed to restore from backup: %v", err),
					http.StatusInternalServerError,
				)
				return
			}
		} else { // Handle dynamic configuration
			profileID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
			if err != nil {
				http.Error(
					w,
					fmt.Sprintf("Failed to parse profile_id: %v", err),
					http.StatusBadRequest,
				)
				return
			}

			dynamic := &dynamic.Configuration{}
			content, err := io.ReadAll(file)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to read file: %v", err), http.StatusInternalServerError)
				return
			}

			switch extension {
			case ".yaml", ".yml":
				if err = yaml.Unmarshal(content, dynamic); err != nil {
					http.Error(w, fmt.Sprintf("Failed to decode YAML file: %v", err), http.StatusInternalServerError)
					return
				}
			case ".json":
				if err = json.Unmarshal(content, dynamic); err != nil {
					http.Error(w, fmt.Sprintf("Failed to decode JSON file: %v", err), http.StatusInternalServerError)
					return
				}
			default:
				http.Error(w, "Invalid file type", http.StatusBadRequest)
				return
			}

			// Write to database
			traefik.DynamicToDB(r.Context(), a.Conn.Q, profileID, dynamic)
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
