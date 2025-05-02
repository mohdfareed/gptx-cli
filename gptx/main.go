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

	var msgPrompt []string
	cmd := &cli.Command{
		Name:  "gptx",
		Usage: "message an OpenAI model",
		Arguments: []cli.Argument{
			&cli.StringArgs{
				Name:        "msg",
				UsageText:   "the message prompt",
				Destination: &msgPrompt,
				Max:         -1,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return MessageModel(strings.Join(msgPrompt, " "))
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
