# GPTx CLI

GPTx CLI is a command-line interface for OpenAI's GPT models built in Go.
It allows you to interact with any GPT model directly from the terminal.

## References

- https://go.dev/doc
- https://github.com/openai/openai-go
- https://cli.urfave.org/v3
- https://github.com/knadh/koanf/v2

## TODO

- Remove `koanf` dependency
- Add chat command
    - Subcommands:
        - Set session history file (default to temporary file)
        - Print session file path, contents, stats
    - Each chat is a json file
- Add msg command
    - Prompts the user for a message, allowing for multiline input
    - Flags for config options
- Support custom model tools
    - Shell commands
