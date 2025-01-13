package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MizuchiLabs/mantrae/internal/db"
)

// Valid sources
const (
	internalSource = "internal"
	externalSource = "external"
)

func CreateTraefikConfig(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var traefik_config db.CreateTraefikConfigParams
		if err := json.NewDecoder(r.Body).Decode(&traefik_config); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if traefik_config.Source != internalSource && traefik_config.Source != externalSource {
			http.Error(w, "invalid source", http.StatusBadRequest)
			return
		}
		if err := q.CreateTraefikConfig(r.Context(), traefik_config); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
func GetTraefikConfig(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		config, err := q.GetTraefikConfig(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(config)
	}
}

func GetTraefikConfigBySource(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profile_id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		config, err := q.GetTraefikConfigBySource(r.Context(), db.GetTraefikConfigBySourceParams{
			ProfileID: profile_id,
			Source:    r.PathValue("source"),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(config)
	}
}

func UpdateTraefikConfig(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var config db.UpdateTraefikConfigParams
		if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if config.Source != internalSource && config.Source != externalSource {
			http.Error(w, "invalid source", http.StatusBadRequest)
			return
		}

		err := q.UpdateTraefikConfig(r.Context(), config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func DeleteTraefikConfig(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profile_id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = q.DeleteTraefikConfig(r.Context(), profile_id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func PublishTraefikConfig(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profile, err := q.GetProfileByName(r.Context(), r.PathValue("name"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		config, err := q.GetTraefikConfigBySource(r.Context(), db.GetTraefikConfigBySourceParams{
			ProfileID: profile.ID,
			Source:    "internal",
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(config.Config)
	}
}
