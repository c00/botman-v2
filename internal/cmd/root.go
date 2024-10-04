package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/c00/botman-v2/clitools"
	"github.com/c00/botman-v2/internal/cmd/mainloop"
	"github.com/c00/botman-v2/internal/config"
	"github.com/c00/botman-v2/internal/history"
	"github.com/c00/botman-v2/internal/storageprovider"
	"github.com/c00/botman-v2/logger"
	"github.com/spf13/cobra"
)

var verbosity int
var versionFlag *bool
var helpFlag *bool
var interactiveFlag *bool
var showConfig *bool
var historyFlag *int
var configFile *string
var continueFlag *bool

var log = logger.New("main")

func init() {
	rootCmd.PersistentFlags().CountVarP(&verbosity, "verbose", "v", "increase verbosity")
	versionFlag = rootCmd.Flags().BoolP("version", "", false, "Prints the version")
	helpFlag = rootCmd.Flags().BoolP("help", "", false, "Prints help")
	interactiveFlag = rootCmd.Flags().BoolP("interactive", "i", false, "Creates an interactive chat session rather than a single response")
	showConfig = rootCmd.Flags().BoolP("show-config", "", false, "Shows current configuration")
	configFile = rootCmd.Flags().StringP("config", "", "", "Use configuration file")
	continueFlag = rootCmd.Flags().BoolP("continue", "c", false, "Continue the last conversation. Does not show the conversation so far. Use -ih 0 for that instead.")
	historyFlag = rootCmd.Flags().IntP("history", "h", -1, "Show historical chat, looking baxk [n] chats. Can be combined with -i to continue conversation")
}

var rootCmd = &cobra.Command{
	Use:   binary,
	Short: fmt.Sprintf("%v is a tool for talking to LLMs.", binary),
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
		if *versionFlag {
			logger.Log("%v %v", binary, version)
			return
		}

		if *helpFlag {
			cmd.Help()
			return
		}

		if *configFile == "" {
			*configFile = config.GetUserConfigFilename()
		}

		if *showConfig {
			content, err := os.ReadFile(*configFile)
			if err != nil {
				log.Error("could not open config file: %v", err)
				return
			}
			log.Log("Config file: %v\n\n%v", *configFile, string(content))
			return
		}

		conf, err := config.Load(*configFile)
		if err != nil {
			log.Error2(err)
			os.Exit(1)
		}

		//Do some magic to inject api keys into other places in the config
		conf.InjectApiKeys()

		chatter, err := getChatter(conf)
		if err != nil {
			log.Error("cannot instantiate chatter: %v", err)
			os.Exit(1)
		}

		histPath := filepath.Join(config.GetUserConfigPath(), "history")
		histKeeper := history.NewYamlHistory(histPath)
		var activeConversation *history.HistoryEntry

		//history
		if *continueFlag {
			*historyFlag = 0
			*interactiveFlag = true
		}
		if *historyFlag > -1 {
			chat, err := histKeeper.LoadChat(*historyFlag)
			if err != nil {
				log.Error("could not load chat: %v", err)
				return
			}

			if !*continueFlag {
				chat.Print()
			}

			if !*interactiveFlag {
				return
			}
			activeConversation = &chat
		}

		pipedIn, err := clitools.GetPipedIn()
		if err != nil {
			log.Error2(err)
			os.Exit(1)
		}

		prompt := strings.TrimSpace(fmt.Sprintf("%v %v", pipedIn, strings.Join(args, " ")))

		if prompt == "" {
			*interactiveFlag = true
		}

		store, err := getStorageProvider(conf.Storage)
		if err != nil {
			log.Error("cannot instantiate storage provider: %v", err)
			os.Exit(1)
		}

		ml := mainloop.New(chatter, histKeeper, store, *interactiveFlag, 0, os.Stdin, os.Stdout)
		if conf.Tools != nil {
			ml.SetTools(conf.Tools)
		}
		if activeConversation != nil {
			ml.SetConversation(*activeConversation)
		}
		ml.Start(prompt)

		log.Debug("Main loop finished.")
	},
}

func getStorageProvider(conf config.StorageConfig) (storageprovider.StorageProvider, error) {
	if conf.Type == "" {
		conf.Type = storageprovider.StorageTypeLocal
	}

	switch conf.Type {
	case storageprovider.StorageTypeLocal:
		store, err := storageprovider.NewLocalStore(filepath.Join(config.GetUserConfigPath(), "generated-images"))
		if err != nil {
			return nil, fmt.Errorf("cannot create local storage provider: %v", store)
		}
		return store, nil
	case storageprovider.StorageTypeMemory:
		return storageprovider.NewMemStore(), nil
	case storageprovider.StorageTypeS3:
		if conf.S3 == nil {
			return nil, fmt.Errorf("missing S3 Storage configuration")
		}
		return storageprovider.NewS3Store(*conf.S3), nil
	}

	return nil, fmt.Errorf("unknown storage provider: %v", conf.Type)
}
