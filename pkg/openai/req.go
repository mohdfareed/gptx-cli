package openai

import (
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/responses"
	"github.com/openai/openai-go/shared"
)

// ModelRequest represents the model request structure.
type ModelRequest = responses.ResponseNewParams

// MARK: Requests
// ============================================================================

func NewRequest(
	model shared.ResponsesModel,
	sysPrompt string, msgs []MsgData, tools []ToolDef,
	maxTokens int64, temp float64, reasoning shared.ReasoningEffort,
) ModelRequest {
	history := responses.ResponseNewParamsInputUnion{OfInputItemList: msgs}
	data := responses.ResponseNewParams{
		Model:           model,
		Input:           history,
		Tools:           tools,
		Instructions:    param.Opt[string]{Value: sysPrompt},
		MaxOutputTokens: param.Opt[int64]{Value: maxTokens},
		Temperature:     param.Opt[float64]{Value: temp},
		Include: []responses.ResponseIncludable{
			responses.ResponseIncludableMessageInputImageImageURL,
		},
		ParallelToolCalls: param.Opt[bool]{Value: true},
	}

	// reasoning
	if reasoning != "" {
		data.Reasoning = shared.ReasoningParam{
			Effort:          reasoning,
			GenerateSummary: shared.ReasoningGenerateSummaryDetailed,
		}
	}
	return data
}

// MARK: Response
// ============================================================================

func parseStream(event responses.ResponseStreamEventUnion) string {
	return event.Delta
}

func parse(response *responses.Response) ([]MsgData, MsgUsage, error) {
	msgs := []MsgData{}
	for _, item := range response.Output {
		msg := MsgData{}
		switch item.AsAny().(type) {
		case responses.ResponseOutputMessage:
			msgData := item.AsMessage().ToParam()
			msg = MsgData{OfOutputMessage: &msgData}
		case responses.ResponseFunctionWebSearch:
			msgData := item.AsWebSearchCall().ToParam()
			msg = MsgData{OfWebSearchCall: &msgData}
		case responses.ResponseFunctionToolCall:
			msgData := item.AsFunctionCall().ToParam()
			msg = MsgData{OfFunctionCall: &msgData}
		case responses.ResponseReasoningItem:
			msgData := item.AsReasoning().ToParam()
			msg = MsgData{OfReasoning: &msgData}
		case responses.ResponseComputerToolCall:
			msgData := item.AsComputerCall().ToParam()
			msg = MsgData{OfComputerCall: &msgData}
		case responses.ResponseFileSearchToolCall:
			msgData := item.AsFileSearchCall().ToParam()
			msg = MsgData{OfFileSearchCall: &msgData}
		}
		msgs = append(msgs, msg)
	}
	return msgs, response.Usage, nil
}
