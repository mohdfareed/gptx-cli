// Package openai implements OpenAI's Responses API integration.
package openai

import (
	"github.com/mohdfareed/gptx-cli/pkg/gptx"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/responses"
	"github.com/openai/openai-go/shared"
)

// NewRequest creates a request for the OpenAI Responses API.
func NewRequest(
	model gptx.Model, msgs []MsgData,
	reasoning shared.ReasoningEffort,
	userID string,
) ModelRequest {
	// Convert internal tool definitions to OpenAI tool format
	var tools []ToolDef
	for _, tool := range model.Tools.GetTools() {
		tools = append(tools, NewTool(tool.Name, tool.Desc, tool.Params))
	}

	// Create the input history from messages
	history := responses.ResponseNewParamsInputUnion{OfInputItemList: msgs}

	// Build the complete request with all configuration options
	data := responses.ResponseNewParams{
		// Core parameters
		Model:        model.Config.Model,         // Which model to use (e.g., "o4-mini")
		Input:        history,                    // Message history
		Tools:        tools,                      // Available tools
		Instructions: param.Opt[string]{Value: model.Config.SysPrompt}, // System prompt
		User:         param.Opt[string]{Value: userID},                 // User identifier

		// Default settings
		Store:             param.Opt[bool]{Value: false}, // Don't store conversations
		ParallelToolCalls: param.Opt[bool]{Value: true},  // Allow parallel tool usage
	}

	// Apply temperature setting if specified
	// Temperature controls randomness: higher values make output more random,
	// lower values make it more deterministic and focused
	if model.Config.Temp != nil {
		data.Temperature = param.Opt[float64]{Value: *model.Config.Temp}
	}

	// Apply token limit if specified
	// This limits the maximum length of the model's response
	if model.Config.Tokens != nil {
		data.MaxOutputTokens = param.Opt[int64]{
			Value: int64(*model.Config.Tokens),
		}
	}

	// Reasoning setting is currently disabled
	// When enabled, this would control how much of the model's reasoning
	// process is exposed in the response
	// if reasoning != "" {
	// 	data.Reasoning = shared.ReasoningParam{
	// 		Effort:          reasoning,
	// 		GenerateSummary: shared.ReasoningGenerateSummaryDetailed,
	// 	}
	// }
	return data
}
