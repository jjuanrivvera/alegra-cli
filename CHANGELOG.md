# Changelog

All notable changes to this project are documented here. The format is based on
[Keep a Changelog](https://keepachangelog.com/) and this project adheres to
[Semantic Versioning](https://semver.org/).

## [Unreleased]

## [0.8.1] - 2026-06-10

### Fixed
- **catalogen refuses to write empty catalogs.** A format change on Alegra's
  country parameter pages that parsed to zero entries would have silently
  hollowed out the embedded per-country catalogs on the next `catalog-sync`;
  the generator now aborts instead. Its fetch→parse→write pipeline is covered
  end-to-end by tests (58% → 84.5%), and regenerated output was verified
  byte-identical against the live pages.

## [0.8.0] - 2026-06-10

Feature release: SAT product-keys catalog for Mexican accounts (full catalog
parity with the official MCP), lossless money amounts, and field-level contract
testing against the documented API.

### Added
- **SAT product-keys catalog (México)** — closes the last catalog gap vs the
  official MCP. `alegra catalog sync-sat` downloads the SAT's `c_ClaveProdServ`
  catalog (~52k keys, sourced from the SAT's published data via the phpcfdi
  mirror) into a shared local cache, and `alegra catalog product-keys [query]`
  searches it offline, accent- and case-insensitively, across keys, names, and
  the SAT's similar-names lists. `alegra init` offers the download when it
  detects a Mexican account (interactive only, never fatal), and `alegra
  doctor` reports the catalog's state on Mexican accounts.
- **Field-level contract tests against the documented API.** `make spec-sync`
  now also harvests the OpenAPI definitions embedded in every REST reference
  page into `testdata/spec/schemas.json`, and a contract test verifies every
  typed struct field exists in the documented response fields with a compatible
  JSON type — catching phantom fields and type drift that resource-level
  matching cannot.
- `api.Refs` flexible type for related-record fields Alegra serializes as
  either one `{id,name}` object or an array of them; `income-debit-notes`'
  `costCenter` (documented in both shapes) now uses it — found by the new
  contract test.
- **Replay-style integration tests**: recorded API fixtures with the real
  API's awkward shapes (string amounts with padded decimals, mixed id types)
  exercised end-to-end through the command tree, asserting rendered output and
  wire behavior (pagination requests, Basic auth).

### Changed
- **`Money` now preserves the exact decimal text the API sent** instead of
  converting through `float64`. Amounts can no longer lose digits (values
  beyond ~15 significant figures were silently rewritten) and are no longer
  normalized (`10.00` stays `10.00`, not `10`). Output note: amount fields the
  API returns as explicit `0` now appear in JSON/YAML output instead of being
  omitted — higher fidelity to what the API actually said. Amounts still render
  as plain JSON numbers, so `jq`-style consumers are unaffected.

### Fixed
- specsync refuses to write an empty endpoint manifest when the upstream
  `llms.txt` format changes, instead of silently disabling drift detection;
  its parser is now covered by tests against real-index fixtures.

## [0.7.1] - 2026-06-10

Hardening release: fiscal-safety fixes around electronic emission, decoder
correctness, and crash-safe persistence. No new commands, no breaking changes.

### Fixed
- **`invoices emit` idempotency guard** — three silent failure modes that could
  lead to double electronic (fiscal) emission:
  - A corrupt or unreadable `emitted.json` was silently treated as an empty
    cache, dropping the guard. `emit` now aborts with the cache path and a
    remediation hint (a missing file is still a normal fresh start).
  - The cache was persisted once after all batches; a crash mid-run lost every
    stamped id. It is now saved after each successful batch, and emission stops
    (listing the affected ids) if the save fails.
  - The cache file was written in place and could be torn by a crash mid-write;
    writes now go through temp-file + rename.
- **`Money` accepted `"NaN"`/`"Inf"` strings**, producing a value that breaks
  any later JSON re-serialization of the document; `Int` converted them through
  an undefined `int64(NaN)`. Both decoders now reject non-finite values.
  (Found by the new value-level fuzz properties within seconds.)
- **`ID` and `Int` lost precision above 2⁵³** by round-tripping integers
  through `float64`; both now decode via `int64` first.
- `import --dry-run` and `create --draft` ignored `json.Marshal` errors,
  risking garbage output or a stale body sent to the API.
- `auth status` ignored config-load errors and could dereference a nil config.
- `config.yaml` is now written atomically (temp + rename), same crash-safety
  treatment as the emission cache.
- `AsAPIError` now delegates to `errors.As`, so `errors.Join` multi-error
  trees are unwrapped correctly.

### Changed
- Table/CSV output: when the auto-detected column set is capped at 10, a note
  on stderr now reports how many columns were dropped and suggests `--columns`
  or `--output json` (stdout stays clean for piping).
- `--all` pagination: the page-cap warning now includes a remediation hint.
- Resource list filters that would collide with built-in flags are recorded at
  registration and asserted empty by a registry test, so a bad resource
  definition fails CI instead of silently losing the filter.
- CONTRIBUTING: resource tests should cover what is unique to the resource
  (not re-test generic CRUD); documented failure-path testing and fuzzing
  expectations.

## [0.7.0] - 2026-06-09

Coverage parity with Alegra's official MCP (except Support Center and the SAT
product-key catalog, which Alegra exposes via neither REST nor a working MCP).

### Added
- **`alegra catalog`** (aliases `catalogs`/`reference`) — offline per-country
  reference catalogs: units of measure and reference enums (identification types,
  tax types, payment methods, document types, régimenes, …) for Colombia, Mexico,
  Costa Rica, Peru, Spain, and Panama. Embedded data generated by `tools/catalogen`
  (`make catalog-sync`) from Alegra's published country parameter pages; works
  offline with `--country`, no login.
- **`alegra items stock <id>`** — per-warehouse stock read from the item
  inventory, with `--date` for a historical snapshot. The inventory
  `warehouses[]` array is now modeled.
- **`bills` sub-actions**: `advances`, `attach`, `attachment-delete`,
  `perceptions`, `retentions`, `comment-update`, `comment-delete`.
- A "vs. the official MCP" comparison page in the docs.

### Changed
- Brought the README, MkDocs site, and agent skill up to date with the full
  v0.6.x/v0.7.0 feature set (bulk import/export, the MCP editor-wiring/stream/tools
  surface, country detection + pre-flight validation, reference catalogs).
- Internal: generic `NewPutActionCmd` builder and `Client.DeleteInto`.

### Fixed
- Docs: `claude mcp add alegra -- alegra mcp` now correctly uses
  `alegra mcp start` (the bare `alegra mcp` prints help and would not serve).

## [0.6.1] - 2026-06-09

### Changed
- Every resource command now ships a `Long` description. Because `ophis`
  surfaces a command's `Long` as its MCP tool description, this brings the
  `alegra mcp` tool docs (and `alegra <resource> --help`) to parity with the
  official Alegra MCP server's per-tool descriptions (#22).

## [0.6.0] - 2026-06-09

API fidelity & coverage release (#22, #27).

### Added
- **`sellers`** (vendedores) and **`webhook-subscriptions`** resources, both
  verified against the live API.
- **Spec tooling**: `make spec-sync` reconstructs a manifest of the documented
  API surface (`testdata/spec/endpoints.json`) from the official docs index;
  `make spec-check` is a network-free CI guardrail asserting every CLI resource
  is documented; a weekly workflow opens a PR when the surface drifts.
- **Live smoke-test suite** (`make smoke`, build-tagged `smoke`): validates the
  structs against what the API returns at runtime, with unknown-field detection
  and an optional create→update→delete write cycle on safe master-data. Runs
  weekly / on demand, never on PRs.
- Richer resource help text (`Long`) that flows to both `--help` and the
  `alegra mcp` tool descriptions.

### Fixed
- **Reconciliations** now uses the documented `/conciliations` path (the old
  `/reconciliations` returned HTTP 403); `doctor`'s plan probe no longer reports
  a false 403.
- Struct↔schema fidelity: seller `identification` decodes as a number; the
  webhook list response's `subscriptions` wrapper is recognized.

### Changed
- Audited the resource structs against the live API (0 decode failures across 37
  resources); documented that contact `identification` is a string on read but
  an object on write, and pinned it with a regression test.

## [0.5.0] - 2026-06-08

Polish & distribution release.

### Added
- **`alegra init`** — guided onboarding: authenticate, auto-detect the account
  country, save the profile, and print next steps.
- **`version --json`** (structured build info) and **`version --check`** (compare
  against the latest GitHub release).
- **Colorized output** for `doctor` (respects `NO_COLOR`; global `--no-color`).
- **Shell completions** (bash/zsh/fish) bundled in release archives, Homebrew,
  Scoop, and the Linux packages.
- **New install channels**: `.deb`/`.rpm`/`.apk` packages, a Docker image on
  **GHCR** (`ghcr.io/jjuanrivvera/alegra-cli`), and a **Scoop** bucket.
- **Supply chain**: SBOM (syft) per archive and **cosign keyless signatures** of
  the checksums file.
- README animated demo (VHS) + DeepWiki badge + Codecov badge; social preview.

### Changed
- README header/badges centered; `canvas-cli` references removed.
- A CI coverage gate now fails the build below 80% (coverage ~86%).

## [0.4.4] - 2026-06-08

Test-quality release — no functional changes to the CLI.

### Changed
- Cover every resource accessor and the custom record/collection actions; total
  statement coverage rises to ~84% (Codecov line coverage clears 80%).

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
