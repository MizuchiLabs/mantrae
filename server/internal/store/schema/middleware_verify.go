package schema

import (
	"fmt"
	"strings"

	"github.com/mizuchilabs/mantrae/pkg/util"
)

func (m HTTPMiddleware) Verify() error {
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
