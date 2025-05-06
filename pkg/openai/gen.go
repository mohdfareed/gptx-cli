package openai

import (
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/responses"
	"github.com/openai/openai-go/shared"
)

// MARK: Requests
// ============================================================================

func newRequest(
	config ModelConfig, msgs []openAIMsg, tools []Tool,
) responses.ResponseNewParams {
	// config
	data := responses.ResponseNewParams{
		Model:        config.Model,
		Instructions: param.Opt[string]{Value: config.SysPrompt},
		Include: []responses.ResponseIncludable{
			responses.ResponseIncludableMessageInputImageImageURL,
		},
		MaxOutputTokens:   param.Opt[int64]{Value: config.MaxTokens},
		ParallelToolCalls: param.Opt[bool]{Value: true},
		Temperature:       param.Opt[float64]{Value: config.Temp},
		User:              param.Opt[string]{Value: userID},
	}

	// reasoning
	if config.ReasonEffort != "" {
		data.Reasoning = shared.ReasoningParam{
			Effort:          config.ReasonEffort,
			GenerateSummary: shared.ReasoningGenerateSummaryDetailed,
		}
	}

	// tools
	var toolDefs []ToolDef
	for _, tool := range tools {
		toolDefs = append(toolDefs, tool.Def)
	}
	data.Tools = toolDefs

	// messages
	reqMsgs := make([]openAIMsg, len(msgs))
	for i, msg := range msgs {
		reqMsgs[i] = msg.Data
	}
	data.Input.OfInputItemList = reqMsgs
	return data
}

// MARK: Response
// ============================================================================

func parse(response *responses.Response) ([]openAIMsg, error) {
	msgs := []openAIMsg{}
	for _, item := range response.Output {
		msg := openAIMsg{}
		switch item.AsAny().(type) {
		case responses.ResponseOutputMessage:
			msgData := item.AsMessage().ToParam()
			msg.Data = openAIMsg{OfOutputMessage: &msgData}
		case responses.ResponseFunctionWebSearch:
			msgData := item.AsWebSearchCall().ToParam()
			msg.Data = openAIMsg{OfWebSearchCall: &msgData}
		case responses.ResponseFunctionToolCall:
			msgData := item.AsFunctionCall().ToParam()
			msg.Data = openAIMsg{OfFunctionCall: &msgData}
		case responses.ResponseReasoningItem:
			msgData := item.AsReasoning().ToParam()
			msg.Data = openAIMsg{OfReasoning: &msgData}
		case responses.ResponseComputerToolCall:
			msgData := item.AsComputerCall().ToParam()
			msg.Data = openAIMsg{OfComputerCall: &msgData}
		case responses.ResponseFileSearchToolCall:
			msgData := item.AsFileSearchCall().ToParam()
			msg.Data = openAIMsg{OfFileSearchCall: &msgData}
		}
		msgs = append(msgs, msg)
	}

	// set usage of last message
	if len(msgs) > 0 {
		msgs[len(msgs)-1].Usage = &response.Usage
	}
	return msgs, nil
}

// MARK: Messages
// ============================================================================

// TextMsg creates a user message with the given text.
func UserMsg(text string) openAIMsg {
	msg := responses.ResponseInputTextParam{Text: text}
	data := openAIMsg{
		OfInputMessage: &responses.ResponseInputItemMessageParam{
			Role: "user", Content: []openAIData{{OfInputText: &msg}},
		},
	}
	return openAIMsg{Data: data}
}
