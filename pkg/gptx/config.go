package gptx

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/joho/godotenv"
)

const AppName string = "gptx"
const DefaultSysPrompt string = `
You are '%s', a CLI app. You are an extension of the command line.
You behave and respond like a command line tool. Be concise.
`

var appDir string = func() string {
	configDir, _ := os.UserConfigDir()
	return filepath.Join(configDir, AppName)
}()

// Config is the model's configuration.
type Config struct {
	APIKey string   `koanf:"api_key"`
	Model  string   `koanf:"model"`
	Files  []string `koanf:"files"`
	Tools  []string `koanf:"tools"`
	Tokens int      `koanf:"max_tokens"`
	Temp   int      `koanf:"temperature"`
	Instr  string   `koanf:"instructions"`
}

// MARK: Defaults & Load
// ============================================================================

// DefaultConfig is the default configuration for the model.
func DefaultConfig() Config {
	return Config{
		Model:  "gpt-4o",
		Tokens: 16,
		Temp:   1,
		Instr:  fmt.Sprintf(DefaultSysPrompt, AppName),
	}
}

// LoadConfig loads the model's configuration from files and env.
// The config is loaded in the following order:
// env, dotenv (dev), cwd, parents, app dir
func LoadConfig() error {
	// load .env files
	err := godotenv.Load(configFIles()...)
	if err != nil {
		return fmt.Errorf("config files: %w", err)
	}
	return nil
}

// MARK: Fields & SysPrompt
// ============================================================================

// Fields returns the fields of the config struct.
func (c Config) Fields() []configField {
	var fields []configField
	t := reflect.TypeOf(c)
	v := reflect.ValueOf(c)

	for i := range t.NumField() {
		ft := t.Field(i)
		fv := v.Field(i)
		fields = append(fields, configField{
			Type: ft.Type, Value: fv,
		})
	}
	return fields
}

// SysPrompt returns the system prompt for the model.
func (c Config) SysPrompt(path string) (string, error) {
	// load prompt from file if path is provided
	if _, err := os.Stat(c.Instr); err == nil {
		data, err := os.ReadFile(c.Instr)
		if err != nil {
			return c.Instr, fmt.Errorf("config prompt: %w", err)
		}
		c.Instr = string(data)
	}
	return c.Instr, nil
}

// MARK: Helpers
// ============================================================================

type configField struct {
	reflect.Type
	reflect.Value
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
	if appDir != "" {
		f := filepath.Join(appDir, "config")
		if _, err := os.Stat(f); err == nil {
			files = append(files, f)
		}
	}
	return files
}
