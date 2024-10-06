package clitools

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInput(t *testing.T) {
	inputVal := "some input"
	reader := bytes.NewReader([]byte(inputVal))
	writer := &stringWriter{}

	value := GetInput("some label", reader, writer)
	assert.Equal(t, inputVal, value)
	assert.Equal(t, "some label: \n", writer.String())
}

func TestSetInput(t *testing.T) {
	value := ""
	reader := bytes.NewReader([]byte("some input\n"))
	writer := &stringWriter{}

	SetInput("some label", &value, reader, writer)
	assert.Equal(t, "some input", value)
	assert.Equal(t, "some label: \n", writer.String())
}

func TestSetInput2(t *testing.T) {
	value := "foo"
	reader := bytes.NewReader([]byte("some input\n"))
	writer := &stringWriter{}

	SetInput("some label", &value, reader, writer)
	assert.Equal(t, "some input", value)
	assert.Equal(t, "some label [foo]: \n", writer.String())
}

func TestSetInput3(t *testing.T) {
	value := "foo"
	reader := bytes.NewReader([]byte("\n"))
	writer := &stringWriter{}

	SetInput("some label", &value, reader, writer)
	assert.Equal(t, "foo", value)
	assert.Equal(t, "some label [foo]: ", writer.String())
}

type stringWriter struct {
	io.Writer
	data []byte
}

func (sw *stringWriter) Write(newData []byte) (n int, err error) {
	sw.data = append(sw.data, newData...)
	return len(newData), nil
}

func (sw *stringWriter) String() string {
	return string(sw.data)
}

func (sw *stringWriter) StringFlush() string {
	data := string(sw.data)
	sw.data = []byte{}
	return data
}
