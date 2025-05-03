package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
)

// The app name.
const AppName string = "gptx"

func main() {
	// load .env file
	_ = godotenv.Load()

	cmd := &cli.Command{
		Name:  "gptx",
		Usage: "message an OpenAI model",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			prompt, err := Editor("")
			if err != nil {
				prompt, err = Terminal()
				if err != nil {
					return err
				}
			}
			return MessageModel(strings.Trim(prompt, "\n"))
		},

		Commands: []*cli.Command{
			{
				Name:      "config",
				Usage:     "show the app's config",
				Arguments: []cli.Argument{},
				Action:    printConfig,
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func printConfig(ctx context.Context, cmd *cli.Command) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	configStr, err := Serialize(config)
	if err != nil {
		return err
	}
	println(configStr)
	return nil
}
