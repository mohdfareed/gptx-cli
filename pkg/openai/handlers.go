// Package openai implements the OpenAI Responses API integration.
package openai

import (
	"context"
	"encoding/json"

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
	rawJSON, err := json.MarshalIndent(usage, "", "  ")
	if err != nil {
		return ""
	}
	return string(rawJSON)
}

func (c *OpenAIClient) handleStreamEvent(
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
			// Signal that a web search is happening
			request.Callbacks.OnWebSearch()

		case responses.ResponseFunctionToolCall:
			// Handle regular tool calls
			toolCall := item.AsFunctionCall()
			if request.ToolHandler != nil {
				result, err := request.ToolHandler(ctx, toolCall.Name, toolCall.Arguments)
				if err != nil {
					// Pass the error up without wrapping it again to avoid duplicate messages
					return err
				}
				if request.Callbacks.OnReply != nil && result != "" {
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
