package gptx

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// ToolDefinition defines a tool that can be used by the model.
type ToolDefinition struct {
	Name        string         // Tool name
	Description string         // Tool description
	Schema      map[string]any // JSON schema for parameters
	Command     string         // Shell command for custom tools
}

// RegisterBuiltinTools adds built-in tools to the model.
func RegisterBuiltinTools(m *Model) {
	// Shell execution tool
	m.RegisterTool("shell", func(ctx context.Context, input []byte) ([]byte, error) {
		var params struct {
			Command string `json:"command"`
		}

		if err := json.Unmarshal(input, &params); err != nil {
			return nil, fmt.Errorf("parse input: %w", err)
		}

		// Execute command
		cmd := exec.CommandContext(ctx, "sh", "-c", params.Command)
		output, err := cmd.CombinedOutput()

		// Prepare result
		result := map[string]string{
			"output": string(output),
		}

		if err != nil {
			result["error"] = err.Error()
		}

		// Marshal to JSON
		response, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("marshal result: %w", err)
		}

		return response, nil
	})
}

// RegisterCustomTool registers a custom tool with the model.
func RegisterCustomTool(m *Model, def ToolDefinition) {
	m.RegisterTool(def.Name, func(ctx context.Context, input []byte) ([]byte, error) {
		// Execute command with input provided to stdin
		cmd := exec.CommandContext(ctx, "sh", "-c", def.Command)
		cmd.Stdin = strings.NewReader(string(input))
		output, err := cmd.CombinedOutput()

		// If output is valid JSON, return it directly
		if json.Valid(output) {
			return output, nil
		}

		// Otherwise wrap in a simple response
		result := map[string]string{
			"output": string(output),
		}

		if err != nil {
			result["error"] = err.Error()
		}

		response, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("marshal result: %w", err)
		}

		return response, nil
	})
}

// GetShellToolSchema returns the schema for the shell tool.
func GetShellToolSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"command": map[string]any{
				"type":        "string",
				"description": "Shell command to execute",
			},
		},
		"required": []string{"command"},
	}
}

// containsTool checks if a tool is in the tools list.
func ContainsTool(tools []string, name string) bool {
	for _, tool := range tools {
		if tool == name {
			return true
		}
	}
	return false
}
