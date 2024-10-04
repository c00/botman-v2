package clitools

import (
	"fmt"
	"io"
	"os"
)

func GetPipedIn() (string, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return "", fmt.Errorf("could not stat stdin: %w", err)
	}

	if info.Mode()&os.ModeNamedPipe == 0 {
		//No stdin
		return "", nil
	}

	content, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("could not read stdin: %w", err)
	}

	return string(content), nil
}
