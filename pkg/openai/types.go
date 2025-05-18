package openai

import (
	"github.com/openai/openai-go/packages/param"
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

// File represents a file to be attached to a message.
type File struct {
	Name    string // File name
	Path    string // File path
	Content string // File content
}

// MARK: Tools
// ============================================================================

// WebSearch is the tool definition for web search.
var WebSearch ToolDef = responses.ToolUnionParam{
	OfWebSearch: &responses.WebSearchToolParam{
		Type: responses.WebSearchToolTypeWebSearchPreview,
	},
}

// NewTool creates a new tool definition.
func NewTool(name, desc string, params map[string]any) ToolDef {
	return responses.ToolUnionParam{
		OfFunction: &responses.FunctionToolParam{
			Name:        name,
			Description: param.Opt[string]{Value: desc},
			Parameters:  params,
			Strict:      true,
		},
	}
}
