package streamrecorder

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStreamRecorder(t *testing.T) {
	writer := &simpleWriter{}
	text := "Are you the muffin man?\nOr are you the garbage man?!"

	r := &StreamRecorder{
		reader: bufio.NewReader(bytes.NewReader([]byte(text))),
		writer: writer,
	}

	data, err := r.ReadBytes('\n')
	assert.Nil(t, err)
	assert.Equal(t, "Are you the muffin man?\n", string(data))

	data, err = r.ReadBytes('\n')
	// assert.Nil(t, err)
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, io.EOF))
	assert.Equal(t, "Or are you the garbage man?!", string(data))

	assert.Equal(t, writer.String(), text)
}

func TestRecordToFile(t *testing.T) {
	text := "Are you the muffin man?\nOr are you the garbage man?!"
	filename := filepath.Join(os.TempDir(), "testrecording.txt")
	r, err := NewFileRecorder(
		filename,
		bufio.NewReader(bytes.NewReader([]byte(text))),
	)
	assert.Nil(t, err)

	for {
		_, err := r.ReadBytes('\n')
		if err != nil {
			break
		}
	}

	f, err := os.ReadFile(filename)
	assert.Nil(t, err)
	assert.Equal(t, string(f), text)

}

type simpleWriter struct {
	data []byte
}

func (w *simpleWriter) Write(p []byte) (n int, err error) {
	if w.data == nil {
		w.data = p
		return len(p), nil
	}

	w.data = append(w.data, p...)
	return len(p), nil
}

func (w *simpleWriter) String() string {
	return string(w.data)
}
