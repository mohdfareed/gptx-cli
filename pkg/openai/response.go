package openai

import (
	"context"

	"github.com/mohdfareed/gptx-cli/internal/events"
	"github.com/mohdfareed/gptx-cli/internal/tools"
	"github.com/openai/openai-go/responses"
)

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

func parseStream(
	event responses.ResponseStreamEventUnion,
	modelEvents events.ModelEvents, ctx context.Context,
) {
	switch variant := event.AsAny().(type) {
	case responses.ResponseTextDeltaEvent:
		modelEvents.Reply.Emit(ctx, variant.Delta)
	case responses.ResponseRefusalDeltaEvent:
		modelEvents.Reply.Emit(ctx, variant.Delta)
	case responses.ResponseFunctionCallArgumentsDeltaEvent:
		modelEvents.Reply.Emit(ctx, "")
	case responses.ResponseWebSearchCallSearchingEvent:
		modelEvents.Reply.Emit(ctx, "")
	}
}
