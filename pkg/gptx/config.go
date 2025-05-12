package gptx

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
)

const CATEGORY = "model config"
const SYS_PROMPT string = `
You are '%s', a CLI app. You are an extension of the command line.
You behave and respond like a command line tool. Be concise.
`

// Config is the model's configuration.
type Config struct {
	APIKey    string   `json:"api_key"`
	Model     string   `json:"model"`
	SysPrompt string   `json:"sys_prompt"`
	Files     []string `json:"files"`
	Tools     []string `json:"tools"`
	Tokens    *int     `json:"max_tokens"`
	Temp      int      `json:"temperature"`
}

// MARK: CLI Flags & Env Vars
// ============================================================================

// Flags returns the CLI flags for the model configuration.
func (c *Config) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name: "api-key", Usage: "OpenAI API key",
			Category: CATEGORY, Destination: &c.APIKey,
			Sources:  cli.EnvVars(EnvVar("API_KEY"), "OPENAI_API_KEY"),
			Required: true,
		},
		&cli.StringFlag{
			Name: "model", Usage: "OpenAI model",
			Category: CATEGORY, Destination: &c.Model,
			Sources: cli.EnvVars(EnvVar("MODEL")),
			Value:   "o4-mini",
			// Required: true,
		},
		&cli.StringFlag{
			Name: "sys-prompt", Usage: "the system prompt",
			Category: CATEGORY, Destination: &c.SysPrompt,
			Sources: cli.EnvVars(EnvVar("INSTRUCTIONS")),
			Value:   fmt.Sprintf(SYS_PROMPT, AppName), Aliases: []string{"s"},
			TakesFile: true, Action: c.readSysPrompt, HideDefault: true,
		},
		&cli.StringSliceFlag{
			Name: "files", Usage: "files to attach",
			Category: CATEGORY, Destination: &c.Files,
			Sources: cli.EnvVars(EnvVar("FILES")),
			Value:   []string{}, Aliases: []string{"f"},
			TakesFile: true, Action: c.resolveFiles,
		},
		&cli.StringSliceFlag{
			Name: "tools", Usage: "tools to load",
			Category: CATEGORY, Destination: &c.Tools,
			Sources: cli.EnvVars(EnvVar("TOOLS")),
			Value:   []string{}, Aliases: []string{"t"},
		},
		&cli.IntFlag{
			Name: "max-tokens", Usage: "max output tokens",
			Category: CATEGORY, Destination: c.Tokens,
			Sources:     cli.EnvVars(EnvVar("MAX_TOKENS")),
			HideDefault: true,
		},
		&cli.IntFlag{
			Name: "temp", Usage: "model temperature",
			Category: CATEGORY, Destination: &c.Temp,
			Sources: cli.EnvVars(EnvVar("TEMPERATURE")),
			Value:   1,
		},
	}
}

func (c *Config) readSysPrompt(
	ctx context.Context, cmd *cli.Command, prompt string,
) error {
	println("sys-prompt:", prompt)
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
		matches, err := filepath.Glob(path)
		if err != nil {
			return fmt.Errorf("file pattern %q: %w", path, err)
		}
		files = append(files, matches...)
	}
	c.Files = files
	return nil
}

// MARK: Config Files
// ============================================================================

func init() {
	godotenv.Load(configFIles()...)
}

func configFIles() []string {
	var files []string // cwd, parents, app dir

	// load files from cwd then its parents
	for dir, _ := os.Getwd(); ; dir = filepath.Dir(dir) {
		f := filepath.Join(dir, "."+AppName)
		if _, err := os.Stat(f); err == nil {
			files = append(files, f)
		}

		if dir == filepath.Dir(dir) {
			break // reached root
		}
	}

	// support $XDG_CONFIG_HOME and %APPDATA%
	if AppDir != "" {
		f := filepath.Join(AppDir, "config")
		if _, err := os.Stat(f); err == nil {
			files = append(files, f)
		}
	}
	return files
}
