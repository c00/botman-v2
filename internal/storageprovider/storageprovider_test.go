package storageprovider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func RunSuite(t *testing.T, pf func() StorageProvider) {
	storeAndRetrieveFile(t, pf())
	loadNonExistingFile(t, pf())
	filenameErrors(t, pf())
}

func storeAndRetrieveFile(t *testing.T, p StorageProvider) {
	filename := "/foo.txt"
	content := "some content for automated tests"

	path, err := p.Save(filename, []byte(content))
	assert.Nil(t, err)
	assert.NotEqual(t, path, "")

	got, err := p.Load(filename)
	assert.Nil(t, err)
	assert.Equal(t, string(got), content)
}

func loadNonExistingFile(t *testing.T, p StorageProvider) {
	filename := "/foo.txt"

	got, err := p.Load(filename)
	assert.NotNil(t, err)
	assert.Len(t, got, 0)
}

func filenameErrors(t *testing.T, p StorageProvider) {
	_, err := p.Save("", []byte{})
	assert.ErrorIs(t, err, ErrEmptyFilename)

	_, err = p.Load("")
	assert.ErrorIs(t, err, ErrEmptyFilename)
}
