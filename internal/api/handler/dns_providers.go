package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/MizuchiLabs/mantrae/internal/config"
	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/internal/dns"
	"github.com/MizuchiLabs/mantrae/internal/util"
)

func ListDNSProviders(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := a.Conn.GetQuery()
		dns_providers, err := q.ListDNSProviders(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(dns_providers); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetDNSProvider(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := a.Conn.GetQuery()
		dns_provider_id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		dns_provider, err := q.GetDNSProvider(r.Context(), dns_provider_id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(dns_provider); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func CreateDNSProvider(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := a.Conn.GetQuery()
		var params db.CreateDNSProviderParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := q.CreateDNSProvider(r.Context(), params); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		util.Broadcast <- util.EventMessage{
			Type:     util.EventTypeCreate,
			Category: util.EventCategoryDNS,
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func UpdateDNSProvider(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := a.Conn.GetQuery()
		var params db.UpdateDNSProviderParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := q.UpdateDNSProvider(r.Context(), params); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		util.Broadcast <- util.EventMessage{
			Type:     util.EventTypeUpdate,
			Category: util.EventCategoryDNS,
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func DeleteDNSProvider(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := a.Conn.GetQuery()
		dns_provider_id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = q.DeleteDNSProvider(r.Context(), dns_provider_id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		util.Broadcast <- util.EventMessage{
			Type:     util.EventTypeDelete,
			Category: util.EventCategoryDNS,
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func ListRouterDNSProviders(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := a.Conn.GetQuery()
		traefik_id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		providers, err := q.ListRouterDNSProvidersByTraefikID(r.Context(), traefik_id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(providers); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func SetRouterDNSProvider(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params struct {
			ProviderIDs []string `json:"providerIds"`
			TraefikID   int64    `json:"traefikId"`
			RouterName  string   `json:"routerName"`
		}
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Fetch current providers for this router
		q := a.Conn.GetQuery()
		currentProviders, err := q.GetRouterDNSProviders(
			r.Context(),
			db.GetRouterDNSProvidersParams{
				TraefikID:  params.TraefikID,
				RouterName: params.RouterName,
			},
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		currentProviderSet := make(map[int64]bool)
		for _, p := range currentProviders {
			currentProviderSet[p.ProviderID] = true
		}

		// Create sets for incoming and existing providers
		newProviderSet := make(map[int64]bool)
		for _, id := range params.ProviderIDs {
			id, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			newProviderSet[id] = true
		}

		// Insert missing providers
		for _, id := range params.ProviderIDs {
			id, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if !currentProviderSet[id] {
				if err := q.AddRouterDNSProvider(r.Context(), db.AddRouterDNSProviderParams{
					TraefikID:  params.TraefikID,
					RouterName: params.RouterName,
					ProviderID: id,
				}); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}

		// Delete providers that were removed
		for _, current := range currentProviders {
			if !newProviderSet[current.ProviderID] {
				err := q.DeleteRouterDNSProvider(r.Context(), db.DeleteRouterDNSProviderParams{
					TraefikID:  params.TraefikID,
					RouterName: params.RouterName,
					ProviderID: current.ProviderID,
				})
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}

		// Update DNS asynchronously
		go func() {
			if err := dns.UpdateDNS(a.Conn.Get()); err != nil {
				slog.Error("Failed to update DNS", "error", err)
			}
		}()

		w.WriteHeader(http.StatusNoContent)
	}
}

func GetRouterDNSProvider(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := a.Conn.GetQuery()
		var params db.GetRouterDNSProvidersParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		providers, err := q.GetRouterDNSProviders(r.Context(), params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(providers); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func DeleteRouterDNSProvider(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := a.Conn.GetQuery()
		var params db.DeleteRouterDNSProviderParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		go func() {
			err := dns.DeleteDNS(a.Conn.Get(), params)
			if err != nil {
				slog.Error("Failed to delete DNS record", "error", err)
			}
		}()
		if err := q.DeleteRouterDNSProvider(r.Context(), params); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
