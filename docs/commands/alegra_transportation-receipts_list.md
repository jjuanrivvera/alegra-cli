---
title: alegra transportation-receipts list
---

## alegra transportation-receipts list

List transportation-receipts

```
alegra transportation-receipts list [flags]
```

### Examples

```
  alegra transportation-receipts list
  alegra transportation-receipts list --limit 30 --all -o json
  alegra transportation-receipts list --count
  alegra transportation-receipts list --status <value>
  alegra transportation-receipts list --param <api_param>=<value>
```

### Options

```
      --all                      Fetch all pages
      --client-id string         Filter by client ID
      --count                    Print only the total number of matching records
      --date string              Filter by date (YYYY-MM-DD)
  -h, --help                     help for list
      --limit int                Max records per page (max 30)
      --order-direction string   Sort direction: ASC or DESC
      --order-field string       Field to sort by (id, date, status)
      --param stringArray        Arbitrary API query parameter: key=value (repeatable; e.g. --param date_after=2026-01-01)
  -q, --query string             Free-text search
      --since string             Start of date range (YYYY-MM-DD, today, this-month, last-month, 7d, 3m, ...)
      --start int                Offset to start from (pagination)
      --status string            Filter by status: open, closed, draft, void
      --until string             End of date range (same formats as --since)
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

* [alegra transportation-receipts](alegra_transportation-receipts.md)	 - Manage transportation receipts (documentos de traslado)

