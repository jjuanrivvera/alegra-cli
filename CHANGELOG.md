# Changelog

All notable changes to this project are documented here. The format is based on
[Keep a Changelog](https://keepachangelog.com/) and this project adheres to
[Semantic Versioning](https://semver.org/).

## [Unreleased]

## [0.4.3] - 2026-06-08

Test/CI accuracy release — no functional changes to the CLI.

### Added
- Test covering the docs generator's `run()` (raises statement coverage to ~82%).
- `codecov.yml`: ignore the binary entry point (`cmd/alegra`) and keep coverage
  status checks informational, so the Codecov badge reflects the testable surface.

## [0.4.2] - 2026-06-08

Documentation/CI polish — no functional changes to the CLI.

### Added
- Codecov coverage badge, backed by a Codecov upload step in the CI coverage job.

### Changed
- README: center the title and badge row.

### Removed
- All `canvas-cli` references across the README, AGENTS.md, and RELEASING.md.

## [0.4.1] - 2026-06-08

Documentation, repository, and test-quality release — no functional changes to
the CLI.

### Added
- README badges (CI, release, pkg.go.dev, Go Report Card, Go version, license,
  Ask DeepWiki).
- Community-health files: `SECURITY.md`, `CODE_OF_CONDUCT.md`, issue forms
  (bug/feature), a pull request template, and Dependabot config (weekly gomod +
  github-actions updates).

### Changed
- `CONTRIBUTING.md` aligned with the repo's Conventional Commits / `develop`
  workflow and a bug/security reporting section.
- Test coverage raised from 39.8% to 80.6% across the testable packages.

## [0.4.0] - 2026-06-07

### Added
- **Auto-detect the account's country (version).** Alegra is one API whose
  required fields, enums, and electronic-invoicing flow are localized per country
  (exposed as the company's `applicationVersion`). `auth login` now detects it
  and caches it on the profile (`profiles.<name>.country`); `doctor` refreshes
  it. Pre-flight validation reads this detected value, so it stays in sync with
  the actual account instead of a hand-set string.

### Changed
- `config set-country` is now an **offline fallback hint** only — the detected
  per-profile country takes precedence. Pre-flight validation country resolution
  is now: `--country` flag > detected profile country > `set-country` hint.
- Skill (`skills/alegra-cli`) documents the per-country API differences and
  instructs detecting the connected version with `alegra company get` before any
  country-specific write (skill `0.4.0`).

## [0.3.2]

### Changed
- Pre-commit hook (`.githooks/pre-commit`) now resolves a **golangci-lint v2**
  binary (PATH or `GOPATH/bin`) — the config is v2 format, which v1 binaries
  reject — and skips the lint step gracefully when no v2 binary is present.
  Also fixed a hang when a commit staged no Go files. Enable with
  `make setup-hooks`.

## [0.3.1]

### Fixed
- `table`/`csv` truncation is now rune-safe: long values with multi-byte UTF-8
  characters (accented Spanish names) are no longer split mid-character into
  invalid output.

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
