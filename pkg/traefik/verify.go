package traefik

import (
	"fmt"
	"net"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/MizuchiLabs/mantrae/pkg/util"
)

func (p *Profile) Verify() error {
	if p.Name == "" {
		return fmt.Errorf("profile name cannot be empty")
	}

	if p.URL != "" {
		if !isValidURL(p.URL) {
			return fmt.Errorf("invalid url")
		}
	} else {
		return fmt.Errorf("url cannot be empty")
	}
	return nil
}

func (r *Router) Verify() error {
	if r.Name == "" {
		return fmt.Errorf("router name cannot be empty")
	}
	if r.Service == "" {
		return fmt.Errorf("service cannot be empty")
	}
	if r.RouterType == "" {
		return fmt.Errorf("router type cannot be empty")
	}
	if r.Rule == "" && r.RouterType != "udp" {
		return fmt.Errorf("rule cannot be empty")
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
	return nil
}

func (m *Middleware) Verify() error {
	if m.Name == "" {
		return fmt.Errorf("middleware name cannot be empty")
	}
	if m.Provider == "" {
		return fmt.Errorf("provider cannot be empty")
	}
	if m.Type == "" {
		return fmt.Errorf("type cannot be empty")
	}

	// Hashes the password strings in the middleware
	if m.BasicAuth != nil {
		for i, u := range m.BasicAuth.Users {
			hash, err := util.HashPassword(u)
			if err != nil {
				return fmt.Errorf("error hashing password: %s", err.Error())
			}
			m.BasicAuth.Users[i] = hash
		}
	}
	if m.DigestAuth != nil {
		for i, u := range m.DigestAuth.Users {
			hash, err := util.HashPassword(u)
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

func isValidURL(u string) bool {
	// If no scheme is provided, prepend "http://"
	if !strings.Contains(u, "://") {
		u = "http://" + u
	}

	parsedURL, err := url.Parse(u)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return false
	}

	host := parsedURL.Hostname()
	port := parsedURL.Port()
	if port != "" {
		if _, err = net.LookupPort("tcp", port); err != nil {
			return false
		}
	}

	// Check if it's an IP address (including loopback)
	ip := net.ParseIP(host)
	if ip != nil {
		return true
	}

	// Check if it's localhost
	if !strings.Contains(host, ".") {
		_, err = strconv.Atoi(host)
		return err != nil // Valid if it's not just a number
	}

	// Check if it's a valid domain name
	domainRegex := `^([a-zA-Z0-9-]+\.)*[a-zA-Z0-9-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(domainRegex, host)
	if err != nil {
		return false
	}
	return matched
}
