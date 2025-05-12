// Package main provides CLI documentation text
package main

const (
	// APP_DESC is the main command description
	APP_DESC = `GPTx is a CLI app for interacting with OpenAI models.

Features:
- Send messages to various OpenAI models
- Attach files and file snippets using the tagging system
- Configure model parameters and behavior
- Support for multiple configuration methods (env vars, config files, flags)
- Extensible with tools and plugins

For detailed usage instructions, run 'gptx help <command>'`

	// MSG_DESC is the description for the msg command
	MSG_DESC = `Send a message to an OpenAI model and get a response.

You can include file contents directly in your message using tags:
  @file(path) - Include the entire file
  @file(path:start-end) - Include specific lines from a file

Examples:
  gptx msg "What does @file(main.go) do?"
  gptx msg "Explain @file(pkg/gptx/config.go:10-30)"`

	// CONFIG_DESC is the description for the config command
	CONFIG_DESC = `Display or modify the current configuration.

The configuration is loaded from multiple sources in the following order:
1. Default values
2. Global config file in application directory
3. .gptx files in the current directory and parent directories
4. Environment variables
5. Command-line flags`

	// DEMO_DESC is the description for the demo command
	DEMO_DESC = `Demonstrate the application UI and logging capabilities.

This command shows how the application looks and feels, including:
- Terminal colors and formatting
- Message prefixes
- Different logging levels`
)
