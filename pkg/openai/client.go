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

// Generate starts a conversation with the model using streaming responses.
// Implements the gptx.Client interface.
func (c *OpenAIClient) Generate(ctx context.Context, request gptx.Request) error {
	// Create user message with prompt and attached files
	msg, err := UserMsg(request.Prompt, request.Config.Files)
	if err != nil {
		return fmt.Errorf("openai: %w", err)
	}
	msgs := []MsgData{msg}

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
	req := NewRequest(request, msgs, c.userID, openAITools)

	// Start the streaming request
	stream := c.client.Responses.NewStreaming(ctx, req)
	defer stream.Close()

	// Signal the start of processing
	if request.Callbacks.OnStart != nil {
		request.Callbacks.OnStart(request.Config)
	}

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

	// Check for errors
	if err := stream.Err(); err != nil {
		if request.Callbacks.OnError != nil {
			request.Callbacks.OnError(err)
		}
		// Signal completion even on error
		if request.Callbacks.OnDone != nil {
			request.Callbacks.OnDone(getUsageJSON(response.Usage))
		}
		return fmt.Errorf("openai: %w", err)
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
		return err
	}

	// Process the complete response
	if err := c.handleCompleteResponse(ctx, &response, request); err != nil {
		return err
	}

	// Signal completion with usage information
	if request.Callbacks.OnDone != nil {
		request.Callbacks.OnDone(getUsageJSON(response.Usage))
	}
	return nil
}
