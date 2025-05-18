import (
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/responses"
)

// ResponsesModel converts a model string to the proper type.
func ResponsesModel(model string) responses.ResponsesModel {
	switch model {
	case "gpt-4o":
		return responses.ResponsesModelGpt4o
	case "gpt-4o-mini":
		return responses.ResponsesModelGpt4oMini
	case "gpt-4":
		return responses.ResponsesModelGpt4
	case "gpt-4-turbo":
		return responses.ResponsesModelGpt4Turbo
	case "gpt-3.5-turbo":
		return responses.ResponsesModelGpt35Turbo
	case "o4", "gpt-4.1":
		return responses.ResponsesModelGpt4o
	case "o4-mini":
		return responses.ResponsesModelGpt4oMini
	default:
		// Default to the model string as-is, for future models
		return responses.ResponsesModel(model)
	}
}

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
			Type:        responses.ToolTypeFunction,
			Strict:      true,
		},
	}
}
