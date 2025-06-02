package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/MizuchiLabs/mantrae/internal/api/agent"
	"github.com/MizuchiLabs/mantrae/internal/config"
	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/internal/settings"
	"github.com/MizuchiLabs/mantrae/internal/util"
	"github.com/google/uuid"
)

func ListAgents(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := a.Conn.GetQuery()
		agents, err := q.ListAgents(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(agents); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func ListAgentsByProfile(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := a.Conn.GetQuery()
		profileID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		agents, err := q.ListAgentsByProfile(r.Context(), profileID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(agents); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetAgent(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := a.Conn.GetQuery()
		agent, err := q.GetAgent(r.Context(), r.PathValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(agent); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func CreateAgent(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := a.Conn.GetQuery()
		profileID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		serverUrl, err := a.SM.Get(r.Context(), settings.KeyServerURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		claims := &agent.AgentClaims{
			AgentID:   uuid.New().String(),
			ProfileID: profileID,
			ServerURL: serverUrl.String("http://127.0.0.1:3000"),
		}

		// Generate a JWT for the agent and let it expire based on the cleanup interval
		agentInterval, err := a.SM.Get(r.Context(), settings.KeyAgentCleanupInterval)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		token, err := claims.EncodeJWT(
			a.Config.Secret,
			agentInterval.Duration(time.Hour*72),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := q.CreateAgent(r.Context(), db.CreateAgentParams{
			ID:        claims.AgentID,
			ProfileID: claims.ProfileID,
			Token:     token,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		util.Broadcast <- util.EventMessage{
			Type:     util.EventTypeCreate,
			Category: util.EventCategoryAgent,
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func UpdateAgentIP(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := a.Conn.GetQuery()
		var agent db.UpdateAgentIPParams
		if err := json.NewDecoder(r.Body).Decode(&agent); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := q.UpdateAgentIP(r.Context(), agent); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		util.Broadcast <- util.EventMessage{
			Type:     util.EventTypeUpdate,
			Category: util.EventCategoryAgent,
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func DeleteAgent(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := a.Conn.GetQuery()
		agent, err := q.GetAgent(r.Context(), r.PathValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := q.DeleteTraefikConfigByAgent(r.Context(), &agent.ID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := q.DeleteAgent(r.Context(), agent.ID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		util.Broadcast <- util.EventMessage{
			Type:     util.EventTypeDelete,
			Category: util.EventCategoryAgent,
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func RotateAgentToken(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := a.Conn.GetQuery()
		dbAgent, err := q.GetAgent(r.Context(), r.PathValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		serverUrl, err := a.SM.Get(r.Context(), settings.KeyServerURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		claims := &agent.AgentClaims{
			AgentID:   dbAgent.ID,
			ProfileID: dbAgent.ProfileID,
			ServerURL: serverUrl.String("http://127.0.0.1:3000"),
		}

		agentInterval, err := a.SM.Get(r.Context(), settings.KeyAgentCleanupInterval)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		token, err := claims.EncodeJWT(
			a.Config.Secret,
			agentInterval.Duration(time.Hour*72),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := q.UpdateAgentToken(r.Context(), db.UpdateAgentTokenParams{
			ID:    dbAgent.ID,
			Token: token,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		util.Broadcast <- util.EventMessage{
			Type:     util.EventTypeUpdate,
			Category: util.EventCategoryAgent,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(token); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
