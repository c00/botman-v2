package chattools

import (
	"github.com/c00/botman-v2/jsonschema"
)

const (
	ToolTypeAddNumbers = "add"
	ToolTypeDalle      = "dalle"
	ToolTypeSdxl       = "sdxl"
)

type ToolDefinition struct {
	ToolType     string      `yaml:"toolType"`
	Name         string      `yaml:"name"`
	Description  string      `yaml:"description"`
	SdxlSettings *SdxlConfig `yaml:"sdxl"`
	// MockSettings *MockSettings
	// DalleSettings *DalleSettings
}

func (d ToolDefinition) Schema() jsonschema.JsonSchema {
	schema := jsonschema.New()

	switch d.ToolType {
	case ToolTypeAddNumbers:
		schema.AddNumber("a", "the first number", true)
		schema.AddNumber("b", "the second number", true)
		return schema
	case ToolTypeSdxl:
		schema.AddString("prompt", "The positive prompt for the image. Be descriptive. Describe at least the scene, the quality, the type of image.", true)
		schema.AddString("negativePrompt", "The negative prompt for the image. Optionally add things you don't want in the prompt. This can be helpful to mitigate common problems with SDXL, such as hands with too many fingers.", false)
	}
	//Todo implement the rest
	return schema
}

// Chatters will return this when they require a tool to be called
type ToolCall struct {
	ID     string
	Name   string
	Params map[string]any
}

type ToolResult struct {
	ID      string
	Name    string
	Content string
	Success bool
	Value   any
}

type SdxlConfig struct {
	NegativePrompt string `json:"negativePrompt"`
	PositivePrompt string `json:"positivePrompt"`
	ApiKey         string `json:"apiKey,omitempty"`
	Model          string `json:"model,omitempty"`
}
