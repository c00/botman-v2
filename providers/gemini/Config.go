package gemini

type Config struct {
	ApiKey       string `yaml:"apiKey"`
	Model        string `yaml:"model"`
	SystemPrompt string `yaml:"systemPrompt"`
}

//https://ai.google.dev/gemini-api/docs/models

var Models = []string{
	"gemini-2.5-pro-exp-03-25",
	"gemini-2.0-flash",
	"gemini-2.0-flash-lite",
	"gemini-1.5-flash",
	"gemini-1.5-flash-8b",
	"gemini-1.5-pro",
}
