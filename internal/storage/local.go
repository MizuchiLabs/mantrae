package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/MizuchiLabs/mantrae/internal/app"
)

// LocalStorage implements StorageBackend for local filesystem
type LocalStorage struct {
	basePath string
}

func NewLocalStorage(path string) (*LocalStorage, error) {
	resolvedPath := app.ResolvePath(path)
	if err := os.MkdirAll(resolvedPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}
	return &LocalStorage{basePath: path}, nil
}

func (ls *LocalStorage) Store(ctx context.Context, name string, data io.Reader) error {
	path := filepath.Join(ls.basePath, name)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, data); err != nil {
		return fmt.Errorf("failed to write backup data: %w", err)
	}
	return nil
}

func (ls *LocalStorage) Retrieve(ctx context.Context, name string) (io.ReadCloser, error) {
	path := filepath.Join(ls.basePath, name)
	return os.Open(path)
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
		files = append(files, StoredFile{
			Name:      entry.Name(),
			Timestamp: info.ModTime(),
			Size:      info.Size(),
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Timestamp.After(files[j].Timestamp)
	})

	return files, nil
}

func (ls *LocalStorage) Delete(ctx context.Context, name string) error {
	path := filepath.Join(ls.basePath, name)
	return os.Remove(path)
}
