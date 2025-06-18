package settings

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/mizuchilabs/mantrae/internal/store/db"
)

func (sm *SettingsManager) validate(ctx context.Context, params *db.UpsertSettingParams) error {
	q := sm.conn.GetQuery()

	// Trim whitespace
	params.Value = strings.TrimSpace(params.Value)

	switch params.Key {
	case KeyServerURL:
		params.Value = cleanURL(params.Value)
		if params.Value == "" {
			return errors.New("server url cannot be empty")
		}

	case KeyS3Endpoint:
		params.Value = cleanURL(params.Value)
		if params.Value == "" {
			return errors.New("S3 endpoint cannot be empty")
		}

	case KeyOIDCIssuerURL:
		params.Value = cleanURL(params.Value)
		if params.Value == "" {
			return errors.New("OIDC issuer URL cannot be empty")
		}

	case KeyEmailPort:
		port, err := strconv.Atoi(params.Value)
		if err != nil || port < 1 || port > 65535 {
			return errors.New("email port must be an integer between 1 and 65535")
		}

	case KeyBackupKeep:
		i, err := strconv.Atoi(params.Value)
		if err != nil || i < 1 {
			return errors.New("backup keep must be an integer greater than 0")
		}

	case KeyPasswordLoginDisabled:
		// Don't allow disabling password login unless OIDC is enabled
		enabled, ok := sm.Get(KeyOIDCEnabled)
		if !ok {
			return errors.New("failed to get OIDC setting")
		}
		if params.Value == "true" && !AsBool(enabled) {
			return errors.New("password login cannot be disabled unless OIDC is enabled")
		}

	case KeyOIDCEnabled:
		// If Password Login is disabled ensure to enable it again if oidc gets disabled
		pwLogin, ok := sm.Get(KeyPasswordLoginDisabled)
		if !ok {
			return errors.New("failed to get password login setting")
		}
		if params.Value == "false" && AsBool(pwLogin) {
			return q.UpsertSetting(ctx, db.UpsertSettingParams{
				Key:   KeyPasswordLoginDisabled,
				Value: "false",
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
