# Changelog

All notable changes to this project are documented here. The format is based on
[Keep a Changelog](https://keepachangelog.com/) and this project adheres to
[Semantic Versioning](https://semver.org/).

## [Unreleased]

## [0.3.1]

### Security
- Neutralize CSV formula injection (CWE-1236): `csv` output cells beginning with
  `=`, `+`, `@`, tab, or CR (and a leading `-` that isn't a real number) are
  prefixed with a quote so spreadsheets treat them as text.
- Clamp `Retry-After: 0` so it no longer triggers an immediate zero-delay retry
  (matching the `X-Rate-Limit-Reset` guard).
- Warn when the configured base URL is not HTTPS.

### Added
- Developer `worktree` skill (`wt` CLI) for building/running/testing each branch
  in an isolated worktree. Dev tooling only — not part of the released binary.

## [0.3.0]

### Added
- **Agent skill** — alegra-cli now ships a `SKILL.md` (at the repo root) that
  teaches AI agents how to drive the CLI. Install across agents with
  `npx skills add jjuanrivvera/alegra-cli`.
- **`alegra skills install`** — write the bundled skill (embedded via `go:embed`)
  into an agent's skills directory, with `--global`, `--agent`
  (claude/cursor/windsurf/codex/gemini/copilot/opencode), `--dir`, and
  `--dry-run`; plus `alegra skills path` and `alegra skills print`.
- **Claude Code plugin** — `.claude-plugin/plugin.json` + `marketplace.json` so
  the skill installs via `/plugin marketplace add jjuanrivvera/alegra-cli`.
- `references/alegra-commands.md` — a condensed command cheatsheet bundled with
  the skill.

## [0.2.0]

### Added
- **`alegra doctor`** — read-only diagnostics for config, credentials, auth,
  company/country, rate-limit budget, numbering resolutions, and plan access.
- **Aliases** — `alegra alias set/list/remove`; a saved alias expands before
  parsing and never shadows a built-in command.
- **Natural date ranges** — `--since`/`--until` on every `list`
  (`this-month`, `last-month`, `7d`, `3m`, `YYYY-MM-DD`, …).
- **CSV import/export** — `<resource> import --file --map --set` (per-row,
  continue-on-error, dotted nested paths) and `<resource> export` (auto-paginated
  CSV/JSON).
- **Electronic-invoice emit lifecycle** — `invoices emit` stamps drafts in
  auto-chunked batches of 10 with a local idempotency guard (`--force` to
  override); `create --draft` keeps a document internal.
- **Country-aware pre-flight validation** — `create` checks the body against
  per-country rules (CO/MX/PE/CR) before sending; `config set-country`,
  `--country`, and `--no-validate` to control it.

### Changed
- **Friendlier errors** — every failure now suggests a fix (e.g. `402` → "plan
  doesn't include this", `429` → rate-limit hint, stamping `AEP*/EPR*` codes →
  remedies).
- **Adaptive rate limiting** — the client reads `X-Rate-Limit-*` headers, slows
  down as the quota drains, and waits the exact reset window on `429`.

## [0.1.0]

### Added
- Initial release of `alegra-cli`.
- Full Alegra v1 resource surface with `list`/`get`/`create`/`update`/`delete`
  plus resource-specific actions (void, open, email, stamp, transfer, comments,
  close, import-by-cufe, …).
- Generic typed API client (`Resource[T]`) with HTTP Basic auth, exponential
  backoff retries, adaptive client-side rate limiting, and offset pagination
  (`start`/`limit`, auto `--all`).
- `table`, `json`, `yaml`, and `csv` output with `--columns` selection.
- Named profiles, environment-variable overrides, and OS keyring token storage
  (`alegra auth login`).
- `--dry-run` mode that prints the equivalent `curl` request.
- Built-in MCP server (`alegra mcp`) exposing the command tree to AI agents.
- MkDocs Material documentation site and GoReleaser-based release pipeline.
