package claude

import (
	"encoding/json"

	"github.com/c00/botman-v2/chatbot"
)

const (
	MsgTypeMessageStart      = "message_start"
	MsgTypeMessageDelta      = "message_delta"
	MsgTypeMessageStop       = "message_stop"
	MsgTypeContentBlockStart = "content_block_start"
	MsgTypeContentBlockDelta = "content_block_delta"
	MsgTypeContentBlockStop  = "content_block_stop"
	MsgTypePing              = "ping"
)

type StreamMessage struct {
	Type string `json:"type"`

	MessageStart *MessageStart
	MessageDelta *MessageDelta
	BlockStart   *BlockStart
	BlockStop    *BlockStop
	BlockDelta   *BlockDelta
}

func (b StreamMessage) MarshalJSON() ([]byte, error) {
	switch b.Type {
	case MsgTypeMessageStart:
		return json.Marshal(b.MessageStart)
	case MsgTypeMessageDelta:
		return json.Marshal(b.MessageDelta)
	case MsgTypeContentBlockStart:
		return json.Marshal(b.BlockStart)
	case MsgTypeContentBlockDelta:
		return json.Marshal(b.BlockDelta)
	case MsgTypeContentBlockStop:
		return json.Marshal(b.BlockStop)
	default:
		return json.Marshal(map[string]any{"type": b.Type})
	}
}

func (b *StreamMessage) UnmarshalJSON(raw []byte) error {
	var data map[string]any
	json.Unmarshal(raw, &data)
	if val, ok := data["type"].(string); ok {
		b.Type = val
	}

	switch b.Type {
	case MsgTypeMessageStart:
		b.MessageStart = &MessageStart{}
		json.Unmarshal(raw, b.MessageStart)
	case MsgTypeMessageDelta:
		b.MessageDelta = &MessageDelta{}
		json.Unmarshal(raw, b.MessageDelta)
	case MsgTypeContentBlockStart:
		b.BlockStart = &BlockStart{}
		json.Unmarshal(raw, b.BlockStart)
	case MsgTypeContentBlockDelta:
		b.BlockDelta = &BlockDelta{}
		json.Unmarshal(raw, b.BlockDelta)
	case MsgTypeContentBlockStop:
		b.BlockStop = &BlockStop{}
		json.Unmarshal(raw, b.BlockStop)
	}
	return nil
}

func (m StreamMessage) Delta() string {
	switch m.Type {
	case MsgTypeContentBlockStart:
		return m.BlockStart.ContentBlock.Delta()
	case MsgTypeContentBlockDelta:
		return m.BlockDelta.Delta.Delta()
	}
	return ""
}

type StreamMessages []StreamMessage

// todo test this
func (pm StreamMessages) ToFinalMessage() ClaudeMessage {
	//I'm assuming that the indexes in Deltas come in the right order.
	//This may be a bad assumption.
	content := []ContentBlock{}

	for _, msg := range pm {
		switch msg.Type {
		case MsgTypeContentBlockStart:
			//This is a new block.
			if msg.BlockStart.Index != len(content) {
				//This is fucked.
				log.Error("[StreamMessages.ToFinalMessage] Block Start index out of sync: Index gotten: %v. Length: %v", msg.BlockStart.Index, len(content))
				panic("[StreamMessages.ToFinalMessage] index out of sync")
			}
			content = append(content, msg.BlockStart.ContentBlock)
		case MsgTypeContentBlockDelta:
			if msg.BlockDelta.Index >= len(content) {
				log.Error("[StreamMessages.ToFinalMessage] Block Delta index out of sync: Index gotten: %v. Length: %v", msg.BlockStart.Index, len(content))
				panic("[StreamMessages.ToFinalMessage] index out of sync")
			}
			content[msg.BlockDelta.Index].Add(msg.BlockDelta.Delta)
		case MsgTypeContentBlockStop:
			if msg.BlockStop.Index >= len(content) {
				log.Error("[StreamMessages.ToFinalMessage] Block Stop index out of sync: Index gotten: %v. Length: %v", msg.BlockStart.Index, len(content))
				panic("[StreamMessages.ToFinalMessage] index out of sync")
			}
			content[msg.BlockStop.Index].Finalize()
		}
	}

	return ClaudeMessage{
		Role:    chatbot.ChatMessageRoleAssistant,
		Content: content,
	}
}

type MessageStart struct {
	Type    string         `json:"type"`
	Message MessageContent `json:"message"`
}

type MessageDelta struct {
	Type  string         `json:"type"`
	Delta MessageContent `json:"delta"`
	Usage Usage          `json:"usage"`
}

type MessageContent struct {
	ID           string         `json:"id"`
	Type         string         `json:"type"`
	Role         string         `json:"role"`
	Model        string         `json:"model"`
	Content      []ContentBlock `json:"content"`
	StopReason   string         `json:"stop_reason"`
	StopSequence string         `json:"stop_sequence"`
	Usage        Usage          `json:"usage"`
}

type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type BlockStart struct {
	Type         string       `json:"type"`
	Index        int          `json:"index"`
	ContentBlock ContentBlock `json:"content_block"`
}

type BlockStop struct {
	Type  string `json:"type"`
	Index int    `json:"index"`
}

type BlockDelta struct {
	Type  string       `json:"type"`
	Index int          `json:"index"`
	Delta ContentBlock `json:"delta"`
}
