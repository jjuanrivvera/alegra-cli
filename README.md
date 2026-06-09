<div align="center">

# alegra-cli

[![CI](https://github.com/jjuanrivvera/alegra-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/jjuanrivvera/alegra-cli/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/jjuanrivvera/alegra-cli)](https://github.com/jjuanrivvera/alegra-cli/releases/latest)
[![codecov](https://codecov.io/gh/jjuanrivvera/alegra-cli/branch/main/graph/badge.svg)](https://codecov.io/gh/jjuanrivvera/alegra-cli)
[![Go Reference](https://pkg.go.dev/badge/github.com/jjuanrivvera/alegra-cli.svg)](https://pkg.go.dev/github.com/jjuanrivvera/alegra-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/jjuanrivvera/alegra-cli)](https://goreportcard.com/report/github.com/jjuanrivvera/alegra-cli)
[![Go version](https://img.shields.io/github/go-mod/go-version/jjuanrivvera/alegra-cli)](go.mod)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/jjuanrivvera/alegra-cli)

<br>

<img src="assets/demo.gif" alt="alegra-cli in action" width="900">

</div>

A fast, scriptable command-line interface for the [Alegra](https://www.alegra.com/)
accounting API — manage contacts, invoices, items, payments, taxes, and the full
Alegra resource surface from your terminal, with `table`/`json`/`yaml`/`csv`
output, named profiles, a dry-run mode, and a built-in MCP server.

Built in Go with a small generic core: a typed API client plus one thin file
per resource.

> Unofficial. Not affiliated with Alegra. Uses the public API at
> `https://api.alegra.com/api/v1`.

## Install

```bash
# Homebrew (macOS/Linux)
brew install jjuanrivvera/alegra-cli/alegra-cli

# Scoop (Windows)
scoop bucket add alegra-cli https://github.com/jjuanrivvera/scoop-alegra-cli
scoop install alegra-cli

# Docker
docker run --rm ghcr.io/jjuanrivvera/alegra-cli:latest version

# From source
go install github.com/jjuanrivvera/alegra-cli/cmd/alegra@latest

# Or build locally
make build && ./bin/alegra --help
```

Linux `.deb`/`.rpm`/`.apk` packages are attached to each
[release](https://github.com/jjuanrivvera/alegra-cli/releases/latest). Release
archives, packages, and the Homebrew/Scoop installs all ship **shell
completions** (bash/zsh/fish). Releases are also **SBOM-attested and the
checksums are signed** with [cosign](https://github.com/sigstore/cosign)
(keyless); see [RELEASING.md](RELEASING.md) for verification.

## Authenticate

Alegra uses HTTP Basic auth with your account **email** and an **API token**
(Alegra → Configuración → Integraciones → API). Two ways:

```bash
# Interactive: token is stored in your OS keyring, never written to disk
alegra auth login

# Or via environment variables (great for CI / scripts)
export ALEGRA_EMAIL="you@example.com"
export ALEGRA_TOKEN="your-api-token"

alegra auth status   # verifies against /users/self
```

## Usage

```bash
# List the first page of contacts as a table
alegra contacts list

# Filter + paginate, fetch every page
alegra contacts list --type client --all --limit 30

# Get one record as JSON
alegra invoices get 12 -o json

# Create from a JSON file (or stdin with -f -)
alegra invoices create -f new-invoice.json

# Create with inline fields (values parsed as JSON when valid)
alegra contacts create --set name="Acme S.A.S" --set 'type=["client"]'

# Custom actions
alegra invoices void 12
alegra invoices email 12 --set 'emails=["client@acme.com"]'
alegra payments stamp 88

# See the exact request without sending it
alegra invoices list --status open --dry-run
```

### Filtering, counting & pagination

```bash
# Natural date ranges (this-month, last-month, 7d, 3m, YYYY-MM-DD, ...)
alegra invoices list --status open --since this-month
alegra payments list --type in --since 7d

# Just the count (uses the API's metadata total — no full fetch)
alegra invoices list --status open --count

# Escape hatch: pass ANY Alegra query parameter the CLI doesn't expose as a flag
alegra contacts list --param identification=900123456 --param order_field=name

# Fetch every page
alegra items list --all
```

### Power features

```bash
alegra doctor                              # diagnose config, auth, plan, rate limit
alegra alias set unpaid "invoices list --status open --all"
alegra unpaid --client-id 12               # run the alias (+ extra args)
alegra contacts import -f clients.csv --map 'Name=name,NIT=identification.number'
alegra invoices export --param status=open > receivables.csv
alegra invoices create -f invoice.json --draft   # country pre-flight validated
alegra invoices emit --all                 # stamp drafts, auto-chunked, idempotent
```

### Output formats

`-o table` (default), `-o json`, `-o yaml`, `-o csv`. Pick columns with
`--columns`:

```bash
alegra items list -o csv --columns id,name,price > items.csv
alegra contacts list --columns id,name,email
```

### Profiles

Manage multiple Alegra accounts:

```bash
alegra config set-profile --name prod --email you@biz.com
alegra auth login --profile prod
alegra --profile prod invoices list
alegra config use prod          # set the default
```

Config lives at `~/.alegra-cli/config.yaml` (tokens stay in the keyring).

### MCP server

Expose the entire CLI to AI agents over the Model Context Protocol:

```bash
alegra mcp           # see MCP subcommands
```

## Resources

Full Alegra v1 surface, each with `list`/`get`/`create`/`update`/`delete`
(plus resource-specific actions):

`contacts`, `items`, `item-categories`, `invoices`, `global-invoices`,
`recurring-invoices`, `credit-notes`, `income-debit-notes`, `estimates`,
`remissions`, `transportation-receipts`, `bills`, `debit-notes`,
`purchase-orders`, `payments`, `recurring-payments`, `taxes`, `retentions`,
`terms`, `currencies`, `number-templates`, `bank-accounts`, `reconciliations`,
`cost-centers`, `journals`, `categories`, `additional-charges`, `company`,
`users`, `warehouses`, `warehouse-transfers`, `inventory-adjustments`,
`inventory-adjustment-numerations`, `price-lists`, `custom-fields`,
`variant-attributes`, `reports`.

Run `alegra <resource> --help` for actions and filters.

## Use it from your AI agent (skill)

alegra-cli ships an **agent skill** that teaches AI coding agents (Claude Code,
Cursor, Codex, Gemini CLI, Windsurf, Copilot, …) how to drive it. Install the
skill across every agent you have with one command:

```bash
npx skills add jjuanrivvera/alegra-cli
```

Or use one of the built-in / native paths:

```bash
alegra skills install --global              # write the bundled skill (no Node needed)
alegra skills install --agent cursor        # target a specific agent

# Native Claude Code plugin:
#   /plugin marketplace add jjuanrivvera/alegra-cli
#   /plugin install alegra-cli@alegra
```

The skill wraps this binary, so install the CLI (above) and authenticate first.
For structured tool access (Claude Desktop, etc.) the CLI is also an MCP server:
`claude mcp add alegra -- alegra mcp`.

## Guides & recipes

The [documentation site](https://jjuanrivvera.github.io/alegra-cli/) has more than
flag references — it's full of real accounting recipes:

- **[Cookbook](https://jjuanrivvera.github.io/alegra-cli/cookbook/)** — copy-paste commands for everyday tasks
- **[Invoice → Cash](https://jjuanrivvera.github.io/alegra-cli/guides/invoice-to-cash/)** — a full sales cycle: client → item → invoice → email → payment
- **[Electronic Invoicing](https://jjuanrivvera.github.io/alegra-cli/guides/electronic-invoicing/)** — the draft→open→stamp lifecycle, `numberTemplate`/`stamp`, and per-country fields (CO/MX/PE/CR)
- **[Expenses & Purchases](https://jjuanrivvera.github.io/alegra-cli/guides/expenses-and-purchases/)** — suppliers, bills, purchase orders, outgoing payments
- **[Reporting & Month-End](https://jjuanrivvera.github.io/alegra-cli/guides/reporting-and-month-end/)** — totals, aging, and closing snapshots
- **[Productivity Features](https://jjuanrivvera.github.io/alegra-cli/guides/productivity/)** — `doctor`, aliases, date ranges, counts, validation
- **[Automation & Scripting](https://jjuanrivvera.github.io/alegra-cli/guides/automation/)** — CSV bulk ops, cron, CI, MCP agents
- **[Errors & Rate Limits](https://jjuanrivvera.github.io/alegra-cli/reference/errors/)** · **[FAQ](https://jjuanrivvera.github.io/alegra-cli/reference/faq/)**

## Configuration reference

| Env var | Meaning |
| --- | --- |
| `ALEGRA_EMAIL` | Account email (Basic auth user) |
| `ALEGRA_TOKEN` | API token (Basic auth password) |
| `ALEGRA_BEARER_TOKEN` | OAuth bearer token (marketplace apps) |
| `ALEGRA_BASE_URL` | Override the API base URL |
| `ALEGRA_PROFILE` | Active profile name |
| `ALEGRA_OUTPUT` | Default output format |
| `ALEGRA_CONFIG` | Override the config file path |

Global flags: `--profile`, `-o/--output`, `--base-url`, `--columns`,
`--requests-per-second`, `--dry-run`, `--show-token`, `-v/--verbose`.

## Architecture

```
cmd/alegra/         entry point
commands/           cobra command tree (one file per resource, init-registered)
  root.go           global flags, client construction, output rendering
  generic.go        generic list/get/create/update/delete + action builders
internal/
  api/              typed client + generic Resource[T] + one file per resource
    client.go       auth, retries, rate limiting, dry-run
    resource.go     generic CRUD over the REST collection
    pagination.go   start/limit offset pagination
  config/           YAML profiles + env overrides
  auth/             OS keyring token storage
  output/           table / json / yaml / csv rendering
  version/          build metadata
tools/gendocs/      command reference generator (MkDocs)
```

Each resource is a thin typed wrapper: a Go struct + a one-line `Client`
accessor + a small command registration. The generic core handles HTTP, auth,
retries, adaptive rate limiting, pagination, and output.

## Development

```bash
make dev            # fmt + vet + build
make test           # run tests
make lint           # golangci-lint
make check          # full local quality gate
make docs-serve     # preview the docs site
```

## License

MIT — see [LICENSE](LICENSE).
