package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) (*LocalStorage, error) {
	abs, err := filepath.Abs(basePath)
	if err != nil {
		return nil, fmt.Errorf("invalid upload path: %w", err)
	}
	if err := os.MkdirAll(abs, 0755); err != nil {
		return nil, fmt.Errorf("create upload dir: %w", err)
	}
	return &LocalStorage{basePath: abs}, nil
}

func (s *LocalStorage) Save(filename string, reader io.Reader) (path string, size int64, err error) {
	ext := filepath.Ext(filename)
	id := uuid.New().String()
	rel := filepath.Join(id[:2], id[2:4], id+ext)
	full := filepath.Join(s.basePath, rel)

	if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
		return "", 0, fmt.Errorf("create subdir: %w", err)
	}

	f, err := os.Create(full)
	if err != nil {
		return "", 0, fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	written, err := io.Copy(f, reader)
	if err != nil {
		os.Remove(full)
		return "", 0, fmt.Errorf("write file: %w", err)
	}

	return rel, written, nil
}

func (s *LocalStorage) Open(path string) (io.ReadCloser, error) {
	full := filepath.Join(s.basePath, path)
	if !strings.HasPrefix(full, filepath.Clean(s.basePath)) {
		return nil, fmt.Errorf("invalid path")
	}
	return os.Open(full)
}

func (s *LocalStorage) Delete(path string) error {
	full := filepath.Join(s.basePath, path)
	if !strings.HasPrefix(full, filepath.Clean(s.basePath)) {
		return fmt.Errorf("invalid path")
	}
	return os.Remove(full)
}
