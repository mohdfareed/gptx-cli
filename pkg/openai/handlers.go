// Package openai implements the OpenAI Responses API integration.
package openai

import (
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

// This space intentionally left blank
// The functionality from handleCompleteResponse has been moved to extractResponseData in client.go
