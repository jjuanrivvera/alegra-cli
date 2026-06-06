---
title: Agent Skill
---

# Use it from your AI agent (skill)

alegra-cli ships an **agent skill** — a `SKILL.md` that teaches AI coding agents
(Claude Code, Cursor, Codex, Gemini CLI, Windsurf, GitHub Copilot, and more) how
and when to drive the `alegra` CLI. Once installed, you can just ask your agent
things like *"create a draft invoice for client 12 and email it"* and it knows
the commands, flags, and safety rules.

## Install the skill

The cross-agent way (recommended) — detects every agent you have and stays up to
date:

```bash
npx skills add jjuanrivvera/alegra-cli
```

Built into the CLI (no Node required) — writes the bundled skill directly:

```bash
alegra skills install                 # current project (./.claude/skills)
alegra skills install --global        # user-level (~/.claude/skills)
alegra skills install --agent cursor --global
alegra skills path                    # show where it would install
alegra skills print                   # print the SKILL.md
```

Native Claude Code plugin:

```text
/plugin marketplace add jjuanrivvera/alegra-cli
/plugin install alegra-cli@alegra
```

## Prerequisites

The skill **wraps** the `alegra` binary — it doesn't bundle it. So:

1. Install the CLI: `brew install jjuanrivvera/alegra-cli/alegra-cli`
   (or `go install github.com/jjuanrivvera/alegra-cli/cmd/alegra@latest`).
2. Authenticate: `alegra auth login` (or set `ALEGRA_EMAIL`/`ALEGRA_TOKEN`).
3. Verify: `alegra doctor`.

## Skill vs MCP

- **Skill** (this page) — markdown instructions for shell-capable agents. The
  agent runs `alegra …` commands directly. Lowest friction; best for CLI use.
- **[MCP server](mcp.md)** — `alegra mcp` exposes every command as a structured
  tool, for clients that prefer tool schemas (Claude Desktop, etc.).

Use whichever your agent supports; they can coexist.

## What the agent learns

The skill teaches the auth → discover → act → verify workflow, the golden rules
(preview writes with `--dry-run`, parse with `-o json`, invoices are
append-only, prefer `--count`), the resource/action map, how to write nested
bodies, electronic-invoicing emission, bulk CSV import/export, and error
handling. A condensed cheatsheet ships alongside it in
`references/alegra-commands.md`.
