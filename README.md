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
  - Support for image attachments (jpg, jpeg, png, svg)

- **Callback System**
  - Simple callback-based event system for model interaction
  - Callbacks for model responses, tool usage, and completions
  - Real-time visibility into model actions with minimal overhead

- **Tools Support**
  - Unified tool registry system for easy extensibility
  - Built-in tools for web search and shell commands
  - Clean API for adding custom tools

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
gptx --shell=auto --web=true msg "Find all files in the current directory and summarize them"
```

View current configuration:
```
gptx cfg
```

## Architecture

The application follows a clean separation of concerns:

- **CLI layer** (`cmd/gptx/`): User-facing commands and flags
- **Core logic** (`internal/`): Business logic, configuration, events, and tools
- **API layer** (`pkg/`): Model interactions and OpenAI Responses API integration

For detailed architecture documentation, see [docs/architecture.md](docs/architecture.md).

## Installation

```
go install github.com/mohdfareed/gptx-cli/cmd/gptx@latest
```

Or clone and build from source:

```
git clone https://github.com/mohdfareed/gptx-cli.git
cd gptx-cli
go build -o gptx ./cmd/gptx
```

## Documentation

- [Architecture](docs/architecture.md): System design, components, and API reference

## References

- [Go Documentation](https://go.dev/doc)
- [OpenAI Responses API](https://platform.openai.com/docs/api-reference/responses)
- [OpenAI Go SDK](https://github.com/openai/openai-go)
- [CLI Framework](https://cli.urfave.org/v3)

## Roadmap

- [ ] Custom user tools
- [ ] Implement chat history
