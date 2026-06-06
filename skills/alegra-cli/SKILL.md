---
name: alegra-cli
description: Manage the Alegra accounting API (https://alegra.com) from the terminal with the `alegra` CLI — contacts/clients, items, sales invoices, payments, bills, credit notes, taxes, reports, and DIAN/SAT electronic invoicing. Use this whenever the user wants to create or list invoices, manage clients or products, record or reconcile payments, run sales reports, bulk-import accounting data, or emit electronic (fiscal) invoices in Colombia, Mexico, Peru, or Costa Rica.
version: 0.3.0
homepage: https://github.com/jjuanrivvera/alegra-cli
license: MIT
allowed-tools: Bash(alegra:*)
metadata: {"openclaw":{"category":"finance","emoji":"🧾","requires":{"bins":["alegra"],"env":["ALEGRA_EMAIL","ALEGRA_TOKEN"]},"install":[{"kind":"brew","formula":"jjuanrivvera/alegra-cli/alegra-cli","bins":["alegra"]},{"kind":"go","package":"github.com/jjuanrivvera/alegra-cli/cmd/alegra@latest","bins":["alegra"]}]}}
---

# Alegra CLI

Drive the [Alegra](https://alegra.com) accounting API through the `alegra`
command-line tool. This skill teaches you how and when to use it.

## Prerequisites

- The `alegra` binary must be on `PATH`. Check with `alegra version`. If missing,
  install it: `brew install jjuanrivvera/alegra-cli/alegra-cli` or
  `go install github.com/jjuanrivvera/alegra-cli/cmd/alegra@latest`.
- Credentials: either run `alegra auth login` once, or set `ALEGRA_EMAIL` and
  `ALEGRA_TOKEN` in the environment. Confirm with `alegra auth status` /
  `alegra doctor`.

## Golden rules (read before acting)

1. **Preview writes.** For any create/update/delete/emit, run it once with
   `--dry-run` first to show the exact request, then run for real after the user
   confirms. Never emit invoices in a blind retry loop — emission is **not
   idempotent**.
2. **Parse with JSON.** Add `-o json` and pipe to `jq` when you need to read
   values; the default `table` output is for humans.
3. **Invoices are append-only.** Never try to edit or delete a stamped invoice —
   issue a credit note (`alegra credit-notes create`).
4. **Count, don't dump.** Use `--count` when the user only needs a number.
5. **Confirm destructive actions** (`delete`, `void`) with the user; `delete`
   prompts unless `-y` is passed.

## Workflow: auth → discover → act → verify

```bash
alegra doctor                      # 1. verify auth, plan, rate-limit, country
alegra <resource> --help           # 2. discover a resource's actions & filters
alegra <resource> list -o json     #    inspect real data/ids
# 3. act (always --dry-run first for writes)
alegra <resource> create -f body.json --dry-run
alegra <resource> create -f body.json
alegra <resource> get <id> -o json # 4. verify the result
```

## Core commands

`alegra <resource> {list|get|create|update|delete|export|import}` plus resource
actions. Resources include: `contacts`, `items`, `invoices`, `credit-notes`,
`debit-notes`, `estimates`, `payments`, `bills`, `purchase-orders`, `taxes`,
`retentions`, `terms`, `price-lists`, `bank-accounts`, `number-templates`,
`warehouses`, `journals`, `categories`, `company`, `users`, `reports`, and more
(run `alegra --help` for the full list).

```bash
# Contacts
alegra contacts list --type client --all
alegra contacts create --set name="Acme S.A.S" \
  --set 'identification={"type":"NIT","number":"901123456"}' \
  --set 'type=["client"]' --set kindOfPerson=LEGAL_ENTITY

# Items / taxes (tax id 3 = IVA 19% in Colombia — confirm with `alegra taxes list`)
alegra items create --set name="Consultoría" --set price=500000 --set 'tax=[{"id":3}]'

# Invoices (see Electronic invoicing below for the body)
alegra invoices list --status open --since this-month
alegra invoices get 1234 -o json | jq '{number,status,total,balance}'

# Payments (allocate to invoices via a JSON body)
alegra payments list --type in --since 7d
alegra payments create -f payment.json
```

## Filtering, counting, output

```bash
alegra invoices list --status open --since last-month --until last-month
alegra invoices list --status open --count            # total via API metadata
alegra contacts list --param identification=901123456 # any raw API query param
alegra items list --all -o csv --columns id,name,price > items.csv
```

- `--since/--until`: `today`, `this-month`, `last-month`, `7d`, `3m`, `YYYY-MM-DD`.
- `--param key=value`: escape hatch for any Alegra query parameter.
- `-o table|json|yaml|csv`, `--columns a,b,c`, `--all` (every page), `--count`.

## Writing bodies

create/update accept the body three ways (combine freely):

```bash
alegra invoices create -f invoice.json          # file (best for nested documents)
alegra invoices create -f -                      # stdin
alegra contacts create --set name="X" --set 'type=["client"]'   # flat fields
```

`--set` sets top-level fields (values parsed as JSON when valid); for nested
documents (invoice `items[]`, payment allocations) use `--file`.

## Electronic invoicing (facturación electrónica)

An invoice becomes fiscal when it has a numbering resolution
(`numberTemplate.id`) and is stamped (`stamp.generateStamp: true`), which returns
a CUFE (CO) / UUID (MX). Find resolutions with `alegra number-templates list`.

Safe flow — create a draft, review, then emit:

```bash
alegra invoices create -f invoice.json --draft   # internal draft (validated per country)
alegra invoices get <id>                          # review
alegra invoices emit <id>                         # stamp; or `emit --all` for every draft
```

`invoices emit` auto-chunks into batches of 10 and keeps a local idempotency
guard (skips already-emitted ids unless `--force`). `create` runs country-aware
pre-flight validation (`alegra config set-country colombia`; `--no-validate` to
skip). `invoice.json` (Colombia):

```json
{
  "client": { "id": 12 },
  "numberTemplate": { "id": 7 },
  "date": "2026-06-06", "dueDate": "2026-06-21", "paymentForm": "CASH",
  "items": [ { "id": 5, "price": 50000, "quantity": 1, "tax": [{ "id": 3 }] } ],
  "stamp": { "generateStamp": true }
}
```

## Bulk operations

```bash
alegra contacts import -f clients.csv \
  --map 'Name=name,NIT=identification.number' \
  --set 'identification.type=NIT' --set 'type=["client"]'
alegra invoices export --param status=open > receivables.csv
```

## Reports

```bash
alegra reports sales-by-client --from 2026-01-01 --to 2026-03-31
alegra reports sales-by-seller --from 2026-01-01 --to 2026-12-31 -o csv
```

## Errors

`alegra` prints `Error: alegra: HTTP <code>: …` with a remediation hint. Common:
`401` bad credentials (`alegra auth login`), `402` the feature isn't in the
account's plan, `403` no permission, `429` rate limit (150/min — the CLI backs
off automatically), validation errors name the missing field. On a stamping
error, run `alegra invoices get <id>` to check whether it actually emitted before
retrying.

## More

Full docs and recipes: https://jjuanrivvera.github.io/alegra-cli/ . A condensed
command cheatsheet ships alongside this skill in `references/alegra-commands.md`.
