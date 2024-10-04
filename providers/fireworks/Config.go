package fireworks

type Config struct {
	ApiKey       string `yaml:"apiKey"`
	Model        string `yaml:"model"`
	SystemPrompt string `yaml:"systemPrompt"`
}

var Models = []string{
	"accounts/fireworks/models/firefunction-v2",
	"accounts/fireworks/models/firellava-13b",
	"accounts/fireworks/models/mixtral-8x7b-instruct",
	"accounts/fireworks/models/mixtral-8x22b-instruct",
	"accounts/fireworks/models/hermes-2-pro-mistral-7b",
	"accounts/fireworks/models/llama-v3-70b-instruct-hf",
	"accounts/fireworks/models/llama-v3-8b-hf",
	"accounts/fireworks/models/mixtral-8x7b-instruct-hf",
	"accounts/fireworks/models/qwen2-72b-instruct",
}
