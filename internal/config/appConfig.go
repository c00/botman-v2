package config

import (
	"github.com/c00/botman-v2/chattools"
	"github.com/c00/botman-v2/internal/storageprovider"
	"github.com/c00/botman-v2/providers/claude"
	"github.com/c00/botman-v2/providers/fireworks"
	"github.com/c00/botman-v2/providers/openai"
)

const LlmProviderOpenAi = "openai"
const LlmProviderFireworksAi = "fireworksai"
const LlmProviderClaude = "claude"
const LlmProviderYappie = "yappie"

// To keep track of breaking changes in the config file
const currentVersion = 1

type BotmanConfig struct {
	Version      int                        `yaml:"version"`
	SaveHistory  bool                       `yaml:"saveHistory"`
	SystemPrompt string                     `yaml:"systemPrompt"`
	LlmProvider  string                     `yaml:"llmProvider"`
	OpenAi       openai.Config              `yaml:"openAi"`
	FireworksAi  fireworks.Config           `yaml:"fireworksAi"`
	Claude       claude.Config              `yaml:"claude"`
	Tools        []chattools.ToolDefinition `yaml:"tools"`
	Storage      StorageConfig              `yaml:"storage"`
}

// Inject API keys as defined in the chatters into tools where needed (e.g. openAi key for Dall-e and Fireworks API key for SDXL)
func (c *BotmanConfig) InjectApiKeys() {
	for idx, t := range c.Tools {
		changes := false

		if t.SdxlSettings != nil && t.SdxlSettings.ApiKey == "" {
			changes = true
			t.SdxlSettings.ApiKey = c.FireworksAi.ApiKey
		}

		if changes {
			c.Tools[idx] = t
		}
	}
}

// Currently only supports s3
type StorageConfig struct {
	Type string                    `yaml:"type"`
	S3   *storageprovider.S3Config `yaml:"s3"`
}
