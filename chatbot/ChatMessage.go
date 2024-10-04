package chatbot

import (
	"fmt"
	"strings"

	"github.com/c00/botman-v2/chattools"
)

type ChatMessage struct {
	Role        string                 `yaml:"role"`
	Content     string                 `yaml:"content"`
	ToolCalls   []chattools.ToolCall   `yaml:"toolCalls,omitempty"`
	ToolResults []chattools.ToolResult `yaml:"toolResults,omitempty"`
}

func (msg ChatMessage) Sprint() string {
	parts := []string{fmt.Sprintf("Role: %v", msg.Role)}

	if msg.Content != "" {
		parts = append(parts, fmt.Sprintf("Content: %v", msg.Content))
	}

	if len(msg.ToolCalls) > 0 {
		for _, tc := range msg.ToolCalls {
			parts = append(parts, fmt.Sprintf("ToolCall: %v (%v): %+v", tc.Name, tc.ID, tc.Params))
		}
	}

	if len(msg.ToolResults) > 0 {
		for _, tr := range msg.ToolResults {
			parts = append(parts, fmt.Sprintf("ToolResult: %v (%v): %v - %v", tr.Name, tr.ID, tr.Content, tr.Value))
		}
	}

	return strings.Join(parts, "\n  ")
}

const (
	ChatMessageRoleUser      = "user"
	ChatMessageRoleAssistant = "assistant"
	ChatMessageRoleTool      = "tool"
)
