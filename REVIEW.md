# GPTx CLI Code Review

## Overview

This document tracks our ongoing code review and simplification efforts for the GPTx CLI project. We'll use this to document our findings, track concerns, and guide refactoring decisions.

## Current Architecture

- **CLI layer** (`cmd/gptx/`): User-facing commands and flags
- **Core logic** (`pkg/gptx/`): Business logic and configuration
- **OpenAI API** (`pkg/openai/`): Thin abstraction over the OpenAI API

## Initial Observations

- Configuration system uses multiple layers (env vars, config files, CLI flags)
- Event system is planned but not yet implemented (empty events.go file)
- Tools framework has minimal implementation so far
- File attachment handling is implemented but may have complexity concerns
- CLI structure follows urfave/cli patterns with commands for messaging, config display, and demo UI
- README.md is minimal and includes TODOs for refactoring config and supporting custom model tools

## Areas to Review

- [x] Configuration complexity (multiple sources: env vars, config files in multiple locations, CLI flags)
- [ ] Events system necessity and implementation
- [ ] Tools framework approach (current implementation is minimal)
- [x] File attachment handling with tagging system
- [ ] Error handling patterns
- [ ] CLI command structure
- [ ] Overall code organization

## Simplification Opportunities

1. **Configuration System**:
   - Consider if all configuration methods are necessary
   - Potential duplication between cmd/gptx/config.go and pkg/gptx/config.go
   - Configuration loading happens in multiple places

2. **Events System**:
   - Evaluate if a simpler approach would suffice
   - Empty events.go file indicates planned but not implemented feature

3. **Tools Framework**:
   - Consider a simpler interface design
   - Current implementation is minimal but could grow in complexity

4. **File Attachments**:
   - Current implementation in pkg/openai/msg.go handles files well
   - Tagging system provides a clean approach to inline file references

5. **CLI Structure**:
   - The msg command is very basic currently
   - Consider if all the planned features are necessary

## In-depth Analysis

### Configuration System
Configuration is loaded from multiple sources with potential duplication between cmd/gptx/config.go and pkg/gptx/config.go. The system loads configuration from:
- Environment variables
- .gptx files in the current directory and all parent directories
- XDG_CONFIG_HOME or APPDATA config file
- CLI flags

**Improvement implemented**: Consolidated configuration loading into pkg/gptx/config.go and improved the clarity of the code with better documentation and organization.

**Decision**: Keep the hierarchical config search (parent directories) as it provides significant value for minimal code complexity. This follows Git-like behavior which is intuitive for users.

### File Context & Tagging System
The tagging system (`@file(path:start-end)`) is an elegant approach for:

1. **Inline file references**: Users can embed file snippets directly in their prompts
2. **Targeted excerpts**: Supports line range specification for focused context
3. **Future extensibility**: Architecture allows for new tag types beyond files

**Implementation plan**:
- Integrate the existing tagging processor into the message handling flow
- Document the feature in user-facing materials (help text, README)
- Consider adding examples in a "getting started" guide

**User discovery considerations**:
- Add documentation in the CLI help text
- Include examples in the README
- Consider adding a hint in response to certain queries (e.g., "How do I include code?")

### Chat History
For maintaining chat context across sessions, we discussed:
1. Storing conversations in a simple format (JSON)
2. Adding a `--continue` or `--session` flag
3. Creating a format for multi-message sequences

**Future work**: Implement a simple chat history mechanism that allows referencing previous conversations.

### OpenAI Integration
The OpenAI integration is well-structured with:
- Clean handling of file attachments
- Stream parsing for real-time responses
- Good abstraction over the Responses API

**Recommendation**: Maintain the current approach as it's well-designed.

### CLI Implementation
Uses urfave/cli with a standard structure. The main commands are:
- msg: Send a message to the model (currently just prints config)
- cfg: Display the current configuration
- demo: Demonstrate the UI elements

**Recommendation**: Focus on completing the msg command with actual model interaction before adding more features.

## Review Progress

| Component          | Reviewed | Findings                                              | Recommendations                                         |
| ------------------ | -------- | ----------------------------------------------------- | ------------------------------------------------------- |
| Configuration      | ✅        | Multiple layers with parent directory search          | Improved organization while retaining functionality     |
| OpenAI Integration | ✅        | Good abstraction over the Responses API               | Maintain current approach                               |
| CLI Implementation | ✅        | Standard urfave/cli usage with incomplete msg command | Complete core functionality first                       |
| Events System      | ✅        | Not yet implemented (empty file)                      | Consider if necessary                                   |
| Tools Framework    | ✅        | Minimal implementation so far                         | Simplify design before expanding                        |
| File Handling      | ✅        | Tag system provides elegant inline references         | Integrate existing tag system with proper documentation |
| README             | ✅        | Minimal with TODOs                                    | Update with clearer project goals                       |

## Next Steps

1. Integrate the tagging system for processing prompt text
   - Ensure proper handling in the message flow
   - Add user-facing documentation and examples

2. Complete core messaging functionality in the msg command
   - Connect to the OpenAI service
   - Implement proper error handling and response formatting

3. Decide on necessity of events system
   - Evaluate if a simpler callback approach would suffice
   - If needed, implement minimal version first

4. Simplify tools framework design
   - Start with basic tools implementation (web search, shell)
   - Define clean interface before adding extensibility

5. Implement chat history mechanism
   - Design simple format for storing conversations
   - Add command flags for continuing previous sessions
