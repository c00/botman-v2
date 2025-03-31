package gemini

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/c00/botman-v2/chatbot"
	"github.com/c00/botman-v2/chattools"
	"github.com/c00/botman-v2/internal/channeltools"
	"github.com/c00/botman-v2/internal/logger"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var _ chatbot.Chatter = (*Gemini)(nil)

var log = logger.New("Gemini")

func New(cfg Config) (*Gemini, error) {
	if cfg.ApiKey == "" {
		return nil, errors.New("missing gemini api key")
	}

	client, err := genai.NewClient(context.Background(), option.WithAPIKey(cfg.ApiKey))
	if err != nil {
		return nil, fmt.Errorf("cannot create gemini client: %w", err)
	}

	log.Debug("Using Gemini model: %v", cfg.Model)

	return &Gemini{
		cfg:     cfg,
		client:  client,
		session: client.GenerativeModel(cfg.Model).StartChat(),
	}, nil
}

// Gemini is a mock Chatbot for testing.
type Gemini struct {
	cfg     Config
	client  *genai.Client
	session *genai.ChatSession
}

// Get a list of features that this chatter supports.
func (c Gemini) SupportedFeatures() []string {
	return []string{}
}

// Set tools for Gemini
func (c *Gemini) SetTools(tools []chattools.ToolDefinition) {
	panic("not supported")
}

func printResponse(resp *genai.GenerateContentResponse, ch chan<- string) string {
	parts := []string{}
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				str := fmt.Sprintf("%v", part)
				parts = append(parts, str)
				ch <- str
			}
		}
	}
	return strings.Join(parts, "")
}

func contentToString(contentParts []genai.Part) string {
	parts := []string{}

	for _, part := range contentParts {
		str := fmt.Sprintf("%v", part)
		parts = append(parts, str)
	}

	return strings.Join(parts, "")
}

func (c *Gemini) GetStreamingResponse(message chatbot.ChatMessage, streamChan chan<- string) (chatbot.ChatMessage, error) {
	defer close(streamChan)

	log.Debug("GetStreamingResponse: %v", message.Sprint())

	parts := []string{}

	iter := c.session.SendMessageStream(context.Background(), genai.Text(message.Content))
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return chatbot.ChatMessage{}, fmt.Errorf("cannot get response: %w", err)
		}
		parts = append(parts, printResponse(resp, streamChan))
	}

	return chatbot.ChatMessage{
		Role:    chatbot.ChatMessageRoleAssistant,
		Content: strings.Join(parts, ""),
	}, nil
}

func (c *Gemini) GetResponse(message chatbot.ChatMessage) (chatbot.ChatMessage, error) {
	ch := channeltools.BlackHoleChannel()
	return c.GetStreamingResponse(message, ch)
}

func toGeminiMessages(in []chatbot.ChatMessage) []*genai.Content {
	result := make([]*genai.Content, 0, len(in))
	for _, msg := range in {
		result = append(result, &genai.Content{
			Role:  msg.Role,
			Parts: []genai.Part{genai.Text(msg.Content)},
		})
	}

	return result
}

func toChatMessages(in []*genai.Content) []chatbot.ChatMessage {
	result := make([]chatbot.ChatMessage, 0, len(in))
	for _, msg := range in {
		result = append(result, chatbot.ChatMessage{
			Role:    msg.Role,
			Content: contentToString(msg.Parts),
		})

	}
	return result
}

func (c *Gemini) AddMessages(messages []chatbot.ChatMessage) {
	c.session.History = append(c.session.History, toGeminiMessages(messages)...)
}

func (c *Gemini) SetMessages(messages []chatbot.ChatMessage) {
	c.session.History = toGeminiMessages(messages)
}

func (c *Gemini) GetMessages() []chatbot.ChatMessage {
	return toChatMessages(c.session.History)
}

func (c *Gemini) SetSystemPrompt(prompt string) {
	c.cfg.SystemPrompt = prompt
}

func (c *Gemini) GetSystemPrompt() string {
	return c.cfg.SystemPrompt
}
