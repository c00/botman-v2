package yappie

import (
	"testing"

	"github.com/c00/botman-v2/chatbot"
	chattertest "github.com/c00/botman-v2/internal/chattertest"
)

func TestChatterSuite(t *testing.T) {
	chattertest.RunSuite(t, func() chatbot.Chatter {
		return &Yappie{}
	})
}
