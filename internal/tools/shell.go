// Package tools implements model tools for system interaction.
package tools

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/mohdfareed/gptx-cli/internal/cfg"
)

// ShellToolDef is the shell tool identifier.
const ShellToolDef = "shell"

// ShellToolDescription describes the shell tool's purpose for the model.
const ShellToolDescription = `Execute shell commands.
Use this for file operations, system information, or any command-line tasks.
`

// NewShellTool creates a shell tool from the given config.
func NewShellTool(config cfg.Config) ToolDef {
	// Determine which shell to use
	shell := getDefaultShell()
	if config.Shell != "auto" {
		shell = config.Shell
	}

	// Create the tool definition
	return ToolDef{
		Name: ShellToolDef,
		Desc: ShellToolDescription,
		Params: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"cmd": map[string]any{
					"type":        "string",
					"description": "The command to execute",
				},
			},
			"required": []string{"cmd"},
		},
		Handler: func(ctx context.Context, params map[string]any) (string, error) {
			// Extract command from params
			cmd, ok := params["cmd"].(string)
			if !ok {
				return "", fmt.Errorf("shell: missing required parameter 'cmd'")
			}

			// Execute the command with the configured shell
			return shellHandler(map[string]any{
				"shell": shell,
				"cmd":   cmd,
			})
		},
	}
}

// shellHandler implements the shell tool functionality.
// It executes a shell command and returns the output or an error.
//
// Parameters:
// - shell: The shell to use (bash, zsh, powershell, etc.)
// - cmd: The command to execute
//
// Returns:
// - The command output as a string
// - An error if the command fails or the shell is not available
func shellHandler(params map[string]any) (string, error) {
	shell := params["shell"].(string)
	cmd := params["cmd"].(string)

	// Check if the shell is available
	if _, err := exec.LookPath(shell); err != nil {
		return "", fmt.Errorf("shell not found: %s", shell)
	}

	// Execute the command
	out, err := exec.Command(shell, "-c", cmd).Output()
	if err != nil {
		return "", fmt.Errorf("command execution failed: %w", err)
	}

	// Return the output as a string
	return string(out), nil
}

// getDefaultShell returns the default shell based on the OS
func getDefaultShell() string {
	// Check for SHELL environment variable
	if shell, ok := os.LookupEnv("SHELL"); ok && shell != "" {
		return filepath.Base(shell) // Extract just the name
	}

	// Default shells based on OS
	if runtime.GOOS == "windows" {
		return "pwsh.exe"
	}

	// Default to bash for Unix-like systems
	return "bash"
}
