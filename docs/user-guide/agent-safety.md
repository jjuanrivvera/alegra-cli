---
title: Agent safety
description: Stop an AI agent from running destructive accounting operations, with hooks and per-host config for Claude Code, Codex, and OpenCode.
---

# Agent safety: gating destructive operations

Letting an AI agent touch your accounting is useful but risky: it can create
invoices, register payments, delete records, and emit electronic invoices to
DIAN/SAT. This page shows how to make sure an agent **cannot** run the
operations you don't want it to — whether it drives `alegra` through a shell or
through the [`alegra mcp`](mcp.md) server.

Be clear about what each layer actually does:

| Layer | Mechanism | What it does | Enforces? |
| --- | --- | --- | --- |
| 1. Tool annotations | `alegra mcp` marks each tool read-only or destructive | Lets an annotation-aware host gate writes on its own | **Advisory** — depends on the host |
| 2. Host config | deny / ask rules **and a PreToolUse hook** | The hook inspects the real command and can't be dodged by shell tricks — the definitive block | **Yes** — this is the real gate |
| 3. CLI built-ins | `--dry-run`, `delete` confirmation | Preview and a manual confirm in a terminal | Shell only (not the MCP surface) |

The honest summary: **layer 1 makes good hosts behave well by default, but it is
not a security boundary. Layer 2 is what actually blocks.** Use both.

!!! tip "The hook is the layer that can't be bypassed"
    Permission `deny`/`ask` globs are easy to read, but a crafted shell command
    can dodge them (quoting, subshells, `&&`). A **PreToolUse hook** runs your
    code against the *actual* command or tool before it executes, so it is the
    layer that **definitively** blocks — on both the shell and MCP surfaces.
    `alegra agent guard --host claude-code` generates it for you. (Hooks are a
    Claude Code feature; Codex blocks with a read-only sandbox, OpenCode with
    `deny` rules.)

## The quick way: `alegra agent guard`

`alegra` generates the layer-2 config for you, with the destructive operations
derived from the live command tree (so the list is always complete):

```bash
alegra agent guard --host claude-code   # prints settings.json + the PreToolUse hook
alegra agent guard --host codex         # prints config.toml
alegra agent guard --host opencode      # prints opencode.json
```

By default it **hard-blocks the irreversible actions** (`delete`, `void`,
`emit`, `stamp`, `close`, and the `*-delete` actions) and makes ordinary writes
(`create`, `update`, `import`) **require approval**; reads stay allowed. Flags:

- `--all-writes` — block every write, not just the irreversible ones.
- `--write` — install the files instead of printing them (it never overwrites an
  existing config; it prints that one for you to merge).

Review the output, drop it into your host, and you're done. The rest of this
page explains what it generates and how to tune it by hand.

## What counts as destructive

`alegra` read operations are `list`, `get`, and `export` (plus `catalog`,
`reports`, `doctor`, `version`). Everything else writes: `create`, `update`,
`delete`, `import`, and the resource actions `void`, `emit`, `stamp`, `email`,
`open`, `close`, `transfer`, …

## Layer 1 — built-in MCP annotations

`alegra mcp` advertises each tool's nature through the standard MCP tool
annotations: read tools carry `readOnlyHint: true`, and writes/actions carry
`destructiveHint: true`. You get this for free — no configuration.

What that buys you depends entirely on the host:

- A host that honors annotations (e.g. **Codex**) will require approval for the
  destructive tools and let the read-only ones run freely, automatically.
- A host that gates by tool name (e.g. **Claude Code**) ignores the annotation
  for permission decisions — you configure rules instead (layer 2).

Annotations are advisory per the MCP spec: a host that ignores them runs
everything. They are never a substitute for layer 2.

!!! note "Tool names"
    The MCP tools are named `alegra_<resource>_<subcommand>` — e.g.
    `alegra_invoices_void`, `alegra_contacts_delete`, `alegra_invoices_list`.
    Hosts namespace them further; in Claude Code a tool is
    `mcp__<server>__alegra_<resource>_<subcommand>`.

## Layer 2 — per-host enforcement

### Claude Code

Claude Code gates by **tool/command name**, and **`deny` always wins over
`allow`**. Put this in your project's `.claude/settings.json` (team-shared) or
`~/.claude/settings.json` (all projects).

**Shell surface** — block destructive `alegra` commands:

```json
{
  "permissions": {
    "deny": [
      "Bash(alegra * delete:*)",
      "Bash(alegra * create:*)",
      "Bash(alegra * update:*)",
      "Bash(alegra invoices emit:*)",
      "Bash(alegra invoices void:*)",
      "Bash(alegra invoices stamp:*)",
      "Bash(alegra invoices email:*)"
    ]
  }
}
```

Use `"ask"` instead of `"deny"` to require approval rather than hard-block. A
denied command never runs; in a headless/CI session `ask` also fails closed
(no human to approve → blocked).

**The definitive block is a PreToolUse hook.** The `deny` globs above help, but
a crafted shell command can dodge them (quoting, subshells, `&&`); a hook runs
your code against the actual command before it executes, so it cannot be
bypassed — and the same hook also covers the MCP surface. This is the most
important piece on Claude Code. Create `.claude/hooks/block-alegra-writes.sh`:

```bash
#!/usr/bin/env bash
# Reads the hook payload from stdin and denies destructive alegra commands.
input=$(cat)
cmd=$(printf '%s' "$input" | jq -r '.tool_input.command // empty')

if printf '%s' "$cmd" | grep -qE 'alegra\b.*\b(delete|create|update|import|emit|void|stamp|email|open|close|transfer)\b'; then
  jq -n '{
    hookSpecificOutput: {
      hookEventName: "PreToolUse",
      permissionDecision: "deny",
      permissionDecisionReason: "Destructive alegra operation blocked by policy."
    }
  }'
fi
exit 0
```

Register it in `.claude/settings.json`:

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          { "type": "command", "command": "${CLAUDE_PROJECT_DIR}/.claude/hooks/block-alegra-writes.sh" }
        ]
      }
    ]
  }
}
```

A PreToolUse hook that prints `permissionDecision: "deny"` (with exit code 0)
blocks the call; exit code 2 also blocks and shows stderr to the model. This is
the strongest option — it runs your logic, not a glob.

**MCP surface** — when alegra runs as a local MCP server
(`claude mcp add alegra -- alegra mcp start`), deny the write tools by name.
Because deny wins over allow, deny the write verbs rather than trying to
allow-list reads:

```json
{
  "permissions": {
    "deny": [
      "mcp__alegra__alegra_.*_(create|update|delete|import|void|emit|stamp|email|open|close|transfer)"
    ]
  }
}
```

MCP permission patterns are matched as regular expressions, so the `.*` and the
`(...)` group work as written; read tools (`..._list`, `..._get`,
`..._export`) are untouched. A PreToolUse hook can also gate MCP tools — set
`"matcher": "mcp__alegra__.*"` and inspect `.tool_name` in the script for
allow-list semantics (deny anything that isn't a read verb).

### Codex

Codex gates with two settings in `config.toml` (`~/.codex/config.toml`):
`sandbox_mode` and `approval_policy`.

**Shell surface** — the strongest, simplest gate is to take away write access:

```toml
sandbox_mode   = "read-only"   # Codex may read; edits/commands/network need approval
approval_policy = "untrusted"  # only known-safe reads auto-run; everything else asks
```

With `sandbox_mode = "read-only"`, Codex cannot run a state-changing `alegra`
command without an explicit approval — there is no per-command allow-list to
maintain. If you want writes to be possible but always reviewed, keep
`approval_policy` on a prompting mode so each write pauses for a human.

**MCP surface** — this is where layer 1 pays off. Codex's docs state that
destructive-annotated MCP tool calls **always require approval**. Since
`alegra mcp` marks `void`/`emit`/`delete`/`update`/… as destructive, Codex
prompts before running them automatically, while read-only tools run without
friction. No per-tool config needed; keep `approval_policy` on a prompting mode
so the prompt is honored.

Honest note: on Codex the gate is **approval** (a human says yes), not a hard
block — fine for interactive use, but it means an unattended Codex run should
use `sandbox_mode = "read-only"` so there is nothing to approve.

### OpenCode

OpenCode gates with a `permission` block in `opencode.json`. Each rule is
`"allow"`, `"ask"`, or `"deny"`, and **the last matching pattern wins**, so put
the catch-all first and the specific denials after.

**Shell surface** — `bash` takes per-pattern globs:

```json
{
  "$schema": "https://opencode.ai/config.json",
  "permission": {
    "bash": {
      "*": "allow",
      "alegra * delete*": "deny",
      "alegra * create*": "ask",
      "alegra * update*": "ask",
      "alegra invoices emit*": "deny",
      "alegra invoices void*": "deny",
      "alegra invoices stamp*": "deny"
    }
  }
}
```

`deny` blocks the command; `ask` prompts for approval before it runs.

**MCP surface** — the same permission keys match MCP tool names as wildcard
patterns, so you can deny a server's write tools:

```json
{
  "permission": {
    "alegra_*_delete": "deny",
    "alegra_*_create": "ask",
    "alegra_*_update": "ask",
    "alegra_invoices_void": "deny",
    "alegra_invoices_emit": "deny",
    "alegra_invoices_stamp": "deny"
  }
}
```

Confirm the exact tool prefix your OpenCode setup uses for the server, then
match it (the tool name is `alegra_<resource>_<subcommand>`).

## Layer 3 — CLI built-ins

When the agent uses a **shell**, two built-ins help even without host config:

- `--dry-run` on any command prints the exact request and sends nothing.
- `delete` asks for confirmation unless `-y` is passed.

These do not cover the **MCP surface**: an MCP tool call has no terminal, so the
interactive `delete` confirmation does not apply. On the MCP surface, rely on
layers 1 and 2.

## Recommended setup

Defense in depth:

1. Let layer 1 do its job — it ships on by default.
2. Add layer 2 for your host: a `deny` (or `ask`) rule **and** a hook on Claude
   Code; `sandbox_mode = "read-only"` on Codex; `permission` denials on OpenCode.
3. Keep `--dry-run` in your own habits and scripts.

Test it: ask the agent to run `alegra invoices void 1` (or call the
`alegra_invoices_void` tool) and confirm it is blocked or paused for approval
before anything reaches the API.
