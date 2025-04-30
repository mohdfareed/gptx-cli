package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// The app's configuration.
type AppConfig struct {
	APIKey    string `json:"OPENAI_API_KEY"`
	Model     string `json:"OPENAI_MODEL"`
	SysPrompt string `json:"OPENAI_SYS_PROMPT"`
}

// The default system prompt.
var DefaultSysPrompt string = `
You are a terminal tool that communicates following shell standards.
Be as concise as a terminal tool is expected to be.
`

// Load the app's configuration.
func LoadConfig() (*AppConfig, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	// return default if no config exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return defaultConfig(&AppConfig{}), nil // config doesn't exist
	} else if err != nil {
		return nil, err
	}

	// open config file
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// decode config json
	var config *AppConfig
	if err := json.NewDecoder(f).Decode(&config); err != nil {
		log.Print(fmt.Errorf("error decoding config: %w", err))
		config = defaultConfig(&AppConfig{})
	} // create default if file is invalid
	return config, nil
}

// Store the app configuration.
func StoreConfig(config AppConfig) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	// create file if no config exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if _, err = os.Create(path); err != nil {
			return err
		}
	}

	// open config file
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// encode app config
	file, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	if _, err := f.Write(append(file, '\n')); err != nil {
		return err
	}
	return nil
}

// The default app configuration.
func defaultConfig(config *AppConfig) *AppConfig {
	if config.APIKey == "" {
		config.APIKey = os.Getenv("OPENAI_API_KEY")
	}
	if config.Model == "" {
		config.Model = os.Getenv("OPENAI_MODEL")
	}
	if config.SysPrompt == "" {
		config.SysPrompt = os.Getenv("OPENAI_SYS_PROMPT")
	}

	if config.Model == "" {
		config.Model = "gpt-4o-mini"
	}
	if config.SysPrompt == "" {
		config.SysPrompt = DefaultSysPrompt
	}
	return config
}

// The configuration file path.
func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// use ~/.chatgpt.json if it exists
	configPath := filepath.Join(home, ".chatgpt.json")
	if _, err := os.Stat(configPath); err == nil {
		return configPath, nil // already exists
	} else if !os.IsNotExist(err) {
		return "", err // use default by ignoring `NotExists` errors
	}

	// default to $XDG_CONFIG_HOME/chatgpt-cli/config.json
	if dir, err := os.UserConfigDir(); err == nil {
		configDir := filepath.Join(dir, "chatgpt-cli")
		if err := os.MkdirAll(configDir, 0700); err == nil {
			configPath = filepath.Join(configDir, "config.json")
		} // fallback to home config path
	}
	return configPath, nil
}
