package handler

import (
	"encoding/json"
	"net/http"

	"github.com/MizuchiLabs/mantrae/internal/db"
)

func GetHTTPRoutersBySource(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var router db.GetHTTPRoutersBySourceParams
		if err := json.NewDecoder(r.Body).Decode(&router); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		routers, err := q.GetHTTPRoutersBySource(r.Context(), router)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routers)
	}
}

func GetTCPRoutersBySource(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var router db.GetTCPRoutersBySourceParams
		if err := json.NewDecoder(r.Body).Decode(&router); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		routers, err := q.GetTCPRoutersBySource(r.Context(), router)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routers)
	}
}

func GetUDPRoutersBySource(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var router db.GetUDPRoutersBySourceParams
		if err := json.NewDecoder(r.Body).Decode(&router); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		routers, err := q.GetUDPRoutersBySource(r.Context(), router)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routers)
	}
}

func GetHTTPRouterByName(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var router db.GetHTTPRouterByNameParams
		if err := json.NewDecoder(r.Body).Decode(&router); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		routers, err := q.GetHTTPRouterByName(r.Context(), router)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routers)
	}
}

func GetTCPRouterByName(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var router db.GetTCPRouterByNameParams
		if err := json.NewDecoder(r.Body).Decode(&router); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		routers, err := q.GetTCPRouterByName(r.Context(), router)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routers)
	}
}

func GetUDPRouterByName(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var router db.GetUDPRouterByNameParams
		if err := json.NewDecoder(r.Body).Decode(&router); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		routers, err := q.GetUDPRouterByName(r.Context(), router)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routers)
	}
}

func UpsertHTTPRouter(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var router db.UpsertHTTPRouterParams
		if err := json.NewDecoder(r.Body).Decode(&router); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := q.UpsertHTTPRouter(r.Context(), router); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func UpsertTCPRouter(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var router db.UpsertTCPRouterParams
		if err := json.NewDecoder(r.Body).Decode(&router); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := q.UpsertTCPRouter(r.Context(), router); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func UpsertUDPRouter(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var router db.UpsertUDPRouterParams
		if err := json.NewDecoder(r.Body).Decode(&router); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := q.UpsertUDPRouter(r.Context(), router); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func DeleteHTTPRouter(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var router db.DeleteHTTPRouterParams
		if err := json.NewDecoder(r.Body).Decode(&router); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := q.DeleteHTTPRouter(r.Context(), router); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func DeleteTCPRouter(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var router db.DeleteTCPRouterParams
		if err := json.NewDecoder(r.Body).Decode(&router); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := q.DeleteTCPRouter(r.Context(), router); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func DeleteUDPRouter(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var router db.DeleteUDPRouterParams
		if err := json.NewDecoder(r.Body).Decode(&router); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := q.DeleteUDPRouter(r.Context(), router); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
