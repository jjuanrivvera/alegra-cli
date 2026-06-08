# AGENTS.md

Guidance for AI agents (Claude Code, Cursor, Copilot, …) working in this repo.

## What this is

`alegra-cli` is a Go CLI for the Alegra accounting API
(`https://api.alegra.com/api/v1`), built with Cobra. The architecture is a
generic typed client plus one thin file per resource.

## Commands

```bash
make build        # build to bin/alegra
make dev          # fmt + vet + build
make test         # go test ./...
make lint         # golangci-lint
make check        # fmt + vet + lint + test
make docs-gen     # regenerate docs/commands from the cobra tree
go run ./cmd/alegra <args>
```

## Architecture

```
cmd/alegra/main.go        entry point
commands/
  root.go                 global flags, getAPIClient(), render()
  generic.go              generic CRUD + action command builders
  <resource>.go           one per resource; self-registers via init()
internal/
  api/
    client.go             auth (Basic email:token / Bearer), retries, rate limit, dry-run
    resource.go           generic Resource[T]: List/Get/Create/Update/Delete/Action
    pagination.go         start/limit offset pagination
    types.go              ID, Ref, Money flexible JSON types
    <resource>.go         one per resource: struct(s) + Client accessor
  config/                 YAML profiles + env overrides
  auth/                   OS keyring token storage
  output/                 table/json/yaml/csv rendering
  version/                build metadata (ldflags)
tools/gendocs/            command reference generator
```

## Key patterns

- **Resource = struct + accessor + registration.** Add a resource with three
  files and zero edits to shared code (see [CONTRIBUTING.md](CONTRIBUTING.md)).
- **Generic core.** `internal/api/resource.go` and `commands/generic.go` provide
  CRUD; resources only declare types, columns, filters, and custom actions.
- **Flexible JSON types.** Use `api.ID` for identifiers, `api.Money` for amounts,
  `*api.Ref` for nested `{id,name}` objects. Unknown JSON fields are ignored, so
  structs need not be exhaustive.
- **Auth.** HTTP Basic with email + API token. Tokens live in the OS keyring;
  email/baseURL in `~/.alegra-cli/config.yaml`. Env overrides: `ALEGRA_EMAIL`,
  `ALEGRA_TOKEN`, `ALEGRA_BASE_URL`, `ALEGRA_PROFILE`, `ALEGRA_OUTPUT`.
- **Pagination.** Alegra uses offset `start`/`limit` (max 30). `--all` walks all
  pages via `Resource.ListAll`.
- **Dry-run.** `--dry-run` prints the equivalent curl and skips the request.

## Testing

Service tests use `httptest.NewServer` via the `newTestClient(t, handler)` helper
in `internal/api/contacts_test.go`. Use `require` for fatal checks, `assert`
otherwise.

## Conventions

- Comments explain WHY, not WHAT.
- `gofmt -s` clean; pass `golangci-lint`.
- Never commit credentials. Tokens belong in the keyring or env, never in code,
  config-in-repo, or commit messages.
- Reference the upstream API docs index at `https://developer.alegra.com/llms.txt`
  (each `/reference/<slug>.md` page embeds the full OpenAPI definition).
