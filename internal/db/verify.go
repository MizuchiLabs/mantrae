package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/MizuchiLabs/mantrae/pkg/util"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

type Server struct {
	URL     string `json:"url,omitempty"`
	Address string `json:"address,omitempty"`
}

func (p *CreateProfileParams) Verify() error {
	if p.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if p.Url == "" {
		return fmt.Errorf("url cannot be empty")
	}
	if !util.IsValidURL(p.Url) {
		return fmt.Errorf("url is not valid")
	}
	return nil
}

func (p *UpdateProfileParams) Verify() error {
	if p.ID == 0 {
		return fmt.Errorf("id cannot be empty")
	}
	if p.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if p.Url == "" {
		return fmt.Errorf("url cannot be empty")
	}
	if !util.IsValidURL(p.Url) {
		return fmt.Errorf("url is not valid")
	}
	return nil
}

func (e *UpsertEntryPointParams) Verify() error {
	if e.ProfileID == 0 {
		return fmt.Errorf("profile id cannot be empty")
	}
	if e.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if e.Address == "" {
		return fmt.Errorf("address cannot be empty")
	}
	e.Http, _ = json.Marshal(e.Http)
	return nil
}

func (r *Router) Verify() error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	if r.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if r.ProfileID == 0 {
		return fmt.Errorf("profile id cannot be empty")
	}
	if r.Protocol == "" {
		return fmt.Errorf("protocol cannot be empty")
	}
	if r.Provider == "" {
		return fmt.Errorf("provider cannot be empty")
	}
	if r.Provider == "http" {
		r.Service = r.Name
	}
	if r.DnsProvider == nil || *r.DnsProvider == 0 {
		r.UpdateError("dns", "")
	}
	r.EntryPoints, _ = json.Marshal(r.EntryPoints)
	r.Middlewares, _ = json.Marshal(r.Middlewares)
	r.Tls, _ = json.Marshal(r.Tls)
	r.Errors, _ = json.Marshal(r.Errors)

	// Check if router name has changed
	oldRouter, err := Query.GetRouterByID(context.Background(), r.ID)
	if err == nil && oldRouter.Name != r.Name {
		if err = Query.DeleteRouterByID(context.Background(), r.ID); err != nil {
			return err
		}
	}
	return nil
}

func (s *UpsertServiceParams) Verify() error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	if s.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if s.ProfileID == 0 {
		return fmt.Errorf("profile id cannot be empty")
	}
	if s.Protocol == "" {
		return fmt.Errorf("protocol cannot be empty")
	}

	loadBalancerBytes, err := json.Marshal(s.LoadBalancer)
	if err != nil {
		return fmt.Errorf("error marshaling LoadBalancer: %s", err.Error())
	}

	// Unmarshal LoadBalancer into a map for processing
	var loadBalancerMap map[string]interface{}
	if err := json.Unmarshal(loadBalancerBytes, &loadBalancerMap); err != nil {
		return fmt.Errorf("error unmarshaling LoadBalancer: %s", err.Error())
	}

	if s.Provider == "http" {
		// Process servers in LoadBalancer
		if servers, found := loadBalancerMap["servers"].([]interface{}); found {
			validServers := make([]Server, 0)
			for _, server := range servers {
				serverMap, ok := server.(map[string]interface{})
				if !ok {
					return fmt.Errorf("invalid server format")
				}

				address, addressOk := serverMap["address"].(string)
				url, urlOk := serverMap["url"].(string)

				// Initialize address and url if they are not present or nil
				if addressOk && address != "" {
					validServers = append(validServers, Server{
						Address: address,
					})
				}
				if urlOk && url != "" {
					validServers = append(validServers, Server{
						URL: url,
					})
				}
			}

			if len(validServers) == 0 {
				return fmt.Errorf("no valid servers found in load balancer")
			}

			// Reassign valid servers back to the LoadBalancer map
			loadBalancerMap["servers"] = validServers
		} else {
			return fmt.Errorf("servers key not found in LoadBalancer")
		}
	}

	s.LoadBalancer, _ = json.Marshal(loadBalancerMap)
	s.Failover, _ = json.Marshal(s.Failover)
	s.Mirroring, _ = json.Marshal(s.Mirroring)
	s.Weighted, _ = json.Marshal(s.Weighted)
	s.ServerStatus, _ = json.Marshal(s.ServerStatus)
	return nil
}

func (m *UpsertMiddlewareParams) Verify() error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	if m.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if m.ProfileID == 0 {
		return fmt.Errorf("profile id cannot be empty")
	}
	if m.Protocol == "" {
		return fmt.Errorf("protocol cannot be empty")
	}

	// Marshal the content field to JSON
	contentBytes, err := json.Marshal(m.Content)
	if err != nil {
		return fmt.Errorf("error marshaling content: %s", err.Error())
	}

	// Unmarshal the JSON bytes into a map for processing based on type
	var contentMap map[string]interface{}
	if err := json.Unmarshal(contentBytes, &contentMap); err != nil {
		return fmt.Errorf("error unmarshaling content: %s", err.Error())
	}

	// Handle BasicAuth and DigestAuth types
	switch m.Type {
	case "basicauth":
		if users, found := contentMap["users"].([]interface{}); found {
			for i, u := range users {
				user, ok := u.(string)
				if !ok {
					return fmt.Errorf("invalid user format in BasicAuth")
				}
				hash, err := util.HashBasicAuth(user)
				if err != nil {
					return fmt.Errorf("error hashing password: %s", err.Error())
				}
				users[i] = hash
			}
		} else {
			return fmt.Errorf("users key not found in BasicAuth")
		}
	case "digestauth":
		if users, found := contentMap["users"].([]interface{}); found {
			for i, u := range users {
				user, ok := u.(string)
				if !ok {
					return fmt.Errorf("invalid user format in DigestAuth")
				}
				hash, err := util.HashBasicAuth(user)
				if err != nil {
					return fmt.Errorf("error hashing password: %s", err.Error())
				}
				users[i] = hash
			}
		} else {
			return fmt.Errorf("users key not found in DigestAuth")
		}
	}

	m.Content, _ = json.Marshal(contentMap)
	return nil
}

func (a *UpsertAgentParams) Verify() error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	if a.Hostname == "" {
		return fmt.Errorf("name cannot be empty")
	}

	a.PrivateIps, _ = json.Marshal(a.PrivateIps)
	a.Containers, _ = json.Marshal(a.Containers)
	return nil
}

func (u *User) Verify() error {
	if u.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if u.Email != nil && *u.Email != "" {
		if !util.IsValidEmail(*u.Email) {
			return fmt.Errorf("email is not valid")
		}
	}
	if u.ID == 0 && u.Password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	if u.Password != "" {
		hash, err := util.HashPassword(strings.TrimSpace(u.Password))
		if err != nil {
			return fmt.Errorf("failed to hash password: %s", err.Error())
		}
		u.Password = hash
	} else {
		user, err := Query.GetUserByID(context.Background(), u.ID)
		if err != nil {
			return fmt.Errorf("failed to get user: %s", err.Error())
		}
		u.Password = user.Password
	}
	return nil
}

func (s *CreateSettingParams) Verify() error {
	if s.Key == "" {
		return fmt.Errorf("key cannot be empty")
	}
	if s.Value == "" {
		return fmt.Errorf("value cannot be empty")
	}
	return nil
}

func (s *UpdateSettingParams) Verify() error {
	if s.Key == "" {
		return fmt.Errorf("key cannot be empty")
	}
	if s.Value == "" {
		return fmt.Errorf("value cannot be empty")
	}

	if s.Key == "backup-schedule" {
		_, err := cron.ParseStandard(s.Value)
		if err != nil {
			return fmt.Errorf("invalid cron expression: %s", err.Error())
		}
	}
	if s.Key == "backup-keep" {
		_, err := strconv.Atoi(s.Value)
		if err != nil {
			return fmt.Errorf("invalid backup-keep value: %s", err.Error())
		}
	}
	return nil
}

func (p *CreateProviderParams) Verify() error {
	if p.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if p.Type == "" {
		return fmt.Errorf("provider type cannot be empty")
	}
	if p.ExternalIp == "" {
		return fmt.Errorf("external ip cannot be empty")
	}
	if p.ApiKey == "" {
		return fmt.Errorf("api key cannot be empty")
	}
	return nil
}

func (p *UpdateProviderParams) Verify() error {
	if p.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if p.Type == "" {
		return fmt.Errorf("provider type cannot be empty")
	}
	if p.ExternalIp == "" {
		return fmt.Errorf("external ip cannot be empty")
	}
	if p.ApiKey == "" {
		return fmt.Errorf("api key cannot be empty")
	}
	return nil
}

// Decoding -------------------------------------------------------------------

// Check if the data is a byte slice and attempt to decode
func decode(data interface{}, target interface{}) error {
	if data == nil {
		return nil
	}

	switch data := data.(type) {
	case []byte:
		return json.Unmarshal(data, target)
	case string:
		return json.Unmarshal([]byte(data), target)
	case *json.RawMessage:
		return json.Unmarshal(*data, target)
	default:
		return nil
	}
}

func (e *Entrypoint) DecodeFields() error {
	if err := decode(e.Http, &e.Http); err != nil {
		return fmt.Errorf("field http: %s", err.Error())
	}
	return nil
}

func (r *Router) DecodeFields() error {
	if err := decode(r.EntryPoints, &r.EntryPoints); err != nil {
		return fmt.Errorf("field entrypoints: %s", err.Error())
	}
	if err := decode(r.Middlewares, &r.Middlewares); err != nil {
		return fmt.Errorf("field middlewares: %s", err.Error())
	}
	if err := decode(r.Tls, &r.Tls); err != nil {
		return fmt.Errorf("field tls: %s", err.Error())
	}
	if err := decode(r.Errors, &r.Errors); err != nil {
		return fmt.Errorf("field errors: %s", err.Error())
	}

	return nil
}

func (s *Service) DecodeFields() error {
	if err := decode(s.LoadBalancer, &s.LoadBalancer); err != nil {
		return fmt.Errorf("field loadbalancer: %s", err.Error())
	}
	if err := decode(s.Failover, &s.Failover); err != nil {
		return fmt.Errorf("field failover: %s", err.Error())
	}
	if err := decode(s.Mirroring, &s.Mirroring); err != nil {
		return fmt.Errorf("field mirroring: %s", err.Error())
	}
	if err := decode(s.Weighted, &s.Weighted); err != nil {
		return fmt.Errorf("field weighted: %s", err.Error())
	}
	if err := decode(s.ServerStatus, &s.ServerStatus); err != nil {
		return fmt.Errorf("field serverstatus: %s", err.Error())
	}
	return nil
}

func (m *Middleware) DecodeFields() error {
	if err := decode(m.Content, &m.Content); err != nil {
		return fmt.Errorf("field content: %s", err.Error())
	}

	return nil
}

func (a *Agent) DecodeFields() error {
	if err := decode(a.Containers, &a.Containers); err != nil {
		return fmt.Errorf("field containers: %s", err.Error())
	}
	if err := decode(a.PrivateIps, &a.PrivateIps); err != nil {
		return fmt.Errorf("field privateips: %s", err.Error())
	}
	return nil
}

func (r *Router) UpdateError(key string, value string) {
	if err := r.DecodeFields(); err != nil {
		slog.Error("Failed to decode router", "error", err)
		return
	}

	if r.Errors == nil {
		r.Errors = make(map[string]interface{})
	} else {
		// Attempt to cast, and reset if the type is incorrect
		if _, ok := r.Errors.(map[string]interface{}); !ok {
			slog.Warn("Invalid errors format detected, resetting to empty map")
			r.Errors = make(map[string]interface{})
		}
	}

	if errorMap, ok := r.Errors.(map[string]interface{}); ok {
		if value == "" {
			delete(errorMap, key)
		} else {
			errorMap[key] = value
		}
	}

	updatedErrors, err := json.Marshal(r.Errors)
	if err != nil {
		slog.Error("Failed to update router", "error", err)
		return
	}

	if _, err := Query.UpsertRouter(context.Background(), UpsertRouterParams{
		Name:      r.Name,
		ProfileID: r.ProfileID,
		Errors:    updatedErrors,
	}); err != nil {
		slog.Error("Failed to update router", "error", err)
	}

	if string(updatedErrors) != "{}" {
		util.Broadcast <- util.EventMessage{
			Type:    "router_updated",
			Message: "Updated router " + r.Name + " with errors: " + string(updatedErrors),
		}
	}
}

func (r *Router) SSLCheck() {
	if err := r.DecodeFields(); err != nil {
		slog.Error("Failed to decode router", "error", err)
		return
	}

	if r.EntryPoints == nil || r.Tls == nil {
		r.UpdateError("ssl", "")
		return
	}

	isHTTPS := false
	for _, ep := range r.EntryPoints.([]interface{}) {
		entrypoint, err := Query.GetEntryPointByName(
			context.Background(),
			GetEntryPointByNameParams{
				ProfileID: r.ProfileID,
				Name:      ep.(string),
			},
		)
		if err != nil {
			slog.Error("Failed to get entry point", "name", ep.(string), "error", err)
			continue
		}
		if entrypoint.Address == "443" {
			isHTTPS = true
			break
		}
	}
	if !isHTTPS {
		slog.Debug("Router is not using HTTPS entrypoint", "name", r.Name)
		r.UpdateError("ssl", "")
		return
	}

	if tlsMap, ok := r.Tls.(map[string]interface{}); ok {
		if tlsMap["certResolver"] == "" {
			slog.Debug("Router is not using a certificate resolver", "name", r.Name)
			r.UpdateError("ssl", "")
			return
		}
	} else {
		slog.Error("Unexpected type for TLS config", "type", fmt.Sprintf("%T", r.Tls))
		return
	}

	// Perform SSL validation and update only the "ssl" error field
	if domain, _ := util.ExtractDomainFromRule(r.Rule); domain != "" {
		if err := util.ValidSSLCert(domain); err != nil {
			r.UpdateError("ssl", err.Error())
		} else {
			r.UpdateError("ssl", "")
		}
	}
}
