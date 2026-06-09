---
title: Creating & Updating
---

# Creating & Updating Records

Alegra request bodies are often deeply nested (an invoice has line items, taxes,
client references, etc.), so create/update commands accept a **JSON body** three
ways. They can be combined: a base body from `--data`/`--file`, with `--set`
overrides merged on top.

## From a file (or stdin)

```bash
alegra invoices create -f invoice.json
cat invoice.json | alegra invoices create -f -
```

## From an inline JSON string

```bash
alegra contacts create -d '{"name":"Acme","type":["client"]}'
```

## From `--set key=value` pairs

Values are parsed as JSON when valid (numbers, booleans, arrays, objects, null),
otherwise treated as strings:

```bash
alegra contacts create \
  --set name="Acme S.A.S" \
  --set 'type=["client"]' \
  --set 'address={"city":"Cali"}'
```

## Updating

`update` sends a partial body (Alegra applies PATCH semantics — only the fields
you send change):

```bash
alegra invoices update 12 --set observations="Pagada por transferencia"
alegra contacts update 99 -f patch.json
```

## Deleting

```bash
alegra contacts delete 99       # asks for confirmation
alegra contacts delete 99 -y    # skip the prompt
```

## Resource actions

Many documents support extra actions beyond CRUD:

```bash
alegra invoices void 12
alegra invoices open 12
alegra invoices email 12 --set 'emails=["client@acme.com"]'
alegra invoices stamp --set 'invoices=[{"id":12},{"id":13}]'
alegra bills close 7 --set date=2026-06-05 --set 'category={"id":5}'
alegra bank-accounts transfer 3 --set 'destinationAccount={"id":4}' --set amount=100000
```

## Bulk import & export (CSV/JSON)

Almost every resource supports bulk operations.

**Import** creates one record per CSV row. The header row names the fields; use
`--map` to rename a column to an API field (dotted paths build nested objects)
and `--set` to apply a constant to every row. Rows are processed independently —
a failed row is reported and does not stop the run:

```bash
alegra contacts import -f clients.csv \
  --map 'Name=name,NIT=identification.number' \
  --set 'identification.type=NIT' --set 'type=["client"]'

alegra items import -f catalog.csv --map 'SKU=reference,Name=name'
```

Preview every row without writing with `--dry-run`.

**Export** fetches *all* pages and writes them to a file or stdout, as CSV
(default) or JSON, honoring `--columns` and any `--param` filter:

```bash
alegra invoices export --param status=open > receivables.csv
alegra items export --format json --out items.json
alegra contacts export --columns id,name,email,status
```

## Dry run

Preview the exact HTTP request (and a copy-pasteable `curl`) without sending it:

```bash
alegra invoices create -f invoice.json --dry-run
```

Add `--show-token` to include the auth header in the printed curl (off by default).

!!! tip "Pre-flight validation"
    `create` validates the body for your detected country (CO/MX/PE/CR) before
    sending and catches common mistakes early. Skip it with `--no-validate`, or
    set the country explicitly with `--country` (see
    [Configuration](configuration.md)).
