package main

import (
	"fmt"
)

func ChatCMD(message string, _ []string) (string, error) {
	// 1) load app config
	config, err := LoadConfig()
	if err != nil {
		return "", fmt.Errorf("error loading config: %w", err)
	}

	// 2) call OpenAI
	client := NewClient(config.APIKey, config.Model)
	resp, err := client.Prompt(message)
	if err != nil {
		return "", fmt.Errorf("chat error: %w", err)
	}
	return fmt.Sprint(resp), nil
}

func ConfigCMD() (string, error) {
	configPath, err := configPath()
	if err != nil {
		return "", err
	}

	config, err := LoadConfig()
	if err != nil {
		return "", err
	}

	return "Path: " + configPath + "\n" + fmt.Sprint(config), nil
}

func CreateConfigCMD() (string, error) {
	config, err := LoadConfig()
	if err != nil {
		return "", err
	}
	err = StoreConfig(*config)
	if err != nil {
		return "", err
	}
	configPath, err := configPath()
	if err != nil {
		return "", err
	}
	return "Config written to: " + configPath, nil
}
