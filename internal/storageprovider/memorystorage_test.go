package storageprovider

import "testing"

func TestMemSuite(t *testing.T) {
	RunSuite(t, func() StorageProvider { return NewMemStore() })
}
