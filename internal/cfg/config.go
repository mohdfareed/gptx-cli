package cfg

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

// CATEGORY is the category for configuration flags.
const CATEGORY = "model config"

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
	Tools     []string `env:"tools"`
	Repo      string   `env:"tools_repo"`
	Shell     string   `env:"tools_shell"`
	Tokens    *int     `env:"max_tokens"`
	Temp      int      `env:"temperature"`
}

// MARK: Flags
// ============================================================================

// Flags returns the CLI flags for the model configuration.
func (c *Config) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name: "api-key", Usage: "Set Platform API key",
			Category: CATEGORY, Destination: &c.APIKey,
			Sources:  cli.EnvVars(EnvVar(c, "APIKey")),
			Required: true,
		},
		&cli.StringFlag{
			Name: "model", Usage: "Select model to use",
			Category: CATEGORY, Destination: &c.Model,
			Sources: cli.EnvVars(EnvVar(c, "Model")),
			Value:   "o4-mini",
		},
		&cli.StringFlag{
			Name: "sys-prompt", Usage: "Set system prompt",
			Category: CATEGORY, Destination: &c.SysPrompt,
			Sources: cli.EnvVars(EnvVar(c, "SysPrompt")),
			Value:   fmt.Sprintf(SYS_PROMPT, AppName), Aliases: []string{"s"},
			TakesFile: true, Action: c.resolveSysPrompt, HideDefault: true,
		},
		&cli.StringSliceFlag{
			Name: "files", Usage: "Attach files to the message",
			Category: CATEGORY, Destination: &c.Files,
			Sources: cli.EnvVars(EnvVar(c, "Files")),
			Value:   []string{}, Aliases: []string{"f"},
			TakesFile: true, Action: c.resolveFiles,
		},
		&cli.StringSliceFlag{
			Name: "tools", Usage: "Enable specific tools",
			Category: CATEGORY, Destination: &c.Tools,
			Sources: cli.EnvVars(EnvVar(c, "Tools")),
			Value:   []string{}, Aliases: []string{"t"},
		},
		&cli.StringFlag{
			Name: "shell", Usage: "Set shell for the model",
			Category: CATEGORY, Destination: &c.Shell,
			Sources: cli.EnvVars(EnvVar(c, "Shell")),
		},
		&cli.StringFlag{
			Name: "repo", Usage: "Root path for repository exploration",
			Category: CATEGORY, Destination: &c.Repo,
			Sources: cli.EnvVars(EnvVar(c, "Repo")),
		},
		&cli.IntFlag{
			Name: "max-tokens", Usage: "Limit response length",
			Category: CATEGORY, Destination: c.Tokens,
			Sources:     cli.EnvVars(EnvVar(c, "Tokens")),
			HideDefault: true,
		},
		&cli.IntFlag{
			Name: "temp", Usage: "Set response randomness (0-100)",
			Category: CATEGORY, Destination: &c.Temp,
			Sources: cli.EnvVars(EnvVar(c, "Temp")),
			Value:   1,
		},
	}
}

// MARK: Helper Actions
// ============================================================================

func (c *Config) resolveSysPrompt(
	ctx context.Context, cmd *cli.Command, prompt string,
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

func (c *Config) resolveFiles(
	ctx context.Context, cmd *cli.Command, paths []string,
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
