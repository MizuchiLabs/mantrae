package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/internal/source"
	"github.com/MizuchiLabs/mantrae/internal/util"
)

func ListProfiles(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := db.New(DB)
		profiles, err := q.ListProfiles(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(profiles)
	}
}

func GetProfile(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := db.New(DB)
		profile_id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		profile, err := q.GetProfile(r.Context(), profile_id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(profile)
	}
}

func CreateProfile(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := db.New(DB)
		var profile db.CreateProfileParams
		if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		profileID, err := q.CreateProfile(r.Context(), profile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create default local config
		if err := q.UpsertTraefikConfig(r.Context(), db.UpsertTraefikConfigParams{
			ProfileID: profileID,
			Source:    source.Local,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		util.Broadcast <- util.EventMessage{
			Type:    util.EventTypeCreate,
			Message: "profile",
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func UpdateProfile(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := db.New(DB)
		var profile db.UpdateProfileParams
		if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := q.UpdateProfile(r.Context(), profile); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		util.Broadcast <- util.EventMessage{
			Type:    util.EventTypeUpdate,
			Message: "profile",
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func DeleteProfile(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := db.New(DB)
		profile_id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := q.DeleteProfile(r.Context(), profile_id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		util.Broadcast <- util.EventMessage{
			Type:    util.EventTypeDelete,
			Message: "profile",
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
