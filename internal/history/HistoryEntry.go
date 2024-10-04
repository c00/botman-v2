package history

import (
	"fmt"
	"time"

	"github.com/c00/botman-v2/chatbot"
)

func NewEntry() HistoryEntry {
	now := time.Now()

	return HistoryEntry{
		Name: now.Format(time.RFC3339),
		Date: now,
	}
}

type HistoryEntry struct {
	Name     string                `yaml:"name"`
	Date     time.Time             `yaml:"date"`
	Messages []chatbot.ChatMessage `yaml:"messages"`
}

func (e HistoryEntry) Print() {
	for _, message := range e.Messages {
		if message.Role == "system" {
			continue
		} else if message.Role == "assistant" {
			fmt.Print(message.Content, "\n\n")
		} else if message.Role == "user" {
			fmt.Print("You: ", message.Content, "\n\n")
		} else {
			fmt.Printf("%v: %v\n", message.Role, message.Content)
		}
	}
}

func (e HistoryEntry) PrintLastMessage() {
	if len(e.Messages) == 0 {
		return
	}

	message := e.Messages[len(e.Messages)-1]

	if message.Role == "assistant" {
		fmt.Print(message.Content, "\n")
	} else {
		fmt.Printf("%v: %v\n", message.Role, message.Content)
	}
}
