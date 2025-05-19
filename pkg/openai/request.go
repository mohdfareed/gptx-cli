// Package openai implements OpenAI's Responses API integration.
package openai

import (
	"github.com/mohdfareed/gptx-cli/internal/tools"
	"github.com/mohdfareed/gptx-cli/pkg/gptx"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/responses"
	"github.com/openai/openai-go/shared"
)

// NewRequest creates a request for the OpenAI Responses API.
func NewRequest(
	request gptx.Request, msgs []MsgData, userID string, tools []ToolDef,
) ModelRequest {
	// Create the input history from messages
	history := responses.ResponseNewParamsInputUnion{OfInputItemList: msgs}

	// Build the complete request with all configuration options
	data := responses.ResponseNewParams{
		// Core parameters
		Model:        request.Config.Model,                               // Which model to use (e.g., "o4-mini")
		Input:        history,                                            // Message history
		Tools:        tools,                                              // Available tools from registry
		Instructions: param.Opt[string]{Value: request.Config.SysPrompt}, // System prompt
		User:         param.Opt[string]{Value: userID},                   // User identifier

		// Default settings
		Store:             param.Opt[bool]{Value: false}, // Don't store conversations
		ParallelToolCalls: param.Opt[bool]{Value: true},  // Allow parallel tool usage
	}

	// Apply temperature setting if specified
	// Temperature controls randomness: higher values make output more random,
	// lower values make it more deterministic and focused
	if request.Config.Temp >= 0 {
		data.Temperature = param.Opt[float64]{Value: request.Config.Temp}
	}

	// Apply token limit if specified
	// This limits the maximum length of the model's response
	if request.Config.Tokens >= 0 {
		data.MaxOutputTokens = param.Opt[int64]{
			Value: int64(request.Config.Tokens),
		}
	}

	// Reasoning setting is currently disabled
	// When enabled, this would control how much of the model's reasoning
	// process is exposed in the response
	if request.Config.Reason {
		data.Reasoning = shared.ReasoningParam{
			Effort:          shared.ReasoningEffortHigh,
			GenerateSummary: shared.ReasoningGenerateSummaryDetailed,
		}
	}
	return data
}

// NewTool creates a new tool definition.
func NewTool(tool tools.ToolDef) ToolDef {
	return responses.ToolUnionParam{
		OfFunction: &responses.FunctionToolParam{
			Name:        tool.Name,
			Description: param.Opt[string]{Value: tool.Desc},
			// Parameters:  tool.Params,
			Parameters: map[string]any{
				"type":                 "object",
				"properties":           tool.Params,
				"required":             tool.Required,
				"additionalProperties": false,
			},
			Strict: true,
		},
	}
}
