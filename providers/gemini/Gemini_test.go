package gemini

import (
	"os"
	"testing"

	"github.com/c00/botman-v2/chatbot"
	chattertest "github.com/c00/botman-v2/internal/chattertest"
	"github.com/c00/botman-v2/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestChatterSuite(t *testing.T) {
	logger.SetLevel(5)
	chattertest.RunSuite(t, func() chatbot.Chatter {
		chatter, err := New(Config{
			ApiKey: os.Getenv("GEMINI_API_KEY"),
			Model:  "gemini-1.5-flash-8b",
		})
		assert.Nil(t, err)
		return chatter
	})
}
