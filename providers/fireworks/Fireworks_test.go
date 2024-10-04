package fireworks

import (
	"os"
	"testing"

	"github.com/c00/botman-v2/chatbot"
	chattertest "github.com/c00/botman-v2/internal/chattertest"
	"github.com/stretchr/testify/assert"
)

func TestChatterSuite(t *testing.T) {
	chattertest.RunSuite(t, func() chatbot.Chatter {
		chatter, err := New(Config{
			ApiKey: os.Getenv("FIREWORKS_API_KEY"),
			Model:  "accounts/fireworks/models/llama-v3p1-8b-instruct",
		})
		assert.Nil(t, err)
		return chatter
	})
}
