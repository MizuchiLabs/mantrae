package db

import (
	"fmt"
	"strings"

	"github.com/mizuchilabs/mantrae/internal/util"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
)

func VerifyMiddleware(m *dynamic.Middleware) error {
	if auth := m.BasicAuth; auth != nil {
		// Process each user entry
		for i, user := range auth.Users {
			// Only process if it looks like cleartext (contains colon)
			if strings.Contains(user, ":") {
				parts := strings.SplitN(user, ":", 2)
				if len(parts) != 2 {
					return fmt.Errorf("invalid user:password format")
				}

				// Only encode if not already in htpasswd format
				if !util.IsHtpasswdFormat(parts[1]) {
					encoded, err := util.HashPassword(parts[1])
					if err != nil {
						return fmt.Errorf("failed to encode password: %s", err.Error())
					}
					auth.Users[i] = parts[0] + ":" + encoded
				}
			}
		}
	}

	if auth := m.DigestAuth; auth != nil {
		// Process each user entry
		for i, user := range auth.Users {
			// Only process if it looks like cleartext (contains colon)
			if strings.Contains(user, ":") {
				parts := strings.SplitN(user, ":", 2)
				if len(parts) != 2 {
					return fmt.Errorf("invalid user:password format")
				}

				// Only encode if not already in htpasswd format
				if !util.IsHtpasswdFormat(parts[1]) {
					encoded, err := util.HashPassword(parts[1])
					if err != nil {
						return fmt.Errorf("failed to encode password: %s", err.Error())
					}
					auth.Users[i] = parts[0] + ":" + encoded
				}
			}
		}
	}

	return nil
}

// func (r *Router) SSLCheck() {
// 	if err := r.DecodeFields(); err != nil {
// 		slog.Error("Failed to decode router", "error", err)
// 		return
// 	}

// 	if r.EntryPoints == nil || r.Tls == nil {
// 		r.UpdateError("ssl", "")
// 		return
// 	}

// 	isHTTPS := false
// 	for _, ep := range r.EntryPoints.([]interface{}) {
// 		entrypoint, err := Query.GetEntryPointByName(
// 			context.Background(),
// 			GetEntryPointByNameParams{
// 				ProfileID: r.ProfileID,
// 				Name:      ep.(string),
// 			},
// 		)
// 		if err != nil {
// 			slog.Error("Failed to get entry point", "name", ep.(string), "error", err)
// 			continue
// 		}
// 		if entrypoint.Address == "443" {
// 			isHTTPS = true
// 			break
// 		}
// 	}
// 	if !isHTTPS {
// 		r.UpdateError("ssl", "")
// 		return
// 	}

// 	if tlsMap, ok := r.Tls.(map[string]interface{}); ok {
// 		if tlsMap["certResolver"] == "" {
// 			slog.Debug("Router is not using a certificate resolver", "name", r.Name)
// 			r.UpdateError("ssl", "")
// 			return
// 		}
// 	} else {
// 		slog.Error("Unexpected type for TLS config", "type", fmt.Sprintf("%T", r.Tls))
// 		return
// 	}

// 	domains, err := util.ExtractDomainFromRule(r.Rule)
// 	if err != nil {
// 		slog.Error("Failed to extract domains from rule", "error", err)
// 		return
// 	}

// 	// Perform SSL validation and update only the "ssl" error field
// 	for _, domain := range domains {
// 		if err := util.ValidSSLCert(domain); err != nil {
// 			r.UpdateError("ssl", err.Error())
// 		} else {
// 			r.UpdateError("ssl", "")
// 		}
// 	}
// }
