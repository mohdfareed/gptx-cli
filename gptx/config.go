package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

// serialize the config
func Serialize(model any) (string, error) {
	parser := koanf.New(".")
	_ = parser.Load(structs.Provider(model, "koanf"), nil)

	data, err := dotenv.Parser().Marshal(parser.All())
	if err != nil {
		return "", fmt.Errorf("config serialization: %w", err)
	}

	str := strings.ReplaceAll(string(data), "\\n", "\n")
	return str, nil
}

// Load the model's configuration in the following order:
// Defaults, XDG, parents, cwd, env vars, .env file.
func LoadConfig() (*ModelConfig, error) {
	// create config parser
	parser, err := createParser()
	if err != nil {
		return nil, fmt.Errorf("config loader: %w", err)
	}

	// deserialize the config
	var config ModelConfig
	if err := parser.Unmarshal("", &config); err != nil {
		return nil, fmt.Errorf("config deserialization: %w", err)
	}
	return &config, nil
}

// create a config parser with the following order:
// Defaults, XDG, parents, cwd, env vars.
func createParser() (*koanf.Koanf, error) {
	parser := koanf.New(".")

	// set defaults
	parser.Set("api_key", os.Getenv("OPENAI_API_KEY"))
	parser.Set("model", "gpt-4o-mini")
	// parser.Set("prompt", strings.Trim(DefaultSysPrompt, "\n"))
	parser.Set("prompt", strings.Trim("", "\n"))

	// load config files
	var files []string = configFIles()
	for i := len(files) - 1; i >= 0; i-- {
		_ = parser.Load(file.Provider(files[i]), dotenv.Parser())
	}

	// support $XDG_CONFIG_HOME
	if configDir, err := os.UserConfigDir(); err == nil {
		f := filepath.Join(configDir, AppName, "config")
		if _, err := os.Stat(f); err == nil {
			_ = parser.Load(file.Provider(f), dotenv.Parser())
		}
	}

	// load environment variables
	var envPrefix = strings.ToUpper(AppName) + "_"
	_ = parser.Load(env.Provider(envPrefix, ".", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, envPrefix))
	}), nil)
	return parser, nil
}

// return the config files in the cwd and its parents, in that order
func configFIles() []string {
	var files []string
	for dir, _ := os.Getwd(); ; dir = filepath.Dir(dir) {
		f := filepath.Join(dir, "."+AppName)

		if _, err := os.Stat(f); err == nil {
			files = append(files, f)
		}

		if dir == filepath.Dir(dir) {
			break // reached root
		}
	}
	return files
}
