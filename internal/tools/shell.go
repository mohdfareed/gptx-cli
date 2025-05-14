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

const ShellToolName = "shell"
const ShellToolDescription = `Execute shell commands.
Use this for file operations, system information, or any command-line tasks.
`

var shellTool = ToolDef{
	Name: ShellToolName,
	Desc: ShellToolDescription,
	Params: map[string]any{
		"shell": "auto",
		"cmd":   "",
	},
	Handler: shellHandler,
}

func ShellTool(config cfg.Config) ToolDef {
	tool := shellTool
	if config.Shell == "auto" {
		tool.Params["shell"] = getDefaultShell()
	} else {
		tool.Params["shell"] = config.Shell
	}
	return tool
}

func shellHandler(ctx context.Context, params map[string]any) (string, error) {
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
