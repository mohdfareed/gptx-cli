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

// Config is the model's configuration.
type Config struct {
	APIKey    string   `env:"api_key"`
	Model     string   `env:"model"`
	SysPrompt string   `env:"sys_prompt"`
	Files     []string `env:"files"`
	WebSearch *bool    `env:"web_search"`
	Shell     *string  `env:"shell_tool"`
	Tokens    *int     `env:"max_tokens"`
	Temp      *float64 `env:"temperature"`
}

// MARK: Flags
// ============================================================================

// Flags returns the CLI flags for the model configuration.
func (c *Config) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name: "key", Usage: "Set Platform API key",
			Category: "config", Destination: &c.APIKey,
			Sources:  cli.EnvVars(EnvVar(c, "APIKey")),
			Required: true,
		},
		&cli.StringFlag{
			Name: "model", Usage: "Select model to use",
			Category: "config", Destination: &c.Model,
			Sources: cli.EnvVars(EnvVar(c, "Model")),
			Value:   "o4-mini",
		},
		&cli.IntFlag{
			Name: "max", Usage: "Limit response length",
			Category: "config", Destination: c.Tokens,
			Sources:     cli.EnvVars(EnvVar(c, "Tokens")),
			HideDefault: true,
		},
		&cli.Float64Flag{
			Name: "temp", Usage: "Set response randomness (0-100)",
			Category: "config", Destination: c.Temp,
			Sources: cli.EnvVars(EnvVar(c, "Temp")),
			Value:   1,
		},
		&cli.StringFlag{
			Name: "prompt", Usage: "Set system prompt",
			Category: "config", Destination: &c.SysPrompt,
			Sources: cli.EnvVars(EnvVar(c, "SysPrompt")),
			Value:   fmt.Sprintf(SYS_PROMPT, AppName), Aliases: []string{"s"},
			TakesFile: true, Action: c.resolveSysPrompt, HideDefault: true,
		},
		&cli.StringSliceFlag{
			Name: "files", Usage: "Attach files to the message",
			Category: "context", Destination: &c.Files,
			Sources: cli.EnvVars(EnvVar(c, "Files")),
			Value:   []string{}, Aliases: []string{"f"},
			TakesFile: true, Action: c.resolveFiles,
		},
		&cli.BoolFlag{
			Name: "web", Usage: "Enable web search",
			Category: "context", Destination: c.WebSearch,
			Sources: cli.EnvVars(EnvVar(c, "WebSearch")),
		},
		&cli.StringFlag{
			Name: "shell", Usage: "Set the shell for the model to use",
			Category: "context", Destination: c.Shell,
			Sources: cli.EnvVars(EnvVar(c, "Shell")),
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
