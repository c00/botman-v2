package claude

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/c00/botman-v2/chatbot"
	"github.com/c00/botman-v2/internal/channeltools"
	"github.com/stretchr/testify/assert"
)

func TestConsumeChatStream(t *testing.T) {
	ch := channeltools.BlackHoleChannel()
	reader := bufio.NewReader(bytes.NewReader([]byte(normalStream)))
	message, err := consumeStream(reader, ch)
	assert.Nil(t, err)

	assert.Len(t, message.Content, 1)
	assert.Equal(t, "Hi!", message.Content[0].TextBlock.Text)
}

func TestConsumeToolStream(t *testing.T) {
	ch := channeltools.BlackHoleChannel()
	reader := bufio.NewReader(bytes.NewReader([]byte(toolUseStream)))
	message, err := consumeStream(reader, ch)
	assert.Nil(t, err)
	assert.Len(t, message.Content, 1)

	expected := ContentBlock{
		Type: ContentTypeToolCall,
		ToolCallBlock: &ToolCallBlock{
			Type:  ContentTypeToolCall,
			ID:    "toolu_01DtNpfALyP25sbvmY4KpJGf",
			Name:  "add_numbers",
			Input: map[string]any{"a": float64(10), "b": float64(10)},
		},
	}

	assert.Equal(t, expected, message.Content[0])
}

func TestConsumeMultiToolStream(t *testing.T) {
	ch := channeltools.BlackHoleChannel()
	reader := bufio.NewReader(bytes.NewReader([]byte(multiToolStream)))
	message, err := consumeStream(reader, ch)
	assert.Nil(t, err)
	assert.Len(t, message.Content, 3)

	expected := ClaudeMessage{
		Role: chatbot.ChatMessageRoleAssistant,
		Content: []ContentBlock{
			{
				Type: ContentTypeText,
				TextBlock: &TextBlock{
					Type: ContentTypeText,
					Text: "Certainly! I'll make two function calls to add the numbers as you've requested, and I'll do both calls at the same time.",
				},
			},
			{
				Type: ContentTypeToolCall,
				ToolCallBlock: &ToolCallBlock{
					Type:  ContentTypeToolCall,
					ID:    "toolu_01TzG4Z4zyvQ2fPLAKm6P5dt",
					Name:  "add_numbers",
					Input: map[string]any{"number_a": float64(3), "number_b": float64(4)},
				},
			},
			{
				Type: ContentTypeToolCall,
				ToolCallBlock: &ToolCallBlock{
					Type:  ContentTypeToolCall,
					ID:    "toolu_01Tg7CyaMV8pGSF6rtL2GFpE",
					Name:  "add_numbers",
					Input: map[string]any{"number_a": float64(8), "number_b": float64(9)},
				},
			},
		},
	}

	assert.Equal(t, expected, message)
}

const normalStream = `event: message_start
data: {"type":"message_start","message":{"id":"msg_011C196hjcEFExDoxPpD5zXV","type":"message","role":"assistant","model":"claude-3-haiku-20240307","content":[],"stop_reason":null,"stop_sequence":null,"usage":{"input_tokens":17,"output_tokens":1}}}

event: content_block_start
data: {"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}      }

event: ping
data: {"type": "ping"}

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Hi"}            }

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"!"}          }

event: content_block_stop
data: {"type":"content_block_stop","index":0          }

event: message_delta
data: {"type":"message_delta","delta":{"stop_reason":"end_turn","stop_sequence":null},"usage":{"output_tokens":5}  }

event: message_stop
data: {"type":"message_stop"  }

`

const toolUseStream = `event: message_start
data: {"type":"message_start","message":{"id":"msg_01R7Dp13khRMm2QSZp3uMAyS","type":"message","role":"assistant","model":"claude-3-haiku-20240307","content":[],"stop_reason":null,"stop_sequence":null,"usage":{"input_tokens":375,"output_tokens":4}}     }

event: content_block_start
data: {"type":"content_block_start","index":0,"content_block":{"type":"tool_use","id":"toolu_01DtNpfALyP25sbvmY4KpJGf","name":"add_numbers","input":{}}            }

event: ping
data: {"type": "ping"}

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"input_json_delta","partial_json":""}          }

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"input_json_delta","partial_json":"{\""}}

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"input_json_delta","partial_json":"a\": 10"}       }

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"input_json_delta","partial_json":", \"b\""}       }

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"input_json_delta","partial_json":": 10}"}               }

event: content_block_stop
data: {"type":"content_block_stop","index":0  }

event: message_delta
data: {"type":"message_delta","delta":{"stop_reason":"tool_use","stop_sequence":null},"usage":{"output_tokens":70}   }

event: message_stop
data: {"type":"message_stop" }

`

const multiToolStream = `event: message_start
data: {"type":"message_start","message":{"id":"msg_013GYyFZMyiCQPovnsbsbhwt","type":"message","role":"assistant","model":"claude-3-5-sonnet-20240620","content":[],"stop_reason":null,"stop_sequence":null,"usage":{"input_tokens":430,"output_tokens":1}}            }

event: content_block_start
data: {"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}               }

event: ping
data: {"type": "ping"}

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Certainly"} }

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"! I"}         }

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"'ll"}}

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":" make two"}              }

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":" function calls to"}               }

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":" add the"}    }

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":" numbers as"}       }

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":" you've requeste"}     }

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"d,"}         }

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":" and I"}    }

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"'ll do both"}}

event: content_block_delta
data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":" calls at the same time."}}

event: content_block_stop
data: {"type":"content_block_stop","index":0      }

event: content_block_start
data: {"type":"content_block_start","index":1,"content_block":{"type":"tool_use","id":"toolu_01TzG4Z4zyvQ2fPLAKm6P5dt","name":"add_numbers","input":{}}         }

event: content_block_delta
data: {"type":"content_block_delta","index":1,"delta":{"type":"input_json_delta","partial_json":""}    }

event: content_block_delta
data: {"type":"content_block_delta","index":1,"delta":{"type":"input_json_delta","partial_json":"{\"number_a\":"}               }

event: content_block_delta
data: {"type":"content_block_delta","index":1,"delta":{"type":"input_json_delta","partial_json":" 3"}     }

event: content_block_delta
data: {"type":"content_block_delta","index":1,"delta":{"type":"input_json_delta","partial_json":", "}}

event: content_block_delta
data: {"type":"content_block_delta","index":1,"delta":{"type":"input_json_delta","partial_json":"\"n"}            }

event: content_block_delta
data: {"type":"content_block_delta","index":1,"delta":{"type":"input_json_delta","partial_json":"umber_b\":"}      }

event: content_block_delta
data: {"type":"content_block_delta","index":1,"delta":{"type":"input_json_delta","partial_json":" 4}"}             }

event: content_block_stop
data: {"type":"content_block_stop","index":1             }

event: content_block_start
data: {"type":"content_block_start","index":2,"content_block":{"type":"tool_use","id":"toolu_01Tg7CyaMV8pGSF6rtL2GFpE","name":"add_numbers","input":{}}       }

event: content_block_delta
data: {"type":"content_block_delta","index":2,"delta":{"type":"input_json_delta","partial_json":""}      }

event: content_block_delta
data: {"type":"content_block_delta","index":2,"delta":{"type":"input_json_delta","partial_json":"{\"number_a\""}           }

event: content_block_delta
data: {"type":"content_block_delta","index":2,"delta":{"type":"input_json_delta","partial_json":": 8"}      }

event: content_block_delta
data: {"type":"content_block_delta","index":2,"delta":{"type":"input_json_delta","partial_json":", \"number"}       }

event: content_block_delta
data: {"type":"content_block_delta","index":2,"delta":{"type":"input_json_delta","partial_json":"_b\":"}             }

event: content_block_delta
data: {"type":"content_block_delta","index":2,"delta":{"type":"input_json_delta","partial_json":" 9}"}}

event: content_block_stop
data: {"type":"content_block_stop","index":2       }

event: message_delta
data: {"type":"message_delta","delta":{"stop_reason":"tool_use","stop_sequence":null},"usage":{"output_tokens":180}          }

event: message_stop
data: {"type":"message_stop"           }
`
