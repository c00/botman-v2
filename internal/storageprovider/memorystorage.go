package storageprovider

import (
	"fmt"
	"path/filepath"

	"github.com/c00/botman-v2/internal/simplekeyvaluestore"
)

var _ StorageProvider = (*MemoryStorage)(nil)

func NewMemStore() *MemoryStorage {
	return &MemoryStorage{
		store: simplekeyvaluestore.NewSimpleKeyValueStore[[]byte](),
		path:  "memory::",
	}
}

type MemoryStorage struct {
	store simplekeyvaluestore.SimpleKeyValueStore[[]byte]
	path  string
}

func (s *MemoryStorage) getFullPath(name string) string {
	return filepath.Join(s.path, name)
}

func (s *MemoryStorage) Save(name string, data []byte) (string, error) {
	if name == "" {
		return "", ErrEmptyFilename
	}

	path := s.getFullPath(name)
	s.store.Set(path, data)
	return path, nil
}

func (s *MemoryStorage) Load(name string) ([]byte, error) {
	if name == "" {
		return nil, ErrEmptyFilename
	}

	path := s.getFullPath(name)
	data, err := s.store.Get(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read from memory store: %w", err)
	}
	return data, nil
}
