// Package main provides the GPTx CLI application entry point.
package main

import (
	"context"
	"os"

	"github.com/mohdfareed/gptx-cli/internal/cfg"
	"github.com/urfave/cli/v3"
)

// main is the application entry point.
func main() {
	config := &cfg.Config{} // Initialize configuration
	cmd := mainCMD()        // Create CLI application

	cmd.Flags = append([]cli.Flag{
		debugFlag,
		silentFlag,
		editorFlag,
	}, config.Flags()...)

	cmd.Commands = []*cli.Command{
		msgCMD(config),
		configCMD(),
		demoCMD(),
	}

	// run the app
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		Error(err)
		os.Exit(1)
	}
}
