// Package openai implements the OpenAI Responses API integration.
package openai

import (
	"context"
	"fmt"

	"github.com/mohdfareed/gptx-cli/pkg/gptx"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/responses"
)

// OpenAIClient implements the gptx.Client interface for the OpenAI API.
type OpenAIClient struct {
	client openai.Client // OpenAI SDK client
	userID string        // User identifier for API tracking
}

// NewOpenAIClient creates a new OpenAI client with the provided API key.
func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		client: openai.NewClient(option.WithAPIKey(apiKey)),
	}
}

// WithUserID sets the user ID for the client and returns the client.
func (c *OpenAIClient) WithUserID(id string) *OpenAIClient {
	c.userID = id
	return c
}

// SendRequest sends a single request to the OpenAI API and returns the response.
// Implements the gptx.Client interface.
func (c *OpenAIClient) SendRequest(ctx context.Context, request gptx.Request) (gptx.Response, error) {
	// Convert gptx messages to OpenAI messages
	var openAIMessages []MsgData

	// Process all messages in the conversation
	for _, msg := range request.Messages {
		// Handle file attachments for user messages
		if msg.Role == "user" && len(request.Config.Files) > 0 {
			// Only attach files to the latest user message
			userMsg, err := UserMsg(msg.Content, request.Config.Files)
			if err != nil {
				return gptx.Response{}, fmt.Errorf("openai: %w", err)
			}
			openAIMessages = append(openAIMessages, userMsg)
		} else {
			// For other messages or user messages without files
			openAIMsg := responses.ResponseInputItemParamOfMessage(msg.Content, responses.EasyInputMessageRole(msg.Role))
			openAIMessages = append(openAIMessages, openAIMsg)
		}
	}

	// Convert tool definitions to OpenAI format
	var openAITools []ToolDef
	if len(request.ToolDefs) > 0 {
		// Convert tools from registry to OpenAI format
		for _, def := range request.ToolDefs {
			openAITools = append(openAITools, NewTool(def))
		}
	}

	if request.Config.WebSearch {
		// Add web search tool if enabled
		openAITools = append(openAITools, WebSearch)
	}

	// Create the request using our helper function
	req := NewRequest(request, openAIMessages, c.userID, openAITools)

	// Signal the start of processing
	if request.Callbacks.OnStart != nil {
		request.Callbacks.OnStart(request.Config)
	}

	// Start the streaming request
	stream := c.client.Responses.NewStreaming(ctx, req)
	defer stream.Close()

	// Stream and process the response
	var response responses.Response
	for stream.Next() {
		data := stream.Current()

		// Check if we have a complete response
		if data.Response.Status == responses.ResponseStatusCompleted {
			response = data.AsResponseCompleted().Response
			break
		}

		// Process streaming events
		c.handleStreamEvent(data, request)
	}

	// Check for errors in the stream
	if err := stream.Err(); err != nil {
		if request.Callbacks.OnError != nil {
			request.Callbacks.OnError(err)
		}
		// Signal completion even on error
		if request.Callbacks.OnDone != nil {
			request.Callbacks.OnDone(getUsageJSON(response.Usage))
		}
		return gptx.Response{}, fmt.Errorf("openai: %w", err)
	}

	if response.IncompleteDetails.Reason != "" {
		err := fmt.Errorf("openai: %s", response.IncompleteDetails.Reason)
		if request.Callbacks.OnError != nil {
			request.Callbacks.OnError(err)
		}
		// Signal completion even on error
		if request.Callbacks.OnDone != nil {
			request.Callbacks.OnDone(getUsageJSON(response.Usage))
		}
		return gptx.Response{}, err
	}

	// Process the complete response and extract data
	responseMessages, hasToolCalls := c.extractResponseData(ctx, &response, request)

	// Signal completion with usage information
	if request.Callbacks.OnDone != nil {
		request.Callbacks.OnDone(getUsageJSON(response.Usage))
	}

	// Return the response
	return gptx.Response{
		Messages:     responseMessages,
		Usage:        getUsageJSON(response.Usage),
		HasToolCalls: hasToolCalls,
	}, nil
}

// extractResponseData processes a completed response and extracts messages and tool call status.
func (c *OpenAIClient) extractResponseData(
	ctx context.Context,
	response *responses.Response,
	request gptx.Request,
) ([]gptx.Message, bool) {
	messages := []gptx.Message{}
	hasToolCalls := false

	// Process each output item
	for _, item := range response.Output {
		switch item.AsAny().(type) {
		case responses.ResponseOutputMessage:
			// Extract text from the message
			for _, content := range item.AsMessage().ToParam().Content {
				if content.OfOutputText != nil {
					// Add assistant message
					messages = append(messages, gptx.Message{
						Role:    "assistant",
						Content: content.OfOutputText.Text,
					})

					// Send text via callback
					if request.Callbacks.OnReply != nil {
						request.Callbacks.OnReply(content.OfOutputText.Text)
					}
				}
				if content.OfRefusal != nil {
					// Add refusal as assistant message
					messages = append(messages, gptx.Message{
						Role:    "assistant",
						Content: content.OfRefusal.Refusal,
					})

					// Send refusal via callback
					if request.Callbacks.OnReply != nil {
						request.Callbacks.OnReply(content.OfRefusal.Refusal)
					}
				}
			}

		case responses.ResponseFunctionWebSearch:
			// Signal that a web search is happening
			request.Callbacks.OnWebSearch()
			hasToolCalls = true

		case responses.ResponseFunctionToolCall:
			// Process tool call
			toolCall := item.AsFunctionCall()
			hasToolCalls = true

			if request.ToolHandler != nil {
				// Execute the tool
				result, err := request.ToolHandler(ctx, toolCall.Name, toolCall.Arguments)
				if err != nil {
					// Handle tool execution error
					errorMsg := fmt.Sprintf("Error executing tool %s: %s", toolCall.Name, err.Error())
					if request.Callbacks.OnError != nil {
						request.Callbacks.OnError(err)
					}

					// Add error as a message
					messages = append(messages, gptx.Message{
						Role:    "system",
						Content: errorMsg,
					})
				} else {
					// Add successful tool result as a message
					messages = append(messages, gptx.Message{
						Role:    "tool",
						Content: result,
						Name:    toolCall.Name,
					})

					// Send result via callback
					if request.Callbacks.OnReply != nil && result != "" {
						request.Callbacks.OnReply(result)
					}
				}
			}

		case responses.ResponseReasoningItem:
			// Handle reasoning output
			reasoning := item.AsReasoning()
			if request.Callbacks.OnReasoning != nil {
				for _, step := range reasoning.Summary {
					request.Callbacks.OnReasoning(step.Text)
				}
			}
		}
	}

	return messages, hasToolCalls
}
