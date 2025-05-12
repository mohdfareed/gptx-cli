package main

import (
	"fmt"
	"os"

	"github.com/mohdfareed/gptx-cli/pkg/gptx"
	"github.com/urfave/cli/v3"
)

// MARK: Exit Codes
// ============================================================================

type ExitCode int

const (
	ExitCodeOK ExitCode = iota
	ExitCodeError
	ExitCodeConfigError
	ExitCodeModelError
)

func exit(code ExitCode) {
	os.Exit(int(code))
}

// MARK: Logging
// ============================================================================

var verbose bool = false // debug output
var silent bool = false  // no output (except errors)

var debugFlag = &cli.BoolFlag{
	Name:        "verbose",
	Usage:       "output debug messages",
	Sources:     cli.EnvVars(gptx.EnvVar("VERBOSE")),
	Destination: &verbose,
}

var quietFlag = &cli.BoolFlag{
	Name:        "silent",
	Usage:       "suppress all output",
	Sources:     cli.EnvVars(gptx.EnvVar("SILENT")),
	Value:       !isTerm,
	Destination: &silent,
}

var (
	errMsgStr   = Bold + R + "error: " + Reset + "%s\n"
	warnMsgStr  = Bold + Y + " warn: " + Reset + "%s\n"
	infoMsgStr  = Bold + B + " info: " + Reset + "%s\n"
	debugMsgStr = Bold + M + "debug: " + Reset + "%s\n"
)

func errMsg(err error) {
	fmt.Fprintf(os.Stderr, errMsgStr, err.Error())
}

func warnMsg(err error) {
	if silent {
		return
	}
	fmt.Fprintf(os.Stderr, warnMsgStr, err.Error())
}

func infoMsg(err error) {
	if silent {
		return
	}
	fmt.Fprintf(os.Stderr, infoMsgStr, err.Error())
}

func debugMsg(err error) {
	if !verbose {
		return
	}
	fmt.Fprintf(os.Stderr, debugMsgStr, err.Error())
}
