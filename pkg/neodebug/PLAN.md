# Neodebug MCP

- Keep MCP in the neodebug codebase, launched as a standalone stdio mode.
- Extract a shared debugging engine for the TUI and MCP; remove Bubble Tea and global configuration from connection handling.
- Store named configurations in `neodebug.json` beside `go.mod`: schema generation commands and connection targets (transporter, URL, options).
- The agent supplies a folder and configuration name. Find the nearest parent `go.mod`, load that root's configuration, and run schema generation from that directory.
- Commands use executable/argument arrays without shell expansion. MCP tools cannot supply or override commands. Running repository code requires developer trust; command changes invalidate approval.
- Expose tools to load a project, inspect routes and presets, open/close sessions, send requests, and read buffered events with cursors and bounded waits.
- One MCP process supports multiple independent sessions, including multiple connections to the same target. Application servers run separately.
- Reuse Neoroute schemas and clients; share schema-aware request conversion between interfaces. Reserve stdout for MCP and send logs to stderr.
- Register the same launcher through each agent's project MCP configuration; there is no universal repository registration file.

## How to structure

- `cmd/neodebug` will be a new package depending on `pkg/neodebug` that handles the command itself
- `pkg/neodebug`: Will only be slightly refactored to get the main connection logic to not depend on bubble tea
- `cmd/neodebug`: Will contain the MCP server and more

## cmd/neodebug

- We use `charm/huh` for forms with custom styling and components we made (like choosing, etc.)
- Find a proper library that can do the sub-commands and arguments properly
- Make sure to handle schema versions properly (reject generating for old schemas / too new ones)
- `neodebug.json` stores all of the config for the project we're in
  - Support mono-repos by storing the projects we're looking at in there as well
  - Store a name for every transporter + its project
- MCP
  - Needs to support multiple sessions (across also multiple paths)
  - Schema should be generated new on every connection to make sure we always have the newest one, also specify in the tool description so the agent knows
  - Minimize the number of tools
- `neodebug` asks the user if they want to set up neodebug, create a new config (`neodebug.json`) or connect directly
- `neodebug [run] <name>` starts the debugger for a specific transporter specified in `neodebug.json`.
  - `run` is optional, if there is a sub-command that clashes, that one is executed
- `neodebug mcp` starts the MCP server.
- `neodebug http/ws <url>` connects to HTTP / WebSocket (or first opens form when not configured with args)
  - Either `--file <schema.json>` or `--project <../some_project>` is required to get the schema
- `neodebug tidy` finds the current project and adds new transporters that aren't initialized
  - In interactive mode: The user will also be asked if they want to delete "old" transporters (meaning ones that don't exist in the schema anymore)
  - In interactive mode: The user is asked to configure the individual transporters
  - This command can be executed again and if there is something that is not filled out (""), we ask the user to still configure
  - If there is no config, we exit with an error
- `neodebug init` initializes `neodebug.json` with a basic project path to the current project
- `neodebug add [path]` adds a new project to the `neodebug.json` config file
  - In interactive mode: Asks to configure just like `neodebug tidy`
  - If path is not set it is the current directory
- `--non-interactive` exists for all commands except connection and makes sure no forms open, instead of the form opening you just get an error (this is for AI specifically).
