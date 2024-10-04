package claude

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContentBlock_JSON(t *testing.T) {
	//This tests marshalls and unmarshalls a claudeMessage and checks that it ends up the same.
	tests := []struct {
		name string
		data ContentBlock
	}{
		{
			name: "text content", data: ContentBlock{
				Type:      ContentTypeText,
				TextBlock: &TextBlock{Type: ContentTypeText, Text: "Hello!"},
			},
		},
		{
			name: "Tool Use / Call content", data: ContentBlock{
				Type:          ContentTypeToolCall,
				ToolCallBlock: &ToolCallBlock{Type: ContentTypeToolCall, ID: "tool_1234", Name: "add_numbers", Input: map[string]any{"foo": "bar"}},
			},
		},
		{
			name: "Tool Result content", data: ContentBlock{
				Type:            ContentTypeToolResult,
				ToolResultBlock: &ToolResultBlock{Type: ContentTypeToolResult, ToolUseId: "tool_1234", Content: "8", IsError: true},
			},
		},
		{
			name: "Text Delta content", data: ContentBlock{
				Type:           ContentTypeTextDelta,
				TextDeltaBlock: &TextDeltaBlock{Type: ContentTypeTextDelta, Text: "!"},
			},
		},
		{
			name: "Input Json Delta content", data: ContentBlock{
				Type:                ContentTypeInputJsonDelta,
				InputJsonDeltaBlock: &InputJsonDeltaBlock{Type: ContentTypeInputJsonDelta, PartialJson: `{"`},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := tt.data
			got, err := json.Marshal(original)
			assert.Nil(t, err)

			unmarshalled := ContentBlock{}
			err = json.Unmarshal(got, &unmarshalled)
			assert.Nil(t, err)

			assert.Equal(t, original, unmarshalled)
		})
	}
}
