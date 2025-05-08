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
		Name:  gptx.AppName,
		Usage: "message an OpenAI model",
		Flags: []cli.Flag{
			colorizeFlag,
			// &cli.StringFlag{
			// 	Name:    "editor",
			// 	Aliases: []string{"e"},
			// 	Usage:   "the prompt editor",
			// },
			// &cli.BoolFlag{
			// 	Name:    "stream",
			// 	Aliases: []string{"s"},
			// 	Usage:   "stream the model's response",
			// 	Value:   true,
			// },
		},
	}
}

func configCMD(config gptx.Config) *cli.Command {
	return &cli.Command{
		Name:  "config",
		Usage: "the app's config",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			data, err := json.MarshalIndent(config, "", "  ")
			if err != nil {
				return err
			}
			json := string(data)

			fmt.Println(json)
			return nil
		},
	}
}
