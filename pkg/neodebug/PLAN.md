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
