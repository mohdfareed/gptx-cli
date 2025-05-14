package gptx

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/mohdfareed/gptx-cli/pkg/openai"
)

// Model handles interactions with OpenAI.
type Model struct {
	Config *Config
	Events *Event
	tools  map[string]ToolFunc
}

// ToolFunc is a function that executes a tool and returns its result.
type ToolFunc func(context.Context, []byte) ([]byte, error)

// NewModel creates a new model instance with the given configuration.
func NewModel(config *Config) *Model {
	return &Model{
		Config: config,
		Events: NewEvent(),
		tools:  make(map[string]ToolFunc),
	}
}

// RegisterTool registers a tool with the model.
func (m *Model) RegisterTool(name string, tool ToolFunc) {
	m.tools[name] = tool
}

// Message sends a message to the model and streams the response.
func (m *Model) Message(ctx context.Context, prompt string, output io.Writer) error {
	prompt, err := ProcessTags(prompt)
	if err != nil {
		return fmt.Errorf("prompt: %w", err)
	}
	// Emit start event
	m.Events.Emit(ctx, EventStart, m.Config)

	// Create request with messages
	req := openai.ModelRequest{
		Model: m.Config.Model,
		Messages: []openai.Message{
			{Role: "system", Content: m.Config.SysPrompt},
			{Role: "user", Content: prompt},
		},
		Temperature: float32(m.Config.Temp) / 100.0,
	}

	// Add max tokens if specified
	if m.Config.Tokens != nil {
		req.MaxTokens = *m.Config.Tokens
	}

	// // Add files if specified
	// if len(m.Config.Files) > 0 {
	// 	files, err := LoadFiles(m.Config.Files)
	// 	if err != nil {
	// 		m.Events.Emit(ctx, EventError, err)
	// 		return fmt.Errorf("load files: %w", err)
	// 	}
	// 	req.Files = files
	// }

	// Add tools if specified
	if len(m.Config.Tools) > 0 {
		req.Tools = m.getTools()
	}

	// Send request to OpenAI
	client := &openai.Client{APIKey: m.Config.APIKey}
	stream, err := client.Generate(ctx, req)
	if err != nil {
		m.Events.Emit(ctx, EventError, err)
		return fmt.Errorf("generate: %w", err)
	}

	// Process the response stream
	return m.processStream(ctx, stream, output)
}

// getTools creates tool definitions from config.
func (m *Model) getTools() []openai.Tool {
	var tools []openai.Tool

	for _, name := range m.Config.Tools {
		switch name {
		case "web_search":
			tools = append(tools, openai.Tool{
				Type: "web_search",
				Name: "web_search",
			})
		default:
			// Add custom tools if registered
			if _, ok := m.tools[name]; ok {
				tools = append(tools, openai.Tool{
					Type: "function",
					Name: name,
				})
			}
		}
	}

	return tools
}

// processStream handles the streaming response.
func (m *Model) processStream(ctx context.Context, stream *openai.StreamResponse, output io.Writer) error {
	done := make(chan error, 1)

	go func() {
		defer close(done)

		for {
			select {
			case text, ok := <-stream.Text:
				if !ok {
					// Channel closed, we're done
					m.Events.Emit(ctx, EventComplete, nil)
					done <- nil
					return
				}

				// Emit reply event and write to output
				m.Events.Emit(ctx, EventReply, text)
				if _, err := fmt.Fprint(output, text); err != nil {
					done <- fmt.Errorf("write output: %w", err)
					return
				}

			case toolCall, ok := <-stream.ToolCalls:
				if !ok {
					continue
				}

				// Emit tool event
				m.Events.Emit(ctx, EventTool, map[string]any{
					"name":   toolCall.Name,
					"params": string(toolCall.Parameters),
				})

				// Execute tool if registered
				if tool, ok := m.tools[toolCall.Name]; ok {
					result, err := tool(ctx, toolCall.Parameters)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Tool error: %v\n", err)
					}

					// Emit tool result
					m.Events.Emit(ctx, EventTool, map[string]any{
						"name":   toolCall.Name,
						"result": string(result),
					})
				}

			case <-ctx.Done():
				done <- ctx.Err()
				return
			}
		}
	}()

	return <-done
}
