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
		DefaultCommand: "msg",
	}
}

func msgCMD(config *gptx.Config) *cli.Command {
	var msg string
	return &cli.Command{
		Name: "msg", Usage: "message an OpenAI model",
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
		Action: func(ctx context.Context, cmd *cli.Command) error {
			modelPrefix("o4-mini", "demo-chat")
			println("Hello, world!")

			errMsg(fmt.Errorf("an error message"))
			warnMsg(fmt.Errorf("a warning message"))
			infoMsg(fmt.Errorf("an info message"))
			debugMsg(fmt.Errorf("a debug message"))
			return nil
		},
	}
}
