// Package api provides handlers for the API
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MizuchiLabs/mantrae/internal/db"
)

// Users ----------------------------------------------------------------------

// GetUsers returns all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := db.Query.ListUsers(r.Context())
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to get users: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	writeJSON(w, users)
}

// GetUser returns a single user
func GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	user, err := db.Query.GetUserByID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("User not found: %s", err.Error()), http.StatusNotFound)
		return
	}
	writeJSON(w, user)
}

// UpsertUser updates a single user
func UpsertUser(w http.ResponseWriter, r *http.Request) {
	var user db.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode user: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if err := user.Verify(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := db.Query.UpsertUser(r.Context(), db.UpsertUserParams{
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
		IsAdmin:  user.IsAdmin,
	})
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to update user: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	writeJSON(w, data)
}

// DeleteUser deletes a single user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	if err := db.Query.DeleteUserByID(r.Context(), id); err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to delete user: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
}
