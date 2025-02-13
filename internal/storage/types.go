package storage

import (
	"context"
	"io"
	"time"
)

// StorageBackend defines interface for different storage solutions
type Backend interface {
	Store(ctx context.Context, name string, data io.Reader) error
	Retrieve(ctx context.Context, name string) (io.ReadCloser, error)
	List(ctx context.Context) ([]StoredFile, error)
	Delete(ctx context.Context, name string) error
}

type StoredFile struct {
	Name      string    `json:"name,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Size      int64     `json:"size,omitempty"`
}
