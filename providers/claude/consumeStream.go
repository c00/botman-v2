package claude

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

func consumeStream(reader *bufio.Reader, ch chan<- string) (ClaudeMessage, error) {
	dataPrefix := []byte("data: ")
	msgs := StreamMessages{}

	for {
		line, err := reader.ReadBytes(byte('\n'))
		if errors.Is(err, io.EOF) {
			responseMessage := msgs.ToFinalMessage()
			return responseMessage, nil
		}

		if err != nil {
			log.Error("Error getting Claude Chat Completion: %v", err)
			return ClaudeMessage{}, fmt.Errorf("error getting Claude Chat Completion: %w", err)
		}

		if len(line) > 5 && bytes.Equal(line[:6], dataPrefix) {
			msg := StreamMessage{}
			err := json.Unmarshal(line[6:], &msg)
			if err != nil {
				return ClaudeMessage{}, fmt.Errorf("cannot parse message: %w", err)
			}

			ch <- msg.Delta()
			msgs = append(msgs, msg)
		}

	}
}
