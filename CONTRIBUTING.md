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

Requires Go 1.25+ (the `toolchain` in `go.mod` pins the patched 1.25.x used by
CI). Linting uses `golangci-lint` **v2** (`go install
github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6`).

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
   `newTestClient` helper). Test what is **unique** to the resource — special
   field types, custom actions, odd response shapes — not the generic CRUD
   plumbing, which is already covered once by the `Resource[T]` tests
   (`resource_test.go`, `client_failures_test.go`). A List/Get happy-path pair
   adds volume, not signal.

Custom actions (e.g. `void`, `email`) go through the `Extra` hook using
`NewActionCmd` / `NewCollectionActionCmd`. Non-CRUD resources (singletons,
reports) build a plain cobra command using `client.GetInto/PostInto/PutInto`.

## Commits & branches

- **Branch from `develop`** (not `main`) and open PRs against `develop`. Use a
  type-prefixed branch name matching the change: `feat/...`, `fix/...`,
  `docs/...`, `test/...`, `chore/...`.
- Write [Conventional Commits](https://www.conventionalcommits.org/), e.g.
  `feat(invoices): add emit --all` or `fix(output): rune-safe truncation`. The
  [CHANGELOG](CHANGELOG.md) follows [Keep a Changelog](https://keepachangelog.com/).
- If you change the command tree (new resource, flag, or help text), regenerate
  the command reference with `make docs-gen` and commit the result.
- Run `make check` before pushing; CI must be green.

## Tests

- Service tests use `httptest.NewServer` (see the `newTestClient` helper in
  `internal/api`).
- Prefer `require` for fatal assertions, `assert` for the rest.
- Keep coverage healthy — the suite sits above 80%; new code should ship tests.
- **Test failure paths, not just happy paths.** Every parse of external state
  (API bodies, config files, caches) needs a test with corrupt input; every
  batch operation needs a partial-failure test asserting counts and a non-zero
  exit. Coverage measures execution, not assertion quality — a swallowed error
  can be 100% "covered" and still hide a bug.
- The flexible JSON types have fuzzers with value-level properties
  (`internal/api/fuzz_test.go`); run them after touching a decoder:
  `go test ./internal/api -fuzz '^FuzzID$' -fuzztime 30s` (likewise `FuzzInt`,
  `FuzzMoney`, `FuzzStringOrSlice`). Counterexamples land in `testdata/fuzz/`
  and become permanent regression cases — commit them.

## Reporting bugs & security issues

- Bugs and feature requests: open an issue (templates guide the details we need).
- **Security vulnerabilities**: do **not** open a public issue — see
  [SECURITY.md](SECURITY.md).
- By participating you agree to our [Code of Conduct](CODE_OF_CONDUCT.md).
