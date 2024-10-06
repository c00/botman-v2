package storagecmd

import (
	"os"

	"github.com/c00/botman-v2/clitools"
	"github.com/c00/botman-v2/internal/config"
	"github.com/c00/botman-v2/internal/logger"
	"github.com/c00/botman-v2/internal/storageprovider"
	"github.com/spf13/cobra"
)

var log = logger.New("storageCmd")

var Command = &cobra.Command{
	Use:   "storage",
	Short: "Manage how botman stores additional files like generated images",
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

func runSetup(conf config.BotmanConfig) {
	currentChoiceIndex := 99
	switch conf.Storage.Type {
	case storageprovider.StorageTypeLocal:
		currentChoiceIndex = 0
	case storageprovider.StorageTypeS3:
		currentChoiceIndex = 1
	}

	choice := clitools.GetChoice([]string{"local", "s3"}, currentChoiceIndex, os.Stdin, os.Stdout)
	switch choice {
	case 0:
		conf.Storage.Type = storageprovider.StorageTypeLocal
		//We don't configure anything for this.
	case 1:
		conf.Storage.Type = storageprovider.StorageTypeS3
		//configure the thing.
		log.Log("To setup S3 you will need: A bucket, an Endpoint, an Access Key, a Secret Key and a region. For more information check the manual here: [todo add link]\n")
		if conf.Storage.S3 == nil {
			s3 := &storageprovider.S3Config{ForcePathStyle: true}
			conf.Storage.S3 = s3
		}

		clitools.SetInput("Endpoint", &conf.Storage.S3.Endpoint, os.Stdin, os.Stdout)
		clitools.SetInput("Bucket", &conf.Storage.S3.Bucket, os.Stdin, os.Stdout)
		clitools.SetInput("Access Key", &conf.Storage.S3.AccessKey, os.Stdin, os.Stdout)
		clitools.SetInput("Secret Key", &conf.Storage.S3.SecretKey, os.Stdin, os.Stdout)
		clitools.SetInput("Region", &conf.Storage.S3.Region, os.Stdin, os.Stdout)
		clitools.SetInput("Base path (optional)", &conf.Storage.S3.BasePath, os.Stdin, os.Stdout)
	}

	err := config.SaveForUser(conf)
	if err != nil {
		log.Error("cannot save config: %v", err)
		os.Exit(1)
	}
	log.Log("Configuration saved.")
}
