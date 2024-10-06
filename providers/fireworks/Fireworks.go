package fireworks

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/c00/botman-v2/chatbot"
	"github.com/c00/botman-v2/chattools"
	"github.com/c00/botman-v2/internal/channeltools"
	"github.com/c00/botman-v2/internal/logger"
)

const apiUrl = "https://api.fireworks.ai/inference/v1/chat/completions"

var log = logger.New("Fireworks")

func New(cfg Config) (*Fireworks, error) {
	if cfg.ApiKey == "" {
		return nil, errors.New("missing claude api key")
	}

	return &Fireworks{
		cfg: cfg,
	}, nil
}

type Fireworks struct {
	cfg      Config
	messages []fireworksMessage
}

type fireworksPostBody struct {
	Model            string             `json:"model"`
	Messages         []fireworksMessage `json:"messages"`
	MaxTokens        string             `json:"max_tokens,omitempty"`
	TopP             int                `json:"top_p,omitempty"`
	TopK             int                `json:"top_k,omitempty"`
	PresencePenalty  int                `json:"presence_penalty,omitempty"`
	FrequencyPenalty int                `json:"frequency_penalty,omitempty"`
	Temperature      float32            `json:"temperature,omitempty"`
	Stream           bool               `json:"stream,omitempty"`
	N                int                `json:"n,omitempty"`
}

type fireworksMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (c *Fireworks) GetStreamingResponse(message chatbot.ChatMessage, streamChan chan<- string) (chatbot.ChatMessage, error) {
	defer close(streamChan)

	log.Debug("GetStreamingResponse Content: %v", message.Content)

	c.messages = append(c.messages, fireworksMessage{Role: message.Role, Content: message.Content})

	postMessages := []fireworksMessage{
		{Role: "system", Content: c.cfg.SystemPrompt},
	}
	postMessages = append(postMessages, c.messages...)

	body := fireworksPostBody{
		Model:    c.cfg.Model,
		Messages: postMessages,
		Stream:   true,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return chatbot.ChatMessage{}, fmt.Errorf("cannot marshall post body: %w", err)
	}

	req, err := http.NewRequest("POST", apiUrl, bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.cfg.ApiKey))
	if err != nil {
		return chatbot.ChatMessage{}, fmt.Errorf("cannot create request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return chatbot.ChatMessage{}, fmt.Errorf("cannot do request: %w", err)
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	if resp.StatusCode != 200 {
		bodyContent, err := io.ReadAll(reader)
		if err != nil {
			return chatbot.ChatMessage{}, fmt.Errorf("got status: %v, error: %v", resp.StatusCode, err)
		}
		return chatbot.ChatMessage{}, fmt.Errorf("got status: %v, %v", resp.StatusCode, string(bodyContent))
	}

	responseContent := make([]string, 0, 50)

	//Read the streaming response.
	//I'm assuming each message will have a newline at the end.
	for {
		line, err := reader.ReadBytes(byte('\n'))
		chunk := parseChunk(line)

		if errors.Is(err, io.EOF) {
			if !chunk.Empty {
				streamChan <- chunk.Delta
				responseContent = append(responseContent, chunk.Delta)
			}

			message := fireworksMessage{Role: chatbot.ChatMessageRoleAssistant, Content: strings.Join(responseContent, "")}
			c.messages = append(c.messages, message)

			return chatbot.ChatMessage{Role: message.Role, Content: message.Content}, nil
		}

		if chunk.Empty {
			continue
		}

		if err != nil {
			return chatbot.ChatMessage{}, fmt.Errorf("error getting FireworksAI Chat Completion: %w", err)
		}

		if chunk.LastMessage {
			continue
		}

		streamChan <- chunk.Delta
		responseContent = append(responseContent, chunk.Delta)
	}
}

func (c *Fireworks) GetResponse(message chatbot.ChatMessage) (chatbot.ChatMessage, error) {
	ch := channeltools.BlackHoleChannel()
	return c.GetStreamingResponse(message, ch)
}

func (c *Fireworks) AddMessages(messages []chatbot.ChatMessage) {
	for _, m := range messages {
		c.messages = append(c.messages, convertMessage(m))
	}
}

func (c *Fireworks) SetMessages(messages []chatbot.ChatMessage) {
	c.messages = make([]fireworksMessage, 0, len(messages))
	for _, m := range messages {
		c.messages = append(c.messages, convertMessage(m))
	}
}

func (c *Fireworks) GetMessages() []chatbot.ChatMessage {
	result := make([]chatbot.ChatMessage, 0, len(c.messages))

	for _, m := range c.messages {
		result = append(result, chatbot.ChatMessage{Role: m.Role, Content: m.Content})
	}

	return result
}

func (c *Fireworks) SetSystemPrompt(prompt string) {
	c.cfg.SystemPrompt = prompt

}

func (c *Fireworks) GetSystemPrompt() string {
	return c.cfg.SystemPrompt
}

// Get a list of features that this chatter supports.
func (c Fireworks) SupportedFeatures() []string {
	return []string{}
}

// Set tools for Yappie
func (c *Fireworks) SetTools(tools []chattools.ToolDefinition) {
	panic("tools not supported")
}

func convertMessage(m chatbot.ChatMessage) fireworksMessage {
	return fireworksMessage{
		Role:    m.Role,
		Content: m.Content,
	}
}
