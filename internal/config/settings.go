package config

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/MizuchiLabs/mantrae/internal/db"
)

type SettingWithDescription struct {
	Value       any     `json:"value"`
	Description *string `json:"description,omitempty"`
}

// Settings defines all application settings
type Settings struct {
	ServerURL            string        `setting:"server_url"             default:"http://localhost:3000" description:"Base URL for the server"`
	BackupEnabled        bool          `setting:"backup_enabled"         default:"true"                  description:"Enable automatic backups"`
	BackupInterval       time.Duration `setting:"backup_interval"        default:"24h"                   description:"Interval between backups"`
	BackupKeep           int           `setting:"backup_keep"            default:"3"                     description:"Number of backups to retain"`
	EmailHost            string        `setting:"email_host"             default:"localhost"             description:"SMTP server hostname"`
	EmailPort            int           `setting:"email_port"             default:"587"                   description:"SMTP server port"`
	EmailUser            string        `setting:"email_user"             default:""                      description:"SMTP username"`
	EmailPassword        string        `setting:"email_password"         default:""                      description:"SMTP password"`
	EmailFrom            string        `setting:"email_from"             default:"mantrae@localhost"     description:"From email address"`
	AgentCleanupEnabled  bool          `setting:"agent_cleanup_enabled"  default:"true"                  description:"Enable automatic agent cleanup"`
	AgentCleanupInterval time.Duration `setting:"agent_cleanup_interval" default:"24h"                   description:"Interval for agent cleanup"`
}

type SettingsManager struct {
	conn     *db.Connection
	defaults *Settings
}

func NewSettingsManager(conn *db.Connection) *SettingsManager {
	return &SettingsManager{
		conn:     conn,
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
	q := sm.conn.GetQuery()
	existingSettings, err := q.ListSettings(ctx)
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
		if err := q.UpsertSetting(ctx, db.UpsertSettingParams{
			Key:   key,
			Value: value,
		}); err != nil {
			return err
		}
	}
	return nil
}

// Modified GetAll to return settings with descriptions
func (sm *SettingsManager) GetAll(ctx context.Context) (map[string]SettingWithDescription, error) {
	q := sm.conn.GetQuery()
	settings := make(map[string]SettingWithDescription)
	v := reflect.ValueOf(sm.defaults).Elem()
	t := v.Type()

	// Get all settings from database
	dbSettings, err := q.ListSettings(ctx)
	if err != nil {
		return nil, err
	}

	// Create map of database settings
	dbMap := make(map[string]db.Setting)
	for _, setting := range dbSettings {
		dbMap[setting.Key] = setting
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		key := field.Tag.Get("setting")
		if key == "" {
			continue
		}

		fieldValue := reflect.New(field.Type).Elem()
		description := field.Tag.Get("description")

		// If value exists in database, use it
		if dbSetting, exists := dbMap[key]; exists {
			if err := parseValue(fieldValue, dbSetting.Value); err != nil {
				return nil, fmt.Errorf("error parsing setting %s: %w", key, err)
			}
			if dbSetting.Description != nil {
				description = *dbSetting.Description
			}
		} else {
			// Use default value
			defaultVal := field.Tag.Get("default")
			if defaultVal != "" {
				if err := parseValue(fieldValue, defaultVal); err != nil {
					return nil, fmt.Errorf("error parsing default value for %s: %w", key, err)
				}
			}
		}

		desc := &description
		if description == "" {
			desc = nil
		}

		// Convert duration fields to formatted strings
		var value any
		if field.Type == reflect.TypeOf(time.Duration(0)) {
			value = fieldValue.Interface().(time.Duration).String()
		} else {
			value = fieldValue.Interface()
		}

		settings[key] = SettingWithDescription{
			Value:       value,
			Description: desc,
		}
	}

	return settings, nil
}

// Modified Get to return setting with description
func (sm *SettingsManager) Get(ctx context.Context, key string) (*SettingWithDescription, error) {
	q := sm.conn.GetQuery()
	setting, err := q.GetSetting(ctx, key)
	if err != nil {
		return nil, err
	}

	v := reflect.ValueOf(sm.defaults).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Tag.Get("setting") == key {
			fieldValue := reflect.New(t.Field(i).Type).Elem()
			if err := parseValue(fieldValue, setting.Value); err != nil {
				return nil, fmt.Errorf("error parsing setting %s: %w", key, err)
			}

			// Convert duration fields to formatted strings
			var value any
			if t.Field(i).Type == reflect.TypeOf(time.Duration(0)) {
				value = fieldValue.Interface().(time.Duration).String()
			} else {
				value = fieldValue.Interface()
			}

			return &SettingWithDescription{
				Value:       value,
				Description: setting.Description,
			}, nil
		}
	}

	return nil, fmt.Errorf("unknown setting key: %s", key)
}

// Set updates a setting with proper type conversion from string input
func (sm *SettingsManager) Set(
	ctx context.Context,
	key, strValue string,
	description *string,
) error {
	q := sm.conn.GetQuery()
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
			params := db.UpsertSettingParams{
				Key:   key,
				Value: strValue,
			}
			if description != nil {
				params.Description = description
			}
			return q.UpsertSetting(ctx, params)
		}
	}

	return fmt.Errorf("unknown setting key: %s", key)
}
