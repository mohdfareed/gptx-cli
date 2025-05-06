package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

// MARK: Logging
// ============================================================================

func exit(err error) {
	print(Theme.Bold + Theme.Red) // style panic prefix
	println(fmt.Errorf("error: %w"+Theme.Reset, err).Error())
	os.Exit(1)
}

func warn(err error) {
	print(Theme.Bold + Theme.Yellow) // style warning prefix
	println(fmt.Errorf("warning: %w"+Theme.Reset, err).Error())
}

// MARK: Views
// ============================================================================

// func printMsg(msg Msg) {
// 	// msg.
// }

// MARK: Model Prefix
// ============================================================================

func ModelPrefix(model string, chat string) string {
	app := Theme.Dim + AppName + Theme.Reset
	model = Theme.Bold + Theme.Green + model + Theme.Reset
	title := Theme.Bold + Theme.Blue + chat + Theme.Reset

	sep := Theme.Dim + "@" + Theme.Reset
	prefix := Theme.Dim + " ~> " + Theme.Reset
	postfix := Theme.Dim + " $ " + Theme.Reset

	if chat == "" {
		return fmt.Sprintf(
			Theme.Bold+"%s%s%s%s"+Theme.Reset, app, sep, model, postfix,
		)
	} else {
		return fmt.Sprintf(Theme.Bold+"%s%s%s%s%s%s"+Theme.Reset,
			app, sep, model, prefix, title, postfix,
		)
	}
}

// MARK: Demo
// ============================================================================

func DemoCMD() *cli.Command {
	return &cli.Command{
		Name: "demo", Usage: "demo the app ui",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			icons := []string{
				Theme.AppIcon, Theme.ChatIcon, Theme.HistoryIcon,
				Theme.SettingsIcon, Theme.ToolIcon, Theme.WebSearchIcon,
				Theme.ShellIcon, Theme.FileIcon, Theme.FolderIcon,
				Theme.RemoteIcon, Theme.On, Theme.Off, Theme.Success,
				Theme.Error, Theme.Warning, Theme.Info, Theme.Unknown,
				Theme.Debug,
			}

			colors := []string{
				Theme.Reset, Theme.Bold, Theme.Dim,
				Theme.Black, Theme.White,
				Theme.Red, Theme.Green, Theme.Yellow,
				Theme.Blue, Theme.Magenta, Theme.Cyan,
			}

			msg := "Hello, world!"
			fmt.Println(ModelPrefix("o4-mini", "first-chat") + msg)
			for _, color := range colors {
				for _, icon := range icons {
					fmt.Print(color + icon + Theme.Reset)
				}
				fmt.Println()
			}
			return nil
		},
	}
}
