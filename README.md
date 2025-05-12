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

- **Event System**
  - Simple channel-based event system for model interaction
  - Events for model responses, tool requests, and completions
  - Easy integration with custom handlers and UI components

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

## Architecture

The application follows a clean separation of concerns:

- **CLI layer** (`cmd/gptx/`): User-facing commands and flags
- **Core logic** (`pkg/gptx/`): Business logic, configuration, and events
- **OpenAI API** (`pkg/openai/`): Thin abstraction over the OpenAI Responses API

## References

- https://go.dev/doc
- https://github.com/openai/openai-go
- https://cli.urfave.org/v3

## Roadmap

- [ ] Complete model interaction through the `msg` command
- [ ] Support model tools
  - [ ] Web search
  - [ ] Shell commands
  - [ ] Custom user tools
- [ ] Implement chat history
