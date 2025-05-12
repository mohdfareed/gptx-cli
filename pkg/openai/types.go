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

// StreamParser is a parser for streaming responses from OpenAI.
type StreamParser struct {
	Text         chan string
	Refusal      chan string
	FunctionCall chan string
	WebSearch    chan struct{}
}
