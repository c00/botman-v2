package chattertest

import (
	"slices"
	"sync"
	"testing"

	"github.com/c00/botman-v2/chatbot"
	"github.com/c00/botman-v2/chattools"
	"github.com/stretchr/testify/assert"
)

func RunSuite(t *testing.T, chatterFactory func() chatbot.Chatter) {
	setAndGetSystemPrompt(t, chatterFactory())
	setAndGetMessages(t, chatterFactory())
	getResponse(t, chatterFactory())
	getStreamingResponse(t, chatterFactory())
	toolCalls(t, chatterFactory())
	toolResults(t, chatterFactory())
}

func toolResults(t *testing.T, chatter chatbot.Chatter) {
	if !slices.Contains(chatter.SupportedFeatures(), "tools") {
		//tool use not supported.
		return
	}

	toolDef := chattools.ToolDefinition{
		ToolType:    chattools.ToolTypeAddNumbers,
		Name:        "add_numbers",
		Description: "Add 2 numbers together",
	}
	chatter.SetTools([]chattools.ToolDefinition{toolDef})

	chatter.AddMessages([]chatbot.ChatMessage{
		{Role: chatbot.ChatMessageRoleUser, Content: "Add the numbers 2 and 3 together. say 'Great!' when you're done."},
		{Role: chatbot.ChatMessageRoleAssistant, ToolCalls: []chattools.ToolCall{
			{ID: "tool_1234", Name: "add_numbers", Params: map[string]any{"a": 2, "b": 3}},
		}},
	})

	response, err := chatter.GetResponse(chatbot.ChatMessage{Role: chatbot.ChatMessageRoleTool, ToolResults: []chattools.ToolResult{
		{ID: "tool_1234", Name: "add_numbers", Content: "5", Success: true, Value: 5},
	}})

	assert.Nil(t, err)
	assert.NotEqual(t, response.Content, "")
}

func toolCalls(t *testing.T, chatter chatbot.Chatter) {
	if !slices.Contains(chatter.SupportedFeatures(), "tools") {
		//tool use not supported.
		return
	}

	toolDef := chattools.ToolDefinition{
		ToolType:    chattools.ToolTypeAddNumbers,
		Name:        "add_numbers",
		Description: "Add 2 numbers together",
	}
	chatter.SetTools([]chattools.ToolDefinition{toolDef})

	response, err := chatter.GetResponse(chatbot.ChatMessage{Role: chatbot.ChatMessageRoleUser, Content: "Add the numbers 10 and 10 together. Use the add_numbers tool for this."})
	assert.Nil(t, err)
	assert.Len(t, response.ToolCalls, 1)

	call := response.ToolCalls[0]
	a, ok := call.Params["a"]
	assert.True(t, ok)
	assert.Equal(t, a, float64(10))

	b, ok := call.Params["b"]
	assert.True(t, ok)
	assert.Equal(t, b, float64(10))
}

func setAndGetSystemPrompt(t *testing.T, chatter chatbot.Chatter) {
	prompt := "You like big butts and you cannot lie"
	chatter.SetSystemPrompt(prompt)
	assert.Equal(t, chatter.GetSystemPrompt(), prompt)
}

func setAndGetMessages(t *testing.T, chatter chatbot.Chatter) {
	assert.Len(t, chatter.GetMessages(), 0)

	chatter.AddMessages([]chatbot.ChatMessage{
		{Role: chatbot.ChatMessageRoleUser, Content: "Just say hi"},
		{Role: chatbot.ChatMessageRoleAssistant, Content: "hi"},
	})

	assert.Len(t, chatter.GetMessages(), 2)

	chatter.AddMessages([]chatbot.ChatMessage{
		{Role: chatbot.ChatMessageRoleUser, Content: "Just say hi"},
		{Role: chatbot.ChatMessageRoleAssistant, Content: "hi"},
	})

	assert.Len(t, chatter.GetMessages(), 4)
}

func getResponse(t *testing.T, chatter chatbot.Chatter) {
	chatter.SetSystemPrompt("You are a helpful chatbot")
	assert.Len(t, chatter.GetMessages(), 0)

	message, err := chatter.GetResponse(chatbot.ChatMessage{
		Role:    chatbot.ChatMessageRoleUser,
		Content: "Just say hi.",
	})

	assert.Nil(t, err)
	assert.NotEmpty(t, message.Content)
	assert.Equal(t, message.Role, chatbot.ChatMessageRoleAssistant)

	assert.Len(t, chatter.GetMessages(), 2)
}

func getStreamingResponse(t *testing.T, chatter chatbot.Chatter) {
	chatter.SetSystemPrompt("You are a helpful chatbot")
	assert.Len(t, chatter.GetMessages(), 0)
	ch := make(chan string)
	count := 0

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(ch chan string) {
		for range ch {
			count++
		}
		wg.Done()
	}(ch)

	message, err := chatter.GetStreamingResponse(chatbot.ChatMessage{
		Role:    chatbot.ChatMessageRoleUser,
		Content: "Just say hi.",
	}, ch)

	//Let the goroutines settle.
	wg.Wait()

	assert.Nil(t, err)
	assert.NotEmpty(t, message.Content)
	assert.Equal(t, message.Role, chatbot.ChatMessageRoleAssistant)
	assert.Greater(t, count, 0)
	assert.Len(t, chatter.GetMessages(), 2)

	//Check that channel is closed.
	_, ok := <-ch
	assert.False(t, ok)
}
