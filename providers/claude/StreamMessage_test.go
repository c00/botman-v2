package claude

import (
	"encoding/json"
	"testing"

	"github.com/c00/botman-v2/chatbot"
	"github.com/stretchr/testify/assert"
)

func TestStreamMessage_JSON(t *testing.T) {
	//This tests marshalls and unmarshalls a claudeMessage and checks that it ends up the same.
	tests := []struct {
		name string
		data StreamMessage
	}{
		{
			name: "empty message", data: StreamMessage{
				Type: MsgTypePing,
			},
		},
		{
			name: "Message Start", data: StreamMessage{
				Type: MsgTypeMessageStart,
				MessageStart: &MessageStart{
					Type:    MsgTypeMessageStart,
					Message: MessageContent{Type: "message", ID: "m1234", Role: "assistant", Model: "haiku", Usage: Usage{InputTokens: 10, OutputTokens: 1}},
				},
			},
		},
		{
			name: "Message Delta", data: StreamMessage{
				Type: MsgTypeMessageDelta,
				MessageDelta: &MessageDelta{
					Type:  MsgTypeMessageDelta,
					Delta: MessageContent{Type: "message", StopReason: "tool_use"},
					Usage: Usage{InputTokens: 100, OutputTokens: 50},
				},
			},
		},
		{
			name: "Content Start", data: StreamMessage{
				Type: MsgTypeContentBlockStart,
				BlockStart: &BlockStart{
					Type:  MsgTypeContentBlockStart,
					Index: 1,
					ContentBlock: ContentBlock{
						Type:      ContentTypeText,
						TextBlock: &TextBlock{Type: ContentTypeText, Text: "ligma"},
					},
				},
			},
		},
		{
			name: "Content Delta", data: StreamMessage{
				Type: MsgTypeContentBlockDelta,
				BlockDelta: &BlockDelta{
					Type:  MsgTypeContentBlockDelta,
					Index: 1,
					Delta: ContentBlock{
						Type:      ContentTypeText,
						TextBlock: &TextBlock{Type: ContentTypeText, Text: "ligma"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := tt.data
			got, err := json.Marshal(original)
			assert.Nil(t, err)

			unmarshalled := StreamMessage{}
			err = json.Unmarshal(got, &unmarshalled)
			assert.Nil(t, err)

			assert.Equal(t, original, unmarshalled)
		})
	}
}

func TestParsedMessages_ToChatMessage(t *testing.T) {
	tests := []struct {
		name string
		pm   StreamMessages
		want ClaudeMessage
	}{
		{
			name: "basic message", pm: StreamMessages{
				{Type: MsgTypeContentBlockStart, BlockStart: &BlockStart{ContentBlock: ContentBlock{Type: ContentTypeText, TextBlock: &TextBlock{Type: "text", Text: "Hi!"}}}},
				{Type: MsgTypeContentBlockDelta, BlockDelta: &BlockDelta{Delta: ContentBlock{Type: ContentTypeTextDelta, TextDeltaBlock: &TextDeltaBlock{Type: "text", Text: " I"}}}},
				{Type: MsgTypeContentBlockDelta, BlockDelta: &BlockDelta{Delta: ContentBlock{Type: ContentTypeTextDelta, TextDeltaBlock: &TextDeltaBlock{Type: "text", Text: " am your"}}}},
				{Type: MsgTypeContentBlockDelta, BlockDelta: &BlockDelta{Delta: ContentBlock{Type: ContentTypeTextDelta, TextDeltaBlock: &TextDeltaBlock{Type: "text", Text: " dad."}}}},
			},
			want: ClaudeMessage{
				Role: chatbot.ChatMessageRoleAssistant,
				Content: []ContentBlock{
					{Type: ContentTypeText, TextBlock: &TextBlock{Type: ContentTypeText, Text: "Hi! I am your dad."}},
				},
			},
		},
		{
			name: "tool use message", pm: StreamMessages{
				{Type: MsgTypeContentBlockStart, BlockStart: &BlockStart{ContentBlock: ContentBlock{Type: ContentTypeToolCall, ToolCallBlock: &ToolCallBlock{Type: "tool_use", ID: "foo123", Name: "add_numbers", Input: map[string]any{}}}}},
				{Type: MsgTypeContentBlockDelta, BlockDelta: &BlockDelta{Delta: ContentBlock{Type: ContentTypeInputJsonDelta, InputJsonDeltaBlock: &InputJsonDeltaBlock{Type: "input_json_delta", PartialJson: `{"`}}}},
				{Type: MsgTypeContentBlockDelta, BlockDelta: &BlockDelta{Delta: ContentBlock{Type: ContentTypeInputJsonDelta, InputJsonDeltaBlock: &InputJsonDeltaBlock{Type: "input_json_delta", PartialJson: `a":`}}}},
				{Type: MsgTypeContentBlockDelta, BlockDelta: &BlockDelta{Delta: ContentBlock{Type: ContentTypeInputJsonDelta, InputJsonDeltaBlock: &InputJsonDeltaBlock{Type: "input_json_delta", PartialJson: `10,`}}}},
				{Type: MsgTypeContentBlockDelta, BlockDelta: &BlockDelta{Delta: ContentBlock{Type: ContentTypeInputJsonDelta, InputJsonDeltaBlock: &InputJsonDeltaBlock{Type: "input_json_delta", PartialJson: `"b": 5}`}}}},
				{Type: MsgTypeContentBlockStop, BlockStop: &BlockStop{}},
			},
			want: ClaudeMessage{
				Role: chatbot.ChatMessageRoleAssistant,
				Content: []ContentBlock{
					{Type: ContentTypeToolCall, ToolCallBlock: &ToolCallBlock{Type: ContentTypeToolCall, ID: "foo123", Name: "add_numbers", Input: map[string]any{"a": float64(10), "b": float64(5)}}},
				},
			},
		},

		{
			name: "tool use and text message", pm: StreamMessages{
				{Type: MsgTypeContentBlockStart, BlockStart: &BlockStart{Index: 0, ContentBlock: ContentBlock{Type: ContentTypeText, TextBlock: &TextBlock{Type: "text", Text: "Hi!"}}}},
				{Type: MsgTypeContentBlockDelta, BlockDelta: &BlockDelta{Index: 0, Delta: ContentBlock{Type: ContentTypeTextDelta, TextDeltaBlock: &TextDeltaBlock{Type: "text", Text: " I"}}}},
				{Type: MsgTypeContentBlockStart, BlockStart: &BlockStart{Index: 1, ContentBlock: ContentBlock{Type: ContentTypeToolCall, ToolCallBlock: &ToolCallBlock{Type: "tool_use", ID: "foo123", Name: "add_numbers", Input: map[string]any{}}}}},
				{Type: MsgTypeContentBlockDelta, BlockDelta: &BlockDelta{Index: 0, Delta: ContentBlock{Type: ContentTypeTextDelta, TextDeltaBlock: &TextDeltaBlock{Type: "text", Text: " am your"}}}},
				{Type: MsgTypeContentBlockDelta, BlockDelta: &BlockDelta{Index: 1, Delta: ContentBlock{Type: ContentTypeInputJsonDelta, InputJsonDeltaBlock: &InputJsonDeltaBlock{Type: "input_json_delta", PartialJson: `{"`}}}},
				{Type: MsgTypeContentBlockDelta, BlockDelta: &BlockDelta{Index: 1, Delta: ContentBlock{Type: ContentTypeInputJsonDelta, InputJsonDeltaBlock: &InputJsonDeltaBlock{Type: "input_json_delta", PartialJson: `a":`}}}},
				{Type: MsgTypeContentBlockDelta, BlockDelta: &BlockDelta{Index: 0, Delta: ContentBlock{Type: ContentTypeTextDelta, TextDeltaBlock: &TextDeltaBlock{Type: "text", Text: " dad."}}}},
				{Type: MsgTypeContentBlockDelta, BlockDelta: &BlockDelta{Index: 1, Delta: ContentBlock{Type: ContentTypeInputJsonDelta, InputJsonDeltaBlock: &InputJsonDeltaBlock{Type: "input_json_delta", PartialJson: `10,`}}}},
				{Type: MsgTypeContentBlockDelta, BlockDelta: &BlockDelta{Index: 1, Delta: ContentBlock{Type: ContentTypeInputJsonDelta, InputJsonDeltaBlock: &InputJsonDeltaBlock{Type: "input_json_delta", PartialJson: `"b": 5}`}}}},
				{Type: MsgTypeContentBlockStop, BlockStop: &BlockStop{Index: 1}},
			},
			want: ClaudeMessage{
				Role: chatbot.ChatMessageRoleAssistant,
				Content: []ContentBlock{
					{Type: ContentTypeText, TextBlock: &TextBlock{Type: ContentTypeText, Text: "Hi! I am your dad."}},
					{Type: ContentTypeToolCall, ToolCallBlock: &ToolCallBlock{Type: ContentTypeToolCall, ID: "foo123", Name: "add_numbers", Input: map[string]any{"a": float64(10), "b": float64(5)}}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pm.ToFinalMessage()
			assert.Equal(t, tt.want, got)
		})
	}
}
