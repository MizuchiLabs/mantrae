package api

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
)

func (p *Profile) Verify() error {
	if p.Name == "" {
		return fmt.Errorf("profile name cannot be empty")
	}

	if p.Instance.URL != "" {
		if !isValidURL(p.Instance.URL) {
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

func isValidURL(u string) bool {
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

	ip := net.ParseIP(host)
	if ip != nil {
		return true
	}

	if host != "localhost" {
		domainRegex := `^([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`
		matched, err := regexp.MatchString(domainRegex, host)
		if err != nil {
			return false
		}
		return matched
	}
	return true
}
