package traefik

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/MizuchiLabs/mantrae/pkg/util"
)

func (r *Router) Verify() error {
	if r.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if r.RouterType == "" {
		return fmt.Errorf("router type cannot be empty")
	}
	if r.Rule == "" && r.RouterType != "udp" {
		return fmt.Errorf("rule cannot be empty")
	}
	if r.Provider == "" {
		return fmt.Errorf("provider cannot be empty")
	}

	r.Rule = strings.TrimSpace(r.Rule)
	r.Name = validateName(r.Name, r.Provider)
	if r.Service == "" {
		r.Service = r.Name
	}
	return nil
}

func (s *Service) Verify() error {
	if s.Name == "" {
		return fmt.Errorf("service name cannot be empty")
	}
	if s.ServiceType == "" {
		return fmt.Errorf("service type cannot be empty")
	}
	if s.Provider == "" {
		return fmt.Errorf("provider cannot be empty")
	}

	s.Name = validateName(s.Name, s.Provider)
	return nil
}

func (m *Middleware) Verify() error {
	if m.Name == "" {
		return fmt.Errorf("middleware name cannot be empty")
	}
	if m.Type == "" {
		return fmt.Errorf("type cannot be empty")
	}
	if m.MiddlewareType == "" {
		return fmt.Errorf("middleware type cannot be empty")
	}
	if m.Provider == "" {
		return fmt.Errorf("provider cannot be empty")
	}

	m.Name = validateName(m.Name, m.Provider)

	// Hashes the password strings in the middleware
	if m.BasicAuth != nil {
		for i, u := range m.BasicAuth.Users {
			hash, err := util.HashBasicAuth(u)
			if err != nil {
				return fmt.Errorf("error hashing password: %s", err.Error())
			}
			m.BasicAuth.Users[i] = hash
		}
	}
	if m.DigestAuth != nil {
		for i, u := range m.DigestAuth.Users {
			hash, err := util.HashBasicAuth(u)
			if err != nil {
				return fmt.Errorf("error hashing password: %s", err.Error())
			}
			m.DigestAuth.Users[i] = hash
		}
	}

	// Make sure the correct fields are set in the new middleware
	newMiddleware := &Middleware{
		Name:           m.Name,
		Provider:       m.Provider,
		Type:           m.Type,
		Status:         m.Status,
		MiddlewareType: m.MiddlewareType,
	}

	// Reflect on the original middleware struct
	v := reflect.ValueOf(m)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	// Find and set the field that matches the Type
	for i := 0; i < v.NumField(); i++ {
		fieldName := v.Type().Field(i).Name
		if strings.EqualFold(fieldName, m.Type) {
			if reflect.ValueOf(newMiddleware).Kind() == reflect.Ptr {
				newField := reflect.ValueOf(newMiddleware).Elem().FieldByName(fieldName)
				if newField.CanSet() {
					newField.Set(v.Field(i))
				}
			} else {
				newField := reflect.ValueOf(newMiddleware).FieldByName(fieldName)
				if newField.CanSet() {
					newField.Set(v.Field(i))
				}
			}
			break
		}
	}
	*m = *newMiddleware
	return nil
}

func validateName(s, p string) string {
	name := strings.ToLower(s)
	parts := strings.Split(name, "@")

	if len(parts) > 1 {
		name = parts[0] + "@" + parts[1]
	} else {
		name = parts[0] + "@" + p
	}

	return name
}
