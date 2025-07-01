package settings

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/mizuchilabs/mantrae/internal/store"
	"github.com/mizuchilabs/mantrae/internal/store/db"
	"github.com/stretchr/testify/assert"
)

func setupTest() (*SettingsManager, func()) {
	conn := store.NewConnection(":memory:")
	sm := NewManager(conn)

	return sm, func() { conn.Close() }
}

func TestNewManager(t *testing.T) {
	conn := store.NewConnection(":memory:")
	defer conn.Close()

	sm := NewManager(conn)
	assert.NotNil(t, sm)
	assert.NotNil(t, sm.conn)
}

func TestGetAndSet(t *testing.T) {
	sm, teardown := setupTest()
	defer teardown()

	ctx := context.Background()

	// Test setting and getting a value
	err := sm.Set(ctx, KeyServerURL, "http://localhost:8080")
	assert.NoError(t, err)

	val, ok := sm.Get(ctx, KeyServerURL)
	assert.True(t, ok)
	assert.Equal(t, "http://localhost:8080", val)

	// Test setting an invalid key
	err = sm.Set(ctx, "invalid_key", "some_value")
	assert.Error(t, err)
}

func TestGetAll(t *testing.T) {
	sm, teardown := setupTest()
	defer teardown()

	ctx := context.Background()
	sm.Start(ctx)

	// Test getting all values
	allSettings := sm.GetAll(ctx)
	assert.NotEmpty(t, allSettings)
	assert.Equal(t, "local", allSettings[KeyStorage])
}

func TestGetMany(t *testing.T) {
	sm, teardown := setupTest()
	defer teardown()

	ctx := context.Background()
	sm.Start(ctx)

	// Test getting many values
	keys := []string{KeyServerURL, KeyStorage}
	manySettings := sm.GetMany(ctx, keys)
	assert.Len(t, manySettings, 2)
	assert.Equal(t, "", manySettings[KeyServerURL])
	assert.Equal(t, "local", manySettings[KeyStorage])
}

func TestStart(t *testing.T) {
	sm, teardown := setupTest()
	defer teardown()

	ctx := context.Background()

	// Set an environment variable
	os.Setenv("SERVER_URL", "http://env.test")
	defer os.Unsetenv("SERVER_URL")

	// Add a value to the database
	err := sm.conn.GetQuery().UpsertSetting(ctx, db.UpsertSettingParams{
		Key:   KeyStorage,
		Value: "db_value",
	})
	assert.NoError(t, err)

	sm.Start(ctx)

	// Check that the environment variable is used
	val, ok := sm.Get(ctx, KeyServerURL)
	assert.True(t, ok)
	assert.Equal(t, "http://env.test", val)

	// Check that the database value is used
	val, ok = sm.Get(ctx, KeyStorage)
	assert.True(t, ok)
	assert.Equal(t, "db_value", val)

	// Check that the default value is used
	val, ok = sm.Get(ctx, KeyBackupEnabled)
	assert.True(t, ok)
	assert.Equal(t, "true", val)
}

func TestValidation(t *testing.T) {
	sm, teardown := setupTest()
	defer teardown()

	ctx := context.Background()

	// Test validation for server url
	err := sm.Set(ctx, KeyServerURL, "  ")
	assert.Error(t, err)

	// Test validation for email port
	err = sm.Set(ctx, KeyEmailPort, "abc")
	assert.Error(t, err)

	err = sm.Set(ctx, KeyEmailPort, "70000")
	assert.Error(t, err)

	// Test validation for backup keep
	err = sm.Set(ctx, KeyBackupKeep, "0")
	assert.Error(t, err)
}

func TestAsHelpers(t *testing.T) {
	testString := "test"
	testStringEmpty := ""
	assert.Equal(t, "test", AsString(&testString))
	assert.Equal(t, "", AsString(&testStringEmpty))
	assert.True(t, AsBool("true"))
	assert.False(t, AsBool("false"))
	assert.False(t, AsBool("invalid"))
	assert.Equal(t, 123, AsInt("123"))
	assert.Equal(t, 0, AsInt("invalid"))
	assert.Equal(t, 123.45, AsFloat64("123.45"))
	assert.Equal(t, 0.0, AsFloat64("invalid"))
	assert.Equal(t, time.Hour, AsDuration("1h"))
	assert.Equal(t, time.Duration(0), AsDuration("invalid"))
}
