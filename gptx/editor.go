package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func Editor(prompt string) (string, error) {
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
	editor := os.Getenv("EDITOR")
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

func Terminal() (string, error) {
	println("press ctrl-D to submit")
	reader := bufio.NewReader(os.Stdin)
	var sb strings.Builder

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("reading input: %w", err)
		}
		sb.WriteString(line)
	}
	return sb.String(), nil
}
