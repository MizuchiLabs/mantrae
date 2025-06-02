package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MizuchiLabs/mantrae/internal/config"
	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/internal/util"
)

func ListUsers(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := a.Conn.GetQuery()
		users, err := q.ListUsers(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(users); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetUser(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := a.Conn.GetQuery()
		user_id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user, err := q.GetUser(r.Context(), user_id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func CreateUser(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := a.Conn.GetQuery()
		var user db.CreateUserParams
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		hash, err := util.HashPassword(user.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user.Password = hash
		if _, err := q.CreateUser(r.Context(), user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		util.Broadcast <- util.EventMessage{
			Type:     util.EventTypeCreate,
			Category: util.EventCategoryUser,
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func UpdateUser(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := a.Conn.GetQuery()
		var user db.UpdateUserParams
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := q.UpdateUser(r.Context(), user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		util.Broadcast <- util.EventMessage{
			Type:     util.EventTypeUpdate,
			Category: util.EventCategoryUser,
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func DeleteUser(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := a.Conn.GetQuery()
		user_id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := q.DeleteUser(r.Context(), user_id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		util.Broadcast <- util.EventMessage{
			Type:     util.EventTypeDelete,
			Category: util.EventCategoryUser,
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func UpdateUserPassword(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := a.Conn.GetQuery()
		var user db.UpdateUserPasswordParams
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Hash
		hash, err := util.HashPassword(user.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user.Password = hash
		if err := q.UpdateUserPassword(r.Context(), user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
