package mainloop

import (
	"reflect"
	"testing"

	"github.com/c00/botman-v2/chattools"
)

// var def = chattools.ToolDefinition{
// 	ToolType: chattools.ToolTypeAddNumbers,
// }

func Test_runAddNumbers(t *testing.T) {
	tests := []struct {
		name string
		call chattools.ToolCall
		want chattools.ToolResult
	}{
		{
			name: "happy path int",
			call: chattools.ToolCall{Name: "Foo", Params: map[string]any{"a": 5, "b": 10}},
			want: chattools.ToolResult{Name: "Foo", Content: "15", Value: 15, Success: true},
		},
		{
			name: "happy path floats",
			call: chattools.ToolCall{Name: "Foo", Params: map[string]any{"a": float64(5), "b": float64(10)}},
			want: chattools.ToolResult{Name: "Foo", Content: "15", Value: 15, Success: true},
		},
		{
			name: "missing fields",
			call: chattools.ToolCall{Name: "Foo", Params: map[string]any{}},
			want: chattools.ToolResult{Name: "Foo", Content: "0", Value: 0, Success: true},
		},
		{
			name: "wrong type fields",
			call: chattools.ToolCall{Name: "Foo", Params: map[string]any{"a": true, "b": "Potato"}},
			want: chattools.ToolResult{Name: "Foo", Content: "0", Value: 0, Success: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := runAddNumbers(tt.call); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("runAddNumbers() = %v, want %v", got, tt.want)
			}
		})
	}
}
