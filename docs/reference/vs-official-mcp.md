---
title: 'Alegra MCP: official server vs. alegra-cli'
description: Alegra MCP compared. The official hosted MCP server vs. alegra-cli's own MCP server and agent-first CLI. Read-only vs. write, OAuth vs. terminal.
---

# alegra-cli vs. the official Alegra MCP server

Alegra publishes a hosted **MCP server** (`mcp.alegra.com`) so AI agents can call the
accounting API over the Model Context Protocol. `alegra-cli` takes a different route: it
is built **agent-first** as a command-line tool that AI agents drive directly — and it
*also* ships its own MCP server (`alegra mcp`) and works for humans in a terminal. This
page compares the two so you can pick the right fit.

!!! note "Snapshot"
    Reflects the official MCP server's **live tool set (66 read-only tools)** and
    `alegra-cli` **v0.9.2**, verified **2026-06-11** by connecting the official MCP through
    claude.ai and inspecting every tool. Both wrap the same Alegra v1 API and both require
    network access to it. Alegra may change its MCP over time — treat this as a
    point-in-time snapshot.

## At a glance

| | Official Alegra MCP | alegra-cli |
| --- | --- | --- |
| Interface | Hosted MCP server (`mcp.alegra.com`, streamable HTTP) | Terminal + agent **skill** + `alegra mcp` server |
| Authentication | OAuth 2.0 (browser login, authorization-code + PKCE; scope `owner`) | Alegra API token (Basic), stored in the **OS keyring** |
| How an agent connects | An interactive MCP host whose callback Alegra allow-listed — **claude.ai** web or **ChatGPT** (it runs the OAuth flow) | One command in **Claude Code**, Cursor, or Codex via a [skill](../user-guide/agent-skill.md); **or** [`alegra mcp`](../user-guide/mcp.md) for MCP hosts |
| Headless / CI / cron | No — OAuth needs an interactive browser login (no `client_credentials` or device grant) | Yes — API token works unattended |
| Install | None — hosted by Alegra | Install a binary (Homebrew, Scoop, Docker, …) |
| Output | Structured JSON tool results | table / JSON / YAML / CSV — composable with `jq` and pipes (fewer tokens for an agent) |
| Agent safety controls | Host-dependent | `--dry-run` on any command, interactive confirm on `delete`, restrictable via shell hooks |
| Human (terminal) use | — | First-class |
| Resource coverage | ~25 areas (66 tools), all reads | 45+ resources |
| Writes (create / update / delete) | Not exposed — read-only | Yes |
| Invoice emission (stamp / void / email) | Not exposed | Yes |

## Built agent-first

`alegra-cli` was designed assuming the primary caller is an AI agent, not a person
typing commands. The rationale is laid out in the post
[CLIs over MCPs](https://jjuanrivvera.com/es/blog/clis-sobre-mcps/); in short:

- **It connects to shell agents in one step.** A coding agent (Claude Code, Cursor,
  Codex) installs the binary, sets an API token, and runs commands — no OAuth flow, no
  browser, no host configuration. The official MCP uses OAuth 2.0 with a **closed
  redirect-URI allow-list**: probing its registration endpoint, only Alegra's
  integration-partner callbacks (`claude.ai` and ChatGPT) are accepted — loopback
  (`localhost`/`127.0.0.1`), custom URI schemes (`vscode://`, `cursor://`), and arbitrary
  HTTPS callbacks are all rejected. So no terminal, CLI, or self-hosted MCP client (Claude
  Code, Cursor, OpenCode, OpenClaw, Hermes, …) can complete the flow. With `alegra-cli`,
  the same API token also works in CI, cron, and scripts, where an interactive OAuth login
  is not possible.
- **Composable output.** An agent can pick `--columns`, pipe to `jq`, or `grep` — sending
  far fewer tokens to the model than a raw JSON tool result.
- **Safety is enforceable.** `--dry-run` previews the exact request on any command,
  `delete` asks for confirmation, and a shell can block destructive commands with hooks —
  so an autonomous agent *cannot* delete even if it tries. See
  [Agent Safety](../user-guide/agent-safety.md) for per-host gating (Claude Code, Codex,
  OpenCode).
- **Two ways in for agents.** A coding agent with shell access uses the **skill** (it
  learns the golden rules and when to call each command); an MCP-protocol host uses
  **`alegra mcp`**, which exposes the same command tree as MCP tools.

!!! note "On credentials"
    Both tools handle credentials responsibly — they are not a point of difference.
    `alegra-cli` keeps a long-lived API token in the OS keyring; the official MCP uses
    short-lived OAuth tokens the model never sees. The real trade-off is reach: an API
    token runs unattended (CI, cron, shell agents), while OAuth requires a one-time
    interactive browser login.

## Resource coverage

Both wrap the same Alegra v1 API but expose different slices of it. The key difference:
the official MCP's 66 tools are **all read-only** (`get…`/`list…`), while `alegra-cli`
reads *and* writes.

### Read by both, written only by alegra-cli
The official MCP reads almost the whole API: contacts; items (with per-warehouse stock)
and the inventory family (item categories, variant attributes, warehouses, transfers,
price lists, adjustments and numerations, custom fields); bank accounts and
reconciliations; accounting (journals, cost centers, ledger categories); invoices;
incoming and outgoing payments; the expenses family (bills, supplier debit notes, purchase
orders); taxes; retentions; currencies; sellers; resolutions/numberings; sales reports;
and the units / SAT product-key reference catalogs. It exposes all of these **read-only**;
`alegra-cli` also creates, updates, and deletes them.

### Only in alegra-cli
Resources the official MCP does not expose at all: estimates, credit notes, customer
(income) debit notes, remissions, transportation receipts, global invoices (CFDI),
recurring invoices and payments, payment terms, additional charges, and webhook
subscriptions. Plus **every write operation** across all resources and electronic invoice
emission (see below).

**SAT product keys** (`claveProdServ`, ~52k Mexico-specific entries): the official MCP
exposes a read-only `config_getProductKeys` tool. `alegra-cli` syncs the SAT's published
catalog on demand (`alegra catalog sync-sat`, offered by `alegra init` on Mexican
accounts) and searches it offline with `alegra catalog product-keys <query>`, no
connection needed.

### Only in the official MCP
- **Support Center** help-desk helpers and **Task Manager** (tasks) — out of scope for an
  accounting CLI.

## Capability differences that matter

**Electronic invoicing.** The official MCP's invoice tools are read-only
(`getInvoices`, `getInvoiceById`/`ByNumber`, plus an exportable-document generator). It
cannot create, update, `stamp`, `void`, or `email` an invoice. Emitting an electronic
invoice to DIAN/SAT, voiding, or emailing it is done with `alegra-cli`
(`alegra invoices stamp|void|email …`) or the REST API directly. If electronic emission
is part of your workflow, the CLI covers the full cycle.

**Payments.** The official MCP reads both incoming (`incomePayments`) and outgoing
(`expenses` outgoing-payments) payments. `alegra-cli` unifies them under one `payments`
resource and adds write, `stamp`, `void`, and `open`.

**Reports.** The two report sets are complementary:

| Report | Official MCP | alegra-cli |
| --- | --- | --- |
| Sales by client / totals | ✓ | ✓ |
| Sales by seller | ✓ | ✓ |
| Sales by seller totals · general sales | ✓ | — |
| Account statement · income statement | — | ✓ |

## Which should I use?

They serve different niches; both are valid.

| Your situation | Best fit |
| --- | --- |
| Read-only consulting of your data from claude.ai or ChatGPT | **Official MCP** (purpose-built for it) |
| Coding agent with shell access (Claude Code, Cursor, Codex) | **alegra-cli + skill** |
| An MCP host the official server doesn't support, or you need write tools | **alegra-cli**'s `alegra mcp` |
| Scripts, CI/CD, cron jobs | **alegra-cli** (pipes, exit codes, env vars) |
| Autonomous agent touching production data | **alegra-cli** (dry-run, confirmations, hooks) |
| Electronic invoice emission | **alegra-cli** |
| Working interactively in a terminal | **alegra-cli** |

Because `alegra-cli`'s MCP server is generated from its command tree, every accounting
operation you can run in the terminal is available to an agent (local setup and credential
commands stay out of the tool surface). See [Agent Skill](../user-guide/agent-skill.md),
[MCP Server](../user-guide/mcp.md), and [Agent Safety](../user-guide/agent-safety.md).
