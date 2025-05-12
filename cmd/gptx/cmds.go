package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mohdfareed/gptx-cli/pkg/gptx"
	"github.com/urfave/cli/v3"
)

func mainCMD() *cli.Command {
	return &cli.Command{
		Name: gptx.AppName, Usage: "OpenAI models CLI",
		Description:    APP_DESC,
		DefaultCommand: "msg",
	}
}

func msgCMD(config *gptx.Config) *cli.Command {
	var msg string
	return &cli.Command{
		Name: "msg", Usage: "message an OpenAI model",
		Description: MSG_DESC,
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name: "message", UsageText: "the message to send",
				Value: "hello world", Destination: &msg,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			cfg, err := json.MarshalIndent(config, "", "  ")
			if err != nil {
				return err
			}
			println(msg)
			println(string(cfg))
			return nil
		},
	}
}

func configCMD(config *gptx.Config) *cli.Command {
	return &cli.Command{
		Name: "cfg", Usage: "the app's config",
		Description: CONFIG_DESC,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			cfg := *config
			if len(cfg.SysPrompt) > 61 {
				cfg.SysPrompt = cfg.SysPrompt[:61-3] + "..."
			} // truncate
			if len(cfg.APIKey) > 65 {
				cfg.APIKey = cfg.APIKey[:65-3] + "..."
			} // truncate
			data, err := json.MarshalIndent(cfg, "", "  ")
			if err != nil {
				return err
			}
			println(string(data))
			return nil
		},
	}
}

func demoCMD() *cli.Command {
	return &cli.Command{
		Name: "demo", Usage: "demo the app ui",
		Description: DEMO_DESC,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			modelPrefix("o4-mini", "demo-chat")
			println("Hello, world!")

			// Show all logging levels with different types of input
			Error("This is an error message")
			Error(fmt.Errorf("this is an error from an error object"))

			Warn("This is a warning message")
			Warn(fmt.Errorf("this is a warning from an error object"))

			Info("This is an info message")
			Info(fmt.Errorf("this is an info from an error object"))

			Debug("This is a debug message (visible with --verbose)")
			Debug(fmt.Errorf("this is a debug from an error object (visible with --verbose)"))

			// Show formatting capabilities
			Info("You can include %s with %d parameters", "formatting", 2)

			return nil
		},
	}
}
