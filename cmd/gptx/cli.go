package main

import (
	"context"
	"os"
	"strings"

	"github.com/mohdfareed/gptx-cli/internal/cfg"
	"github.com/mohdfareed/gptx-cli/internal/events"
	"github.com/mohdfareed/gptx-cli/internal/tools"
	"golang.org/x/term"
)

var isTerm bool = term.IsTerminal(int(os.Stdout.Fd()))

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

func init() {
	var colors = []*string{
		&Reset, &Bold, &Dim, &Black, &R, &G, &Y, &B, &M, &C, &White,
	}

	if !isTerm || os.Getenv("NO_COLOR") != "" {
		for _, color := range colors {
			*color = ""
		}
	}
}

// MARK: Configuration ========================================================

func formatKeyValue(key string, value string) string {
	quote := Y + "\"" + Reset
	equal := M + "=" + Reset
	keyName := Bold + Dim + key + Reset

	// Format value based on content
	if strings.Contains(value, "\n") {
		// Multiline values need quotes (escape existing quotes)
		value = quote + strings.ReplaceAll(value, "\"", "\\\"") + quote
	}
	return keyName + equal + value
}

func printModelEvent(mgr events.ModelEvents) {
	mgr.Start.Subscribe(context.TODO(), func(data cfg.Config) {
		Debug("Model started. Config: %v", data)
	})
	mgr.Reply.Subscribe(context.TODO(), func(data string) {
		Print(data)
	})
	mgr.Reasoning.Subscribe(context.TODO(), func(data string) {
		Info("Reasoning: %s", data)
	})
	mgr.ToolCall.Subscribe(context.TODO(), func(data tools.ToolCall) {
		Info("Tool call: %s", data)
	})
	mgr.ToolResult.Subscribe(context.TODO(), func(data string) {
		Info("Tool result: %s", data)
	})
	mgr.Error.Subscribe(context.TODO(), func(err error) {
		Error("Model error: %s", err)
	})
	mgr.Done.Subscribe(context.TODO(), func(usage string) {
		Debug("Model done. Usage: %s", usage)
	})
}
