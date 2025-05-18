package openai

import (
	"github.com/mohdfareed/gptx-cli/pkg/gptx"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/responses"
	"github.com/openai/openai-go/shared"
)

// NewRequest creates a new model request with the given parameters.
func NewRequest(
	model gptx.Model, msgs []MsgData,
	reasoning shared.ReasoningEffort,
	userID string,
) ModelRequest {
	var tools []ToolDef
	for _, tool := range model.Tools.GetTools() {
		tools = append(tools, NewTool(tool.Name, tool.Desc, tool.Params))
	}

	history := responses.ResponseNewParamsInputUnion{OfInputItemList: msgs}
	data := responses.ResponseNewParams{
		Model:        model.Config.Model,
		Input:        history,
		Tools:        tools,
		Instructions: param.Opt[string]{Value: model.Config.SysPrompt},
		User:         param.Opt[string]{Value: userID},
		// defaults
		Store:             param.Opt[bool]{Value: false},
		ParallelToolCalls: param.Opt[bool]{Value: true},
	}

	// temperature
	if model.Config.Temp != nil {
		data.Temperature = param.Opt[float64]{Value: *model.Config.Temp}
	}

	// max tokens
	if model.Config.Tokens != nil {
		data.MaxOutputTokens = param.Opt[int64]{
			Value: int64(*model.Config.Tokens),
		}
	}

	// // reasoning
	// if reasoning != "" {
	// 	data.Reasoning = shared.ReasoningParam{
	// 		Effort:          reasoning,
	// 		GenerateSummary: shared.ReasoningGenerateSummaryDetailed,
	// 	}
	// }
	return data
}
