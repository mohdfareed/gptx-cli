// Package openai implements OpenAI's Responses API integration.
package openai

import (
	"github.com/mohdfareed/gptx-cli/internal/cfg"
	"github.com/mohdfareed/gptx-cli/pkg/gptx"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/responses"
	"github.com/openai/openai-go/shared"
)

// NewRequest creates a request for the OpenAI Responses API.
func NewRequest(
	request gptx.Request, msgs []MsgData, userID string,
) ModelRequest {
	// Create the input history from messages
	history := responses.ResponseNewParamsInputUnion{OfInputItemList: msgs}

	// Build the complete request with all configuration options
	data := responses.ResponseNewParams{
		// Core parameters
		Model:        request.Config.Model,                               // Which model to use (e.g., "o4-mini")
		Input:        history,                                            // Message history
		Tools:        modelToolsToOpenAI(request.Config),                 // Available tools
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
func NewTool(name, desc string, params map[string]any) ToolDef {
	return responses.ToolUnionParam{
		OfFunction: &responses.FunctionToolParam{
			Name:        name,
			Description: param.Opt[string]{Value: desc},
			Parameters:  params,
			Strict:      true,
		},
	}
}

// WebSearchToolDef is the identifier for the web search tool.
const WebSearchToolDef = "web-search"

// modelToolsToOpenAI converts tool configuration to OpenAI API format.
// This isolates the OpenAI-specific tool formatting from the rest of the app.
func modelToolsToOpenAI(config cfg.Config) []responses.ToolUnionParam {
	var toolDefs []responses.ToolUnionParam

	// Add web search if enabled
	if config.WebSearch {
		toolDefs = append(toolDefs, WebSearch)
	}

	// Add shell tool if enabled
	if config.Shell != "" {
		shellTool := responses.ToolUnionParam{
			OfFunction: &responses.FunctionToolParam{
				Name:        "shell",
				Description: param.Opt[string]{Value: "Execute shell commands."},
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"cmd": map[string]any{
							"type":        "string",
							"description": "The command to execute",
						},
					},
					"required": []string{
						"cmd",
					},
					"additionalProperties": false,
				},
				Strict: true,
			},
		}
		toolDefs = append(toolDefs, shellTool)
	}
	return toolDefs
}
