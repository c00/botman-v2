package openai

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/c00/botman-v2/chatbot"
	"github.com/c00/botman-v2/chattools"
	"github.com/c00/botman-v2/internal/channeltools"
	"github.com/c00/botman-v2/logger"
	openai "github.com/sashabaranov/go-openai"
)

var log = logger.New("Openai")

func New(cfg Config) (*OpenAi, error) {

	return &OpenAi{
		client: openai.NewClient(cfg.ApiKey),
		cfg:    cfg,
	}, nil
}

type OpenAi struct {
	client   *openai.Client
	cfg      Config
	messages []openai.ChatCompletionMessage
}

func (c *OpenAi) GetStreamingResponse(message chatbot.ChatMessage, streamChan chan<- string) (chatbot.ChatMessage, error) {
	defer close(streamChan)

	log.Debug("GetStreamingResponse Content: %v", message.Content)

	c.messages = append(c.messages, openai.ChatCompletionMessage{Role: message.Role, Content: message.Content})

	postMessages := []openai.ChatCompletionMessage{
		{Role: "system", Content: c.cfg.SystemPrompt},
	}
	postMessages = append(postMessages, c.messages...)

	stream, err := c.client.CreateChatCompletionStream(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    c.cfg.Model,
			Messages: postMessages,
		},
	)

	if err != nil {
		return chatbot.ChatMessage{}, fmt.Errorf("error getting OpenAi Chat Completion: %v", err)
	}
	defer stream.Close()

	responseContent := make([]string, 0, 50)

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			message := openai.ChatCompletionMessage{
				Role:    chatbot.ChatMessageRoleAssistant,
				Content: strings.Join(responseContent, ""),
			}
			c.messages = append(c.messages, message)
			return chatbot.ChatMessage{Role: message.Role, Content: message.Content}, nil
		}

		if err != nil {
			return chatbot.ChatMessage{}, fmt.Errorf("stream error: %v", err)
		}

		streamChan <- response.Choices[0].Delta.Content

		responseContent = append(responseContent, response.Choices[0].Delta.Content)
	}
}

func (c *OpenAi) GetResponse(message chatbot.ChatMessage) (chatbot.ChatMessage, error) {
	ch := channeltools.BlackHoleChannel()
	return c.GetStreamingResponse(message, ch)
}

func (c *OpenAi) AddMessages(messages []chatbot.ChatMessage) {
	for _, m := range messages {
		c.messages = append(c.messages, convertMessage(m))
	}
}

func (c *OpenAi) SetMessages(messages []chatbot.ChatMessage) {
	c.messages = make([]openai.ChatCompletionMessage, 0, len(messages))
	for _, m := range messages {
		c.messages = append(c.messages, convertMessage(m))
	}
}

func (c *OpenAi) GetMessages() []chatbot.ChatMessage {
	result := make([]chatbot.ChatMessage, 0, len(c.messages))

	for _, m := range c.messages {
		result = append(result, chatbot.ChatMessage{Role: m.Role, Content: m.Content})
	}

	return result
}

func (c *OpenAi) SetSystemPrompt(prompt string) {
	c.cfg.SystemPrompt = prompt

}

func (c *OpenAi) GetSystemPrompt() string {
	return c.cfg.SystemPrompt
}

// Get a list of features that this chatter supports.
func (c OpenAi) SupportedFeatures() []string {
	return []string{}
}

// Set tools for Yappie
func (c *OpenAi) SetTools(tools []chattools.ToolDefinition) {
	panic("tools not supported")
}

func convertMessage(m chatbot.ChatMessage) openai.ChatCompletionMessage {
	return openai.ChatCompletionMessage{
		Role:    m.Role,
		Content: m.Content,
	}
}
