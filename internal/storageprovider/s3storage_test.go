package storageprovider

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestS3Suite(t *testing.T) {
	godotenv.Load("../../test.env")

	RunSuite(t, func() StorageProvider {
		s3Folder := "botman-tests"
		store := NewS3Store(S3Config{
			AccessKey:      os.Getenv("TEST_S3_ACCESS_KEY"),
			SecretKey:      os.Getenv("TEST_S3_SECRET_KEY"),
			Endpoint:       os.Getenv("TEST_S3_ENDPOINT"),
			Region:         os.Getenv("TEST_S3_REGION"),
			Bucket:         os.Getenv("TEST_S3_BUCKET"),
			BasePath:       s3Folder,
			ForcePathStyle: true,
		})

		store.DeleteFolder("")

		return store
	})
}
