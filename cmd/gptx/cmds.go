// Package main implements the GPTx CLI commands.
package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/mohdfareed/gptx-cli/internal/cfg"
	"github.com/urfave/cli/v3"
)

// mainCMD creates the top-level CLI command.
func mainCMD() *cli.Command {
	return &cli.Command{
		Name: cfg.AppName, Usage: "Interact with an LLM models",
		Description:    APP_DESC,
		DefaultCommand: "msg",
	}
}

// msgCMD creates the message command for model interaction.
func msgCMD(config *cfg.Config) *cli.Command {
	var msg []string
	return &cli.Command{
		Name: "msg", Usage: "Send a message to a model",
		Description: MSG_DESC,
		Arguments: []cli.Argument{
			&cli.StringArgs{
				Name: "prompt", UsageText: "Message to send",
				Value: "hello world", Destination: &msg, Max: -1, // Accept multiple arguments
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Get the user prompt from command line args, stdin, or editor
			prompt, err := PromptUser(*config, msg)
			if err != nil {
				return fmt.Errorf("prompt: %w", err)
			}

			// Run the model with the prompt
			if err := runModel(ctx, *config, prompt); err != nil {
				return err
			}
			return nil
		},
	}
}

func configCMD() *cli.Command {
	return &cli.Command{
		Name: "cfg", Usage: "Show current configuration",
		Description: CONFIG_DESC,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Convert config to a map for display
			configMap := cfg.EnvMap()

			// Print each key-value pair
			for _, key := range configMap {
				keyName := Bold + Dim + key + Reset
				value := configMap[key]

				// Multiline values need quotes (escape existing quotes)
				if strings.Contains(value, "\n") {
					quote := Y + "\"" + Reset
					str := strings.ReplaceAll(value, "\"", "\\\"")
					value = quote + str + quote
				}

				Print(keyName + M + "=" + Reset + value + "\n")
			}

			// Show source files
			PrintErr(Dim + Bold + "Config Files:\n" + Reset)
			for _, file := range cfg.ConfigFiles() {
				PrintErr("- %s\n", file)
			}
			return nil
		},
	}
}

func demoCMD() *cli.Command {
	return &cli.Command{
		Name: "demo", Usage: "Show UI demonstration",
		Description: DEMO_DESC,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			modelPrefix("o4-mini", "demo-chat")
			Print("Hello, world!\n")
			Print("This is a demo of the gptx CLI.\n")

			Error("this is an error message")
			Warn("this is a warning message")
			Info("this is an info message")
			Debug("this is a debug message (visible with --verbose)")
			return nil
		},
	}
}
