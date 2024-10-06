package yappie

import (
	"fmt"
	"strings"

	"github.com/c00/botman-v2/chatbot"
	"github.com/c00/botman-v2/chattools"
	"github.com/c00/botman-v2/internal/channeltools"
	"github.com/c00/botman-v2/internal/logger"
)

const defaultResponse = "Belloo! poopayee bappleees tank yuuu! Chasy potatoooo tulaliloo belloo! Belloo! baboiii hana dul sae jiji daa po kass. Hahaha hahaha uuuhhh chasy jeje. Butt baboiii poulet tikka masala pepete jeje hana dul sae. Bee do bee do bee do daa wiiiii tank yuuu! Potatoooo gelatooo po kass poopayee. Daa jiji tank yuuu! Uuuhhh bappleees ti aamoo! Gelatooo gelatooo. Tatata bala tu hahaha me want bananaaa! Bananaaaa wiiiii me want bananaaa! Wiiiii tatata bala tu."

var log = logger.New("Yappie")

// Yappie is a mock Chatbot for testing.
type Yappie struct {
	SystemPrompt string
	messages     []chatbot.ChatMessage
	tools        []chattools.ToolDefinition
	//Will return a tool use for the tool at this index if tool exists at index
	UseToolIndex int
}

// Get a list of features that this chatter supports.
func (c Yappie) SupportedFeatures() []string {
	return []string{"tools"}
}

// Set tools for Yappie
func (c *Yappie) SetTools(tools []chattools.ToolDefinition) {
	c.tools = tools
}

func (c *Yappie) GetStreamingResponse(newMessage chatbot.ChatMessage, streamChan chan<- string) (chatbot.ChatMessage, error) {
	log.Debug("GetStreamingResponse Content: %v", newMessage.Content)
	c.messages = append(c.messages, newMessage)

	for _, part := range strings.Split(defaultResponse, " ") {
		streamChan <- part + " "
	}

	close(streamChan)

	response := chatbot.ChatMessage{Role: chatbot.ChatMessageRoleAssistant, Content: defaultResponse}
	c.messages = append(c.messages, response)

	if len(c.messages) == 2 && c.UseToolIndex >= 0 && c.UseToolIndex < len(c.tools) {
		//Create tool call
		tool := c.tools[c.UseToolIndex]
		schema := tool.Schema()
		params := schema.Generate()

		paramMap, ok := params.(map[string]any)
		if !ok {
			return chatbot.ChatMessage{}, fmt.Errorf("generated params isn't a map")
		}

		toolCall := chattools.ToolCall{
			ID:     "random-id-1",
			Name:   c.tools[c.UseToolIndex].Name,
			Params: paramMap,
		}
		response.ToolCalls = []chattools.ToolCall{toolCall}
	}

	return response, nil
}

func (c *Yappie) GetResponse(message chatbot.ChatMessage) (chatbot.ChatMessage, error) {
	ch := channeltools.BlackHoleChannel()
	return c.GetStreamingResponse(message, ch)
}

func (c *Yappie) AddMessages(messages []chatbot.ChatMessage) {
	c.messages = append(c.messages, messages...)
}

func (c *Yappie) SetMessages(messages []chatbot.ChatMessage) {
	c.messages = messages
}

func (c *Yappie) GetMessages() []chatbot.ChatMessage {
	return c.messages
}

func (c *Yappie) SetSystemPrompt(prompt string) {
	c.SystemPrompt = prompt

}

func (c *Yappie) GetSystemPrompt() string {
	return c.SystemPrompt
}
