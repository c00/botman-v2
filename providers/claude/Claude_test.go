package claude

import (
	"os"
	"testing"

	"github.com/c00/botman-v2/chatbot"
	chattertest "github.com/c00/botman-v2/internal/chattertest"
	"github.com/c00/botman-v2/logger"
	"github.com/stretchr/testify/assert"
)

func TestChatterSuite(t *testing.T) {
	logger.SetLevel(5)
	chattertest.RunSuite(t, func() chatbot.Chatter {
		chatter, err := New(Config{
			ApiKey:    os.Getenv("CLAUDE_API_KEY"),
			Model:     "claude-3-haiku-20240307",
			MaxTokens: 100,
		})
		assert.Nil(t, err)
		return chatter
	})
}
