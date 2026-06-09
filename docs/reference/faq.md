---
title: FAQ & Troubleshooting
---

# FAQ & Troubleshooting

### How do I authenticate in a script / CI?
Set environment variables — no keyring needed:
```bash
export ALEGRA_EMAIL="you@biz.com" ALEGRA_TOKEN="…"
alegra auth status
```
Env vars override the config file and keyring. See [Configuration](../user-guide/configuration.md).

### Where do I get my API token?
Alegra → **Configuración → API / Integraciones con otros sistemas**. The page
shows the email to use and lets you generate a token. Auth is HTTP Basic
(`base64(email:token)`) — the CLI handles the encoding.

### `alegra invoices create` fails with HTTP 400. Why?
Almost always a missing country-required field. Colombia needs a
`numberTemplate.id` (resolution) and a client with the right `identification`
type (NIT for companies); Mexico needs régimen/uso CFDI and a matching CP. Run
with `--dry-run` to inspect the body and see [Electronic Invoicing](../guides/electronic-invoicing.md).

### How do I pass a filter the CLI doesn't have a flag for?
Use the escape hatch:
```bash
alegra invoices list --param client_name="Acme" --param numberTemplate_fullNumber="FE-123"
```
Any Alegra query parameter works with `--param key=value` (repeatable).

### Why does `list` only return 30 rows?
Alegra caps `limit` at **30** per page. Use `--all` to fetch every page
automatically, or `--start`/`--limit` to page manually. For a total without
downloading, use `--count`.

### How do I get just a count?
```bash
alegra invoices list --status open --count
```
It uses the API's metadata total instead of downloading rows.

### Can I edit, void, or delete an invoice?
You can `update` a draft, but **stamped fiscal documents are immutable** by
design. To reverse one:

- **Void / reopen**: `alegra invoices void <id>` (annul) or `alegra invoices open <id>`.
  Payments have the same: `alegra payments void <id>` / `alegra payments open <id>`.
- **Correct the amount**: issue a credit note (`alegra credit-notes create`).

The CLI follows accounting's append-only discipline.

### A command returns HTTP 402 — is it broken?
No — `402` means your **Alegra plan doesn't include** that feature (some reports
and add-on/marketplace features). It's an account/plan limit, not a CLI bug. Run
`alegra doctor` to see which endpoints your plan can reach. See the
[Error reference](errors.md).

### `currencies get USD` — why a code, not a number?
Currencies are keyed by their ISO **code**, not a numeric id. `get`/`update`/
`delete` take the code: `alegra currencies get USD`.

### How do I switch between client accounts?
Use [profiles](../user-guide/configuration.md):
```bash
alegra config use clienteA
alegra --profile clienteB invoices list
```

### The output is hard to read / I need a spreadsheet.
Pick a format: `-o table` (default), `-o json`, `-o yaml`, `-o csv`, and narrow
columns with `--columns`:
```bash
alegra contacts list -o csv --columns id,name,email > contacts.csv
```

### Is there a sandbox?
The App API runs against your real account. Use a dedicated **test/demo Alegra
account** (a separate profile) for experiments, and `--dry-run` to preview
writes without sending them.

### How do I report a bug or request a feature?
Open an issue: <https://github.com/jjuanrivvera/alegra-cli/issues>.
