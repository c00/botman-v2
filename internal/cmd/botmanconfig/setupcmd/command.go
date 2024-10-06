package setupcmd

import (
	"github.com/c00/botman-v2/internal/config"
	"github.com/c00/botman-v2/internal/logger"
	"github.com/spf13/cobra"
)

var log = logger.New("setupCmd")

var Command = &cobra.Command{
	Use:   "setup",
	Short: "Set up or update botman configuration",
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
		//Get current config
		conf := config.LoadFromUser()
		runSetup(conf)
	},
}
