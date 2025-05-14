package main

import (
	"context"
	"os"

	"github.com/mohdfareed/gptx-cli/internal/cfg"
	"github.com/urfave/cli/v3"
)

func main() {
	config := &cfg.Config{} // configuration
	cmd := mainCMD()        // cli app

	cmd.Flags = append([]cli.Flag{
		colorizeFlag,
		debugFlag,
		silentFlag,
		editorFlag,
	}, config.Flags()...)

	cmd.Commands = []*cli.Command{
		msgCMD(config),
		configCMD(config),
		demoCMD(),
	}

	// run the app
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		Error(err)
		exit(ModelErrorCode)
	}
}
