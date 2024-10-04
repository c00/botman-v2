package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/c00/botman-v2/chatbot"
	"github.com/c00/botman-v2/internal/config"
	"github.com/c00/botman-v2/providers/claude"
	"github.com/c00/botman-v2/providers/fireworks"
	"github.com/c00/botman-v2/providers/openai"
	"github.com/c00/botman-v2/providers/yappie"
)

func getChatter(conf config.BotmanConfig) (chatbot.Chatter, error) {
	//Add the current date and time.
	conf.SystemPrompt = fmt.Sprintf("The current date and time is %v. %v", time.Now().Format(time.RFC1123Z), conf.SystemPrompt)

	switch conf.LlmProvider {
	case config.LlmProviderClaude:
		conf.Claude.SystemPrompt = strings.TrimSpace(fmt.Sprintf("%v %v", conf.SystemPrompt, conf.Claude.SystemPrompt))
		return claude.New(conf.Claude)
	case config.LlmProviderFireworksAi:
		conf.FireworksAi.SystemPrompt = strings.TrimSpace(fmt.Sprintf("%v %v", conf.SystemPrompt, conf.FireworksAi.SystemPrompt))
		return fireworks.New(conf.FireworksAi)
	case config.LlmProviderOpenAi:
		conf.OpenAi.SystemPrompt = strings.TrimSpace(fmt.Sprintf("%v %v", conf.SystemPrompt, conf.OpenAi.SystemPrompt))
		return openai.New(conf.OpenAi)
	case config.LlmProviderYappie:
		return &yappie.Yappie{SystemPrompt: conf.SystemPrompt}, nil
	}

	return nil, fmt.Errorf("unknown llm provider: %v", conf.LlmProvider)
}
