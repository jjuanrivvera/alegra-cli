# alegra-cli

A fast, scriptable command-line interface for the [Alegra](https://www.alegra.com/)
accounting API — manage contacts, invoices, items, payments, taxes, and the full
Alegra resource surface from your terminal, with `table`/`json`/`yaml`/`csv`
output, named profiles, a dry-run mode, and a built-in MCP server.

Built in Go, mirroring the architecture of
[canvas-cli](https://github.com/jjuanrivvera/canvas-cli).

> Unofficial. Not affiliated with Alegra. Uses the public API at
> `https://api.alegra.com/api/v1`.

## Install

```bash
# Homebrew (once the tap is published)
brew install jjuanrivvera/alegra-cli/alegra-cli

# From source
go install github.com/jjuanrivvera/alegra-cli/cmd/alegra@latest

# Or build locally
make build && ./bin/alegra --help
```

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
