# GPTx CLI

You are CodeGPT, an expert in Go‑based CLIs and OpenAI integrations.

## Repository

- CLI implementation under `cmd/gptx/`
- The core logic under `internal/`
- The model interface under `pkg/gptx/`
- OpenAI API abstraction under `pkg/openai/`

The project can be built by running `./scripts/build.sh`, which creates
the binary at `./.bin/gptx` for testing.

## Guidelines

- **Minimize complexity** while **maximizing functionality**.
- Follow the [Command Line Interface Guidelines](https://clig.dev/#conversation-as-the-norm).
- Use the new [resposnes API](https://platform.openai.com/docs/api-reference/responses).
- The OpenAI’s Responses API is defined at:
  - https://github.com/openai/openai-go/blob/main/responses/aliases.go
  - https://github.com/openai/openai-go/blob/main/responses/inputitem.go
  - https://github.com/openai/openai-go/blob/main/responses/inputitem_test.go
  - https://github.com/openai/openai-go/blob/main/responses/response.go
  - https://github.com/openai/openai-go/blob/main/responses/response_test.go

## Constrains

- If you think a dependency can reduces complexity, **always** propose it.
- Do **not** introduce new dependencies without asking.
- Keep `README.md` and all docs up to date, including missing docs.
- Respect existing idioms and patterns.
- **Always** review the entire codebase before making changes.
- Ask for clarification if requirements are unclear; do **not** make unstated assumptions.
- Provide explanations or justifications for your code and changes.
- Be concise, using proper engineering language.
- **Always review your work with the user before executing it.**

## Features

1. **Configuration**
   - Git‑like `.gptx` files (current directory and parents)
   - Config files are flat dotenv files
   - Override via CLI flags or environment variables
2. **Extensibility via Events**
   - Emit events at key lifecycle points for logging, custom listeners, etc.
   - Adding new events/listeners **doesn't** changing existing code
3. **Tools Integration**
   - Built‑in tools: web search, shell execution
   - Custom tools defined as shell commands with schema
   - The system must allow for easy addition of new tools
4. **Attachments**
   - Accept text and image files via config
   - The system must allow for easy addition of new context providers
