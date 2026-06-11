---
title: MCP Server
---

# MCP Server

alegra-cli can expose its entire command tree as a
[Model Context Protocol](https://modelcontextprotocol.io/) server, so AI agents
(Claude, etc.) can drive Alegra through well-described tools — one per command.

```bash
alegra mcp --help
```

Each CLI command becomes an MCP tool named `alegra_<resource>_<action>`
(e.g. `alegra_invoices_list`, `alegra_contacts_create`) whose input schema is
derived from that command's flags and arguments.

Credentials are resolved exactly as for the CLI (profile → keyring/env), and the
sensitive `--show-token` and `--profile` flags are excluded from the exposed
tool schemas.

!!! warning
    The MCP server can create, modify, and delete real accounting records. Run it
    against a sandbox profile while developing, and review the actions your agent
    is allowed to take.

## Wire it into your editor automatically

Rather than editing config files by hand, let the CLI write them for you. There
are subcommands for the common MCP hosts:

```bash
alegra mcp claude enable     # add the server to Claude Desktop's config
alegra mcp cursor enable     # Cursor
alegra mcp vscode enable     # VS Code

alegra mcp claude list       # show configured servers
alegra mcp claude disable    # remove it again
```

Each `enable` writes the host's MCP config to launch `alegra mcp start` (the
stdio server) with your current setup.

## Transports

| Command | Transport | Use |
| --- | --- | --- |
| `alegra mcp start` | stdio | The default; what the editor wiring launches. Also: `claude mcp add alegra -- alegra mcp start`. |
| `alegra mcp stream` | HTTP | Long-running server for remote/shared agents. Flags: `--host`, `--port` (default 8080), `--log-level`. |
| `alegra mcp tools` | — | Print the full tool schema as JSON (one tool per command) — handy for inspection or feeding another system. |

```bash
# Run the HTTP server on a custom port
alegra mcp stream --port 9000 --log-level debug

# Inspect the exposed tools
alegra mcp tools | jq '.[].name'
```

Then ask the agent to, for example, "list this month's open invoices in Alegra".

## Keep it safe

Every exposed tool is annotated read-only or destructive, so a host that honors
MCP annotations gates writes on its own. To actually block destructive tools
(`void`, `emit`, `delete`, …) on the agent side, see
[Agent Safety](agent-safety.md).
