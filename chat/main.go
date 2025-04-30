package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cmd := &cli.Command{
		Name:     "chat",
		Usage:    "chat with an OpenAI model",
		Commands: commands(),
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func commands() []*cli.Command {

	var message string
	var attachments []string // TODO: implement

	prompt := cli.Command{
		Name:    "prompt",
		Aliases: []string{"p"},
		Usage:   "send a prompt to the model",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:        "message",
				UsageText:   "the prompt message",
				Destination: &message,
			},
		},
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:        "attachments",
				Usage:       "files to attach to the message",
				Destination: &attachments,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			result, err := ChatCMD(message, attachments)
			if err != nil {
				return err
			}
			println(result)
			return nil
		},
	}

	config := cli.Command{
		Name:      "config",
		Aliases:   []string{"c"},
		Usage:     "show the app's config",
		Arguments: []cli.Argument{},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			result, err := ConfigCMD()
			if err != nil {
				return err
			}
			println(result)
			return nil
		},
	}

	writeConfig := cli.Command{
		Name:    "write-config",
		Aliases: []string{"w"},
		Usage:   "write the app config to file",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			result, err := CreateConfigCMD()
			if err != nil {
				return err
			}
			println(result)
			return nil
		},
	}

	return []*cli.Command{
		&prompt,
		&config,
		&writeConfig,
	}
}
