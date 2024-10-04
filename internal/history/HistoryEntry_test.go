package history

import (
	"testing"
	"time"

	"github.com/c00/botman-v2/chatbot"
	"github.com/stretchr/testify/assert"
)

func RunSuite(t *testing.T, keeperFactory func() HistoryKeeper) {
	saveAndLoad(t, keeperFactory())
}

func saveAndLoad(t *testing.T, keeper HistoryKeeper) {
	list, err := keeper.List()
	assert.Nil(t, err)
	assert.Len(t, list, 0)

	entry := HistoryEntry{
		Name: "entry 1",
		Date: mustParse("2024-01-01T00:00:00Z"),
		Messages: []chatbot.ChatMessage{
			{Role: "system", Content: "You like big butts and you cannot lie"},
			{Role: "user", Content: "what do you like?"},
			{Role: "assistant", Content: "Big butts"},
		},
	}

	entry, err = keeper.SaveChat(entry)
	assert.Nil(t, err)
	assert.Len(t, entry.Messages, 3)

	list, err = keeper.List()
	assert.Nil(t, err)
	assert.Len(t, list, 1)

	got, err := keeper.LoadChat(0)
	assert.Nil(t, err)
	assert.Equal(t, got, entry)

	entry, err = keeper.SaveChat(HistoryEntry{
		Name: "entry 2",
		Date: mustParse("2024-01-02T00:00:00Z"),
		Messages: []chatbot.ChatMessage{
			{Role: "system", Content: "Oh Canada"},
		},
	})

	assert.Nil(t, err)
	assert.Len(t, entry.Messages, 1)

	list, err = keeper.List()
	assert.Nil(t, err)
	assert.Len(t, list, 2)

	entry, err = keeper.LoadChat(0)
	assert.Nil(t, err)
	assert.Len(t, entry.Messages, 1)

	entry, err = keeper.LoadChat(1)
	assert.Nil(t, err)
	assert.Len(t, entry.Messages, 3)

	//Update existing
	entry, err = keeper.SaveChat(HistoryEntry{
		Name: "entry 2",
		Date: mustParse("2024-01-02T00:00:00Z"),
		Messages: []chatbot.ChatMessage{
			{Role: "system", Content: "Oh Canada"},
			{Role: "user", Content: "I guess"},
		},
	})
	assert.Nil(t, err)
	assert.Len(t, entry.Messages, 2)

	got, err = keeper.LoadChat(0)
	assert.Nil(t, err)
	assert.Equal(t, got, entry)

	list, err = keeper.List()
	assert.Nil(t, err)
	assert.Len(t, list, 2)
}

func mustParse(input string) time.Time {
	date, err := time.Parse(time.RFC3339, input)
	if err != nil {
		panic("cannot parse date")
	}
	return date
}
