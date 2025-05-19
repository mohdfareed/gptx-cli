// Package cfg handles configuration from CLI flags, env vars, and .gptx files.
// Order: CLI flags > env vars > local .gptx > parent .gptx > global config
package cfg

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

// SYS_PROMPT is the default system prompt template.
const SYS_PROMPT string = `
You are '%s', a CLI app. You are an extension of the command line.
You behave and respond like a command line tool. Be concise.
`

// Config stores application configuration settings.
type Config struct {
	APIKey    string   // OpenAI API key
	Model     string   // Model name
	SysPrompt string   // System prompt
	Files     []string // Attached files
	WebSearch bool     // Enable web search
	Shell     string   // Shell command
	Reason    bool     // Enable reasoning
	Tokens    int      // Max tokens
	Temp      float64  // Temperature (controls randomness)
}

// MARK: Flags
// ============================================================================

// Flags returns the CLI flags for the model configuration.
func (c *Config) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name: "key", Usage: "Set Platform API key",
			Category: "config", Destination: &c.APIKey,
			Sources:  cli.EnvVars(EnvVarPrefix + "API_KEY"),
			Required: true,
		},
		&cli.StringFlag{
			Name: "model", Usage: "Select model to use",
			Category: "config", Destination: &c.Model,
			Sources: cli.EnvVars(EnvVarPrefix + "MODEL"),
			Value:   "o4-mini",
		},
		// CONFIG
		&cli.BoolFlag{
			Name: "reason", Usage: "Allow the model to reason",
			Category: "config", Destination: &c.Reason,
			Sources:     cli.EnvVars(EnvVarPrefix + "REASON"),
			HideDefault: true,
		},
		// CONFIG
		&cli.IntFlag{
			Name: "max", Usage: "Limit response length",
			Category: "config", Destination: &c.Tokens,
			Sources:     cli.EnvVars(EnvVarPrefix + "MAX_TOKENS"),
			HideDefault: true,
		},
		&cli.Float64Flag{
			Name: "temp", Usage: "Set response randomness (0-100)",
			Category: "config", Destination: &c.Temp,
			Sources: cli.EnvVars(EnvVarPrefix + "TEMP"),
			Value:   1,
		},
		// CONTEXT
		&cli.StringFlag{
			Name: "prompt", Usage: "Set system prompt",
			Category: "config", Destination: &c.SysPrompt,
			Sources: cli.EnvVars(EnvVarPrefix + "INSTRUCTIONS"),
			Value:   fmt.Sprintf(SYS_PROMPT, AppName), Aliases: []string{"s"},
			TakesFile: true, Action: c.resolveSysPrompt, HideDefault: true,
		},
		&cli.StringSliceFlag{
			Name: "files", Usage: "Attach files to the message",
			Category: "context", Destination: &c.Files,
			Sources: cli.EnvVars(EnvVarPrefix + "FILES"),
			Value:   []string{}, Aliases: []string{"f"},
			TakesFile: true, Action: c.resolveFiles,
		},
		// TOOLS
		&cli.BoolFlag{
			Name: "web", Usage: "Enable web search",
			Category: "context", Destination: &c.WebSearch,
			Sources: cli.EnvVars(EnvVarPrefix + "WEB_SEARCH"),
		},
		&cli.StringFlag{
			Name: "shell", Usage: "Set the shell for the model to use",
			Category: "context", Destination: &c.Shell,
			Sources: cli.EnvVars(EnvVarPrefix + "SHELL"),
		},
	}
}

// MARK: Helper Actions
// ============================================================================

// Support reading a file for the system prompt.
func (c *Config) resolveSysPrompt(
	_ context.Context, cmd *cli.Command, prompt string,
) error {
	// load prompt from file if path is provided
	if _, err := os.Stat(prompt); err == nil {
		file, err := os.ReadFile(prompt)
		if err != nil {
			return fmt.Errorf("system prompt: %w", err)
		}
		c.SysPrompt = string(file)
	}
	return nil
}

// Support path globbing for file attachments.
func (c *Config) resolveFiles(
	_ context.Context, cmd *cli.Command, paths []string,
) error {
	var files []string
	for _, path := range paths {
		// Handle file globbing for the path part
		matches, err := filepath.Glob(path)
		if err != nil {
			return fmt.Errorf("file pattern %q: %w", path, err)
		}
		files = append(files, matches...)
	}
	c.Files = files
	return nil
}
