// Package main implements the command-line interface for GPTx.
// This file provides logging functionality with different severity levels
// and formatting options for consistent user feedback throughout the application.
package main

import (
	"fmt"
	"os"

	"github.com/mohdfareed/gptx-cli/internal/cfg"
	"github.com/urfave/cli/v3"
)

// MARK: Logging System
// ============================================================================

// Format strings for log messages with color coding for different severity levels.
// Each message type is formatted with a colored prefix to make it visually
// distinguishable in the terminal output. Colors are applied only when
// the output is a terminal (checked in cli.go).
var (
	errMsgStr   = Bold + R + "error: " + Reset + "%s\n" // Red for errors
	warnMsgStr  = Bold + Y + " warn: " + Reset + "%s\n" // Yellow for warnings
	infoMsgStr  = Bold + B + " info: " + Reset + "%s\n" // Blue for information
	debugMsgStr = Bold + M + "debug: " + Reset + "%s\n" // Magenta for debug messages
)

// Print logs a message to stdout without any formatting or severity prefix.
// It accepts either a string format with args or an error object.
// This is typically used for the main output of commands that are meant
// to be consumed by both users and other programs (e.g., model responses).
func Print(msg any, args ...any) {
	fmt.Fprintf(os.Stdout, "%s", logMsg(msg, args))
}

// PrintErr logs a message to stderr without any formatting or severity prefix.
// It accepts either a string format with args or an error object.
// This is useful for error messages that might be parsed by other programs
// without the color coding and prefixes of the Error() function.
func PrintErr(msg any, args ...any) {
	fmt.Fprintf(os.Stderr, "%s", logMsg(msg, args))
}

// Error logs an error message to stderr with red "error:" prefix.
// These messages are always displayed regardless of verbosity settings
// since they represent critical issues that must be addressed.
// It accepts either a string format with args or an error object.
func Error(msg any, args ...any) {
	fmt.Fprintf(os.Stderr, errMsgStr, logMsg(msg, args))
}

// Warn logs a warning message to stderr with yellow "warn:" prefix.
// Warning messages are suppressed when silent mode is enabled.
// These indicate potential issues that aren't blocking execution
// but may lead to unexpected behavior or reduced functionality.
// It accepts either a string format with args or an error object.
func Warn(msg any, args ...any) {
	if silent {
		return
	}
	fmt.Fprintf(os.Stderr, warnMsgStr, logMsg(msg, args))
}

// Info logs an informational message to stderr with blue "info:" prefix.
// Info messages are suppressed when silent mode is enabled.
// These provide additional context about what's happening during execution,
// such as configuration details or processing steps.
// It accepts either a string format with args or an error object.
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
			cfg.EnvVarPrefix+"VERBOSE", cfg.EnvVarPrefix+"DEBUG",
		),
		Destination: &verbose,
	}

	silentFlag = &cli.BoolFlag{
		Name:    "quiet",
		Usage:   "Show only error messages",
		Aliases: []string{"silent", "q"},
		Sources: cli.EnvVars(
			cfg.EnvVarPrefix+"QUIET", cfg.EnvVarPrefix+"SILENT",
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
