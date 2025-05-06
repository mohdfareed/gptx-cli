package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/openai/openai-go/responses"
	"github.com/urfave/cli/v3"
)

// MARK: Models
// ============================================================================

type MsgData = responses.ResponseInputItemUnionParam

// Msg represents a single turn in the conversation.
type Msg struct {
	Data  MsgData                  `json:"data"`
	Usage *responses.ResponseUsage `json:"usage,omitempty"`
}

// History holds an ordered list of Messages.
type History struct {
	Title string `json:"title"`
	Msgs  []Msg  `json:"messages"`
	path  string
}

// MARK: History
// ============================================================================

// LoadChat loads the chat history from the given path.
func LoadChat(path string) (History, error) {
	title := filepath.Base(path)
	chat := History{Title: title, Msgs: []Msg{}, path: path}

	// create the history file if it doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := chat.Save(); err != nil {
			return History{}, fmt.Errorf("history: %w", err)
		}
	}

	// load the history file
	if err := chat.load(); err != nil {
		return History{}, fmt.Errorf("history: %w", err)
	}
	return chat, nil
}

// Appends a new message with the given role and content.
func (h *History) Add(msg Msg) {
	h.Msgs = append(h.Msgs, msg)
}

// The last n messages in the history.
func (h *History) Last(n int) []Msg {
	if n > len(h.Msgs) {
		n = len(h.Msgs)
	} // check if n is greater than the length of the history
	return h.Msgs[len(h.Msgs)-n:]
}

// MARK: File I/O
// ============================================================================

// Writes the history out as JSON (indented) to the given file.
func (h *History) Save() error {
	if h.path == "" {
		return nil
	} // check if the path is empty

	if err := os.MkdirAll(filepath.Dir(h.path), 0o755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	} // create path struct
	f, err := os.Create(h.path)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	} // create history file
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(h)
}

func (h *History) load() error {
	if h.path == "" {
		return nil
	} // check if the path is empty

	// read history file
	f, err := os.Open(h.path)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	// deserialize history
	if err := json.NewDecoder(f).Decode(h); err != nil {
		return fmt.Errorf("loading history: %w", err)
	}
	return nil
}

// MARK: CLI
// ============================================================================

func EditChatCMD(config Config) *cli.Command {
	return &cli.Command{
		Name: "edit", Usage: "edit the chat history",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if config.Chat == "" {
				return fmt.Errorf("no chat history file")
			} // check if the path is empty
			if config.Editor == "" {
				println("config path:" + Theme.Bold + config.Chat + Theme.Reset)
				return nil
			} // check if editor is provided

			// launch editor
			println(
				Theme.Dim+"running: "+Theme.Reset,
				Theme.Green+Theme.Bold+config.Editor+Theme.Reset,
				Theme.Bold+config.Chat+Theme.Reset,
			)
			editor := exec.Command(config.Editor, config.Chat)
			editor.Stdout = os.Stdout
			editor.Stderr = os.Stderr
			editor.Stdin = os.Stdin
			return editor.Run()
		},
	}
}
