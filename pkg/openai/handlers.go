// Package openai implements the OpenAI Responses API integration.
package openai

import (
	"context"
	"fmt"

	"github.com/mohdfareed/gptx-cli/pkg/gptx"
	"github.com/openai/openai-go/responses"
)

// getUsageJSON returns usage information as a JSON string.
// Returns "{}" if usage information isn't available.

// handleStreamEvent processes streaming events from the OpenAI API.
// getUsageJSON returns usage information as a JSON string.
// Returns "{}" if usage information isn't available.
func getUsageJSON(usage responses.ResponseUsage) string {
	// Since we can't directly access the fields easily, just return
	// the raw JSON if it's available or empty JSON if it's not
	rawJSON := usage.RawJSON()
	if rawJSON == "" {
		return "{}"
	}
	return rawJSON
}

func (c *OpenAIClient) handleStreamEvent(
	ctx context.Context,
	event responses.ResponseStreamEventUnion,
	request gptx.Request,
) {
	switch variant := event.AsAny().(type) {
	case responses.ResponseTextDeltaEvent:
		// Handle incremental text responses
		if request.Callbacks.OnReply != nil {
			request.Callbacks.OnReply(variant.Delta)
		}
	case responses.ResponseRefusalDeltaEvent:
		// Handle model refusals
		if request.Callbacks.OnReply != nil {
			request.Callbacks.OnReply(variant.Delta)
		}
	case responses.ResponseFunctionCallArgumentsDeltaEvent:
		// Function call arguments are streamed incrementally
		// We don't handle these at the stream level - wait for complete function calls
	case responses.ResponseWebSearchCallSearchingEvent:
		// Web search status events are used to indicate that a search is happening
		// We could emit a special message here if needed
	}
}

// handleCompleteResponse processes a completed response.
// This handles any tool calls that need to be processed.
func (c *OpenAIClient) handleCompleteResponse(
	ctx context.Context,
	response *responses.Response,
	request gptx.Request,
) error {
	for _, item := range response.Output {
		switch item.AsAny().(type) {
		case responses.ResponseOutputMessage:
			// Handle completed text response
			for _, content := range item.AsMessage().ToParam().Content {
				if content.OfOutputText != nil && request.Callbacks.OnReply != nil {
					request.Callbacks.OnReply(content.OfOutputText.Text)
				}
				if content.OfRefusal != nil && request.Callbacks.OnReply != nil {
					request.Callbacks.OnReply(content.OfRefusal.Refusal)
				}
			}

		case responses.ResponseFunctionWebSearch:
			// Handle web search results
			_ = item.AsWebSearchCall() // We access the data but don't need to use it directly
			if request.ToolHandler != nil {
				// Extract the search query from the function call
				// The web search data doesn't directly expose a Query field in the response
				// Instead, we pass an empty string to the handler which will handle this special case
				result, err := request.ToolHandler(ctx, WebSearchToolDef, "")
				if err != nil {
					return fmt.Errorf("web search: %w", err)
				}
				if request.Callbacks.OnReply != nil {
					request.Callbacks.OnReply(result)
				}
			}

		case responses.ResponseFunctionToolCall:
			// Handle regular tool calls
			toolCall := item.AsFunctionCall()
			if request.ToolHandler != nil {
				result, err := request.ToolHandler(ctx, toolCall.Name, toolCall.Arguments)
				if err != nil {
					return fmt.Errorf("tool %s: %w", toolCall.Name, err)
				}
				if request.Callbacks.OnReply != nil {
					request.Callbacks.OnReply(result)
				}
			}

		case responses.ResponseReasoningItem:
			// Handle reasoning output
			reasoning := item.AsReasoning()
			if request.Callbacks.OnReasoning != nil {
				for _, step := range reasoning.Summary {
					request.Callbacks.OnReasoning(step.Text)
				}
			}
		}
	}

	return nil
}
