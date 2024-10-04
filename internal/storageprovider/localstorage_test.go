package storageprovider

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalSuite(t *testing.T) {
	RunSuite(t, func() StorageProvider {
		path := filepath.Join(os.TempDir(), "botman-test")
		os.RemoveAll(path)
		store, err := NewLocalStore(path)
		assert.Nil(t, err)
		return store
	})
}
