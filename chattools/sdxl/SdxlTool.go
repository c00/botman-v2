package sdxl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/c00/botman-v2/chattools"
	"github.com/c00/botman-v2/internal/storageprovider"
	"github.com/c00/botman-v2/logger"
)

var log = logger.New("SdxlTool")

const defaultModel = "stable-diffusion-xl-1024-v1-0"

func New(def chattools.ToolDefinition, store storageprovider.StorageProvider) (*SdxlTool, error) {
	if def.SdxlSettings == nil {
		return nil, fmt.Errorf("missing sdxl settings in tool definition")
	}

	if def.SdxlSettings.Model == "" {
		def.SdxlSettings.Model = defaultModel
	}

	return &SdxlTool{
		Definition: def,
		Store:      store,
	}, nil
}

type SdxlTool struct {
	Definition chattools.ToolDefinition
	Store      storageprovider.StorageProvider
}

type SdxlResult struct {
	Data []byte `yaml:"-"`
	Path string `yaml:"path"`
}

func (sr SdxlResult) String() string {
	return fmt.Sprintf("SdxlResult{ Data: len(%v), Path: %v}", len(sr.Data), sr.Path)
}

func (t *SdxlTool) Run(call chattools.ToolCall) chattools.ToolResult {
	prompt, ok := call.Params["prompt"].(string)
	if !ok {
		log.Error("missing prompt in SDXL call %+v", call.Params)
		return chattools.ToolResult{
			Success: false,
			Content: "missing prompt",
		}
	}

	negPrompt, _ := call.Params["negativePrompt"].(string)

	data, err := t.getImage(prompt, negPrompt)
	if err != nil {
		log.Error("could not generate image: %v", err)
		return chattools.ToolResult{
			Success: false,
			Content: fmt.Sprintf("error while generating image: %v", err),
		}
	}

	filename := fmt.Sprintf("image-%v-%v.png", time.Now().Format(time.RFC3339), rand.Intn(1000))

	path, err := t.Store.Save(filename, data)
	if err != nil {
		log.Error("error while saving image to %v: %v", path, err)
		return chattools.ToolResult{
			Success: false,
			Content: fmt.Sprintf("error while saving image to %v: %v", path, err),
		}
	}

	val := SdxlResult{
		Data: data,
		Path: path,
	}

	return chattools.ToolResult{Success: true, Content: fmt.Sprintf("Image saved at: %v", path), Value: val}
}

func (t *SdxlTool) getImage(prompt, negativePrompt string) ([]byte, error) {
	url := fmt.Sprintf("https://api.fireworks.ai/inference/v1/image_generation/accounts/fireworks/models/%s", t.Definition.SdxlSettings.Model)

	//Full prompts
	prompt = strings.TrimSpace(fmt.Sprintf("%v %v", t.Definition.SdxlSettings.PositivePrompt, prompt))
	negativePrompt = strings.TrimSpace(fmt.Sprintf("%v %v", t.Definition.SdxlSettings.NegativePrompt, negativePrompt))

	log.Debug("Image prompt: %v", prompt)
	if negativePrompt != "" {
		log.Debug("Image negative prompt: %v", prompt)
	}

	payload := map[string]interface{}{
		"cfg_scale":       7,
		"height":          1024,
		"width":           1024,
		"sampler":         nil,
		"samples":         1,
		"steps":           30,
		"seed":            0,
		"style_preset":    nil,
		"safety_check":    false,
		"prompt":          prompt,
		"negative_prompt": negativePrompt,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "image/jpeg")
	req.Header.Set("Authorization", "Bearer "+t.Definition.SdxlSettings.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("sdxl api responded with " + strconv.Itoa(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
