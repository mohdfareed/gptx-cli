package main

import (
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
func Serialize(model any) string {
	parser := koanf.New(".")
	parser.Load(structs.Provider(model, "koanf"), nil)
	return parser.Sprint()
}

// create a config parser with the following order:
// Defaults, XDG, parents, cwd, env vars.
func createParser() (*koanf.Koanf, error) {
	parser := koanf.New(".")

	// set defaults
	parser.Set("api_key", os.Getenv("OPENAI_API_KEY"))
	parser.Set("model", "gpt-4o-mini")
	parser.Set("prompt", strings.Trim(DefaultSysPrompt, "\n"))

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
		f := filepath.Join(dir, ".gptx")

		if _, err := os.Stat(f); err == nil {
			files = append(files, f)
		}

		if dir == filepath.Dir(dir) {
			break // reached root
		}
	}
	return files
}
