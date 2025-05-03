package main

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// The OpenAI client.
type Model struct {
	cli    *openai.Client
	config ModelConfig
}

// Create a new OpenAI model client.
func CreateModel(config ModelConfig) *Model {
	client := openai.NewClient(option.WithAPIKey(config.APIKey))
	return &Model{
		cli:    &client,
		config: config,
	}
}

// Generate or stream a model reply.
func (m *Model) Prompt(message string) error {
	// history := NewChat().Load(m.config.Chat)
	if m.config.Stream {
		return m.stream(message)
	}
	return m.generate(message)
}
