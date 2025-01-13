package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/pkg/util"
)

// Agents ---------------------------------------------------------------------

// GetAgents retrieves all agents for a given profile by its ID.
func GetAgents(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}

	agents, err := db.Query.ListAgentsByProfileID(context.Background(), id)
	if err != nil {
		http.Error(w, "Failed to get agents", http.StatusInternalServerError)
		return
	}

	for i := range agents {
		if err := agents[i].DecodeFields(); err != nil {
			slog.Error("Failed to decode agent", "name", agents[i].ID, "error", err)
		}
	}
	writeJSON(w, agents)
}

// UpsertAgent inserts or updates an agent in the database based on the provided data.
func UpsertAgent(w http.ResponseWriter, r *http.Request) {
	var agent db.Agent
	if err := json.NewDecoder(r.Body).Decode(&agent); err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to decode agent: %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}

	if err := agent.Verify(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := db.Query.UpsertAgent(context.Background(), db.UpsertAgentParams(agent))
	if err != nil {
		http.Error(w, "Failed to upsert agent", http.StatusInternalServerError)
		return
	}

	if err := data.DecodeFields(); err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to decode service: %s", err.Error()),
			http.StatusInternalServerError,
		)
	}

	writeJSON(w, data)
}

// DeleteAgent removes a specific agent from the database using its ID and type (hard/soft).
func DeleteAgent(w http.ResponseWriter, r *http.Request) {
	agentID := r.PathValue("id")

	if err := db.Query.DeleteAgentByID(context.Background(), agentID); err != nil {
		http.Error(w, "Failed to delete agent", http.StatusInternalServerError)
		return
	}

	// Delete all connected routers
	routers, err := db.Query.ListRoutersByAgentID(context.Background(), &agentID)
	if err != nil {
		http.Error(w, "Failed to get routers", http.StatusInternalServerError)
		return
	}
	for _, router := range routers {
		if err := db.Query.DeleteRouterByID(context.Background(), router.ID); err != nil {
			http.Error(w, "Failed to delete router", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

// RegenerateAgentToken generates a new JWT token for an agent
func RegenerateAgentToken(w http.ResponseWriter, r *http.Request) {
	agent, err := db.Query.GetAgentByID(context.Background(), r.PathValue("id"))
	if err != nil {
		http.Error(w, "Failed to get agent", http.StatusInternalServerError)
		return
	}

	setting, err := db.Query.GetSettingByKey(context.Background(), "server-url")
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to get settings: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	token, err := util.EncodeAgentJWT(agent.ProfileID, setting.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"token": token})
}
