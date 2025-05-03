package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Message represents a single turn in the conversation.
type Message struct {
	Sender  string `json:"sender"`
	Content string `json:"content"`
}

// History holds an ordered list of Messages.
type History struct {
	Title    string    `json:"title"`
	Messages []Message `json:"messages"`
	Size     int       `json:"size"`
}

// Create an empty chat History.
func NewChat() *History {
	return &History{Messages: []Message{}}
}

// Appends a new message with the given role and content.
func (h *History) Add(sender, content string) {
	h.Messages = append(h.Messages, Message{
		Sender:  sender,
		Content: content,
	})
}

// Returns up to the last n messages (or all if len < n || n <= 0).
func (h *History) Last(n int) []Message {
	total := len(h.Messages)
	if n <= 0 || n >= total {
		return h.Messages
	}
	return h.Messages[total-n:]
}

// Loads history from the given JSON file.
func (h *History) Load(path string) error {
	// read history file
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}
	defer f.Close()

	// deserialize history
	if err := json.NewDecoder(f).Decode(h); err != nil {
		return err
	}
	return nil
}

// Writes the history out as JSON (indented) to the given file.
func (h *History) Save(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	} // create path struct
	f, err := os.Create(path)
	if err != nil {
		return err
	} // create history file
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(h)
}
