package openai

import (
	"github.com/mohdfareed/gptx-cli/pkg/llm"
	"github.com/openai/openai-go/responses"
	"github.com/openai/openai-go/shared"
)

const WebSearch llm.Tool = "web_search"

// WebSearchTool is a tool that allows the model to search the web for information.
var WebSearchTool = llm.ModelTool[struct{}, struct{}, struct{}]{
	Name: "web_search",
	Desc: "search the web for information",
	Call: func(a struct{}) (struct{}, error) { return a, nil },
	Init: func() error { return nil },
}

type OpenAIConfig struct {
	APIKey    string                 `json:"api_key"`
	Model     shared.ChatModel       `json:"model"`
	Reasoning shared.ReasoningEffort `json:"reason"`
	Temp      float64                `json:"temp"`
	MaxTokens int64                  `json:"max_tokens"`
	Config    llm.ModelConfig        `json:"config"`
}

type openAIMsg = responses.ResponseInputItemUnionParam
type openAIData = responses.ResponseInputContentUnionParam
type openAIUsage = responses.ResponseUsage
