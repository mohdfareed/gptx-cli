# gptx

NAME:
   gptx - Interact with an LLM models

USAGE:
   gptx [global options] [command [command options]]

DESCRIPTION:
   Interact with LLM models from your terminal.

   Features:
   - Send messages to LLM models
   - Include file contents using the @file tag system
   - Configure model parameters
   - Use multiple configuration methods
   - Extend with tools and plugins

   Learn more about a command:
       gptx help <command>

COMMANDS:
   msg      Send a message to a model
   cfg      Show current configuration
   demo     Show UI demonstration
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --editor string, -e string  Use specified text editor for input [$GPTX_EDITOR, $EDITOR]
   --help, -h                  show help
   --quiet, --silent, -q       Show only error messages (default: true) [$GPTX_QUIET, $GPTX_SILENT]
   --verbose, -v               Show debug messages (default: false) [$GPTX_VERBOSE, $GPTX_DEBUG]

   config

   --key string                Set Platform API key [$GPTX_API_KEY]
   --max int                   Limit response length [$GPTX_MAX_TOKENS]
   --model string              Select model to use (default: "o4-mini") [$GPTX_MODEL]
   --prompt string, -s string  Set system prompt [$GPTX_SYS_PROMPT]
   --temp float                Set response randomness (0-100) (default: 1) [$GPTX_TEMPERATURE]

   context

   --files string, -f string [ --files string, -f string ]  Attach files to the message [$GPTX_FILES]
   --shell string                                           Set the shell for the model to use [$GPTX_SHELL_TOOL]
   --web                                                    Enable web search (default: false) [$GPTX_WEB_SEARCH]

