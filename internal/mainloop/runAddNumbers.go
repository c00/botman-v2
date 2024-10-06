package mainloop

import (
	"strconv"

	"github.com/c00/botman-v2/chattools"
)

func runAddNumbers(call chattools.ToolCall) chattools.ToolResult {
	a := getInt(call.Params, "a")
	b := getInt(call.Params, "b")

	result := a + b

	return chattools.ToolResult{
		Name:    call.Name,
		ID:      call.ID,
		Content: strconv.Itoa(result),
		Success: true,
		Value:   result,
	}
}
