package main

import (
	"os"

	"github.com/c00/botman-v2/internal/logger"
)

const binary = "botman"

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
