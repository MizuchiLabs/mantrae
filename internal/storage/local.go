// Package storage provides a generic interface for storage backends.
package storage

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sort"

	"github.com/mizuchilabs/mantrae/internal/util"
)

// LocalStorage implements StorageBackend for local filesystem
type LocalStorage struct {
	basePath string
}

func NewLocalStorage(path string) (*LocalStorage, error) {
	resolvedPath := util.ResolvePath(path)
	if err := os.MkdirAll(resolvedPath, 0o750); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}
	return &LocalStorage{basePath: resolvedPath}, nil
}

func (ls *LocalStorage) Store(ctx context.Context, name string, data io.Reader) error {
	path := filepath.Join(ls.basePath, name)
	f, err := os.Create(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			slog.Error("failed to close backup file", "error", err)
		}
	}()

	if _, err := io.Copy(f, data); err != nil {
		return fmt.Errorf("failed to write backup data: %w", err)
	}
	return nil
}

func (ls *LocalStorage) Retrieve(ctx context.Context, name string) (io.ReadCloser, error) {
	path := filepath.Join(ls.basePath, name)
	return os.Open(filepath.Clean(path))
}

func (ls *LocalStorage) List(ctx context.Context) ([]StoredFile, error) {
	var files []StoredFile
	entries, err := os.ReadDir(ls.basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read storage directory: %w", err)
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		ts := info.ModTime()
		files = append(files, StoredFile{
			Name:      entry.Name(),
			Size:      info.Size(),
			Timestamp: &ts,
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Timestamp.After(*files[j].Timestamp)
	})

	return files, nil
}

func (ls *LocalStorage) Delete(ctx context.Context, name string) error {
	path := filepath.Join(ls.basePath, name)
	return os.Remove(path)
}
