// Package openai implements OpenAI's Responses API integration.
package openai

import (
	"context"

	"github.com/mohdfareed/gptx-cli/internal/events"
	"github.com/mohdfareed/gptx-cli/internal/tools"
	"github.com/openai/openai-go/responses"
)

// parse processes API responses and emits appropriate events.
// Handles text responses, tool calls, web search, and reasoning.
func parse(
	response *responses.Response, events events.ModelEvents,
	ctx context.Context,
) ([]MsgData, MsgUsage, error) {
	msgs := []MsgData{}
	for _, item := range response.Output {
		msg := MsgData{}
		switch item.AsAny().(type) {
		case responses.ResponseOutputMessage:
			msgData := item.AsMessage().ToParam()
			msg = MsgData{OfOutputMessage: &msgData}
			for _, msg := range msgData.Content {
				if msg.OfOutputText != nil {
					events.Reply.Emit(ctx, msg.OfOutputText.Text)
				} else if msg.OfRefusal != nil {
					events.Reply.Emit(ctx, msg.OfRefusal.Refusal)
				}
			}
		case responses.ResponseFunctionWebSearch:
			msgData := item.AsWebSearchCall().ToParam()
			msg = MsgData{OfWebSearchCall: &msgData}
			events.ToolCall.Emit(ctx, tools.ToolCall{
				Name: tools.WebSearchToolDef, Params: "",
			})
		case responses.ResponseFunctionToolCall:
			msgData := item.AsFunctionCall().ToParam()
			msg = MsgData{OfFunctionCall: &msgData}
			events.ToolCall.Emit(ctx, tools.ToolCall{
				Name: msgData.Name, Params: msgData.Arguments,
			})
		case responses.ResponseReasoningItem:
			msgData := item.AsReasoning().ToParam()
			msg = MsgData{OfReasoning: &msgData}
			for _, msg := range msgData.Summary {
				events.Reply.Emit(ctx, msg.Text)
			}
		}
		msgs = append(msgs, msg)
	}
	return msgs, response.Usage, nil
}

// parseStream processes individual stream events from the OpenAI API's
// streaming response and emits appropriate events. This function
// handles the real-time incremental updates from the API, allowing
// for responsive user interfaces that show progress as it happens.
//
// The function handles different types of streaming events:
// - Text deltas (incremental pieces of the response text)
// - Refusal deltas (when the model refuses to respond)
// - Function call argument deltas (tool call parameter updates)
// - Web search status events (when search is being performed)
//
// This streaming approach provides a better user experience by showing
// responses as they are generated rather than waiting for the complete response.
func parseStream(
	event responses.ResponseStreamEventUnion,
	modelEvents events.ModelEvents, ctx context.Context,
) {
	switch variant := event.AsAny().(type) {
	case responses.ResponseTextDeltaEvent:
		// Emit each text fragment as it arrives
		modelEvents.Reply.Emit(ctx, variant.Delta)
	case responses.ResponseRefusalDeltaEvent:
		// Emit refusal messages (when model declines to respond)
		modelEvents.Reply.Emit(ctx, variant.Delta)
	case responses.ResponseFunctionCallArgumentsDeltaEvent:
		// Function call arguments are built up incrementally
		// We don't emit anything here since we wait for the complete call
		modelEvents.Reply.Emit(ctx, "")
	case responses.ResponseWebSearchCallSearchingEvent:
		// Indicate that web search is in progress
		// We could improve this to show a "searching..." message
		modelEvents.Reply.Emit(ctx, "")
	}
}
