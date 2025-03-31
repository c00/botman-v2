package claude

type Config struct {
	ApiKey       string `yaml:"apiKey"`
	Model        string `yaml:"model"`
	SystemPrompt string `yaml:"systemPrompt"`
	MaxTokens    int    `yaml:"maxTokens"`
}
