package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/internal/source"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
)

type UpsertRouterParams struct {
	Name       string              `json:"name"`
	Type       string              `json:"type"`
	Router     *dynamic.Router     `json:"router"`
	TCPRouter  *dynamic.TCPRouter  `json:"tcpRouter"`
	UDPRouter  *dynamic.UDPRouter  `json:"udpRouter"`
	Service    *dynamic.Service    `json:"service"`
	TCPService *dynamic.TCPService `json:"tcpService"`
	UDPService *dynamic.UDPService `json:"udpService"`
}

// UpsertRouter handles both creation and updates of router/service pairs
func UpsertRouter(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params UpsertRouterParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := validateRouterParams(&params); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Get existing config
		profileID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid profile ID", http.StatusBadRequest)
			return
		}

		existingConfig, err := q.GetTraefikConfigBySource(
			r.Context(),
			db.GetTraefikConfigBySourceParams{
				ProfileID: profileID,
				Source:    source.Local,
			},
		)
		if err != nil {
			http.Error(
				w,
				"Failed to get existing config: "+err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		// Initialize maps if nil
		if existingConfig.Config.Routers == nil {
			existingConfig.Config.Routers = make(map[string]*dynamic.Router)
		}
		if existingConfig.Config.Services == nil {
			existingConfig.Config.Services = make(map[string]*dynamic.Service)
		}
		if existingConfig.Config.TCPRouters == nil {
			existingConfig.Config.TCPRouters = make(map[string]*dynamic.TCPRouter)
		}
		if existingConfig.Config.TCPServices == nil {
			existingConfig.Config.TCPServices = make(map[string]*dynamic.TCPService)
		}
		if existingConfig.Config.UDPRouters == nil {
			existingConfig.Config.UDPRouters = make(map[string]*dynamic.UDPRouter)
		}
		if existingConfig.Config.UDPServices == nil {
			existingConfig.Config.UDPServices = make(map[string]*dynamic.UDPService)
		}

		// Ensure name has @http suffix
		if !strings.HasSuffix(params.Name, "@http") {
			params.Name = fmt.Sprintf("%s@http", strings.Split(params.Name, "@")[0])
		}

		// Update configuration based on type
		switch params.Type {
		case "http":
			existingConfig.Config.Routers[params.Name] = params.Router
			existingConfig.Config.Services[params.Name] = params.Service
		case "tcp":
			existingConfig.Config.TCPRouters[params.Name] = params.TCPRouter
			existingConfig.Config.TCPServices[params.Name] = params.TCPService
		case "udp":
			existingConfig.Config.UDPRouters[params.Name] = params.UDPRouter
			existingConfig.Config.UDPServices[params.Name] = params.UDPService
		default:
			http.Error(w, "invalid router type: must be http, tcp, or udp", http.StatusBadRequest)
			return
		}

		err = q.UpdateTraefikConfig(r.Context(), db.UpdateTraefikConfigParams{
			ID:     existingConfig.ID,
			Source: source.Local,
			Config: existingConfig.Config,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return the updated configuration
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(existingConfig.Config)
	}
}

// DeleteRouter handles the removal of router/service pairs
func DeleteRouter(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profileID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid profile ID", http.StatusBadRequest)
			return
		}

		routerName := r.PathValue("name")
		routerType := r.PathValue("type")

		if routerName == "" || routerType == "" {
			http.Error(w, "Missing router name or type", http.StatusBadRequest)
			return
		}

		// Ensure name has @http suffix
		if !strings.HasSuffix(routerName, "@http") {
			http.Error(w, "Invalid router provider", http.StatusBadRequest)
			return
		}

		existingConfig, err := q.GetTraefikConfigBySource(
			r.Context(),
			db.GetTraefikConfigBySourceParams{
				ProfileID: profileID,
				Source:    source.Local,
			},
		)
		if err != nil {
			http.Error(
				w,
				"Failed to get existing config: "+err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		// Remove router and service based on type
		switch routerType {
		case "http":
			delete(existingConfig.Config.Routers, routerName)
			delete(existingConfig.Config.Services, routerName)
		case "tcp":
			delete(existingConfig.Config.TCPRouters, routerName)
			delete(existingConfig.Config.TCPServices, routerName)
		case "udp":
			delete(existingConfig.Config.UDPRouters, routerName)
			delete(existingConfig.Config.UDPServices, routerName)
		default:
			http.Error(w, "invalid router type: must be http, tcp, or udp", http.StatusBadRequest)
			return
		}

		err = q.UpdateTraefikConfig(r.Context(), db.UpdateTraefikConfigParams{
			ID:     existingConfig.ID,
			Source: source.Local,
			Config: existingConfig.Config,
		})
		if err != nil {
			http.Error(w, "Failed to update config: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func validateRouterParams(params *UpsertRouterParams) error {
	var (
		ErrInvalidRouterType = errors.New("invalid router type: must be http, tcp, or udp")
		ErrMissingRouter     = errors.New("missing router configuration")
		ErrMissingService    = errors.New("missing service configuration")
		ErrInvalidName       = errors.New("invalid router name")
	)

	if params.Name == "" {
		return ErrInvalidName
	}

	switch params.Type {
	case "http":
		if params.Router == nil {
			return ErrMissingRouter
		}
		if params.Service == nil {
			return ErrMissingService
		}
		// Validate HTTP specific fields
		if params.Router.Rule == "" {
			return errors.New("http router requires a rule")
		}
		if len(params.Router.EntryPoints) == 0 {
			return errors.New("http router requires at least one entrypoint")
		}
	case "tcp":
		if params.TCPRouter == nil {
			return ErrMissingRouter
		}
		if params.TCPService == nil {
			return ErrMissingService
		}
		// Validate TCP specific fields
		if params.TCPRouter.Rule == "" {
			return errors.New("tcp router requires a rule")
		}
	case "udp":
		if params.UDPRouter == nil {
			return ErrMissingRouter
		}
		if params.UDPService == nil {
			return ErrMissingService
		}
		// Validate UDP specific fields
		if len(params.UDPRouter.EntryPoints) == 0 {
			return errors.New("udp router requires at least one entrypoint")
		}
	default:
		return ErrInvalidRouterType
	}

	return nil
}
