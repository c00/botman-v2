package storageprovider

import (
	"fmt"
	"os"
	"path/filepath"
)

var _ StorageProvider = (*LocalStorage)(nil)

func NewLocalStore(path string) (*LocalStorage, error) {
	os.Mkdir(path, 0755)
	_, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("cannot create local store: %v: %w", path, err)
	}

	return &LocalStorage{
		Path: path,
	}, nil
}

type LocalStorage struct {
	Path string
}

func (s *LocalStorage) getFullPath(name string) string {
	return filepath.Join(s.Path, name)
}

func (s *LocalStorage) Save(name string, data []byte) (string, error) {
	if name == "" {
		return "", ErrEmptyFilename
	}

	path := s.getFullPath(name)
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		return "", fmt.Errorf("could not write to file: %w", err)
	}
	return path, nil
}

func (s *LocalStorage) Load(name string) ([]byte, error) {
	if name == "" {
		return nil, ErrEmptyFilename
	}

	path := s.getFullPath(name)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read from file: %w", err)
	}

	return data, nil
}
