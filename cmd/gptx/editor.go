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

// MARK: CLI
// ============================================================================

var editorFlag = &cli.StringFlag{
	Name:        "editor",
	Usage:       "Use specified text editor for input",
	Aliases:     []string{"e"},
	Sources:     cli.EnvVars(cfg.EnvVar(nil, "EDITOR"), "EDITOR"),
	Destination: &editor,
}
var editor string

// MARK: Prompt
// ============================================================================

// PromptUser gets user message. Input is retrieved in the following order:
// user input -> editor -> terminal -> no input
func PromptUser(config cfg.Config, args []string) (string, error) {
	if len(args) > 0 { // user input provided
		return strings.Join(args, " "), nil
	} else if editor != "" { // editor specified
		return editorPrompt(editor)
	} else if isTerm { // running in terminal
		return terminalPrompt(config)
	}
	return "", nil
}

// MARK: Editor
// ============================================================================

func editorPrompt(editor string) (string, error) {
	// create a temp file
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

func terminalPrompt(config cfg.Config) (string, error) {
	modelPrefix(config.Model, "")
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		lines = append(lines, line)
	} // FIXME: can't exit
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
