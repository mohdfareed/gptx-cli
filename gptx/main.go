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

	// load the config
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	// create the model
	model := CreateModel(config)

	// create the app
	cmd := &cli.Command{
		Name:  "gptx",
		Usage: "message an OpenAI model",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return msgModel(model)
		},

		Commands: []*cli.Command{
			{
				Name:      "config",
				Usage:     "show the app's config",
				Arguments: []cli.Argument{},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return printConfig(config)
				},
			},
		},
		EnableShellCompletion: true,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func msgModel(model *Model) error {
	prompt, err := Prompt("", model.config.Editor)
	if err != nil {
		return err
	}
	return model.Prompt(strings.Trim(prompt, "\n"))
}

func printConfig(config ModelConfig) error {
	configStr, err := Serialize(config)
	if err != nil {
		return err
	}
	println(configStr)
	return nil
}
