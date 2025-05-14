package main

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/mohdfareed/gptx-cli/pkg/gptx"
	"github.com/urfave/cli/v3"
)

func mainCMD() *cli.Command {
	return &cli.Command{
		Name: gptx.AppName, Usage: "Interact with OpenAI models",
		Description:    APP_DESC,
		DefaultCommand: "msg",
	}
}

func msgCMD(config *gptx.Config) *cli.Command {
	var msg string
	return &cli.Command{
		Name: "msg", Usage: "Send a message to a model",
		Description: MSG_DESC,
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name: "message", UsageText: "Message to send",
				Value: "hello world", Destination: &msg,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Create the model
			model := gptx.NewModel(config)
			for event := range model.Events.Channel {
				printModelEvent(event)
			}
			println() // Add newline after response

			// Listen for events
			go func(event *gptx.Event) {
				for event := range event.Channel {
					printModelEvent(event)
				}
			}(model.Events)

			// Send the message
			err := model.Message(ctx, msg, os.Stdout)
			if err != nil {
				return fmt.Errorf("model: %w", err)
			}
			return nil
		},
	}
}

func configCMD(config *gptx.Config) *cli.Command {
	return &cli.Command{
		Name: "cfg", Usage: "Show current configuration",
		Description: CONFIG_DESC,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Convert config to a map for display
			configMap := config.ToEnvMap()
			configMap["API Key"] = maskAPIKey(config.APIKey)
			configMap["System Prompt"] = shortenText(config.SysPrompt, 40)

			// Sort keys for consistent output
			keys := make([]string, 0, len(configMap))
			for k := range configMap {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			// Print each key-value pair
			for _, key := range keys {
				Print(formatKeyValue(key, configMap[key]) + "\n")
			}

			// Show source files
			PrintErr("\nConfiguration Files:\n")
			for _, file := range gptx.ConfigFiles() {
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
