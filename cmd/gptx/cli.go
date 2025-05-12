package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mohdfareed/gptx-cli/pkg/gptx"
	"github.com/urfave/cli/v3"
	"golang.org/x/term"
)

var isTerm bool = term.IsTerminal(int(os.Stdout.Fd()))

type ExitCode int

const (
	ErrorCode ExitCode = iota
	ConfigErrorCode
	ModelErrorCode
)

func exit(code ExitCode) {
	os.Exit(int(code))
}

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

func deColorize() {
	colors := []*string{
		&Reset, &Bold, &Dim, &Black, &R, &G, &Y, &B, &M, &C, &White,
	}
	for _, color := range colors {
		*color = ""
	}
}

// MARK: CLI ==================================================================

var colorizeFlag = &cli.StringFlag{
	Name:    "color",
	Usage:   "colorize output, one of: auto, always, never",
	Value:   "auto",
	Sources: cli.EnvVars(gptx.EnvVar(nil, "COLORIZE"), "NO_COLOR"),
	Validator: func(value string) error {
		switch value {
		case "auto":
			if !(isTerm && os.Getenv("NO_COLOR") == "") {
				deColorize()
			}
		case "never":
			deColorize()
		case "always":
		default:
			Error(fmt.Errorf("invalid color option: %s", value))
			exit(ErrorCode)
		}
		return nil
	},
	ValidateDefaults: true,
}

// MARK: Configuration ========================================================

// FormatKeyValue formats a key-value pair with appropriate styling.
func FormatKeyValue(key string, value string) string {
	quote := Y + "\"" + Reset

	// Format value based on content
	if strings.Contains(value, "\n") {
		// Multiline values need quotes (escape existing quotes)
		value = quote + strings.ReplaceAll(value, "\"", "\\\"") + quote
	}

	// Return the formatted key=value pair
	return fmt.Sprintf("%s=%s", Bold+Dim+key+Reset, value+Reset)
}
