package storage

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/mizuchilabs/mantrae/internal/settings"
)

type BackendType string

const (
	BackendTypeLocal BackendType = "local"
	BackendTypeS3    BackendType = "s3"
)

// Backend defines interface for different storage solutions
type Backend interface {
	Store(ctx context.Context, name string, data io.Reader) error
	Retrieve(ctx context.Context, name string) (io.ReadCloser, error)
	List(ctx context.Context) ([]StoredFile, error)
	Delete(ctx context.Context, name string) error
}

type StoredFile struct {
	Name      string    `json:"name,omitempty"`
	Size      int64     `json:"size,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

func GetBackend(
	ctx context.Context,
	sm *settings.SettingsManager,
	path string,
) (Backend, error) {
	backendSetting, ok := sm.Get(settings.KeyStorage)
	if !ok {
		return nil, errors.New("failed to get storage backend")
	}

	switch BackendType(backendSetting) {
	case BackendTypeLocal:
		return NewLocalStorage(path)
	case BackendTypeS3:
		return NewS3Storage(ctx, sm)
	default:
		return nil, errors.New("invalid storage backend")
	}
}
