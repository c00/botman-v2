package sdxl

import (
	"os"
	"testing"

	"github.com/c00/botman-v2/chattools"
	"github.com/c00/botman-v2/internal/storageprovider"
	"github.com/stretchr/testify/assert"
)

func TestSdxlTool_getImage(t *testing.T) {
	tool, err := New(chattools.ToolDefinition{
		ToolType: chattools.ToolTypeSdxl,
		SdxlSettings: &chattools.SdxlConfig{
			ApiKey:         os.Getenv("FIREWORKS_API_KEY"),
			PositivePrompt: "Dragonball z anime style.",
		},
	}, nil)
	assert.Nil(t, err)

	data, err := tool.getImage("an image of a cat with too much fur", "")
	assert.Nil(t, err)
	assert.NotNil(t, data)
}

func TestSdxlTool_Run(t *testing.T) {
	store, err := storageprovider.NewLocalStore(os.TempDir())
	assert.Nil(t, err)

	tool, err := New(
		chattools.ToolDefinition{
			ToolType: chattools.ToolTypeSdxl,
			SdxlSettings: &chattools.SdxlConfig{
				ApiKey:         os.Getenv("FIREWORKS_API_KEY"),
				PositivePrompt: "Dragonball z anime style.",
			},
		},
		store,
	)
	assert.Nil(t, err)

	call := chattools.ToolCall{
		Params: map[string]any{"prompt": "a really fluffy kitten"},
	}
	msg := tool.Run(call)
	assert.True(t, msg.Success)

	val, ok := msg.Value.(SdxlResult)
	assert.True(t, ok)
	//Check the path to see the test image. if you feel like it.
	assert.NotEqual(t, "", val.Path)
}
