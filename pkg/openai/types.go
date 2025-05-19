package openai

import (
	"github.com/openai/openai-go/responses"
)

// MsgData represents the data structure for a message.
type MsgData = responses.ResponseInputItemUnionParam

// FileData represents the data structure for a file attachment.
type FileData = responses.ResponseInputContentUnionParam

// ToolDef represents the definition of a model tool.
type ToolDef = responses.ToolUnionParam

// MsgUsage represents the usage information for a response message.
type MsgUsage = responses.ResponseUsage

// ModelRequest represents a request to an OpenAI model.
type ModelRequest = responses.ResponseNewParams

// MARK: Tools
// ============================================================================

// WebSearch is the tool definition for web search.
var WebSearch ToolDef = responses.ToolUnionParam{
	OfWebSearch: &responses.WebSearchToolParam{
		Type: responses.WebSearchToolTypeWebSearchPreview,
	},
}
