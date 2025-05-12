package main

import (
	"context"
	"encoding/json"
	"fmt"
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
			cfg, err := json.MarshalIndent(config, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(msg)     // outputs to stdout
			println(string(cfg)) // print to console
			return nil
		},
	}
}

func configCMD(config *gptx.Config) *cli.Command {
	return &cli.Command{
		Name: "cfg", Usage: "Show current configuration",
		Description: CONFIG_DESC,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			envMap := config.ToEnvMap()

			// Sort keys for consistent output
			keys := make([]string, 0, len(envMap))
			for k := range envMap {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			// Print each key-value pair formatted
			for _, key := range keys {
				fmt.Println(FormatKeyValue(key, envMap[key]))
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
			println("Hello, world!")
			fmt.Println("This is a demo of the gptx CLI.")

			Error("this is an error message")
			Warn("this is a warning message")
			Info("this is an info message")
			Debug("this is a debug message (visible with --verbose)")
			return nil
		},
	}
}
