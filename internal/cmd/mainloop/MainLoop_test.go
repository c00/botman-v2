package mainloop

import (
	"io"
	"testing"

	"github.com/c00/botman-v2/chatbot"
	"github.com/c00/botman-v2/chattools"
	"github.com/c00/botman-v2/internal/history"
	"github.com/c00/botman-v2/internal/storageprovider"
	"github.com/c00/botman-v2/providers/yappie"
	"github.com/stretchr/testify/assert"
)

func TestMainLoopNonInteractive_Run(t *testing.T) {
	chatter := &yappie.Yappie{}
	userInput := &stringReader{}
	output := &stringWriter{}
	hist := &history.InMemoryHistory{}
	store := storageprovider.NewMemStore()

	ml := New(chatter, hist, store, false, 0, userInput, output)

	err := ml.Start("hey")
	assert.Nil(t, err)
	assert.Equal(t, 1, ml.CurrentRun)
	assert.Len(t, ml.Chatter.GetMessages(), 2)

	list, err := hist.List()
	assert.Nil(t, err)
	assert.Len(t, list, 1)
}

func TestMainLoopNonInteractiveNoPrompt_Run(t *testing.T) {
	chatter := &yappie.Yappie{}
	userInput := &stringReader{}
	output := &stringWriter{}
	hist := &history.InMemoryHistory{}
	store := storageprovider.NewMemStore()

	ml := New(chatter, hist, store, false, 0, userInput, output)

	err := ml.Start("")
	assert.Nil(t, err)
	assert.Equal(t, 0, ml.CurrentRun)
	assert.Len(t, ml.Chatter.GetMessages(), 0)

	list, err := hist.List()
	assert.Nil(t, err)
	assert.Len(t, list, 0)
}

func TestMainLoopInteractive_Run(t *testing.T) {
	chatter := &yappie.Yappie{}
	userInput := &stringReader{}
	output := &stringWriter{}

	userInput.Add("one answer")
	userInput.Add("second answer")
	hist := &history.InMemoryHistory{}
	store := storageprovider.NewMemStore()

	ml := New(chatter, hist, store, true, 0, userInput, output)

	err := ml.Start("hey")
	assert.Nil(t, err)
	assert.Equal(t, 3, ml.CurrentRun)
	assert.Len(t, ml.Chatter.GetMessages(), 6)

	list, err := hist.List()
	assert.Nil(t, err)
	assert.Len(t, list, 1)
}

func TestMainLooptools_Run(t *testing.T) {
	chatter := &yappie.Yappie{}
	userInput := &stringReader{}
	output := &stringWriter{}
	hist := &history.InMemoryHistory{}
	store := storageprovider.NewMemStore()

	ml := New(chatter, hist, store, false, 0, userInput, output)
	ml.SetTools([]chattools.ToolDefinition{
		{ToolType: chattools.ToolTypeAddNumbers, Name: "add_numbers", Description: "Add two numbers"},
	})

	err := ml.Start("hey")
	assert.Nil(t, err)
	assert.Equal(t, 2, ml.CurrentRun)
	assert.Len(t, ml.Chatter.GetMessages(), 4)

	toolMessage := ml.Chatter.GetMessages()[2]
	assert.Equal(t, toolMessage.Role, chatbot.ChatMessageRoleTool)
}

type stringReader struct {
	data []string
	pos  int
}

func (sr *stringReader) Read(p []byte) (n int, err error) {
	if sr.pos >= len(sr.data) {
		return 0, io.EOF
	}
	n = copy(p, sr.data[sr.pos])
	sr.pos++
	return n, nil
}

func (sr *stringReader) Add(s string) {
	sr.data = append(sr.data, s+"\n")
}

type stringWriter struct {
	io.Writer
	data []byte
}

func (sw *stringWriter) Write(newData []byte) (n int, err error) {
	sw.data = append(sw.data, newData...)
	return len(newData), nil
}

func (sw *stringWriter) String() string {
	return string(sw.data)
}

func (sw *stringWriter) StringFlush() string {
	data := string(sw.data)
	sw.data = []byte{}
	return data
}
