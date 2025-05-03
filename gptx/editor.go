package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func Prompt(prompt string, editor string) (string, error) {
	if editor == "" {
		prompt, err := terminalPrompt()
		if err != nil {
			return "", err
		}
		return prompt, nil
	}

	prompt, err := editorPrompt(prompt, editor)
	if err != nil {
		return "", err
	}
	return prompt, nil
}

func terminalPrompt() (string, error) {
	fmt.Print("Prompt: ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("reading input: %w", err)
	}
	return line, nil
}

func editorPrompt(prompt string, editor string) (string, error) {
	// create a temp file
	tmpDir := os.TempDir()
	tmp, err := os.CreateTemp(tmpDir, "chat-input-*.md")
	if err != nil {
		return "", fmt.Errorf("creating temp file: %w", err)
	}
	defer os.Remove(tmp.Name())

	// write prompt
	if _, err := tmp.WriteString(prompt); err != nil {
		return "", fmt.Errorf("writing initial prompt: %w", err)
	}
	tmp.Close()

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
	return string(raw), nil
}
