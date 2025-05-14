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

    %% Events and Tools
    Core -->|emits| Events[Events]
    Events -->|listened by| Listeners[Event Listeners]
    Core -->|manages| Tools[Tools]
    Tools -->|includes| BuiltInTools["Built-in Tools\n(web search, shell execution)"]
    Tools -->|includes| CustomTools["Custom Tools\n(defined as shell commands)"]

    %% Attachments
    Core -->|manages| Attachments[Attachments]
    Attachments -->|supports| Files["Text Files"]
    Attachments -->|uses| Tags["Tagging System\n(for file line ranges)"]

    %% Build and Deploy
    Scripts["Build Scripts"] -.->|builds| CLI

    classDef primary fill:#d0e0ff,stroke:#3080ff,stroke-width:2px;
    classDef secondary fill:#e0f0e0,stroke:#30a030,stroke-width:2px;
    classDef tertiary fill:#f0e0e0,stroke:#a03030,stroke-width:1px;
    classDef future fill:#f0f0f0,stroke:#808080,stroke-width:1px,stroke-dasharray: 5 5;

    class CLI,Core,API primary;
    class Config,Events,Tools,Attachments secondary;
    class ConfigFiles,EnvVars,CLIFlags,Listeners,BuiltInTools,CustomTools,Files,Tags,Scripts tertiary;
```

## Package Structure

```mermaid
graph LR
    Main["main (cmd/gptx/main.go)"] -->|entry point| CLI["CLI (cmd/gptx/*.go)"]
    CLI -->|uses| Core["Core Logic (pkg/gptx/*.go)"]
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
        Core_config["config.go\n(Config loading)"]
        Core_env["env.go\n(Environment handling)"]
        Core_events["events.go\n(Event system)"]
        Core_model["model.go\n(Model interactions)"]
        Core_tags["tags.go\n(File tagging system)"]
        Core_tools["tools.go\n(Tool integration)"]
    end

    subgraph "pkg/openai"
        API_client["client.go\n(OpenAI API client)"]
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
    Core --- Core_config & Core_env & Core_events & Core_model & Core_tags & Core_tools
    OpenAI --- API_client

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
    class Core_config,Core_env,Core_events,Core_model,Core_tags,Core_tools core;
    class API_client api;
    class Build_sh,Build_ps1,Release_sh,Release_ps1,Run_sh,Run_ps1 scripts;
```

## Event System

```mermaid
sequenceDiagram
    participant User
    participant CLI
    participant Core
    participant Events
    participant Tools
    participant API
    participant OpenAI

    User->>CLI: Initiates command
    CLI->>Core: Processes command

    Core->>Events: EventStart
    Note over Events: Emit configuration

    Core->>API: Create OpenAI request
    API->>OpenAI: Send API request (Responses API)

    alt Tool is needed
        OpenAI->>API: Request tool execution
        API->>Core: Forward tool request
        Core->>Events: EventTool
        Note over Events: Emit tool details
        Core->>Tools: Execute tool
        Tools->>Core: Return tool result
        Core->>Events: EventTool
        Note over Events: Emit tool results
        Core->>API: Send tool results
        API->>OpenAI: Forward tool results
    end

    OpenAI->>API: Return response
    API->>Core: Forward response
    Core->>Events: EventReply
    Note over Events: Emit response
    Core->>CLI: Output response
    CLI->>User: Display response

    Core->>Events: EventComplete
    Note over Events: Emit completion details

    alt Error occurs
        Core->>Events: EventError
        Note over Events: Emit error details
    end
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

    receiveToolCall --> emitToolEvent["Emit EventTool event with params"]

    emitToolEvent --> executeTool["Execute registered tool function"]

    executeTool --> collectResults["Collect tool results"]

    collectResults --> emitResultEvent["Emit EventTool event with results"]

    emitResultEvent --> returnResults["Return results to OpenAI"]

    returnResults --> end[End]

    classDef process fill:#d0e0ff,stroke:#3080ff,stroke-width:2px;
    classDef event fill:#e0f0e0,stroke:#30a030,stroke-width:2px;
    classDef terminal fill:#f0e0e0,stroke:#a03030,stroke-width:2px;

    class start,end terminal;
    class receiveToolCall,executeTool,collectResults,returnResults process;
    class emitToolEvent,emitResultEvent event;
```

## Responses API Integration

```mermaid
sequenceDiagram
    participant CLI as CLI
    participant Core as Core Logic
    participant API as OpenAI API Layer
    participant Responses as OpenAI Responses API

    CLI->>Core: Send user input
    Core->>API: Process request

    API->>Responses: Create streaming request

    loop Stream Response
        Responses->>API: Stream text or tool calls
        API->>Core: Process stream data

        alt Text Content
            Core->>CLI: Display text
            CLI->>User: Show text to user
        else Tool Call Content
            Core->>Tools: Execute tool based on call
            Tools->>Core: Return tool result
            Core->>API: Send tool result
            API->>Responses: Submit tool output
        end
    end

    Responses->>API: Complete response
    API->>Core: Signal completion
    Core->>CLI: Finalize interaction

    Note over API,Responses: Using streaming pattern with channels for\nboth text and tool calls
```
