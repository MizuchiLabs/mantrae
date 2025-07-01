package settings

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/mizuchilabs/mantrae/internal/store/db"
	"github.com/mizuchilabs/mantrae/pkg/util"
)

func (sm *SettingsManager) validate(ctx context.Context, params *db.UpsertSettingParams) error {
	q := sm.conn.GetQuery()

	// Trim whitespace
	params.Value = strings.TrimSpace(params.Value)

	switch params.Key {
	case KeyServerURL:
		if params.Value == "" {
			return errors.New("server url cannot be empty")
		}
		params.Value = util.CleanURL(params.Value)

	case KeyS3Endpoint:
		if params.Value == "" {
			return errors.New("S3 endpoint cannot be empty")
		}
		params.Value = util.CleanURL(params.Value)

	case KeyOIDCIssuerURL:
		if params.Value == "" {
			return errors.New("OIDC issuer URL cannot be empty")
		}
		params.Value = util.CleanURL(params.Value)

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

	case KeyPasswordLoginEnabled:
		// Don't allow disabling password login unless OIDC is oidcEnabled
		oidcEnabled, ok := sm.Get(ctx, KeyOIDCEnabled)
		if !ok {
			return errors.New("failed to get OIDC setting")
		}
		if params.Value == "false" && !AsBool(oidcEnabled) {
			return errors.New("password login cannot be disabled unless OIDC is enabled")
		}

	case KeyOIDCEnabled:
		// If Password Login is disabled ensure to enable it again if oidc gets disabled
		pwLogin, ok := sm.Get(ctx, KeyPasswordLoginEnabled)
		if !ok {
			return errors.New("failed to get password login setting")
		}
		if params.Value == "false" && !AsBool(pwLogin) {
			return q.UpsertSetting(ctx, db.UpsertSettingParams{
				Key:   KeyPasswordLoginEnabled,
				Value: "true",
			})
		}
	}

	return nil
}
