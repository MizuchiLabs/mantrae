package storage

import (
	"context"
	"io"
	"time"
)

type BackendType string

const (
	BackendTypeLocal BackendType = "local"
	BackendTypeS3    BackendType = "s3"
	// BackendTypeGit   BackendType = "git" //TODO: For future implementation
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

func (t BackendType) Valid() bool {
	switch t {
	case BackendTypeLocal, BackendTypeS3:
		return true
	default:
		return false
	}
}
