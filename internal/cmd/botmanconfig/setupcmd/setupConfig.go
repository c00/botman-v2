package setupcmd

import (
	"fmt"
	"os"

	"github.com/c00/botman-v2/clitools"
	"github.com/c00/botman-v2/internal/config"
)

func runSetup(conf config.BotmanConfig) {
	log.Log("Botman Setup\n")

	//Preferred provider
	log.Log("Select your preferred LLM Provider")
	currentChoiceIndex := 0
	if conf.LlmProvider == config.LlmProviderFireworksAi {
		currentChoiceIndex = 1
	} else if conf.LlmProvider == config.LlmProviderClaude {
		currentChoiceIndex = 2
	}

	choice := clitools.GetChoice([]string{"Open AI", "Fireworks AI", "Claude"}, currentChoiceIndex, os.Stdin, os.Stdout)
	if choice == 0 {
		conf.LlmProvider = config.LlmProviderOpenAi
		setupApiKey(&conf.OpenAi.ApiKey, "OpenAI")
		chooseModel(&conf.OpenAi.Model, OpenAiModels)
	} else if choice == 1 {
		conf.LlmProvider = config.LlmProviderFireworksAi
		setupApiKey(&conf.FireworksAi.ApiKey, "Fireworks AI")
		chooseModel(&conf.FireworksAi.Model, FireworksAIModels)
	} else if choice == 2 {
		conf.LlmProvider = config.LlmProviderClaude
		setupApiKey(&conf.Claude.ApiKey, "Claude")
		chooseModel(&conf.Claude.Model, ClaudeModels)

		// Set defaults
		if conf.Claude.MaxTokens == 0 {
			conf.Claude.MaxTokens = 1024
		}
	}

	//todo setup tools

	log.Log("")
	err := config.SaveForUser(conf)
	if err != nil {
		log.Log("could not update the configuration: %v", err)
		os.Exit(1)
	}

	log.Log("Configuration has been updated")
}

// Setup API key and return true if changes were made.
func setupApiKey(key *string, name string) bool {
	if *key == "" {
		input := clitools.GetInput(fmt.Sprintf("Enter your %v API key", name), os.Stdin, os.Stdout)
		*key = input
		return true
	} else {
		fmt.Printf("Current %v API key: %v\n", name, *key)
		input := clitools.GetInput(fmt.Sprintf("Enter your new %v API key, or press [enter] to keep the current one", name), os.Stdin, os.Stdout)
		if input != "" {
			*key = input
			return true
		}
	}

	return false
}

// Select a model
func chooseModel(currentModel *string, models []string) bool {

	log.Log("\nChoose a model:")

	index := indexOf(models, *currentModel)
	chosen := clitools.GetChoice(models, index, os.Stdin, os.Stdout)
	if chosen == -1 {
		return false
	}

	newModel := models[chosen]
	if newModel == *currentModel {
		return false
	} else {
		*currentModel = newModel
		return true
	}
}

func indexOf(s []string, searchString string) int {
	for i, v := range s {
		if v == searchString {
			return i
		}
	}
	return -1
}
