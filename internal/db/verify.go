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

func (r *UpsertRouterParams) Verify() error {
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
		r.Errors, _ = SetError(r.Errors, "dns", "")
	}
	r.EntryPoints, _ = json.Marshal(r.EntryPoints)
	r.Middlewares, _ = json.Marshal(r.Middlewares)
	r.Tls, _ = json.Marshal(r.Tls)
	r.Errors, _ = json.Marshal(r.Errors)
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
	// if s.Provider == "http" {
	// 	if s.LoadBalancer != nil {
	// 		validServers := make([]traefik.Server, 0)
	// 		for _, server := range s.LoadBalancer.Servers {
	// 			if server.Address != "" || server.URL != "" {
	// 				validServers = append(validServers, server)
	// 			}
	// 		}
	//
	// 		if len(validServers) == 0 {
	// 			return fmt.Errorf("no valid servers found in load balancer")
	// 		}
	//
	// 		s.LoadBalancer.Servers = validServers
	// 	} else {
	// 		return fmt.Errorf("load balancer cannot be nil")
	// 	}
	// }

	s.LoadBalancer, _ = json.Marshal(s.LoadBalancer)
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

	m.Content, _ = json.Marshal(m.Content)

	// 	if m.Type == "basicauth" {
	// 	for i, u := range m.Content["BasicAuth"].(map[string]interface{}).Users {
	// 		hash, err := util.HashBasicAuth(u)
	// 		if err != nil {
	// 			return fmt.Errorf("error hashing password: %s", err.Error())
	// 		}
	// 		m.BasicAuth.Users[i] = hash
	// 	}
	// }
	// if m.DigestAuth != nil {
	// 	for i, u := range m.DigestAuth.Users {
	// 		hash, err := util.HashBasicAuth(u)
	// 		if err != nil {
	// 			return fmt.Errorf("error hashing password: %s", err.Error())
	// 		}
	// 		m.DigestAuth.Users[i] = hash
	// 	}
	// }

	return nil
}

func (u *CreateUserParams) Verify() error {
	if u.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if u.Password == "" {
		return fmt.Errorf("password cannot be empty")
	}
	if u.Type == "" {
		return fmt.Errorf("user type cannot be empty")
	}
	if u.Email != nil && *u.Email != "" {
		if !util.IsValidEmail(*u.Email) {
			return fmt.Errorf("email is not valid")
		}
	}

	hash, err := util.HashPassword(strings.TrimSpace(u.Password))
	if err != nil {
		return fmt.Errorf("failed to hash password: %s", err.Error())
	}
	u.Password = hash
	return nil
}

func (u *UpdateUserParams) Verify() error {
	if u.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if u.Password == "" {
		return fmt.Errorf("password cannot be empty")
	}
	if u.Type == "" {
		return fmt.Errorf("user type cannot be empty")
	}
	if u.Email != nil && *u.Email != "" {
		if !util.IsValidEmail(*u.Email) {
			return fmt.Errorf("email is not valid")
		}
	}

	hash, err := util.HashPassword(strings.TrimSpace(u.Password))
	if err != nil {
		return fmt.Errorf("failed to hash password: %s", err.Error())
	}
	u.Password = hash
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
func (r *Router) DecodeFields() error {
	if r.EntryPoints != nil {
		if err := json.Unmarshal(r.EntryPoints.([]byte), &r.EntryPoints); err != nil {
			return err
		}
	}
	if r.Middlewares != nil {
		if err := json.Unmarshal(r.Middlewares.([]byte), &r.Middlewares); err != nil {
			return err
		}
	}
	if r.Tls != nil {
		if err := json.Unmarshal(r.Tls.([]byte), &r.Tls); err != nil {
			return err
		}
	}
	if r.Errors != nil {
		if err := json.Unmarshal(r.Errors.([]byte), &r.Errors); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) DecodeFields() error {
	if s.LoadBalancer != nil {
		if err := json.Unmarshal(s.LoadBalancer.([]byte), &s.LoadBalancer); err != nil {
			return err
		}
	}
	if s.Failover != nil {
		if err := json.Unmarshal(s.Failover.([]byte), &s.Failover); err != nil {
			return err
		}
	}
	if s.Mirroring != nil {
		if err := json.Unmarshal(s.Mirroring.([]byte), &s.Mirroring); err != nil {
			return err
		}
	}
	if s.Weighted != nil {
		if err := json.Unmarshal(s.Weighted.([]byte), &s.Weighted); err != nil {
			return err
		}
	}
	if s.ServerStatus != nil {
		if err := json.Unmarshal(s.ServerStatus.([]byte), &s.ServerStatus); err != nil {
			return err
		}
	}
	return nil
}

func (m *Middleware) DecodeFields() error {
	if m.Content != nil {
		if err := json.Unmarshal(m.Content.([]byte), &m.Content); err != nil {
			return err
		}
	}

	return nil
}

// SetError adds an error message with a custom key to the Errors map in the provided router.
func SetError(errors interface{}, key string, message string) ([]byte, error) {
	errMap := make(map[string]string)
	switch e := errors.(type) {
	case map[string]interface{}:
		for k, v := range e {
			if strVal, ok := v.(string); ok {
				errMap[k] = strVal
			}
		}
	case string:
		if err := json.Unmarshal([]byte(e), &errMap); err != nil {
			return nil, err
		}
	}

	if message == "" {
		delete(errMap, key)
	} else {
		errMap[key] = message
	}

	updatedErrors, err := json.Marshal(errMap)
	if err != nil {
		return nil, err
	}
	return updatedErrors, nil
}

func (r *Router) SSLCheck() {
	tlsMap := make(map[string]any)
	if r.Tls != nil {
		if existingMap, ok := r.Tls.(map[string]interface{}); ok {
			tlsMap = existingMap
		} else {
			slog.Error("Unexpected type for TLS config", "type", fmt.Sprintf("%T", r.Tls))
			return
		}
	}
	if tlsMap["certResolver"] == "" {
		return
	}

	// Perform SSL validation and update only the "ssl" error field
	if domain, _ := util.ExtractDomainFromRule(r.Rule); domain != "" {
		if err := util.ValidSSLCert(domain); err != nil {
			r.Errors, err = SetError(r.Errors, "ssl", err.Error())
			if err != nil {
				slog.Error("Failed to update router", "error", err)
				return
			}
		} else {
			r.Errors, err = SetError(r.Errors, "ssl", "")
			if err != nil {
				slog.Error("Failed to update router", "error", err)
				return
			}
		}
	}
	if _, err := Query.UpsertRouter(context.Background(), UpsertRouterParams{
		Name:      r.Name,
		ProfileID: r.ProfileID,
		Errors:    r.Errors,
	}); err != nil {
		slog.Error("Failed to update router", "error", err)
	}
}
