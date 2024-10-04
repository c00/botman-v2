package mainloop

import (
	"fmt"
	"io"
	"sync"

	"github.com/c00/botman-v2/chatbot"
	"github.com/c00/botman-v2/chattools"
	"github.com/c00/botman-v2/clitools"
	"github.com/c00/botman-v2/internal/history"
	"github.com/c00/botman-v2/internal/storageprovider"
	"github.com/c00/botman-v2/logger"
)

var log = logger.New("MainLoop")

func New(
	chatter chatbot.Chatter,
	histKeeper history.HistoryKeeper,
	storage storageprovider.StorageProvider,
	interactive bool,
	maxRuns int,
	stdIn io.Reader,
	stdOut io.Writer,
) MainLoop {
	return MainLoop{
		Chatter:      chatter,
		interactive:  interactive,
		maxRuns:      maxRuns,
		stdIn:        stdIn,
		stdOut:       stdOut,
		history:      histKeeper,
		conversation: history.NewEntry(),
		storage:      storage,
	}
}

type MainLoop struct {
	Chatter      chatbot.Chatter
	CurrentRun   int
	interactive  bool
	maxRuns      int
	stdIn        io.Reader
	stdOut       io.Writer
	history      history.HistoryKeeper
	storage      storageprovider.StorageProvider
	conversation history.HistoryEntry
	tools        []chattools.ToolDefinition
}

func (l *MainLoop) SetConversation(conv history.HistoryEntry) {
	l.conversation = conv
	l.Chatter.SetMessages(conv.Messages)
}

func (l *MainLoop) getToolDef(name string) (chattools.ToolDefinition, error) {
	for _, td := range l.tools {
		if td.Name == name {
			return td, nil
		}
	}

	return chattools.ToolDefinition{}, fmt.Errorf("tool %v not found", name)
}

func (l *MainLoop) SetTools(tools []chattools.ToolDefinition) {
	l.tools = tools
	l.Chatter.SetTools(tools)
}

func (l *MainLoop) Start(prompt string) error {
	if prompt == "" {
		if l.interactive {
			prompt = clitools.GetInput("You", l.stdIn, l.stdOut)
			if prompt == "" {
				log.Debug("No input from user.")
				return nil
			}
		} else {
			log.Debug("No prompt, not interactive.")
			return nil
		}
	}

	return l.run(chatbot.ChatMessage{Role: chatbot.ChatMessageRoleUser, Content: prompt})
}

// run an interation of the loop
func (l *MainLoop) run(newMsg chatbot.ChatMessage) error {
	l.CurrentRun++
	if l.maxRuns != 0 && l.CurrentRun > l.maxRuns {
		return fmt.Errorf("max runs exceeded")
	}

	l.conversation.Messages = append(l.conversation.Messages, newMsg)

	//Create channel for streaming output
	wg := &sync.WaitGroup{}
	ch := stdOutChannel(wg, l.stdOut)

	// Prompt it
	msg, err := l.Chatter.GetStreamingResponse(newMsg, ch)
	if err != nil {
		return fmt.Errorf("getting streaming response failed: %w", err)
	}

	//Give channels time to flush their last shit to stdout
	wg.Wait()

	fmt.Println("")
	log.Debug("Gotten message Role: %v, ToolCalls: %v, Content: %v", msg.Role, len(msg.ToolCalls), msg.Content)

	l.conversation.Messages = append(l.conversation.Messages, msg)
	_, err = l.history.SaveChat(l.conversation)
	if err != nil {
		return fmt.Errorf("could not save chat: %w", err)
	}

	//Run tool calls
	//todo make goroutines
	toolResults := []chattools.ToolResult{}
	for _, call := range msg.ToolCalls {
		log.Debug("Executing call: %v, Params: %v", call.Name, call.Params)

		def, err := l.getToolDef(call.Name)
		if err != nil {
			log.Error("cannot get tool def: %v", err)
			toolResults = append(toolResults, chattools.ToolResult{ID: call.ID, Name: call.Name, Success: false, Content: err.Error()})
			continue
		}

		toolResult := l.runTool(call, def)
		log.Debug("Tool %v Result: %v", call.Name, toolResult.Content)
		toolResults = append(toolResults, toolResult)
	}

	if len(toolResults) > 0 {
		//if any tool results have been added, restart loop
		return l.run(chatbot.ChatMessage{Role: chatbot.ChatMessageRoleTool, ToolResults: toolResults})
	}

	//If interactive, Ask for input and restart loop
	if l.interactive {
		prompt := clitools.GetInput("You", l.stdIn, l.stdOut)
		if prompt == "" {
			log.Debug("No input from user.")
			return nil
		}

		return l.run(chatbot.ChatMessage{Role: chatbot.ChatMessageRoleUser, Content: prompt})
	}

	return nil
}
