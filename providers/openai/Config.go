package openai

type Config struct {
	ApiKey       string `yaml:"apiKey"`
	Model        string `yaml:"model"`
	SystemPrompt string `yaml:"systemPrompt"`
}
