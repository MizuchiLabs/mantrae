package config

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/MizuchiLabs/mantrae/internal/db"
)

// Settings defines all application settings
type Settings struct {
	ServerURL            string        `setting:"server_url" default:"http://localhost:3000"`
	BackupEnabled        bool          `setting:"backup_enabled" default:"true"`
	BackupInterval       time.Duration `setting:"backup_interval" default:"24h"`
	BackupKeep           int           `setting:"backup_keep" default:"3"`
	EmailHost            string        `setting:"email_host" default:"localhost"`
	EmailPort            int           `setting:"email_port" default:"587"`
	EmailUser            string        `setting:"email_user" default:""`
	EmailPass            string        `setting:"email_pass" default:""`
	EmailFrom            string        `setting:"email_from" default:"mantrae@localhost"`
	AgentCleanupEnabled  bool          `setting:"agent_cleanup_enabled" default:"true"`
	AgentCleanupInterval time.Duration `setting:"agent_cleanup_interval" default:"24h"`
}

type SettingsManager struct {
	q        *db.Queries
	defaults *Settings
}

func NewSettingsManager(q *db.Queries) *SettingsManager {
	return &SettingsManager{
		q:        q,
		defaults: getDefaults(),
	}
}

// parseValue converts a string value to the appropriate type based on the field
func parseValue(field reflect.Value, strValue string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(strValue)
	case reflect.Bool:
		b, err := strconv.ParseBool(strValue)
		if err != nil {
			return err
		}
		field.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Type() == reflect.TypeOf(time.Duration(0)) {
			d, err := time.ParseDuration(strValue)
			if err != nil {
				return err
			}
			field.SetInt(int64(d))
		} else {
			i, err := strconv.ParseInt(strValue, 10, 64)
			if err != nil {
				return err
			}
			field.SetInt(i)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(strValue, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(i)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(strValue, 64)
		if err != nil {
			return err
		}
		field.SetFloat(f)
	case reflect.Ptr:
		// Handle pointer types
		if strValue == "" {
			field.Set(reflect.Zero(field.Type()))
			return nil
		}

		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		return parseValue(field.Elem(), strValue)
	}
	return nil
}

// valueToString converts a value to its string representation
func valueToString(value reflect.Value) string {
	switch value.Kind() {
	case reflect.String:
		return value.String()
	case reflect.Bool:
		return strconv.FormatBool(value.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value.Type() == reflect.TypeOf(time.Duration(0)) {
			return value.Interface().(time.Duration).String()
		}
		return strconv.FormatInt(value.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'f', -1, 64)
	case reflect.Ptr:
		if value.IsNil() {
			return ""
		}
		return valueToString(value.Elem())
	default:
		return fmt.Sprint(value.Interface())
	}
}

// getDefaults creates a Settings instance with default values from struct tags
func getDefaults() *Settings {
	s := &Settings{}
	v := reflect.ValueOf(s).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		defaultVal := field.Tag.Get("default")
		if defaultVal == "" {
			continue
		}

		if err := parseValue(v.Field(i), defaultVal); err != nil {
			// Log error or handle it as appropriate for your application
			continue
		}
	}
	return s
}

// Initialize ensures all settings exist with default values
func (sm *SettingsManager) Initialize(ctx context.Context) error {
	// First get existing settings
	existingSettings, err := sm.q.ListSettings(ctx)
	if err != nil {
		return err
	}

	// Create map of existing settings for quick lookup
	existing := make(map[string]struct{})
	for _, s := range existingSettings {
		existing[s.Key] = struct{}{}
	}

	// Only initialize settings that don't exist
	v := reflect.ValueOf(sm.defaults).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		key := field.Tag.Get("setting")
		if key == "" {
			continue
		}

		// Skip if setting already exists
		if _, exists := existing[key]; exists {
			continue
		}

		// Only set default value if setting doesn't exist
		value := valueToString(v.Field(i))
		if err := sm.q.UpsertSetting(ctx, db.UpsertSettingParams{
			Key:   key,
			Value: value,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (sm *SettingsManager) GetAll(ctx context.Context) (*Settings, error) {
	settings := &Settings{} // Start with empty settings
	v := reflect.ValueOf(settings).Elem()
	t := v.Type()

	// Create map of field info
	fieldMap := make(map[string]int)
	for i := 0; i < t.NumField(); i++ {
		if key := t.Field(i).Tag.Get("setting"); key != "" {
			fieldMap[key] = i
		}
	}

	// Get all settings from database
	dbSettings, err := sm.q.ListSettings(ctx)
	if err != nil {
		return nil, err
	}

	// Create map of database settings for easier lookup
	dbMap := make(map[string]string)
	for _, setting := range dbSettings {
		dbMap[setting.Key] = setting.Value
	}

	// Iterate through fields and set values
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		key := field.Tag.Get("setting")
		if key == "" {
			continue
		}

		// If value exists in database, use it
		if value, exists := dbMap[key]; exists {
			if err := parseValue(v.Field(i), value); err != nil {
				return nil, fmt.Errorf("error parsing setting %s: %w", key, err)
			}
		} else {
			// If no value in database, use default
			defaultVal := field.Tag.Get("default")
			if defaultVal != "" {
				if err := parseValue(v.Field(i), defaultVal); err != nil {
					return nil, fmt.Errorf("error parsing default value for %s: %w", key, err)
				}
			}
		}
	}

	return settings, nil
}

func (sm *SettingsManager) Get(ctx context.Context, key string) (interface{}, error) {
	setting, err := sm.q.GetSetting(ctx, key)
	if err != nil {
		return nil, err
	}

	// Find the corresponding field and its type in Settings struct
	v := reflect.ValueOf(sm.defaults).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Tag.Get("setting") == key {
			// Create a new value of the correct type
			fieldValue := reflect.New(t.Field(i).Type).Elem()

			// Parse the string value to the correct type
			if err := parseValue(fieldValue, setting.Value); err != nil {
				return nil, fmt.Errorf("error parsing setting %s: %w", key, err)
			}

			// Return the interface{} value
			return fieldValue.Interface(), nil
		}
	}

	return nil, fmt.Errorf("unknown setting key: %s", key)
}

// Set updates a setting with proper type conversion from string input
func (sm *SettingsManager) Set(ctx context.Context, key string, strValue string) error {
	// Find the corresponding field to validate and convert the type
	v := reflect.ValueOf(sm.defaults).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Tag.Get("setting") == key {
			// Create a new value of the correct field type
			fieldValue := reflect.New(t.Field(i).Type).Elem()

			// Parse the string value to the correct type
			if err := parseValue(fieldValue, strValue); err != nil {
				return fmt.Errorf("invalid value for setting %s: %w", key, err)
			}

			return sm.q.UpsertSetting(ctx, db.UpsertSettingParams{
				Key:   key,
				Value: strValue,
			})
		}
	}

	return fmt.Errorf("unknown setting key: %s", key)
}
