package settings

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/mizuchilabs/mantrae/internal/store/db"
)

func (sm *SettingsManager) validateSetting(
	ctx context.Context,
	params *db.UpsertSettingParams,
) error {
	q := sm.conn.GetQuery()

	// Trim whitespace
	params.Value = strings.TrimSpace(params.Value)

	switch params.Key {
	case KeyServerURL:
		params.Value = cleanURL(params.Value)
		if params.Value == "" {
			return fmt.Errorf("Server URL cannot be empty")
		}

	case KeyS3Endpoint:
		params.Value = cleanURL(params.Value)
		if params.Value == "" {
			return fmt.Errorf("OIDC issuer URL cannot be empty")
		}

	case KeyOIDCIssuerURL:
		params.Value = cleanURL(params.Value)
		if params.Value == "" {
			return fmt.Errorf("OIDC issuer URL cannot be empty")
		}

	case KeyEmailPort:
		port, err := strconv.Atoi(params.Value)
		if err != nil || port < 1 || port > 65535 {
			return fmt.Errorf("email port must be a valid TCP port (1-65535)")
		}

	case KeyBackupKeep:
		i, err := strconv.Atoi(params.Value)
		if err != nil || i < 1 {
			return fmt.Errorf("backup_keep must be an integer >= 1")
		}

	case KeyPasswordLoginDisabled:
		// Don't allow disabling password login unless OIDC is enabled
		enabled, err := sm.Get(ctx, KeyOIDCEnabled)
		if err != nil {
			return fmt.Errorf("failed to get OIDC setting: %w", err)
		}
		if params.Value == "true" && !enabled.Bool(false) {
			return fmt.Errorf("cannot disable password login when OIDC is not enabled")
		}

	case KeyOIDCEnabled:
		// If Password Login is disabled ensure to enable it again if oidc gets disabled
		pwLogin, err := sm.Get(ctx, KeyPasswordLoginDisabled)
		if err != nil {
			return fmt.Errorf("failed to get OIDC setting: %w", err)
		}
		if params.Value == "false" && pwLogin.Bool(false) {
			return q.UpsertSetting(ctx, db.UpsertSettingParams{
				Key:         KeyPasswordLoginDisabled,
				Value:       "false",
				Description: pwLogin.Description,
			})
		}
	}

	return nil
}

// Ensure trimmed and valid URL format
func cleanURL(url string) string {
	url = strings.TrimSuffix(url, "/")
	if !strings.HasPrefix(url, "http://") &&
		!strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}
	return url
}
