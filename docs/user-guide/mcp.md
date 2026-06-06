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

## Example: Claude Code

Register the server (stdio transport):

```bash
claude mcp add alegra -- alegra mcp
```

Then ask the agent to, for example, "list this month's open invoices in Alegra".
