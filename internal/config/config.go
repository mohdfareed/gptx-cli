package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/openai/openai-go/shared"
	"github.com/urfave/cli/v3"
)

var AppConfigDir string = func() string {
	if configDir, err := os.UserConfigDir(); err == nil {
		return filepath.Join(configDir, AppName)
	}
	return "." + AppName
}()

const AppConfigFile string = "." + AppName + ".env"

// The default system prompt.
const DefaultSysPrompt string = `
You are '%s', a CLI app. You are an extension of the command line.
You behave and respond like a command line tool. Be concise.
`

// The model's configuration.
type Config struct {
	APIKey       string                 `koanf:"api_key"`
	Model        shared.ChatModel       `koanf:"model"`
	SysPrompt    string                 `koanf:"prompt"`
	ReasonEffort shared.ReasoningEffort `koanf:"reason"`
	Tools        []ModelTool            `koanf:"tools"`
	Files        []string               `koanf:"files"`
	Chat         string                 `koanf:"chat"`
	Temp         float64                `koanf:"temp"`
	MaxTokens    int64                  `koanf:"max_tokens"`
	Stream       bool                   `koanf:"stream"`
	Editor       string                 `koanf:"editor"`
	Color        bool                   `koanf:"color"`
}

// Load the model's configuration in the following order:
// Defaults, XDG, parents, cwd, env vars, .env file.
func LoadConfig() (Config, error) {
	// create config parser
	parser, err := createParser()
	if err != nil {
		return Config{}, fmt.Errorf("config loader: %w", err)
	}

	// deserialize the config
	var config Config
	if err := parser.Unmarshal("", &config); err != nil {
		return Config{}, fmt.Errorf("config deserialization: %w", err)
	}
	config.SysPrompt = strings.Trim(config.SysPrompt, "\n")

	// combine tool names
	toolNames := make([]string, len(config.Tools))
	for i, tool := range config.Tools {
		toolNames[i] = string(tool)
	}
	tools := strings.Split(string(strings.Join(toolNames, " ")), " ")

	// parse the tools
	config.Tools = make([]ModelTool, len(tools))
	for i, tool := range tools {
		config.Tools[i] = ModelTool(tool)
	}

	// set theme
	if config.Color {
		Theme = RichTheme
	} else {
		Theme = PlainTheme
	}
	return config, nil
}

// create a config parser with the following order:
// defaults, xdg, parents, cwd, env.
func createParser() (*koanf.Koanf, error) {
	parser := koanf.New(".")

	// set defaults
	parser.Set("api_key", os.Getenv("OPENAI_API_KEY"))
	parser.Set("model", shared.ChatModelGPT4oMini)
	parser.Set("prompt", fmt.Sprintf(DefaultSysPrompt, AppName))
	parser.Set("stream", true)
	parser.Set("editor", os.Getenv("EDITOR"))

	// load config files
	var files []string = configFIles()
	for i := len(files) - 1; i >= 0; i-- {
		_ = parser.Load(file.Provider(files[i]), dotenv.Parser())
	}

	// load environment variables
	var envPrefix = strings.ToUpper(AppName) + "_"
	_ = parser.Load(env.Provider(envPrefix, ".", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, envPrefix))
	}), nil)
	return parser, nil
}

// Return the app config files in the following order: cwd, parents, xdg.
func configFIles() []string {
	var files []string
	for dir, _ := os.Getwd(); ; dir = filepath.Dir(dir) {
		f := filepath.Join(dir, "."+AppName+".env")

		if _, err := os.Stat(f); err == nil {
			files = append(files, f)
		}

		if dir == filepath.Dir(dir) {
			break // reached root
		}
	}

	// support $XDG_CONFIG_HOME and %APPDATA%
	if AppConfigDir != "" {
		f := filepath.Join(AppConfigDir, "config")
		if _, err := os.Stat(f); err == nil {
			files = append(files, f)
		}
	}
	return files
}

// MARK: CLI
// ============================================================================

func ConfigCMD(config Config) *cli.Command {
	return &cli.Command{
		Name: "config", Usage: "the app's config",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			parser := koanf.New(".")
			_ = parser.Load(structs.Provider(config, "koanf"), nil)

			data, err := dotenv.Parser().Marshal(parser.All())
			if err != nil {
				return fmt.Errorf("serialization: %w", err)
			}

			str := strings.ReplaceAll(string(data), "\\n", "\n")
			println(str)
			return nil
		},
	}
}
