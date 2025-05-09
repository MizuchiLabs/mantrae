package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/MizuchiLabs/mantrae/internal/config"
	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/internal/source"
	"github.com/MizuchiLabs/mantrae/internal/util"
	"github.com/traefik/traefik/v3/pkg/config/runtime"
)

type UpsertRouterParams struct {
	Name       string                  `json:"name"`
	Protocol   string                  `json:"protocol"`
	Router     *runtime.RouterInfo     `json:"router"`
	TCPRouter  *runtime.TCPRouterInfo  `json:"tcpRouter"`
	UDPRouter  *runtime.UDPRouterInfo  `json:"udpRouter"`
	Service    *db.ServiceInfo         `json:"service"`
	TCPService *runtime.TCPServiceInfo `json:"tcpService"`
	UDPService *runtime.UDPServiceInfo `json:"udpService"`
}

type DeleteRouterParams struct {
	ProfileID int64  `json:"profileId"`
	Name      string `json:"name"`
	Protocol  string `json:"protocol"`
}

type BulkDeleteRouterParams struct {
	ProfileID int64 `json:"profileId"`
	Items     []struct {
		Name     string `json:"name"`
		Protocol string `json:"protocol"`
	} `json:"items"`
}

// UpsertRouter handles both creation and updates of router/service pairs
func UpsertRouter(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := a.Conn.GetQuery()
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

		existingConfig, err := q.GetLocalTraefikConfig(r.Context(), profileID)
		if err != nil {
			http.Error(
				w,
				"Failed to get existing config: "+err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		// Initialize maps if nil
		if existingConfig.Config == nil {
			existingConfig.Config = &db.TraefikConfiguration{}
		}
		if existingConfig.Config.Routers == nil {
			existingConfig.Config.Routers = make(map[string]*runtime.RouterInfo)
		}
		if existingConfig.Config.Services == nil {
			existingConfig.Config.Services = make(map[string]*db.ServiceInfo)
		}
		if existingConfig.Config.TCPRouters == nil {
			existingConfig.Config.TCPRouters = make(map[string]*runtime.TCPRouterInfo)
		}
		if existingConfig.Config.TCPServices == nil {
			existingConfig.Config.TCPServices = make(map[string]*runtime.TCPServiceInfo)
		}
		if existingConfig.Config.UDPRouters == nil {
			existingConfig.Config.UDPRouters = make(map[string]*runtime.UDPRouterInfo)
		}
		if existingConfig.Config.UDPServices == nil {
			existingConfig.Config.UDPServices = make(map[string]*runtime.UDPServiceInfo)
		}

		// Ensure name has no @
		params.Name = strings.Split(params.Name, "@")[0]

		// Get default dns provider
		dnsProvider, _ := q.GetActiveDNSProvider(r.Context())

		// Update configuration based on type
		switch params.Protocol {
		case "http":
			if !strings.HasSuffix(params.Router.Service, "@http") {
				params.Router.Service = fmt.Sprintf(
					"%s@http",
					strings.Split(params.Router.Service, "@")[0],
				)
			}
			// Check if router is new and add default dns provider (if there is one)
			if _, ok := existingConfig.Config.Routers[params.Name]; !ok && dnsProvider.ID != 0 {
				if err = q.AddRouterDNSProvider(r.Context(), db.AddRouterDNSProviderParams{
					TraefikID:  existingConfig.ID,
					RouterName: params.Name,
					ProviderID: dnsProvider.ID,
				}); err != nil {
					slog.Error("Failed to add router dns provider", "error", err)
				}
			}
			existingConfig.Config.Routers[params.Name] = params.Router
			existingConfig.Config.Services[params.Name] = params.Service
		case "tcp":
			if !strings.HasSuffix(params.TCPRouter.Service, "@http") {
				params.TCPRouter.Service = fmt.Sprintf(
					"%s@http",
					strings.Split(params.TCPRouter.Service, "@")[0],
				)
			}
			existingConfig.Config.TCPRouters[params.Name] = params.TCPRouter
			existingConfig.Config.TCPServices[params.Name] = params.TCPService
		case "udp":
			if !strings.HasSuffix(params.UDPRouter.Service, "@http") {
				params.UDPRouter.Service = fmt.Sprintf(
					"%s@http",
					strings.Split(params.UDPRouter.Service, "@")[0],
				)
			}
			existingConfig.Config.UDPRouters[params.Name] = params.UDPRouter
			existingConfig.Config.UDPServices[params.Name] = params.UDPService
		default:
			http.Error(w, "invalid router type: must be http, tcp, or udp", http.StatusBadRequest)
			return
		}

		err = q.UpsertTraefikConfig(r.Context(), db.UpsertTraefikConfigParams{
			ProfileID: existingConfig.ProfileID,
			Source:    source.Local,
			Config:    existingConfig.Config,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		util.Broadcast <- util.EventMessage{
			Type:     util.EventTypeUpdate,
			Category: util.EventCategoryTraefik,
		}

		// Return the updated configuration
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(existingConfig.Config); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// DeleteRouter handles the removal of router/service pairs
func DeleteRouter(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params DeleteRouterParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if params.ProfileID == 0 || params.Name == "" || params.Protocol == "" {
			http.Error(w, "Missing router name or protocol", http.StatusBadRequest)
			return
		}

		q := a.Conn.GetQuery()
		existingConfig, err := q.GetLocalTraefikConfig(r.Context(), params.ProfileID)
		if err != nil {
			http.Error(
				w,
				"Failed to get existing config: "+err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		// Remove router and service based on type
		switch params.Protocol {
		case "http":
			delete(existingConfig.Config.Routers, params.Name)
			delete(existingConfig.Config.Services, params.Name)
		case "tcp":
			delete(existingConfig.Config.TCPRouters, params.Name)
			delete(existingConfig.Config.TCPServices, params.Name)
		case "udp":
			delete(existingConfig.Config.UDPRouters, params.Name)
			delete(existingConfig.Config.UDPServices, params.Name)
		default:
			http.Error(w, "invalid router type: must be http, tcp, or udp", http.StatusBadRequest)
			return
		}

		err = q.UpsertTraefikConfig(r.Context(), db.UpsertTraefikConfigParams{
			ProfileID: existingConfig.ProfileID,
			Source:    source.Local,
			Config:    existingConfig.Config,
		})
		if err != nil {
			http.Error(w, "Failed to update config: "+err.Error(), http.StatusInternalServerError)
			return
		}

		util.Broadcast <- util.EventMessage{
			Type:     util.EventTypeDelete,
			Category: util.EventCategoryTraefik,
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func BulkDeleteRouter(a *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params BulkDeleteRouterParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if params.ProfileID == 0 {
			http.Error(w, "Missing router name or protocol", http.StatusBadRequest)
			return
		}

		q := a.Conn.GetQuery()
		existingConfig, err := q.GetLocalTraefikConfig(r.Context(), params.ProfileID)
		if err != nil {
			http.Error(
				w,
				"Failed to get existing config: "+err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		// Remove router and service based on type
		for _, item := range params.Items {
			if item.Name == "" || item.Protocol == "" {
				continue // Skip invalid entries
			}

			switch item.Protocol {
			case "http":
				delete(existingConfig.Config.Routers, item.Name)
				delete(existingConfig.Config.Services, item.Name)
			case "tcp":
				delete(existingConfig.Config.TCPRouters, item.Name)
				delete(existingConfig.Config.TCPServices, item.Name)
			case "udp":
				delete(existingConfig.Config.UDPRouters, item.Name)
				delete(existingConfig.Config.UDPServices, item.Name)
			default:
				http.Error(
					w,
					"invalid router type: must be http, tcp, or udp",
					http.StatusBadRequest,
				)
				return
			}
		}

		err = q.UpsertTraefikConfig(r.Context(), db.UpsertTraefikConfigParams{
			ProfileID: existingConfig.ProfileID,
			Source:    source.Local,
			Config:    existingConfig.Config,
		})
		if err != nil {
			http.Error(w, "Failed to update config: "+err.Error(), http.StatusInternalServerError)
			return
		}

		util.Broadcast <- util.EventMessage{
			Type:     util.EventTypeDelete,
			Category: util.EventCategoryTraefik,
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

	switch params.Protocol {
	case "http":
		if params.Router == nil {
			return ErrMissingRouter
		}
		if params.Service == nil {
			return ErrMissingService
		}
		// Validate HTTP specific fields
		if params.Router.Rule == "" {
			return errors.New("HTTP router requires a rule")
		}
		if len(params.Router.EntryPoints) == 0 {
			return errors.New("HTTP router requires at least one entrypoint")
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
			return errors.New("TCP router requires a rule")
		}
		if params.TCPRouter.EntryPoints == nil {
			return errors.New("TCP router requires at least one entrypoint")
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
			return errors.New("UDP router requires at least one entrypoint")
		}
	default:
		return ErrInvalidRouterType
	}

	return nil
}
