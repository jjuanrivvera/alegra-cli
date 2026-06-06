---
title: alegra warehouse-transfers list
---

## alegra warehouse-transfers list

List warehouse-transfers

```
alegra warehouse-transfers list [flags]
```

### Examples

```
  alegra warehouse-transfers list
  alegra warehouse-transfers list --limit 30 --all -o json
  alegra warehouse-transfers list --count
  alegra warehouse-transfers list --date <value>
  alegra warehouse-transfers list --param <api_param>=<value>
```

### Options

```
      --all                      Fetch all pages
      --count                    Print only the total number of matching records
      --date string              Filter by date (YYYY-MM-DD)
  -h, --help                     help for list
      --item-id string           Filter transfers involving a specific item ID
      --limit int                Max records per page (max 30)
      --order-direction string   Sort direction: ASC or DESC
      --order-field string       Field to sort by (id, date)
      --param stringArray        Arbitrary API query parameter: key=value (repeatable; e.g. --param date_after=2026-01-01)
  -q, --query string             Free-text search
      --since string             Start of date range (YYYY-MM-DD, today, this-month, last-month, 7d, 3m, ...)
      --start int                Offset to start from (pagination)
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

* [alegra warehouse-transfers](alegra_warehouse-transfers.md)	 - Manage inventory transfers between warehouses

