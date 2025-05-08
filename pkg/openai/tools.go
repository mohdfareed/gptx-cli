package openai

import "github.com/openai/openai-go/responses"

type ToolDef = responses.ToolUnionParam

var WebSearch ToolDef = responses.ToolParamOfWebSearch(
	responses.WebSearchToolTypeWebSearchPreview,
)
