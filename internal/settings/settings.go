// Package settings provides functionality for managing application settings.
package settings

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/mizuchilabs/mantrae/internal/store"
	"github.com/mizuchilabs/mantrae/internal/store/db"
)

// Settings defines all application settings
type Settings struct {
	ServerURL            string        `setting:"server_url"             default:""`
	Storage              string        `setting:"storage_select"         default:"local"`
	BackupEnabled        bool          `setting:"backup_enabled"         default:"true"`
	BackupInterval       time.Duration `setting:"backup_interval"        default:"24h"`
	BackupKeep           int           `setting:"backup_keep"            default:"3"`
	S3Endpoint           string        `setting:"s3_endpoint"            default:""`
	S3Bucket             string        `setting:"s3_bucket"              default:"mantrae"`
	S3Region             string        `setting:"s3_region"              default:"us-east-1"`
	S3AccessKey          string        `setting:"s3_access_key"          default:""`
	S3SecretKey          string        `setting:"s3_secret_key"          default:""`
	S3UsePathStyle       bool          `setting:"s3_use_path_style"      default:"false"`
	EmailHost            string        `setting:"email_host"             default:""`
	EmailPort            int           `setting:"email_port"             default:"587"`
	EmailUser            string        `setting:"email_user"             default:""`
	EmailPassword        string        `setting:"email_password"         default:""`
	EmailFrom            string        `setting:"email_from"             default:"mantrae@localhost"`
	PasswordLoginEnabled bool          `setting:"password_login_enabled" default:"true"`
	OIDCEnabled          bool          `setting:"oidc_enabled"           default:"false"`
	OIDCClientID         string        `setting:"oidc_client_id"         default:""`
	OIDCClientSecret     string        `setting:"oidc_client_secret"     default:""`
	OIDCIssuerURL        string        `setting:"oidc_issuer_url"        default:""`
	OIDCProviderName     string        `setting:"oidc_provider_name"     default:""`
	OIDCScopes           string        `setting:"oidc_scopes"            default:""`
	OIDCPKCE             bool          `setting:"oidc_pkce"              default:"false"`
	AgentCleanupEnabled  bool          `setting:"agent_cleanup_enabled"  default:"true"`
	AgentCleanupInterval time.Duration `setting:"agent_cleanup_interval" default:"24h"`
	TraefikSyncInterval  time.Duration `setting:"traefik_sync_interval"  default:"20s"`
	DNSSyncInterval      time.Duration `setting:"dns_sync_interval"      default:"3m"`
	AgentCheckInterval   time.Duration `setting:"agent_check_interval"   default:"5m"`
}

type SettingsManager struct {
	conn    *store.Connection
	structT reflect.Type
}

func NewManager(conn *store.Connection) *SettingsManager {
	if conn == nil {
		log.Fatal("No database connection provided")
	}
	return &SettingsManager{
		conn:    conn,
		structT: reflect.TypeOf(Settings{}),
	}
}

// Start loads settings from ENV > DB > default tags
func (sm *SettingsManager) Start(ctx context.Context) {
	q := sm.conn.Query
	existing, err := q.ListSettings(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Build current DB map
	dbMap := map[string]string{}
	for _, s := range existing {
		dbMap[s.Key] = s.Value
	}

	for _, field := range reflect.VisibleFields(sm.structT) {
		key := field.Tag.Get("setting")
		def := field.Tag.Get("default")

		var val string

		if envVal, ok := os.LookupEnv(strings.ToUpper(key)); ok {
			val = envVal
		} else if dbVal, ok := dbMap[key]; ok {
			val = dbVal
		} else {
			val = def
		}

		if _, exists := dbMap[key]; !exists {
			err = q.UpsertSetting(ctx, &db.UpsertSettingParams{Key: key, Value: val})
			if err != nil {
				log.Fatal(fmt.Errorf("failed to upsert setting %s: %w", key, err))
			}
		}
	}

	// Clean up deprecated DB settings
	validKeys := sm.validKeys()
	for k := range dbMap {
		if _, ok := validKeys[k]; !ok {
			if err = q.DeleteSetting(ctx, k); err != nil {
				slog.Error("failed to delete deprecated setting", "key", k, "error", err)
			}
		}
	}
}

func (sm *SettingsManager) validKeys() map[string]struct{} {
	keys := make(map[string]struct{})
	for _, field := range reflect.VisibleFields(sm.structT) {
		if k := field.Tag.Get("setting"); k != "" {
			keys[k] = struct{}{}
		}
	}
	return keys
}

func (sm *SettingsManager) Get(ctx context.Context, key string) (string, bool) {
	setting, err := sm.conn.Query.GetSetting(ctx, key)
	if err != nil {
		return "", false
	}
	return setting.Value, true
}

func (sm *SettingsManager) Set(ctx context.Context, key, val string) error {
	if _, ok := sm.validKeys()[key]; !ok {
		return fmt.Errorf("invalid setting key: %s", key)
	}

	params := &db.UpsertSettingParams{Key: key, Value: val}
	if err := sm.validate(ctx, params); err != nil {
		return err
	}

	// Update database
	if err := sm.conn.Query.UpsertSetting(ctx, params); err != nil {
		return fmt.Errorf("failed to update setting in database: %w", err)
	}

	return nil
}

func (sm *SettingsManager) GetAll(ctx context.Context) map[string]string {
	settings, err := sm.conn.Query.ListSettings(ctx)
	if err != nil {
		return make(map[string]string)
	}

	result := make(map[string]string, len(settings))
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	return result
}

func (sm *SettingsManager) GetMany(ctx context.Context, keys []string) map[string]string {
	result := make(map[string]string, len(keys))
	for _, k := range keys {
		if val, ok := sm.Get(ctx, k); ok {
			result[k] = val
		}
	}
	return result
}

// Helper

func AsString(val *string) string {
	if val == nil || *val == "" {
		return ""
	}
	return *val
}

func AsBool(val string) bool {
	if v, err := strconv.ParseBool(val); err == nil {
		return v
	}
	return false
}

func AsInt(val string) int {
	if v, err := strconv.Atoi(val); err == nil {
		return v
	}
	return 0
}

func AsDuration(val string) time.Duration {
	if d, err := time.ParseDuration(val); err == nil {
		return d
	}
	return 0
}
