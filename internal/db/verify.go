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
	r.Service = r.Name
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
	s.LoadBalancer, _ = json.Marshal(s.LoadBalancer)
	s.Failover, _ = json.Marshal(s.Failover)
	s.Mirroring, _ = json.Marshal(s.Mirroring)
	s.Weighted, _ = json.Marshal(s.Weighted)
	s.ServerStatus, _ = json.Marshal(s.ServerStatus)
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

	// Unmarshal existing errors if any
	errorsMap := make(map[string]any)
	if r.Errors != nil {
		if existingMap, ok := r.Errors.(map[string]interface{}); ok {
			errorsMap = existingMap
		} else {
			slog.Error("Unexpected type for errors config", "type", fmt.Sprintf("%T", r.Errors))
			return
		}
	}

	// Perform SSL validation and update only the "ssl" error field
	if domain, _ := util.ExtractDomainFromRule(r.Rule); domain != "" {
		slog.Info("Checking SSL", "domain", domain)
		if err := util.ValidSSLCert(domain); err != nil {
			errorsMap["ssl"] = err.Error() // Set "ssl" error message
		} else {
			delete(errorsMap, "ssl") // Remove "ssl" error if SSL is valid
		}
	}

	// Marshal back to JSON and update `Errors`
	if len(errorsMap) > 0 {
		marshalledErrors, err := json.Marshal(errorsMap)
		if err != nil {
			return
		}
		r.Errors = marshalledErrors
	} else {
		r.Errors = nil // Clear `Errors` if no errors remain
	}

	_, err := Query.UpsertRouter(context.Background(), UpsertRouterParams{
		Name:      r.Name,
		ProfileID: r.ProfileID,
		Errors:    r.Errors,
	})
	if err != nil {
		slog.Error("Failed to update router", "error", err)
	}
}
