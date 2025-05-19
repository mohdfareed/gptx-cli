package main

import (
	"context"
	"fmt"

	"github.com/mohdfareed/gptx-cli/internal/cfg"
	"github.com/mohdfareed/gptx-cli/internal/events"
	"github.com/mohdfareed/gptx-cli/internal/tools"
	"github.com/mohdfareed/gptx-cli/pkg/gptx"
	"github.com/mohdfareed/gptx-cli/pkg/openai"
)

// setupCallbacks configures the event callbacks for the CLI.
func setupCallbacks() events.Callbacks {
	return events.NewCallbacks().
		// Model lifecycle events
		WithStartHandler(func(config cfg.Config) {
			Debug("Model started")
		}).
		WithErrorHandler(func(err error) {
			Error("Model error: %s\n", err)
		}).
		WithDoneHandler(func(usage string) {
			Info("Model done. Usage: %s\n", usage)
		}).
		// Output events
		WithReplyHandler(func(text string) {
			Print(text)
		}).
		WithReasoningHandler(func(text string) {
			PrintErr(M+"Reasoning: %s\n"+Reset, text)
		}).
		// Tool events
		WithToolCallHandler(func(call tools.ToolCall) {
			PrintErr(M+"Tool call: %s(%s)\n"+Reset, call.Name, call.Params)
		}).
		WithToolResultHandler(func(result string) {
			PrintErr(M+"Tool result: %s\n"+Reset, result)
		}).
		Build()
}

// setupTools sets up the tool registry with built-in tools.
func setupTools(config cfg.Config, registry *tools.Registry) {
	// Register shell tool if enabled
	if config.Shell != "" {
		registry.Register(tools.NewShellTool(config))
	}
}

// createModel creates a new model with the given configuration.
func createModel(config cfg.Config) *gptx.Model {
	// Create the callbacks manager
	callbacks := setupCallbacks()

	// Create the tool registry
	registry := tools.NewRegistry()
	setupTools(config, registry)

	// Create and configure the client
	client := openai.NewOpenAIClient(config.APIKey)

	// Create the model
	model := gptx.NewModel(
		config,
		gptx.WithClient(client),
		gptx.WithCallbacks(callbacks),
	)

	return model
}

// runModel runs a conversation with the given model and prompt.
func runModel(ctx context.Context, config cfg.Config, prompt string) error {
	model := createModel(config)
	if err := model.Message(ctx, prompt); err != nil {
		return fmt.Errorf("model error: %w", err)
	}
	return nil
}
