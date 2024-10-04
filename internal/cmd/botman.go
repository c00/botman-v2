package main

import (
	"os"

	"github.com/c00/botman-v2/logger"
)

const version = "2.0.3"
const binary = "botman"

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
