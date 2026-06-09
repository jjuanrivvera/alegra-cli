# alegra-cli vs. the official Alegra MCP server

Alegra publishes a hosted **MCP server** (`mcp.alegra.com`) so AI agents can call the
accounting API over the Model Context Protocol. `alegra-cli` takes a different route: it
is built **agent-first** as a command-line tool that AI agents drive directly — and it
*also* ships its own MCP server (`alegra mcp`) and works for humans in a terminal. This
page compares the two so you can pick the right fit.

!!! note "Snapshot"
    Reflects the official MCP server's publicly documented tool set and `alegra-cli`
    **v0.6.1**, verified **2026-06-09**. Both wrap the same Alegra v1 API and both require
    network access to it. Alegra may add tools over time — treat coverage figures as a
    point-in-time snapshot.

## At a glance

| | Official Alegra MCP | alegra-cli |
| --- | --- | --- |
| Interface | MCP tools (hosted, HTTP/SSE) | Terminal + agent **skill** + `alegra mcp` server |
| How an agent connects | MCP protocol (e.g. Claude Desktop) | A [skill](../user-guide/agent-skill.md) for shell-capable agents (Claude Code, Cursor, Codex) **or** [`alegra mcp`](../user-guide/mcp.md) for MCP hosts |
| Install | None — hosted by Alegra | Install a binary (Homebrew, Scoop, Docker, …) |
| Credentials | Configured in the MCP host/client | Alegra token in your **OS keyring** (profiles + env), kept out of agent config files |
| Output | Structured JSON tool results | table / JSON / YAML / CSV — composable with `jq` and pipes (fewer tokens for an agent) |
| Agent safety controls | — | `--dry-run` on any command, interactive confirm on `delete`, restrictable via shell hooks |
| Human (terminal) use | — | First-class |
| Resource coverage | ~24 resource areas | 37 resources |
| Invoice emission (stamp / void / email) | Not exposed | Yes |

## Built agent-first

`alegra-cli` was designed assuming the primary caller is an AI agent, not a person
typing commands. The rationale is laid out in the post
[CLIs over MCPs](https://jjuanrivvera.com/es/blog/clis-sobre-mcps/); in short:

- **Credentials stay in the OS keyring**, not in an agent's config file or a plaintext
  secret the model can read back.
- **Composable output.** An agent can `--fields`, pipe to `jq`, or `grep` — sending far
  fewer tokens to the model than a raw JSON tool result.
- **Safety is enforceable.** `--dry-run` previews the exact request on any command,
  `delete` asks for confirmation, and a shell can block destructive commands with hooks —
  so an autonomous agent *cannot* delete even if it tries.
- **Two ways in for agents.** A coding agent with shell access uses the **skill** (it
  learns the golden rules and when to call each command); an MCP-protocol host uses
  **`alegra mcp`**, which exposes the same command tree as MCP tools.

## Resource coverage

Both wrap the same Alegra v1 API but expose different slices of it.

### Covered by both
Contacts; items (with per-warehouse stock via `items stock`) and the inventory family
(item categories, variant attributes, warehouses, transfers, price lists, adjustments and
numerations, custom fields); bank accounts; reconciliations; journals; cost centers;
taxes; retentions; currencies; sellers; invoices (read/write); bills (with attachments,
comments, advances, and perceptions/retentions); supplier debit notes; purchase orders;
payments; and sales reports.

### Only in alegra-cli
Estimates, credit notes, customer debit notes, remissions, transportation receipts,
global invoices (CFDI), recurring invoices and payments, document numberings, payment
terms, additional charges, and webhook subscriptions — plus electronic invoice emission
(see below).

### Only in the official MCP
- **Catalog/reference lookups** — product keys, units, and reference enums
  (Mexico/country catalogs). These map only to the MCP's own endpoint with no documented
  public REST path, so the CLI does not expose them yet.
- **Support Center** help-desk tickets (out of scope for an accounting CLI).

## Capability differences that matter

**Electronic invoicing.** The official MCP's invoice tools cover create/read/update/
delete; they do not expose `stamp`, `void`, `email`, `open`, or `preview`. Emitting an
electronic invoice to DIAN/SAT, voiding, or emailing it is done with `alegra-cli`
(`alegra invoices stamp|void|email …`) or the REST API directly. If electronic emission
is part of your workflow, the CLI covers the full cycle.

**Payments.** The official MCP splits payments into incoming and outgoing tool families;
`alegra-cli` unifies them under one `payments` resource with a `type` flag, plus
`stamp`/`void`/`open`.

**Reports.** The two report sets are complementary:

| Report | Official MCP | alegra-cli |
| --- | --- | --- |
| Sales by client / totals | ✓ | ✓ |
| Sales by seller | ✓ | ✓ |
| Sales by seller totals · general sales · clients with items | ✓ | — |
| Account statement · income statement | — | ✓ |

## Which should I use?

They serve different niches; both are valid.

| Your situation | Best fit |
| --- | --- |
| Coding agent with shell access (Claude Code, Cursor, Codex) | **alegra-cli + skill** |
| Claude Desktop or another MCP host without a shell | Official MCP, or `alegra mcp` |
| Scripts, CI/CD, cron jobs | **alegra-cli** (pipes, exit codes, env vars) |
| Autonomous agent touching production data | **alegra-cli** (dry-run, confirmations, hooks) |
| Electronic invoice emission | **alegra-cli** |
| Working interactively in a terminal | **alegra-cli** |

Because `alegra-cli`'s MCP server is generated from its command tree, its agent-facing
coverage equals the CLI's: anything you can run in the terminal is available to an agent.
See [Agent Skill](../user-guide/agent-skill.md) and [MCP Server](../user-guide/mcp.md).
