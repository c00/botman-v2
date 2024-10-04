package claude

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/c00/botman-v2/chatbot"
	"github.com/c00/botman-v2/chattools"
	"github.com/c00/botman-v2/internal/channeltools"
	"github.com/c00/botman-v2/jsonschema"
	"github.com/c00/botman-v2/logger"
)

const apiUrl = "https://api.anthropic.com/v1/messages"

var log = logger.New("Claude")

func New(cfg Config) (*Claude, error) {
	if cfg.ApiKey == "" {
		return nil, errors.New("missing claude api key")
	}

	return &Claude{
		cfg: cfg,
	}, nil
}

type Claude struct {
	cfg      Config
	messages []ClaudeMessage
	tools    []chattools.ToolDefinition
}

type claudePostBody struct {
	Model     string          `json:"model"`
	Messages  []ClaudeMessage `json:"messages"`
	MaxTokens int             `json:"max_tokens"`
	Stream    bool            `json:"stream,omitempty"`
	System    string          `json:"system"`
	Tools     []claudeToolDef `json:"tools,omitempty"`
}

type claudeToolDef struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	InputSchema jsonschema.JsonSchema `json:"input_schema"`
}

func (c *Claude) GetStreamingResponse(message chatbot.ChatMessage, streamChan chan<- string) (chatbot.ChatMessage, error) {
	defer close(streamChan)

	log.Debug("GetStreamingResponse: %v", message.Sprint())

	c.messages = append(c.messages, chatMessageToClaudeMessage(message))

	body := claudePostBody{
		Model:     c.cfg.Model,
		Messages:  c.messages,
		System:    c.cfg.SystemPrompt,
		Stream:    true,
		MaxTokens: c.cfg.MaxTokens,
	}

	//Add tools
	if len(c.tools) > 0 {
		body.Tools = []claudeToolDef{}
		for _, t := range c.tools {
			body.Tools = append(body.Tools, claudeToolDef{
				Name:        t.Name,
				Description: t.Description,
				InputSchema: t.Schema(),
			})
		}
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return chatbot.ChatMessage{}, fmt.Errorf("cannot marshal claude post body: %w", err)
	}

	// indented, _ := json.MarshalIndent(body, "", "  ")
	// log.Debug("\n\nBODY: %v\n\n", string(indented))

	req, err := http.NewRequest("POST", apiUrl, bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("anthropic-version", "2023-06-01") //https://docs.anthropic.com/en/api/versioning
	req.Header.Set("x-api-key", c.cfg.ApiKey)
	if err != nil {
		return chatbot.ChatMessage{}, fmt.Errorf("cannot create request for claude: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return chatbot.ChatMessage{}, fmt.Errorf("cannot marshal claude post body: %w", err)
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	if resp.StatusCode != 200 {
		bodyContent, err := io.ReadAll(reader)
		if err != nil {
			log.Error("failed to get response: %v, %v", resp.StatusCode, err)
			return chatbot.ChatMessage{}, fmt.Errorf("got status: %v, error: %w", resp.StatusCode, err)
		}
		log.Error("failed to get response: %v, %v", string(bodyContent), err)
		return chatbot.ChatMessage{}, fmt.Errorf("got status: %v, %v", resp.StatusCode, string(bodyContent))
	}

	responseMessage, err := consumeStream(reader, streamChan)
	if err != nil {
		return chatbot.ChatMessage{}, fmt.Errorf("stream consume error: %w", err)
	}
	c.messages = append(c.messages, responseMessage)
	return responseMessage.ToChatMessage(), nil
}

func (c *Claude) GetResponse(message chatbot.ChatMessage) (chatbot.ChatMessage, error) {
	ch := channeltools.BlackHoleChannel()
	return c.GetStreamingResponse(message, ch)
}

func (c *Claude) AddMessages(messages []chatbot.ChatMessage) {
	for _, m := range messages {
		c.messages = append(c.messages, chatMessageToClaudeMessage(m))
	}
}

func (c *Claude) SetMessages(messages []chatbot.ChatMessage) {
	c.messages = make([]ClaudeMessage, 0, len(messages))
	for _, m := range messages {
		c.messages = append(c.messages, chatMessageToClaudeMessage(m))
	}
}

func (c *Claude) GetMessages() []chatbot.ChatMessage {
	result := make([]chatbot.ChatMessage, 0, len(c.messages))

	for _, cm := range c.messages {
		result = append(result, cm.ToChatMessage())
	}

	return result
}

func (c *Claude) SetSystemPrompt(prompt string) {
	c.cfg.SystemPrompt = prompt

}

func (c *Claude) GetSystemPrompt() string {
	return c.cfg.SystemPrompt
}

// Get a list of features that this chatter supports.
func (c Claude) SupportedFeatures() []string {
	return []string{"tools"}
}

// Set tools
func (c *Claude) SetTools(tools []chattools.ToolDefinition) {
	c.tools = tools
}
