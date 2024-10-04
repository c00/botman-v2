package storageprovider

import "errors"

const (
	StorageTypeLocal  = "local"
	StorageTypeMemory = "memory"
	StorageTypeS3     = "s3"
)

var ErrEmptyFilename = errors.New("empty filename")

type StorageProvider interface {
	Save(name string, data []byte) (string, error)
	Load(name string) ([]byte, error)
}
