---
title: Output Formats
---

# Output Formats

Choose with `-o`/`--output` (or `ALEGRA_OUTPUT`, or `settings.defaultOutputFormat`).

| Format | Use |
| --- | --- |
| `table` | Human-readable (default). A list renders as a grid; a single record as key/value. |
| `json` | Pretty-printed JSON — ideal for `jq`. |
| `yaml` | YAML. |
| `csv` | Spreadsheet-friendly; one row per record. |

```bash
alegra invoices list -o json
alegra items list -o csv > items.csv
alegra contacts get 5 -o yaml
```

## Selecting columns

`--columns` overrides the default columns (by JSON key) for `table` and `csv`:

```bash
alegra contacts list --columns id,name,email,status
alegra items list -o csv --columns id,name,price
```

Each resource ships sensible default columns; `json`/`yaml` always include the
full record.

## Scripting tips

```bash
# Count open invoices
alegra invoices list --status open --all -o json | jq 'length'

# Extract IDs
alegra contacts list --all -o json | jq -r '.[].id'
```
