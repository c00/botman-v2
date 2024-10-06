package main

import (
	"fmt"
	"os"

	"github.com/c00/botman-v2/internal/cmd/botmanconfig/setupcmd"
	"github.com/c00/botman-v2/internal/cmd/botmanconfig/showcmd"
	"github.com/c00/botman-v2/internal/cmd/botmanconfig/storagecmd"
	"github.com/c00/botman-v2/internal/config"
	"github.com/c00/botman-v2/internal/logger"
	"github.com/spf13/cobra"
)

const binary = "botman-conf"

var log = logger.New("main")

var rootCmd = &cobra.Command{
	Use:   binary,
	Short: fmt.Sprintf("%v is the configuration tool for botman", binary),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		verboseFlags, err := cmd.Flags().GetCount("verbose")
		if err != nil {
			return
		}

		if verboseFlags == 0 {
			return
		}

		if verboseFlags > 5 {
			verboseFlags = 5
		}

		logger.IncreaseLevel(verboseFlags)
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Log("Using config: %v\nTo show configuration run botman-config show\nTo set up or update configuration run botman-config setup", config.GetUserConfigFilename())
	},
}

func init() {
	rootCmd.AddCommand(setupcmd.Command, showcmd.Command, storagecmd.Command)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
