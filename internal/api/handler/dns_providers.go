package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/internal/dns"
)

func ListDNSProviders(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dns_providers, err := q.ListDNSProviders(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(dns_providers)
	}
}

func GetDNSProvider(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		json.NewEncoder(w).Encode(dns_provider)
	}
}

func CreateDNSProvider(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dns_provider db.CreateDNSProviderParams
		if err := json.NewDecoder(r.Body).Decode(&dns_provider); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := q.CreateDNSProvider(r.Context(), dns_provider); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func UpdateDNSProvider(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dns_provider db.UpdateDNSProviderParams
		if err := json.NewDecoder(r.Body).Decode(&dns_provider); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := q.UpdateDNSProvider(r.Context(), dns_provider); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func DeleteDNSProvider(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusNoContent)
	}
}

func ListRouterDNSProviders(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		json.NewEncoder(w).Encode(providers)
	}
}

func SetRouterDNSProvider(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dns_provider db.UpsertRouterDNSProviderParams
		if err := json.NewDecoder(r.Body).Decode(&dns_provider); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if router exists
		config, err := q.GetTraefikConfig(r.Context(), dns_provider.TraefikID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		router := config.Config.Routers[dns_provider.RouterName]
		if router == nil {
			http.Error(w, "Router not found", http.StatusBadRequest)
			return
		}

		if err := q.UpsertRouterDNSProvider(r.Context(), dns_provider); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Update DNS
		go func() {
			err := dns.UpdateDNS(q)
			if err != nil {
				slog.Error("Failed to delete DNS record", "error", err)
			}
		}()
		w.WriteHeader(http.StatusNoContent)
	}
}

func GetRouterDNSProvider(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dns_provider db.GetRouterDNSProviderParams
		if err := json.NewDecoder(r.Body).Decode(&dns_provider); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		providers, err := q.GetRouterDNSProvider(r.Context(), dns_provider)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(providers)
	}
}

func DeleteRouterDNSProvider(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dns_provider db.DeleteRouterDNSProviderParams
		if err := json.NewDecoder(r.Body).Decode(&dns_provider); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		go func() {
			err := dns.DeleteDNS(q, dns_provider.TraefikID, dns_provider.RouterName)
			if err != nil {
				slog.Error("Failed to delete DNS record", "error", err)
			}
		}()
		if err := q.DeleteRouterDNSProvider(r.Context(), dns_provider); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
