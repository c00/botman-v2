package setupcmd

var OpenAiModels = []string{"gpt-4o", "gpt-4-turbo", "gpt-4", "gpt-3.5-turbo"}

// https://fireworks.ai/models?type=text
var FireworksAIModels = []string{
	"accounts/fireworks/models/deepseek-v3",
	"accounts/fireworks/models/deepseek-r1",
	"accounts/fireworks/models/firefunction-v2",
	"accounts/fireworks/models/firellava-13b",
	"accounts/fireworks/models/mixtral-8x7b-instruct",
	"accounts/fireworks/models/mixtral-8x22b-instruct",
	"accounts/fireworks/models/hermes-2-pro-mistral-7b",
	"accounts/fireworks/models/llama-v3-70b-instruct-hf",
	"accounts/fireworks/models/llama-v3p1-405b-instruct",
	"accounts/fireworks/models/llama-v3-8b-hf",
	"accounts/fireworks/models/mixtral-8x7b-instruct-hf",
	"accounts/fireworks/models/qwen2-72b-instruct",
}

// https://docs.anthropic.com/en/docs/about-claude/models/all-models
var ClaudeModels = []string{
	"claude-3-7-sonnet-latest",
	"claude-3-5-haiku-latest",
	"claude-3-5-sonnet-latest",
	"claude-3-5-sonnet-20240620",
	"claude-3-opus-20240229",
	"claude-3-sonnet-20240229",
	"claude-3-haiku-20240307",
}

// https://ai.google.dev/gemini-api/docs/models
var GeminiModels = []string{
	"gemini-2.5-pro-exp-03-25",
	"gemini-2.0-flash",
	"gemini-2.0-flash-lite",
	"gemini-1.5-flash",
	"gemini-1.5-flash-8b",
	"gemini-1.5-pro",
}
