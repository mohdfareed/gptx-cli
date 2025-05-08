package main

import (
	"context"
	"os"

	"github.com/mohdfareed/gptx-cli/pkg/gptx"
	"github.com/urfave/cli/v3"
)

func main() {
	config := gptx.DefaultConfig()

	// load config files
	err := gptx.LoadConfig()
	if err != nil {
		errMsg(err)
		os.Exit(1)
	}

	// // setup tools
	// tools := gptx.ModelTools{
	// 	config: config,
	// 	tools: map[gptx.Tool]gptx.ModelTool{
	// 		openai.WebSearchTool.Name: WebSearchTool,
	// 	},
	// }

	// // create model
	// model := gptx.CreateModel(config, tools)

	// commands
	cmd := mainCMD()
	cmd.Commands = []*cli.Command{
		// validateCMD(&model),
		configCMD(config),
		demoCMD(),
	}
	cmd.EnableShellCompletion = true

	// run the app
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		errMsg(err)
		os.Exit(1)
	}
}
