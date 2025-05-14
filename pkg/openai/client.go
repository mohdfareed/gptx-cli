// Package openai provides an OpenAI-specific client for the GPTx CLI
package openai

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mohdfareed/gptx-cli/pkg/gptx"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/responses"
	"github.com/openai/openai-go/shared"
)

// Client implements the gptx.Client interface using the OpenAI API
type Client struct {
	client openai.Client
}

// responseEvent implements the gptx.ResponseEvent interface
type responseEvent struct {
	eventType  string
	content    string
	toolCall   *gptx.ToolCall
	streamData *responses.ResponseStreamingResponseParam
}

// responseStream implements the gptx.ResponseStream interface
type responseStream struct {
	stream       *responses.ResponseStreamingResponseIterator
	client       *Client
	currentEvent *responseEvent
}

// New creates a new OpenAI client
func New(apiKey string, organizationID string) *Client {
	// Create the OpenAI client
	client := openai.NewClient(apiKey)

	// Set organization ID if provided
	if organizationID != "" {
		client = client.WithOrganization(organizationID)
	}

	return &Client{
		client: client,
	}
}

// Generate sends a request to the OpenAI API and returns a stream of events
func (c *Client) Generate(ctx context.Context, req *gptx.Request) (gptx.ResponseStream, error) {
	// Convert request to OpenAI format
	openaiReq, err := c.convertRequest(req)
	if err != nil {
		return nil, fmt.Errorf("convert request: %w", err)
	}

	// Start streaming
	stream := c.client.Responses.NewStreaming(ctx, openaiReq)

	// Return our wrapped stream
	return &responseStream{
		stream: stream,
		client: c,
	}, nil
}

// ProcessStream processes a response stream, handling text and tool calls
func (c *Client) ProcessStream(
	ctx context.Context,
	stream gptx.ResponseStream,
	textHandler func(string),
	toolHandlers map[string]gptx.ToolHandler,
) error {
	// Process stream events
	for stream.HasNext() {
		event := stream.Next()

		switch event.GetType() {
		case "text":
			// Handle text content
			content := event.GetContent()
			if content != "" && textHandler != nil {
				textHandler(content)
			}

		case "tool_call":
			// Handle tool call
			toolCall := event.GetToolCall()
			if toolCall == nil {
				continue
			}

			// Skip if handler not available
			handler, ok := toolHandlers[toolCall.Name]
			if !ok {
				continue
			}

			// Execute tool
			result, err := handler(ctx, toolCall.Name, toolCall.Arguments)
			if err != nil {
				return fmt.Errorf("execute tool %s: %w", toolCall.Name, err)
			}

			// Submit tool output
			err = stream.SubmitToolOutputs(ctx, []gptx.ToolOutput{
				{
					ID:     toolCall.ID,
					Name:   toolCall.Name,
					Output: string(result),
				},
			})

			if err != nil {
				return fmt.Errorf("submit tool output: %w", err)
			}
		}
	}

	// Check for stream errors
	if err := stream.Err(); err != nil {
		return fmt.Errorf("stream error: %w", err)
	}

	return nil
}

// convertRequest converts a gptx.Request to an OpenAI API request
func (c *Client) convertRequest(req *gptx.Request) (responses.ResponseNewParams, error) {
	// Create messages array
	inputItems := []responses.ResponseInputItemUnionParam{}

	// Add system prompt if provided
	if req.SystemPrompt != "" {
		inputItems = append(inputItems, responses.ResponseInputItemUnionParam{
			OfMessageParam: &responses.MessageParam{
				Role:    "system",
				Content: responses.ContentUnionParam{OfString: req.SystemPrompt},
			},
		})
	}

	// Convert messages
	for _, msg := range req.Messages {
		// Convert content based on whether we have files
		var content responses.ContentUnionParam

		// Simple text message
		content = responses.ContentUnionParam{OfString: msg.Content}

		// Add message to input items
		inputItems = append(inputItems, responses.ResponseInputItemUnionParam{
			OfMessageParam: &responses.MessageParam{
				Role:    msg.Role,
				Content: content,
			},
		})
	}

	// Convert tools
	var tools []responses.ToolUnionParam
	for _, tool := range req.Tools {
		var openaiTool responses.ToolUnionParam

		switch tool.Type {
		case "function":
			// Create function tool
			openaiTool = responses.ToolUnionParam{
				OfFunction: &responses.FunctionToolParam{
					Type:     "function",
					Function: convertFunctionParams(tool),
				},
			}
		case "web_search":
			// Create web search tool
			openaiTool = responses.ToolUnionParam{
				OfWebSearch: &responses.WebSearchToolParam{
					Type: responses.WebSearchToolTypeWebSearchPreview,
				},
			}
		}

		tools = append(tools, openaiTool)
	}

	// Prepare file attachments
	fileAttachments := []responses.ResponseInputContentUnionParam{}
	for _, filePath := range req.Files {
		attachment, err := createFileAttachment(filePath)
		if err != nil {
			return responses.ResponseNewParams{}, fmt.Errorf("create file attachment: %w", err)
		}
		fileAttachments = append(fileAttachments, attachment)
	}

	// Create input union
	input := responses.ResponseNewParamsInputUnion{
		OfInputItemList: inputItems,
	}

	// Convert to ChatGPT model format
	var model shared.ChatModel
	if strings.HasPrefix(req.Model, "gpt-4") {
		model = shared.ChatModelGPT4
	} else if strings.HasPrefix(req.Model, "gpt-3.5") {
		model = shared.ChatModelGPT3_5Turbo
	} else {
		// Default to latest model
		model = shared.ChatModelGPT4o
	}

	// Create request
	requestParams := responses.ResponseNewParams{
		Model:             model,
		Input:             input,
		Tools:             tools,
		Temperature:       param.Opt[float64]{Value: float64(req.Temperature)},
		ParallelToolCalls: param.Opt[bool]{Value: true},
	}

	// Set max tokens if specified
	if req.MaxTokens > 0 {
		requestParams.MaxOutputTokens = param.Opt[int64]{Value: int64(req.MaxTokens)}
	}

	// Set user if specified
	if req.User != "" {
		requestParams.User = param.Opt[string]{Value: req.User}
	}

	return requestParams, nil
}

// createFileAttachment creates a file attachment from a file path
func createFileAttachment(filePath string) (responses.ResponseInputContentUnionParam, error) {
	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return responses.ResponseInputContentUnionParam{}, fmt.Errorf("read file: %w", err)
	}

	// Get file type
	fileType := "text/plain"
	switch {
	case strings.HasSuffix(filePath, ".go"):
		fileType = "text/plain"
	case strings.HasSuffix(filePath, ".md"):
		fileType = "text/markdown"
	case strings.HasSuffix(filePath, ".json"):
		fileType = "application/json"
	case strings.HasSuffix(filePath, ".png"):
		fileType = "image/png"
	case strings.HasSuffix(filePath, ".jpg"), strings.HasSuffix(filePath, ".jpeg"):
		fileType = "image/jpeg"
	}

	// Create attachment
	return responses.ResponseInputContentUnionParam{
		OfImageParam: &responses.ImageParam{
			Type:   fileType,
			Source: responses.ImageSourceUnionParam{OfBase64: string(data)},
		},
	}, nil
}

// convertFunctionParams converts function parameters
func convertFunctionParams(tool common.Tool) responses.FunctionObjectParam {
	return responses.FunctionObjectParam{
		Name:        tool.Name,
		Description: param.Opt[string]{Value: tool.Description},
		Parameters:  tool.Parameters,
	}
}

// HasNext returns true if there are more events
func (s *responseStream) HasNext() bool {
	return s.stream.Next()
}

// Next returns the next event
func (s *responseStream) Next() gptx.ResponseEvent {
	data := s.stream.Current()

	// Create event based on data type
	s.currentEvent = &responseEvent{
		streamData: &data,
	}

	switch {
	case data.AsResponseDelta != nil:
		delta := data.AsResponseDelta
		s.currentEvent.eventType = "text"

		// Extract text content if available
		if delta.Delta.Content != nil {
			s.currentEvent.content = *delta.Delta.Content
		}

	case data.AsResponseToolCallDelta != nil:
		toolCall := data.AsResponseToolCallDelta
		s.currentEvent.eventType = "tool_call"

		// Create tool call object
		s.currentEvent.toolCall = &gptx.ToolCall{
			ID:   toolCall.ToolCallID,
			Name: toolCall.Delta.Name,
		}

		// Add function arguments if available
		if toolCall.Delta.Function != nil && toolCall.Delta.Function.Arguments != nil {
			s.currentEvent.toolCall.Arguments = []byte(*toolCall.Delta.Function.Arguments)
		}
	}

	return s.currentEvent
}

// Close closes the stream
func (s *responseStream) Close() {
	s.stream.Close()
}

// Err returns any error that occurred during streaming
func (s *responseStream) Err() error {
	return s.stream.Err()
}

// SubmitToolOutputs submits tool outputs to the model
func (s *responseStream) SubmitToolOutputs(ctx context.Context, outputs []gptx.ToolOutput) error {
	// Convert outputs to OpenAI format
	var toolOutputs []responses.ResponseSubmitToolOutputsParamsToolOutputParam

	for _, output := range outputs {
		toolOutputs = append(toolOutputs, responses.ResponseSubmitToolOutputsParamsToolOutputParam{
			ToolCallID: output.ID,
			Output:     output.Output,
		})
	}

	// Submit tool outputs
	_, err := s.client.client.Responses.SubmitToolOutputs(ctx, responses.ResponseSubmitToolOutputsParams{
		ThreadID:    s.stream.Current().ThreadID,
		ToolOutputs: toolOutputs,
	})

	return err
}

// GetType returns the event type
func (e *responseEvent) GetType() string {
	return e.eventType
}

// GetContent returns the text content if available
func (e *responseEvent) GetContent() string {
	return e.content
}

// GetToolCall returns tool call info if available
func (e *responseEvent) GetToolCall() *gptx.ToolCall {
	return e.toolCall
}
