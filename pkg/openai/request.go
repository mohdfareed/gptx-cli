package openai

import (
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/responses"
	"github.com/openai/openai-go/shared"
)

// NewRequest creates a new model request with the given parameters.
func NewRequest(
	model shared.ChatModel,
	instructs string, msgs []MsgData, tools []ToolDef,
	maxTokens *int64, temp float64, reasoning shared.ReasoningEffort,
	userID string,
) ModelRequest {
	history := responses.ResponseNewParamsInputUnion{OfInputItemList: msgs}
	data := responses.ResponseNewParams{
		Model:        model,
		Input:        history,
		Tools:        tools,
		Instructions: param.Opt[string]{Value: instructs},
		Temperature:  param.Opt[float64]{Value: temp},
		User:         param.Opt[string]{Value: userID},
		// defaults
		Store:             param.Opt[bool]{Value: false},
		ParallelToolCalls: param.Opt[bool]{Value: true},
	}

	// max tokens
	if maxTokens != nil {
		data.MaxOutputTokens = param.Opt[int64]{Value: *maxTokens}
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
