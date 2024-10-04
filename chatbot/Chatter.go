package chatbot

import "github.com/c00/botman-v2/chattools"

// Chatter is an interface to an LLM provider
type Chatter interface {
	SupportedFeatures() []string
	SetTools([]chattools.ToolDefinition)

	GetResponse(ChatMessage) (ChatMessage, error)
	GetStreamingResponse(ChatMessage, chan<- string) (ChatMessage, error)

	//Append messages
	AddMessages(messages []ChatMessage)
	//Replace messages
	SetMessages(messages []ChatMessage)
	GetMessages() []ChatMessage

	SetSystemPrompt(prompt string)
	GetSystemPrompt() string
}
