package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"
)

// The app name.
const AppName string = "gptx"

func main() {
	// load config
	config, err := LoadConfig()
	if err != nil {
		exit(err)
	}

	// setup tools
	tools := ModelTools{
		config: config,
		tools: map[ModelTool]Tool{
			WebSearchTool.Name: WebSearchTool,
		},
	}

	// create model
	model := CreateModel(config, tools)
	cmd := MsgCMD(&model)
	cmd.Commands = []*cli.Command{
		ValidateCMD(&model),
		ConfigCMD(config),
		EditChatCMD(config),
		ToolsCMD(tools),
		UsageCMD(config),
		DemoCMD(),
	}
	cmd.EnableShellCompletion = true

	// run the app
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		exit(err)
	}
}
