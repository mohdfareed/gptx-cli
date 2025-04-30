# ChatGPT CLI

ChatGPT CLI is a command-line interface for OpenAI's ChatGPT built in Go.
It allows you to interact with the ChatGPT model directly from the terminal.

## References

- https://go.dev/doc
- https://cli.urfave.org/v3

## TODO

- Add scoped data management
    - Per invocation
    - Per session
    - Per user
- Add config command
    - Show config path
    - Show current config
    - Overwrite a config item
- Add chat history
    - History is per session
    - Start with an empty history each session
    - Store history in user data
    - Add command to manage history
        - Show history path
        - List stored sessions
        - Create a new session and load it (lazy file creation)
        - Load a previous session history into the current session
        - Delete a session or all sessions
