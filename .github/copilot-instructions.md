# GPTx CLI

This is a command line interface (CLI) for OpenAI models built in Go. The project aims to **minimize complexity** while maximizing functionality. Follow the following in the design of the CLI interface:
[Command Line Interface Guidelines](https://clig.dev/#conversation-as-the-norm)

The project uses OpenAI's new [resposnes API](https://platform.openai.com/docs/api-reference/responses) at:

- https://github.com/openai/openai-go/blob/main/responses/aliases.go
- https://github.com/openai/openai-go/blob/main/responses/inputitem.go
- https://github.com/openai/openai-go/blob/main/responses/inputitem_test.go
- https://github.com/openai/openai-go/blob/main/responses/response.go
- https://github.com/openai/openai-go/blob/main/responses/response_test.go

## Structure

The project has the following structure:

- `cmd/gptx` - The command line interface. This is where the CLI is defined and the commands are implemented.
- `pkg/gptx` - The core package that implements the model. This is where the model's logic is implemented, along with any buisness logic.
- `pkg/openai` - An openai api abstraction layer. This should be a thin layer that provides a simple interface to the OpenAI API. It should not contain any business logic. The goal is to make it easy to switch to a different API in the future, if needed.

The `cmd` should handle anything user facing, like a view in MVC. The `pkg` should handle all logic and features, like a controller in MVC. The `pkg/openai` should handle all API calls, like a service.

## Features

### Configuration

The CLI provides a command for sending a preconfigured model a message. The model's configuration behaves similarly to git, where it is loaded from `.gptx` files at the current directory or one of its parents. The configuration provides details such as the model, the API key, the model's tools, instructions, attachments, etc. All configuration data can also be provided as options to the command, to apply them to the current run only.

### Extensibility

The model should follow the Open-Closed Principle, meaning that it should be easy to add to the model's functionality without modifying the existing code. This can be achieved using events.

The model can emit events at various points in its lifecycle. The events can be used to trigger actions, such as logging. Through adding more events and listeners, the model can be extended to perform more complex tasks. The events are emitted in a way that allows for easy addition of new events and listeners without modifying the existing code.

### Tools

The model can be extended with tools. Tools are external programs that can be called from the model. Tools are configured in the model's configuration. There are built-in tools, such as web search or executing shell commands. The model can also be extended with custom tools by defining their inputs/outputs and the command to invoke them (along with descriptions).

Tools usage are part of the model's lifecycle. The model can emit events when a tool is used, and the tool can emit events when it is finished. The model runs in a loop until **the model** decides to stop. The model can be stopped by sending a signal, such as Ctrl-C.

### Attachments

Files can be provided in the model's configuration to attach to the model's input. The model for now takes in text and image files.

I was hoping for a way to easily reference line ranges in files and insert it to the message. Something like a tagging system.

Something that would be nice is to allow the user to provide a path and have its tree included as context, providing the model with a tool to search the path's files, like codebases.
