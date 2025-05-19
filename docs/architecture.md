# GPTx CLI Architecture

This document outlines the architecture of the GPTx CLI using Mermaid diagrams.

## Component Architecture

```mermaid
graph TD
    User[User] -->|interacts with| CLI["CLI (cmd/gptx)"]
    CLI -->|uses| Core["Core Logic (pkg/gptx)"]
    Core -->|calls| API["OpenAI API Layer (pkg/openai)"]
    API -->|communicates with| OpenAI[OpenAI Services]

    %% Configuration flow
    Config[Configuration] -->|loaded by| Core
    Config -->|from| ConfigFiles[".gptx files"]
    Config -->|from| EnvVars["Environment Variables"]
    Config -->|from| CLIFlags["CLI Flags"]

    %% Callbacks and Tools
    Core -->|triggers| Callbacks[Callbacks]
    Callbacks -->|handled by| Handlers[Event Handlers]
    Core -->|manages| Tools[Tools Registry]
    Tools -->|includes| BuiltInTools["Built-in Tools\n(web search, shell execution)"]
    Tools -->|includes| CustomTools["Custom Tools\n(defined as shell commands)"]

    %% Attachments
    Core -->|manages| Attachments[Attachments]
    Attachments -->|supports| Files["Text Files"]

    %% Build and Deploy
    Scripts["Build Scripts"] -.->|builds| CLI

    classDef primary fill:#d0e0ff,stroke:#3080ff,stroke-width:2px;
    classDef secondary fill:#e0f0e0,stroke:#30a030,stroke-width:2px;
    classDef tertiary fill:#f0e0e0,stroke:#a03030,stroke-width:1px;
    classDef future fill:#f0f0f0,stroke:#808080,stroke-width:1px,stroke-dasharray: 5 5;

    class CLI,Core,API primary;
    class Config,Callbacks,Tools,Attachments secondary;
    class ConfigFiles,EnvVars,CLIFlags,Handlers,BuiltInTools,CustomTools,Files,Scripts tertiary;
```

## Package Structure

GPTx CLI follows a clean architecture with separation of concerns between the CLI interface, core business logic, and API integrations. The codebase is organized into three main layers:

1. **CLI Layer** (`cmd/gptx`): User-facing commands and terminal interactions
2. **Core Logic** (`internal/`): Configuration, events, and tools management
3. **API Layer** (`pkg/`): Model interfaces and OpenAI API integration

```mermaid
graph LR
    Main["main (cmd/gptx/main.go)"] -->|entry point| CLI["CLI (cmd/gptx/*.go)"]
    CLI -->|uses| Core["Core Logic (pkg/gptx/*.go)"]
    Core -->|uses| Internal["Internal Components (internal/*)"]
    Core -->|depends on| OpenAI["OpenAI API (pkg/openai/*.go)"]

    subgraph "cmd/"
        subgraph "cmd/gptx"
            CLI_cmds["cmds.go\n(CLI commands)"]
            CLI_cli["cli.go\n(CLI setup)"]
            CLI_editor["editor.go\n(Editor integration)"]
            CLI_help["help.go\n(Help text)"]
            CLI_logging["logging.go\n(Logging)"]
            CLI_main["main.go\n(Entry point)"]
        end
    end

    subgraph "pkg/gptx"
        Core_client["client.go\n(Client interface)"]
        Core_model["model.go\n(Model controller)"]
    end

    subgraph "internal"
        Internal_cfg["cfg/\n(Configuration)"]
        Internal_callbacks["events/\n(Callbacks)"]
        Internal_tools["tools/\n(Tool registry)"]
    end

    subgraph "pkg/openai"
        API_client["client.go\n(OpenAI client)"]
        API_handlers["handlers.go\n(Event handlers)"]
        API_request["request.go\n(Request builder)"]
        API_types["types.go\n(Type definitions)"]
    end

    subgraph "scripts"
        Build_sh["build.sh\n(Unix build)"]
        Build_ps1["Build.ps1\n(Windows build)"]
        Release_sh["release.sh\n(Unix release)"]
        Release_ps1["Release.ps1\n(Windows release)"]
        Run_sh["run.sh\n(Unix run)"]
        Run_ps1["Run.ps1\n(Windows run)"]
    end

    CLI --- CLI_cmds & CLI_cli & CLI_editor & CLI_help & CLI_logging & CLI_main
    Core --- Core_model & Core_client
    Internal --- Internal_cfg & Internal_callbacks & Internal_tools
    OpenAI --- API_client & API_handlers & API_request & API_types

    %% Script connections
    Build_sh & Build_ps1 -.->|builds| CLI
    Run_sh & Run_ps1 -.->|runs| CLI
    Release_sh & Release_ps1 -.->|packages| CLI

    classDef primary fill:#d0e0ff,stroke:#3080ff,stroke-width:2px;
    classDef cmd fill:#f0f0d0,stroke:#a0a030,stroke-width:1px;
    classDef core fill:#e0f0e0,stroke:#30a030,stroke-width:1px;
    classDef api fill:#f0e0e0,stroke:#a03030,stroke-width:1px;
    classDef scripts fill:#e0e0f0,stroke:#3030a0,stroke-width:1px;

    class Main,CLI,Core,OpenAI primary;
    class CLI_cmds,CLI_cli,CLI_editor,CLI_help,CLI_logging,CLI_main cmd;
    class Core_config,Core_env,Core_events,Core_model,Core_tools core;
    class API_client api;
    class Build_sh,Build_ps1,Release_sh,Release_ps1,Run_sh,Run_ps1 scripts;
```

## Callback System

```mermaid
sequenceDiagram
    participant User
    participant CLI
    participant Model
    participant Tools
    participant Client
    participant OpenAI

    User->>CLI: Initiates command
    CLI->>Model: Message(prompt)

    Model->>Client: Generate(request)
    Note right of Client: Request includes callbacks

    Client->>Model: OnStart(config)
    Client->>OpenAI: Send streaming request

    loop For each stream event
        OpenAI->>Client: Stream event
        alt Text response
            Client->>Model: OnReply(text)
            Model->>CLI: Display text to user
        else Tool call
            Client->>Model: OnToolCall(toolCall)
            Model->>Tools: Execute tool
            Tools->>Model: Return result
            Model->>Client: Continue with tool result
        end
    end

    OpenAI->>Client: Complete response
    Client->>Model: OnDone(usage)
    Model->>CLI: Complete
    CLI->>User: Display final results
```

## Configuration Flow

```mermaid
flowchart TD
    start[Start] --> findConfig["Find .gptx config files\n(current dir and parents)"]
    findConfig --> loadEnv["Load environment variables"]
    loadEnv --> parseCLI["Parse CLI flags"]

    parseCLI --> mergeConfig["Merge configurations\n(CLI flags override env vars\nenv vars override config files)"]

    mergeConfig --> useConfig[Use configuration for OpenAI API]

    classDef process fill:#d0e0ff,stroke:#3080ff,stroke-width:2px;
    classDef terminal fill:#e0f0e0,stroke:#30a030,stroke-width:2px;

    class start terminal;
    class findConfig,loadEnv,parseCLI,mergeConfig,useConfig process;
```

## Tool Execution Flow

```mermaid
flowchart TD
    start[Start] --> receiveToolCall["Receive tool call from OpenAI"]

    receiveToolCall --> callHandlers["Call OnToolCall handler"]

    callHandlers --> lookupTool["Look up tool in registry"]

    lookupTool --> executeTool["Execute tool with parameters"]

    executeTool --> collectResults["Collect tool results"]

    collectResults --> callResultHandler["Call OnToolResult handler"]

    callResultHandler --> returnResults["Return results to OpenAI"]

    returnResults --> end[End]

    classDef process fill:#d0e0ff,stroke:#3080ff,stroke-width:2px;
    classDef callback fill:#e0f0e0,stroke:#30a030,stroke-width:2px;
    classDef terminal fill:#f0e0e0,stroke:#a03030,stroke-width:2px;

    class start,end terminal;
    class receiveToolCall,lookupTool,executeTool,collectResults,returnResults process;
    class callHandlers,callResultHandler callback;
```

## Responses API Integration

```mermaid
sequenceDiagram
    participant CLI as CLI
    participant Model as Model
    participant Client as OpenAI Client
    participant Responses as OpenAI Responses API
    participant Tools as Tool Registry

    CLI->>Model: Message(prompt)
    Model->>Client: Generate(request)
    Note right of Client: Request includes callbacks

    Client->>Responses: Create streaming request
    Client->>Model: OnStart callback

    loop Stream Response
        Responses->>Client: Stream event

        alt Text Content
            Client->>Model: OnReply callback
            Model->>CLI: Display text to user
        else Tool Call Content
            Client->>Model: OnToolCall callback
            Model->>Tools: Look up & execute tool
            Tools->>Model: Return tool result
            Model->>Client: Continue with result
            Client->>Responses: Submit tool output
        end
    end

    Responses->>Client: Complete response
    Client->>Model: OnDone callback with usage
    Model->>CLI: Signal completion
    CLI->>User: Complete interaction

    Note over Client,Responses: Using direct callbacks instead of\nchannels for simpler, more efficient processing
```

## Tool Integration

The new tool system provides a unified registry that makes it easy to add custom tools:

```mermaid
flowchart TD
    start[Start] --> createRegistry["Create tool registry"]
    createRegistry --> registerBuiltIn["Register built-in tools\nbased on config"]
    registerBuiltIn --> registerCustom["Register custom tools\n(future extension)"]

    registerCustom --> modelSetup["Set up model with registry"]

    modelSetup --> receiveToolCall["Receive tool call during generation"]

    receiveToolCall --> lookupTool["Look up tool in registry"]

    lookupTool --> executeTool["Execute tool with parameters"]

    executeTool --> returnResult["Return result to client"]

    returnResult --> continueGeneration["Continue model generation"]

    continueGeneration --> end[End]

    classDef process fill:#d0e0ff,stroke:#3080ff,stroke-width:2px;
    classDef terminal fill:#e0f0e0,stroke:#30a030,stroke-width:2px;

    class start,end terminal;
    class createRegistry,registerBuiltIn,registerCustom,modelSetup,receiveToolCall,lookupTool,executeTool,returnResult,continueGeneration process;
```

## Future Extensibility

The new architecture is designed for easy extension in several areas:

```mermaid
flowchart LR
    Core["Core Model"] --> Client1["OpenAI Client"]
    Core --> Client2["Claude Client (future)"]
    Core --> Client3["Local Model Client (future)"]

    Tools["Tool Registry"] --> Tool1["Shell Tool"]
    Tools --> Tool2["Web Search Tool"]
    Tools --> Tool3["SQL Tool (future)"]
    Tools --> Tool4["Custom Tool (future)"]

    Context["Context Sources"] --> Context1["File Attachments"]
    Context --> Context2["Environment Info"]
    Context --> Context3["Chat History (future)"]
    Context --> Context4["Database (future)"]

    classDef core fill:#d0e0ff,stroke:#3080ff,stroke-width:2px;
    classDef current fill:#e0f0e0,stroke:#30a030,stroke-width:2px;
    classDef future fill:#f0f0f0,stroke:#808080,stroke-width:1px,stroke-dasharray: 5 5;

    class Core,Tools,Context core;
    class Client1,Tool1,Tool2,Context1,Context2 current;
    class Client2,Client3,Tool3,Tool4,Context3,Context4 future;
```

The refactored architecture enables:

1. **Custom Tools**: Register custom tools through the unified tool registry
2. **Chat History**: Can be added as a client-side feature without changing core interfaces
3. **Context Providers**: New sources of context (beyond files) can be added through the model configuration
4. **Client Implementations**: Alternative clients can implement the simple client interface
