# GPTx CLI

GPTx CLI is a command-line interface for OpenAI's GPT models built in Go.
It allows you to interact with any GPT model directly from the terminal.

## Features

- **Configuration**
  - Git-like `.gptx` files in current directory and parent directories
  - Global config in application directory
  - Configuration via environment variables and CLI flags
  - Easy to share configurations between projects

- **File Integration**
  - Attach entire files via `--files` flag with glob pattern support
  - Include file snippets directly in your prompts with tags:
    - `@file(path)` - Include the entire file
    - `@file(path:start-end)` - Include specific lines from a file

- **Event System**
  - Simple channel-based event system for model interaction
  - Events for model responses, tool usage, and completions
  - See what the model is doing in real-time

- **Tools Support**
  - Built-in tools for web search and shell commands
  - Custom tool support via shell scripts
  - Event emissions for tool usage and results

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

Use tools:
```
gptx --tools="shell,web_search" msg "Find all files in the current directory and summarize them"
```

View current configuration:
```
gptx cfg
```

## Architecture

The application follows a clean separation of concerns:

- **CLI layer** (`cmd/gptx/`): User-facing commands and flags
- **Core logic** (`pkg/gptx/`): Business logic, configuration, and events
- **OpenAI API** (`pkg/openai/`): Thin abstraction over the OpenAI Responses API

## Installation

```
go install github.com/mohdfareed/chatgpt-cli/cmd/gptx@latest
```

Or clone and build from source:

```
git clone https://github.com/mohdfareed/chatgpt-cli.git
cd chatgpt-cli
go build -o gptx ./cmd/gptx
```

## References

- [Go Documentation](https://go.dev/doc)
- [OpenAI Go SDK](https://github.com/openai/openai-go)
- [CLI Framework](https://cli.urfave.org/v3)

## Roadmap

- [x] Complete model interaction through the `msg` command
- [x] Support model tools
  - [x] Web search
  - [x] Shell commands
  - [x] Custom user tools
- [ ] Implement chat history
- [ ] Support for image attachments
- [ ] Conversation management with stateful sessions
