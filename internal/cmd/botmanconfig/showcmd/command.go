package showcmd

import (
	"os"

	"github.com/c00/botman-v2/internal/config"
	"github.com/c00/botman-v2/internal/logger"
	"github.com/spf13/cobra"
)

var log = logger.New("showCmd")

var Command = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
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
		filename := config.GetUserConfigFilename()
		data, err := os.ReadFile(filename)
		if err != nil {
			log.Error("could not read config file: %v", err)
			os.Exit(1)
		}

		log.Log("%v", string(data))
	},
}
