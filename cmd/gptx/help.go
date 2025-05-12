// Package main provides CLI documentation text
package main

const (
	// APP_DESC is the main command description
	APP_DESC = `Interact with OpenAI models from your terminal.

Features:
- Send messages to OpenAI models
- Include file contents using the @file tag system
- Configure model parameters
- Use multiple configuration methods
- Extend with tools and plugins

Learn more about a command:
    gptx help <command>`

	// MSG_DESC is the description for the msg command
	MSG_DESC = `Send a message to an OpenAI model.

Include file contents with tags:
    @file(path)            Include entire file
    @file(path:start-end)  Include specific lines

Examples:
    gptx msg "What does @file(main.go) do?"
    gptx msg "Explain @file(config.go:10-30)"`

	// CONFIG_DESC is the description for the config command
	CONFIG_DESC = `Display the configuration.

Configuration sources (in order of precedence):
1. Command-line flags
2. Environment variables
3. .gptx files (current directory, then parents)
4. Global config file
5. Default values

Output is in dotenv format suitable for config files.

Examples:
    # Save current configuration to a file
    gptx config > .gptx

    # Create a configuration with specific files
    gptx --files="*.go" config > .gptx

    # Create a project-specific configuration
    gptx --model="gpt-4o" --files="project/*.go" config > project/.gptx`

	// DEMO_DESC is the description for the demo command
	DEMO_DESC = `Demonstrate the UI and logging capabilities.

Shows:
- Terminal formatting
- Message prefixes
- Logging levels`
)
