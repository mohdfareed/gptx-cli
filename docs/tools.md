# Tools System Design

## Overview

The tool system in GPTx CLI provides an extensible way to add capabilities to the models. Tools can be:

1. Built-in (shell execution, web search)
2. Custom (defined by users)

## Tool Registration

Tools are centrally managed in the `tools.Registry` component. This provides a unified system for:
- Registering tools
- Executing tool calls
- Looking up tool definitions

## Tool Definition

Each tool has the following components:
- Name: A unique identifier (e.g., "shell", "web-search")
- Description: Human-readable description for the model
- Parameters: JSON Schema defining the expected inputs
- Handler: Optional implementation function for local execution

```go
type ToolDef struct {
    Name    string             // Tool identifier
    Desc    string             // Description
    Params  map[string]any     // Parameter schema
    Handler Tool               // Implementation (optional)
}
```

## Built-in Tools

The system includes two built-in tools:

### Shell Tool
- Executes shell commands on the local system
- Parameters: Command to execute
- Implementation: Uses Go's os/exec package

### Web Search Tool
- Performs web searches using OpenAI's Responses API
- Parameters: Search query
- Implementation: Handled directly by the OpenAI client

## Tool Execution Flow

1. User prompt is sent to the model
2. Model decides to use a tool and specifies parameters
3. Tool handler is called with the parameters
4. Results are returned to the model
5. Model incorporates the results in its response

## Adding Custom Tools

New tools can be added by:

1. Creating a `ToolDef` structure
2. Implementing the handler function
3. Registering the tool with the registry

## Future Extensibility

The tool system is designed to be extended in the future:
- Custom tool registration without modifying existing code
- Support for more complex tools with state
- Integration with external APIs and services
