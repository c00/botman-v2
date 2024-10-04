package mainloop

import (
	"fmt"

	"github.com/c00/botman-v2/chattools"
	"github.com/c00/botman-v2/chattools/sdxl"
	"github.com/c00/botman-v2/internal/storageprovider"
)

func (ml *MainLoop) runTool(call chattools.ToolCall, def chattools.ToolDefinition) chattools.ToolResult {
	result := chattools.ToolResult{
		Content: fmt.Sprintf("tool %v not implemented", def.ToolType),
		Success: false,
	}

	switch def.ToolType {
	case chattools.ToolTypeAddNumbers:
		result = runAddNumbers(call)
	case chattools.ToolTypeSdxl:
		result = runSdxlTool(call, def, ml.storage)
	}

	//Set this here to be sure they get set.
	result.ID = call.ID
	result.Name = call.Name

	return result
}

func runSdxlTool(call chattools.ToolCall, def chattools.ToolDefinition, store storageprovider.StorageProvider) chattools.ToolResult {
	result := chattools.ToolResult{}

	//todo get the storage provider here.
	// store, err := storageprovider.NewLocalStore(filepath.Join(config.GetUserConfigPath(), "generated-images"))
	// if err != nil {
	// 	result.Content = fmt.Sprintf("cannot create storage provider: %v", err)
	// 	return result
	// }
	tool, err := sdxl.New(def, store)
	if err != nil {
		result.Content = fmt.Sprintf("cannot create new sdxl tool: %v", err)
		return result
	}

	return tool.Run(call)
}

func getInt(params map[string]any, key string) int {
	switch val := params[key].(type) {
	case nil:
		return 0
	case int:
		return val
	case int32:
		return int(val)
	case int64:
		return int(val)
	case float32:
		return int(val)
	case float64:
		return int(val)
	}

	return 0
}

// func getString(params map[string]any, key string) string {
// 	if val, ok := params[key]; ok {
// 		if num, ok := val.(string); ok {
// 			return num
// 		}
// 	}
// 	return ""
// }
