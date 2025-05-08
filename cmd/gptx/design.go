package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
	"golang.org/x/term"
)

// MARK: Colors ===============================================================

var (
	Reset = "\033[0m"
	Bold  = "\033[1m"
	Dim   = "\033[2m"
	Black = "\033[30m"
	R     = "\033[31m"
	G     = "\033[32m"
	Y     = "\033[33m"
	B     = "\033[34m"
	M     = "\033[35m"
	C     = "\033[36m"
	White = "\033[37m"
)

// MARK: Symbols ==============================================================

var (
	AppIcon       = " "
	ChatIcon      = "󰭹 "
	HistoryIcon   = " "
	SettingsIcon  = " "
	ToolIcon      = " "
	WebSearchIcon = " "
	ShellIcon     = " "
	FileIcon      = " "
	FolderIcon    = " "
	RemoteIcon    = " "
)

// MARK: Status ===============================================================

var (
	On      = Bold + " " + Reset
	Off     = Dim + " " + Reset
	Success = G + " " + Reset
	Error   = R + " " + Reset
	Warning = Y + " " + Reset
	Info    = B + " " + Reset
	Debug   = M + " " + Reset
	Unknown = White + Dim + " " + Reset
)

// MARK: Definitions ==========================================================

var appColors = []*string{
	&Reset, &Bold, &Dim, &Black, &R, &G, &Y, &B, &M, &C, &White,
}

var appIcons = []*string{
	&AppIcon, &ChatIcon, &HistoryIcon, &SettingsIcon, &ToolIcon,
	&WebSearchIcon, &ShellIcon, &FileIcon, &FolderIcon, &RemoteIcon,
	&On, &Off, &Success, &Error, &Warning, &Info, &Debug, &Unknown,
}

func deColorize() {
	for _, color := range appColors {
		*color = ""
	}
	for _, icon := range appIcons {
		*icon = ""
	}

	On = "[ON] "
	Off = "[OFF] "
	Success = "[OK] "
	Error = "[ERROR] "
	Warning = "[WARN] "
	Info = "[INFO] "
	Debug = "[DEBUG] "
	Unknown = "[?] "
}

// MARK: CLI ==================================================================

// colorizeFlag is a flag that controls the output theme.
var colorizeFlag = &cli.StringFlag{
	Name:  "color",
	Usage: "colorize output, one of: auto, always, never",
	Value: "auto",
	Validator: func(value string) error {
		switch value {
		case "auto":
			isTerm := term.IsTerminal(int(os.Stdout.Fd()))
			if !(isTerm && os.Getenv("NO_COLOR") == "") {
				deColorize()
			}
		case "never":
			deColorize()
		case "always":
		default:
			return cli.Exit("invalid color value", 1)
		}
		return nil
	},
}

func demoCMD() *cli.Command {
	return &cli.Command{
		Name: "demo", Usage: "demo the app ui",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			msg := "Hello, world!"
			fmt.Println(modelPrefix("o4-mini", "first-chat") + msg)
			warnMsg(fmt.Errorf("a warning message"))
			errMsg(fmt.Errorf("an error message"))
			doneMsg("cli demo")
			if Reset == "" && AppIcon == "" {
				return nil
			}

			for _, color := range appColors {
				for _, icon := range appIcons {
					fmt.Print(*color + *icon + Reset)
				}
				println()
			}
			return nil
		},
	}
}
