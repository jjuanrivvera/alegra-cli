<!--
Target `develop`, not `main`. Use a Conventional Commit title
(e.g. `feat(invoices): ...`, `fix(output): ...`, `docs: ...`).
-->

## What & why

<!-- What does this change and why? Link issues: "Closes #123". -->

## How

<!-- Notable implementation details, trade-offs, or alternatives. -->

## Checklist

- [ ] Targets the `develop` branch
- [ ] `make check` passes locally (fmt + vet + lint + test)
- [ ] Added/updated tests (coverage stays healthy)
- [ ] Ran `make docs-gen` if the command tree changed (and committed the result)
- [ ] Updated `CHANGELOG.md` under `[Unreleased]` for user-facing changes
- [ ] No credentials/tokens in code, tests, fixtures, or commit messages
