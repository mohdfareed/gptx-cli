package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mohdfareed/gptx-cli/pkg/gptx"
)

// MARK: Prompt
// ============================================================================

// PromptUser gets user message. Input is retrieved in the following order:
// user input -> editor -> terminal
func PromptUser(
	config gptx.Config, msgs []string, editor string) (string, error) {
	var prompt string
	var err error

	if len(msgs) > 0 { // user input provided
		prompt = strings.Join(msgs, " ")
	} else if editor != "" { // editor specified
		prompt, err = editorPrompt(editor)
	} else {
		prompt, err = terminalPrompt(config)
	}
	return strings.TrimSpace(prompt), err
}

// MARK: Terminal
// ============================================================================

func terminalPrompt(config gptx.Config) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	print(modelPrefix(config.Model, ""))

	prompt, err := reader.ReadString('\n')
	if err != nil {
		println("Error reading prompt:", err)
		return prompt, err
	} // FIXME: handle shift-enter
	return prompt, nil
}

// MARK: Editor
// ============================================================================

func editorPrompt(editor string) (string, error) {
	// create a temp file
	tmpDir := os.TempDir()
	tmp, err := os.CreateTemp(tmpDir, "chat-input-*.md")
	if err != nil {
		return "", fmt.Errorf("creating temp file: %w", err)
	}
	defer os.Remove(tmp.Name())

	// launch editor
	cmd := exec.Command(editor, tmp.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("running editor %q: %w", editor, err)
	}

	// read the edited prompt
	raw, err := os.ReadFile(tmp.Name())
	if err != nil {
		return "", fmt.Errorf("reading temp file: %w", err)
	}
	return strings.TrimSpace(string(raw)), nil
}
