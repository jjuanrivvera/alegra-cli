---
title: alegra recurring-invoices export
---

## alegra recurring-invoices export

Export all recurring-invoices to CSV or JSON

### Synopsis

Fetch every page of recurring-invoices and write them to a file (or stdout).

```
alegra recurring-invoices export [flags]
```

### Examples

```
  alegra recurring-invoices export > recurring-invoices.csv
  alegra recurring-invoices export --format json --out recurring-invoices.json
  alegra recurring-invoices export --param status=open
```

### Options

```
      --format string       Export format: csv or json (default "csv")
  -h, --help                help for export
      --out string          Write to this file (default: stdout)
      --param stringArray   Filter by an API query parameter: key=value (repeatable)
  -q, --query string        Free-text search
```

### Options inherited from parent commands

```
      --base-url string             Override the API base URL (env: ALEGRA_BASE_URL)
      --columns strings             Comma-separated columns for table/csv output
      --dry-run                     Print the equivalent curl request without sending it
      --no-color                    Disable colored output (also respects the NO_COLOR env var)
  -o, --output string               Output format: table, json, yaml, csv (env: ALEGRA_OUTPUT)
      --profile string              Configuration profile to use (env: ALEGRA_PROFILE)
      --requests-per-second float   Client-side rate limit (default from config)
      --show-token                  In --dry-run, do not redact the Authorization header
  -v, --verbose                     Enable verbose (debug) logging to stderr
```

### SEE ALSO

* [alegra recurring-invoices](alegra_recurring-invoices.md)	 - Manage recurring invoices (facturas recurrentes)

