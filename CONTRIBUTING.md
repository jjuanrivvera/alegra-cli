# Contributing to alegra-cli

Thanks for your interest in improving alegra-cli!

## Development setup

```bash
git clone https://github.com/jjuanrivvera/alegra-cli
cd alegra-cli
make setup-hooks   # install the pre-commit hook
make dev           # fmt + vet + build
make check         # fmt + vet + lint + test
```

Requires Go 1.25+. Linting uses `golangci-lint`.

## Architecture

See the [README](README.md#architecture). The key idea: each Alegra resource is
a thin typed wrapper over a generic core.

### Adding a resource

A resource needs three files and no edits to shared code:

1. `internal/api/<resource>.go` — the typed struct(s) and a `Client` accessor:

   ```go
   package api

   type Widget struct {
       ID   ID     `json:"id,omitempty"`
       Name string `json:"name,omitempty"`
   }

   func (c *Client) Widgets() *Resource[Widget] {
       return NewResource[Widget](c, "widgets")
   }
   ```

2. `commands/<resource>.go` — register it (it self-attaches via `init()`):

   ```go
   package commands

   import "github.com/jjuanrivvera/alegra-cli/internal/api"

   func init() {
       registerResource(resourceSpec[api.Widget]{
           Use:     "widgets",
           Short:   "Manage widgets",
           New:     func(c *api.Client) *api.Resource[api.Widget] { return c.Widgets() },
           Columns: []string{"id", "name"},
       })
   }
   ```

3. `internal/api/<resource>_test.go` — an httptest-based service test (reuse the
   `newTestClient` helper).

Custom actions (e.g. `void`, `email`) go through the `Extra` hook using
`NewActionCmd` / `NewCollectionActionCmd`. Non-CRUD resources (singletons,
reports) build a plain cobra command using `client.GetInto/PostInto/PutInto`.

## Commits & branches

- Branch from `develop`: `feature/...` or `fix/...`.
- Keep commit messages clear and descriptive.
- Run `make check` before pushing; CI must be green.

## Tests

- Service tests use `httptest.NewServer`.
- Prefer `require` for fatal assertions, `assert` for the rest.
