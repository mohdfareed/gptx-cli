// Package main implements the GPTx CLI.
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mohdfareed/gptx-cli/internal/cfg"
	"github.com/urfave/cli/v3"
)

// MARK: CLI Flags
// ============================================================================

// editorFlag for specifying text editor for input composition.
var editorFlag = &cli.StringFlag{
	Name:        "editor",
	Usage:       "Use specified text editor for input",
	Aliases:     []string{"e"},
	Sources:     cli.EnvVars(cfg.EnvVarPrefix+"EDITOR", "EDITOR"),
	Destination: &editor,
}

// editor command name
var editor string

// MARK: Prompt Handling
// ============================================================================

// PromptUser gets user input from args, editor, or terminal in that order.
func PromptUser(model string, args []string) (string, error) {
	if len(args) > 0 { // Message provided as command line arguments
		return strings.Join(args, " "), nil
	} else if editor != "" { // Editor specified, open it for composition
		return editorPrompt(editor)
	} else if isTerm { // Running in terminal, prompt interactively
		return terminalPrompt(model)
	}
	// No input method available
	return "", nil
}

// MARK: Editor Integration
// ============================================================================

// editorPrompt opens an external text editor for the user to compose a message.
// It creates a temporary file, launches the specified editor with that file,
// and then reads the contents after the editor closes.
//
// This function allows for a more comfortable editing experience when composing
// longer or more complex messages, taking advantage of the user's preferred
// text editor with all its features (syntax highlighting, keyboard shortcuts, etc.)
func editorPrompt(editor string) (string, error) {
	// Create a temporary file for the editor to use
	tmpDir := os.TempDir()
	tmp, err := os.CreateTemp(tmpDir, "chat-input-*.md")
	if err != nil {
		return "", fmt.Errorf("editor temp file: %w", err)
	}
	defer os.Remove(tmp.Name())

	// launch editor
	cmd := exec.Command(editor, tmp.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("editor %q: %w", editor, err)
	}

	// read the edited prompt
	raw, err := os.ReadFile(tmp.Name())
	if err != nil {
		return "", fmt.Errorf("editor temp file: %w", err)
	}
	return strings.TrimSpace(string(raw)), nil
}

// MARK: Terminal
// ============================================================================

func terminalPrompt(model string) (string, error) {
	modelPrefix(model, "")
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		} // exit on empty line (double enter)
		lines = append(lines, line)
	}
	prompt := strings.Join(lines, "\n")
	return prompt, nil
}

func modelPrefix(model string, chat string) {
	app := Dim + cfg.AppName + Reset
	model = Bold + G + model + Reset
	title := Bold + B + chat + Reset

	sep := Dim + "@" + Reset
	prefix := Dim + " ~> " + Reset
	postfix := Dim + " $ " + Reset

	if chat != "" {
		PrintErr(app + sep + model + prefix + title + postfix)
	} else {
		PrintErr(app + sep + model + postfix)
	}
}
