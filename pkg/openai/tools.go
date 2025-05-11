package openai

import (
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/responses"
)

// WebSearch is the tool definition for web search.
var WebSearch ToolDef = responses.ToolParamOfWebSearch(
	responses.WebSearchToolTypeWebSearchPreview,
)

// NewTool creates a new tool definition.
func NewTool(name, desc string, params map[string]any) ToolDef {
	tool := responses.ToolParamOfFunction(
		name, params, true,
	)
	tool.OfFunction.Description = param.Opt[string]{Value: desc}
	return tool
}
