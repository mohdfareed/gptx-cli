package main

import (
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// MARK: Definitions
// ============================================================================

// The OpenAI client.
type Model struct {
	cli    *openai.Client
	config ModelConfig
}

// The model's configuration.
type ModelConfig struct {
	// The OpenAI API key.
	APIKey string `koanf:"api_key"`
	// The OpenAI model to use.
	Model string `koanf:"model"`
	// The system prompts to use. Combined with other sys prompts.
	SysPrompt string `koanf:"prompt"`
	// The paths to the files to attach to the message.
	Files []string `koanf:"files"`
	// The chat history path.
	Chat string `koanf:"chat"`
	// Whether to stream the response.
	Stream bool `koanf:"stream"`
}

// MARK: Interfaces
// ============================================================================

// Message an OpenAI model and return the response.
func MessageModel(message string) error {
	// load the model's configuration
	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	// send the message to the model
	model := CreateModel(*config)
	model.Send(message)
	return nil
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
func (m *Model) Send(message string) error {
	if m.config.Stream {
		return m.stream(message)
	}
	return m.generate(message)
}
