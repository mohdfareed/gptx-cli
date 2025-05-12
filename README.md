# GPTx CLI

GPTx CLI is a command-line interface for OpenAI's GPT models built in Go.
It allows you to interact with any GPT model directly from the terminal.

## Features

- **Configuration**
  - Git-like `.gptx` files in current directory and parent directories
  - Global config in application directory
  - Configuration via environment variables and CLI flags
  - Export configurations in dotenv format for easy sharing and reuse

- **File Integration**
  - Attach entire files via `--files` flag
  - Include file snippets directly in your prompts with tags:
    - `@file(path)` - Include the entire file
    - `@file(path:start-end)` - Include specific lines from a file

- **Editor Support**
  - Use your favorite editor for writing prompts with `--editor`
  - Supports standard `EDITOR` environment variable

## Examples

Message with a file snippet:
```
gptx msg "What does @file(main.go:10-30) do?"
```

Attach multiple files:
```
gptx --files="*.go" msg "Explain this codebase"
```

Save your configuration to a file:
```
gptx --model="gpt-4o" --files="main.go,helper.go" config > .gptx
```

Share configurations between projects:
```
# Create a project-specific config
gptx --files="project/*.go" config > project/.gptx

# View current effective configuration
gptx config
```

## References

- https://go.dev/doc
- https://github.com/openai/openai-go
- https://cli.urfave.org/v3

## TODO

- Support custom model tools
  - Web search
  - Shell commands
- Implement session history
