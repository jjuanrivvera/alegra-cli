---
title: alegra income-debit-notes export
---

## alegra income-debit-notes export

Export all income-debit-notes to CSV or JSON

### Synopsis

Fetch every page of income-debit-notes and write them to a file (or stdout).

```
alegra income-debit-notes export [flags]
```

### Examples

```
  alegra income-debit-notes export > income-debit-notes.csv
  alegra income-debit-notes export --format json --out income-debit-notes.json
  alegra income-debit-notes export --param status=open
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
  -o, --output string               Output format: table, json, yaml, csv (env: ALEGRA_OUTPUT)
      --profile string              Configuration profile to use (env: ALEGRA_PROFILE)
      --requests-per-second float   Client-side rate limit (default from config)
      --show-token                  In --dry-run, do not redact the Authorization header
  -v, --verbose                     Enable verbose (debug) logging to stderr
```

### SEE ALSO

* [alegra income-debit-notes](alegra_income-debit-notes.md)	 - Manage customer debit notes

