package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mohdfareed/gptx-cli/pkg/gptx"
	"github.com/urfave/cli/v3"
)

// MARK: Exit Codes
// ============================================================================

// ExitCode defines application exit codes.
type ExitCode int

const (
	// ExitCodeOK indicates successful execution.
	ExitCodeOK ExitCode = iota
	// ExitCodeError indicates a general error.
	ExitCodeError
	// ExitCodeConfigError indicates a configuration error.
	ExitCodeConfigError
	// ExitCodeModelError indicates an error from the AI model.
	ExitCodeModelError
)

// exit terminates the application with the specified exit code.
func exit(code ExitCode) {
	os.Exit(int(code))
}

// MARK: Logging
// ============================================================================

// LogLevel defines the verbosity level of logging.
type LogLevel int

const (
	// LogLevelError only shows errors.
	LogLevelError LogLevel = iota
	// LogLevelWarn shows errors and warnings.
	LogLevelWarn
	// LogLevelInfo shows errors, warnings, and info messages.
	LogLevelInfo
	// LogLevelDebug shows all messages including debug information.
	LogLevelDebug
)

var (
	// currentLogLevel controls the current logging verbosity.
	currentLogLevel LogLevel = LogLevelInfo
	// silent suppresses all non-error output when true.
	silent bool = false
)

// CLI flags for controlling log output.
var (
	debugFlag = &cli.BoolFlag{
		Name:    "verbose",
		Usage:   "output debug messages",
		Sources: cli.EnvVars(gptx.EnvVar("VERBOSE")),
		Action: func(ctx context.Context, cmd *cli.Command, v bool) error {
			if v {
				currentLogLevel = LogLevelDebug
			}
			return nil
		},
	}

	quietFlag = &cli.BoolFlag{
		Name:        "silent",
		Usage:       "suppress all output except errors",
		Sources:     cli.EnvVars(gptx.EnvVar("SILENT")),
		Value:       !isTerm,
		Destination: &silent,
	}
)

// Format strings for log messages
var (
	errMsgStr   = Bold + R + "error: " + Reset + "%s\n"
	warnMsgStr  = Bold + Y + " warn: " + Reset + "%s\n"
	infoMsgStr  = Bold + B + " info: " + Reset + "%s\n"
	debugMsgStr = Bold + M + "debug: " + Reset + "%s\n"
)

// Error logs an error message. These are always displayed.
// Accepts both string format with args or an error object.
func Error(msg any, args ...any) {
	var text string
	switch m := msg.(type) {
	case error:
		text = m.Error()
	case string:
		text = fmt.Sprintf(m, args...)
	default:
		text = fmt.Sprint(msg)
	}
	fmt.Fprintf(os.Stderr, errMsgStr, text)
}

// Warn logs a warning message if the current log level permits.
// Accepts both string format with args or an error object.
func Warn(msg any, args ...any) {
	if silent || currentLogLevel < LogLevelWarn {
		return
	}

	var text string
	switch m := msg.(type) {
	case error:
		text = m.Error()
	case string:
		text = fmt.Sprintf(m, args...)
	default:
		text = fmt.Sprint(msg)
	}
	fmt.Fprintf(os.Stderr, warnMsgStr, text)
}

// Info logs an informational message if the current log level permits.
// Accepts both string format with args or an error object.
func Info(msg any, args ...any) {
	if silent || currentLogLevel < LogLevelInfo {
		return
	}

	var text string
	switch m := msg.(type) {
	case error:
		text = m.Error()
	case string:
		text = fmt.Sprintf(m, args...)
	default:
		text = fmt.Sprint(msg)
	}
	fmt.Fprintf(os.Stderr, infoMsgStr, text)
}

// Debug logs a debug message if the current log level permits.
// Accepts both string format with args or an error object.
func Debug(msg any, args ...any) {
	if silent || currentLogLevel < LogLevelDebug {
		return
	}

	var text string
	switch m := msg.(type) {
	case error:
		text = m.Error()
	case string:
		text = fmt.Sprintf(m, args...)
	default:
		text = fmt.Sprint(msg)
	}
	fmt.Fprintf(os.Stderr, debugMsgStr, text)
}
