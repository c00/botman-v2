package streamrecorder

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

func NewFileRecorder(path string, source *bufio.Reader) (*StreamRecorder, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, fmt.Errorf("could not open file %v: %w", path, err)
	}

	return &StreamRecorder{
		reader: source,
		writer: file,
	}, nil
}

type StreamRecorder struct {
	reader *bufio.Reader
	writer io.Writer
}

func (sr *StreamRecorder) ReadBytes(delim byte) (n []byte, err error) {
	data, err := sr.reader.ReadBytes(delim)

	if len(data) > 0 {
		sr.writer.Write(data)
	}

	if errors.Is(err, io.EOF) {
		if closer, ok := sr.writer.(io.Closer); ok {
			//Close writer if it needs closing.
			closer.Close()
		}
	}

	return data, err
}

// func (sr *StreamRecorder) Read(p []byte) (n int, err error) {

// 	n, err = sr.reader.Read(p)

// 	if errors.Is(err, io.EOF) {
// 		if closer, ok := sr.writer.(io.Closer); ok {
// 			//Close writer if it needs closing.
// 			closer.Close()
// 		}
// 	}

// 	if err == nil {
// 		//Write to file
// 		sr.writer.Write(p)
// 	}

// 	return n, err
// }
