package claude

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/c00/botman-v2/chatbot"
	"github.com/c00/botman-v2/chattools"
)

const (
	ContentTypeText           = "text"
	ContentTypeTextDelta      = "text_delta"
	ContentTypeToolCall       = "tool_use"
	ContentTypeInputJsonDelta = "input_json_delta"
	ContentTypeToolResult     = "tool_result"
)

// Modeled after what to POST and what we would receive if we'd NOT
// stream the response.
type ClaudeMessage struct {
	//assistant or user
	Role    string         `json:"role"`
	Content []ContentBlock `json:"content"`
}

type ContentBlock struct {
	Type string

	TextBlock           *TextBlock
	TextDeltaBlock      *TextDeltaBlock
	ToolCallBlock       *ToolCallBlock
	InputJsonDeltaBlock *InputJsonDeltaBlock
	ToolResultBlock     *ToolResultBlock
}

func (b ContentBlock) Delta() string {
	switch b.Type {
	case ContentTypeText:
		return b.TextBlock.Text
	case ContentTypeTextDelta:
		return b.TextDeltaBlock.Text
	}

	return ""
}

func (b *ContentBlock) Add(deltaBlock ContentBlock) {
	switch deltaBlock.Type {
	case ContentTypeTextDelta:
		if b.Type != ContentTypeText {
			log.Warn("cannot add text content to block type %v", b.Type)
			return
		}
		b.TextBlock.Text += deltaBlock.TextDeltaBlock.Text
	case ContentTypeInputJsonDelta:
		if b.Type != ContentTypeToolCall {
			log.Warn("cannot add json content to block type %v", b.Type)
			return
		}
		if b.ToolCallBlock.rawInput == nil {
			b.ToolCallBlock.rawInput = []byte{}
		}
		b.ToolCallBlock.rawInput = append(b.ToolCallBlock.rawInput, []byte(deltaBlock.InputJsonDeltaBlock.PartialJson)...)
	default:
		log.Warn("cannot add block of type %v to block type %v", deltaBlock.Type, b.Type)
	}
}

func (b *ContentBlock) Finalize() {
	switch b.Type {
	case ContentTypeToolCall:
		b.ToolCallBlock.Input = map[string]any{}
		err := json.Unmarshal(b.ToolCallBlock.rawInput, &b.ToolCallBlock.Input)
		if err != nil {
			log.Warn("could not parse inputs: %v", err)
		} else {
			b.ToolCallBlock.rawInput = nil
		}

	}
}

func (b ContentBlock) MarshalJSON() ([]byte, error) {
	switch b.Type {
	case ContentTypeText:
		return json.Marshal(b.TextBlock)
	case ContentTypeToolCall:
		return json.Marshal(b.ToolCallBlock)
	case ContentTypeToolResult:
		return json.Marshal(b.ToolResultBlock)
	case ContentTypeTextDelta:
		return json.Marshal(b.TextDeltaBlock)
	case ContentTypeInputJsonDelta:
		return json.Marshal(b.InputJsonDeltaBlock)
	default:
		return []byte{}, fmt.Errorf("unknown content type: %v", b.Type)
	}
}

func (b *ContentBlock) UnmarshalJSON(raw []byte) error {
	var data map[string]any
	json.Unmarshal(raw, &data)
	if val, ok := data["type"].(string); ok {
		b.Type = val
	}

	switch b.Type {
	case ContentTypeText:
		b.TextBlock = &TextBlock{}
		json.Unmarshal(raw, b.TextBlock)
	case ContentTypeToolCall:
		b.ToolCallBlock = &ToolCallBlock{}
		json.Unmarshal(raw, b.ToolCallBlock)
	case ContentTypeToolResult:
		b.ToolResultBlock = &ToolResultBlock{}
		json.Unmarshal(raw, b.ToolResultBlock)
	case ContentTypeTextDelta:
		b.TextDeltaBlock = &TextDeltaBlock{}
		json.Unmarshal(raw, b.TextDeltaBlock)
	case ContentTypeInputJsonDelta:
		b.InputJsonDeltaBlock = &InputJsonDeltaBlock{}
		json.Unmarshal(raw, b.InputJsonDeltaBlock)
	default:
		return fmt.Errorf("unknown content type: %v", b.Type)
	}

	return nil
}

// Todo test
func (cm ClaudeMessage) ToChatMessage() chatbot.ChatMessage {
	texts := []string{}
	toolCalls := []chattools.ToolCall{}

	for _, c := range cm.Content {
		switch c.Type {
		case ContentTypeText:
			texts = append(texts, c.TextBlock.Text)
		case ContentTypeToolCall:
			toolCalls = append(toolCalls, chattools.ToolCall{
				ID:     c.ToolCallBlock.ID,
				Name:   c.ToolCallBlock.Name,
				Params: c.ToolCallBlock.Input,
			})
		case ContentTypeToolResult:
			log.Warn("Why are we calling toChatMessage on a ToolResult?")
		}
	}

	msg := chatbot.ChatMessage{
		Role:    cm.Role,
		Content: strings.Join(texts, " "),
	}
	if len(toolCalls) > 0 {
		msg.ToolCalls = toolCalls
	}
	return msg
}

func chatMessageToClaudeMessage(msg chatbot.ChatMessage) ClaudeMessage {
	cm := ClaudeMessage{
		Role:    msg.Role,
		Content: []ContentBlock{},
	}

	if msg.Content != "" {
		cm.Content = append(cm.Content, ContentBlock{Type: ContentTypeText, TextBlock: &TextBlock{Type: ContentTypeText, Text: msg.Content}})
	}

	if msg.Role == chatbot.ChatMessageRoleTool {
		cm.Role = chatbot.ChatMessageRoleUser
	}
	if len(msg.ToolResults) > 0 {
		for _, tr := range msg.ToolResults {
			cm.Content = append(cm.Content, ContentBlock{Type: ContentTypeToolResult, ToolResultBlock: &ToolResultBlock{
				Type:      ContentTypeToolResult,
				ToolUseId: tr.ID,
				Content:   tr.Content,
				IsError:   !tr.Success,
			}})
		}
	}

	if len(msg.ToolCalls) > 0 {
		for _, tr := range msg.ToolCalls {
			cm.Content = append(cm.Content, ContentBlock{Type: ContentTypeToolCall, ToolCallBlock: &ToolCallBlock{
				Type:  ContentTypeToolCall,
				ID:    tr.ID,
				Name:  tr.Name,
				Input: tr.Params,
			}})
		}
	}

	return cm
}

// type: text
type TextBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type TextDeltaBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type InputJsonDeltaBlock struct {
	Type        string `json:"type"`
	PartialJson string `json:"partial_json"`
}

// type: tool_use
type ToolCallBlock struct {
	Type  string         `json:"type"`
	ID    string         `json:"id"`
	Name  string         `json:"name"`
	Input map[string]any `json:"input"`
	//Raw JSON
	rawInput []byte `json:"-"`
}

// type: tool_result
type ToolResultBlock struct {
	Type      string `json:"type"`
	ToolUseId string `json:"tool_use_id"`
	Content   string `json:"content"`
	IsError   bool   `json:"is_error,omitempty"`
}
