package openai

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// TODO: validate model against [responses.ChatModel]

// Client represents an OpenAI API client
type Client struct {
	APIKey string
}

// Message represents a chat message
type Message struct {
	Role    string // "system", "user", "assistant"
	Content string
}

// FileAttachment represents a file attached to a message
type FileAttachment struct {
	Path    string
	Content string
}

// Tool represents a tool definition
type Tool struct {
	Type        string // "function" or "web_search"
	Name        string
	Description string
	Parameters  map[string]any
}

// ModelRequest represents a request to the OpenAI API
type ModelRequest struct {
	Model       string
	Messages    []Message
	Files       []FileAttachment
	Tools       []Tool
	Temperature float32
	MaxTokens   int
}

// ToolCall represents a tool call from the model
type ToolCall struct {
	ID         string
	Name       string
	Parameters json.RawMessage
}

// StreamResponse provides channels for streaming response content
type StreamResponse struct {
	Text      chan string
	ToolCalls chan ToolCall
	Done      chan struct{}
}

// Generate sends a request to OpenAI and returns a streaming response
func (c *Client) Generate(ctx context.Context, req ModelRequest) (*StreamResponse, error) {
	// Create stream channels
	stream := &StreamResponse{
		Text:      make(chan string, 10),
		ToolCalls: make(chan ToolCall, 5),
		Done:      make(chan struct{}),
	}

	// Prepare the request body
	apiReq := map[string]interface{}{
		"model": req.Model,
	}

	// Add messages
	var messages []map[string]interface{}
	for _, msg := range req.Messages {
		messages = append(messages, map[string]interface{}{
			"role":    msg.Role,
			"content": msg.Content,
		})
	}
	apiReq["messages"] = messages

	// Add file content to the last user message if needed
	if len(req.Files) > 0 && len(messages) > 0 {
		// Find the last user message
		var lastUserIdx int = -1
		for i := len(messages) - 1; i >= 0; i-- {
			if messages[i]["role"] == "user" {
				lastUserIdx = i
				break
			}
		}

		if lastUserIdx >= 0 {
			// Convert the content to an array
			userMsg := messages[lastUserIdx]
			userText := userMsg["content"].(string)

			content := []map[string]interface{}{
				{"type": "text", "text": userText},
			}

			// Add file content
			for _, file := range req.Files {
				content = append(content, map[string]interface{}{
					"type": "text",
					"text": file.Content,
				})
			}

			// Update the message with content array
			delete(userMsg, "content")
			userMsg["content"] = content
		}
	}

	// Add tools if specified
	if len(req.Tools) > 0 {
		var tools []map[string]interface{}

		for _, tool := range req.Tools {
			if tool.Type == "web_search" {
				tools = append(tools, map[string]interface{}{
					"type": "web_search",
				})
			} else if tool.Type == "function" {
				tools = append(tools, map[string]interface{}{
					"type": "function",
					"function": map[string]interface{}{
						"name":        tool.Name,
						"description": tool.Description,
						"parameters":  tool.Parameters,
					},
				})
			}
		}

		if len(tools) > 0 {
			apiReq["tools"] = tools
		}
	}

	// Add other parameters
	if req.Temperature > 0 {
		apiReq["temperature"] = req.Temperature
	}
	if req.MaxTokens > 0 {
		apiReq["max_tokens"] = req.MaxTokens
	}

	// Add streaming parameter
	apiReq["stream"] = true

	// Construct HTTP request
	reqBody, err := json.Marshal(apiReq)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://api.openai.com/v1/chat/completions",
		strings.NewReader(string(reqBody)),
	)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)

	// Start processing in a goroutine
	go func() {
		defer close(stream.Text)
		defer close(stream.ToolCalls)
		defer close(stream.Done)

		// Send the request
		client := &http.Client{}
		resp, err := client.Do(httpReq)
		if err != nil {
			fmt.Printf("error sending request: %v\n", err)
			return
		}
		defer resp.Body.Close()

		// Check for error response
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("error from OpenAI API: %s\n", string(body))
			return
		}

		// Process the stream
		reader := bufio.NewReader(resp.Body)
		toolCallsMap := make(map[string]*ToolCall) // Track partial tool calls by ID

		for {
			// Check if context is canceled
			select {
			case <-ctx.Done():
				return
			default:
				// Continue processing
			}

			// Read a line
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					fmt.Printf("error reading stream: %v\n", err)
				}
				return
			}

			// Skip empty lines
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			// Check for done signal
			if line == "data: [DONE]" {
				stream.Done <- struct{}{}
				return
			}

			// Parse the data
			if strings.HasPrefix(line, "data: ") {
				data := line[6:] // Remove "data: " prefix

				// Parse the JSON
				var chunk map[string]interface{}
				if err := json.Unmarshal([]byte(data), &chunk); err != nil {
					fmt.Printf("error parsing chunk: %v\n", err)
					continue
				}

				// Process choices
				if choices, ok := chunk["choices"].([]interface{}); ok && len(choices) > 0 {
					choice := choices[0].(map[string]interface{})

					// Process delta
					if delta, ok := choice["delta"].(map[string]interface{}); ok {
						// Handle text content
						if content, ok := delta["content"].(string); ok && content != "" {
							select {
							case stream.Text <- content:
							default: // Non-blocking
							}
						}

						// Handle tool calls
						if toolCalls, ok := delta["tool_calls"].([]interface{}); ok {
							for _, tc := range toolCalls {
								toolCall := tc.(map[string]interface{})

								// Get tool call ID and index
								id, _ := toolCall["id"].(string)
								index, _ := toolCall["index"].(float64)

								// Process function call parts
								if function, ok := toolCall["function"].(map[string]interface{}); ok {
									name, nameOk := function["name"].(string)
									args, argsOk := function["arguments"].(string)

									// Create or update tool call in our map
									if _, exists := toolCallsMap[id]; !exists {
										toolCallsMap[id] = &ToolCall{
											ID:         id,
											Name:       "",
											Parameters: nil,
										}
									}

									// Update the tool call with any new data
									if nameOk && name != "" {
										toolCallsMap[id].Name = name
									}

									if argsOk && args != "" {
										// Append to existing parameters if any
										if toolCallsMap[id].Parameters != nil {
											currentArgs := string(toolCallsMap[id].Parameters)
											toolCallsMap[id].Parameters = json.RawMessage(currentArgs + args)
										} else {
											toolCallsMap[id].Parameters = json.RawMessage(args)
										}

										// Check if tool call is complete
										if isValidJSON(string(toolCallsMap[id].Parameters)) &&
											toolCallsMap[id].Name != "" && index == 0 {

											// Send complete tool call
											select {
											case stream.ToolCalls <- *toolCallsMap[id]:
												// Remove from map after sending
												delete(toolCallsMap, id)
											default: // Non-blocking
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}()

	return stream, nil
}

// isValidJSON checks if a string is valid JSON
func isValidJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}
