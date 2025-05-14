package main

import (
	"fmt"
	"os"

	"github.com/mohdfareed/gptx-cli/internal/cfg"
	"github.com/urfave/cli/v3"
)

// MARK: Logging
// ============================================================================

// Format strings for log messages
var (
	errMsgStr   = Bold + R + "error: " + Reset + "%s\n"
	warnMsgStr  = Bold + Y + " warn: " + Reset + "%s\n"
	infoMsgStr  = Bold + B + " info: " + Reset + "%s\n"
	debugMsgStr = Bold + M + "debug: " + Reset + "%s\n"
)

// Print logs a message to stdout.
// Accepts both string format with args or an error object.
func Print(msg any, args ...any) {
	fmt.Fprintf(os.Stdout, "%s", logMsg(msg, args))
}

// Print logs a message to stderr.
// Accepts both string format with args or an error object.
func PrintErr(msg any, args ...any) {
	fmt.Fprintf(os.Stderr, "%s", logMsg(msg, args))
}

// Error logs an error message. These are always displayed.
// Accepts both string format with args or an error object.
func Error(msg any, args ...any) {
	fmt.Fprintf(os.Stderr, errMsgStr, logMsg(msg, args))
}

// Warn logs a warning message if the current log level permits.
// Accepts both string format with args or an error object.
func Warn(msg any, args ...any) {
	if silent {
		return
	}
	fmt.Fprintf(os.Stderr, warnMsgStr, logMsg(msg, args))
}

// Info logs an informational message if the current log level permits.
// Accepts both string format with args or an error object.
func Info(msg any, args ...any) {
	if silent {
		return
	}
	fmt.Fprintf(os.Stderr, infoMsgStr, logMsg(msg, args))
}

// Debug logs a debug message if the current log level permits.
// Accepts both string format with args or an error object.
func Debug(msg any, args ...any) {
	if silent || !verbose {
		return
	}
	fmt.Fprintf(os.Stderr, debugMsgStr, logMsg(msg, args))
}

// MARK: Flags
// ============================================================================

var (
	// currentLogLevel controls the current logging verbosity.
	verbose bool = false
	// silent suppresses all non-error output when true.
	silent bool = false
)

// CLI flags for controlling log output.
var (
	debugFlag = &cli.BoolFlag{
		Name:    "verbose",
		Usage:   "Show debug messages",
		Aliases: []string{"v"},
		Sources: cli.EnvVars(
			cfg.EnvVar(nil, "VERBOSE"), cfg.EnvVar(nil, "DEBUG"),
		),
		Destination: &verbose,
	}

	silentFlag = &cli.BoolFlag{
		Name:    "quiet",
		Usage:   "Show only error messages",
		Aliases: []string{"silent", "q"},
		Sources: cli.EnvVars(
			cfg.EnvVar(nil, "QUIET"), cfg.EnvVar(nil, "SILENT"),
		),
		Value:       !isTerm,
		Destination: &silent,
	}
)

// MARK: Helpers
// ============================================================================

func logMsg(msg any, args []any) string {
	var text string
	switch m := msg.(type) {
	case error:
		text = m.Error()
	case string:
		text = fmt.Sprintf(m, args...)
	default:
		text = fmt.Sprint(msg, args)
	}
	return text
}
