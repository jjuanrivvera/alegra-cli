---
title: alegra warehouse-transfers list
---

## alegra warehouse-transfers list

List warehouse-transfers

```
alegra warehouse-transfers list [flags]
```

### Options

```
      --all                      Fetch all pages
      --date string              Filter by date (YYYY-MM-DD)
  -h, --help                     help for list
      --item-id string           Filter transfers involving a specific item ID
      --limit int                Max records per page (max 30)
      --order-direction string   Sort direction: ASC or DESC
      --order-field string       Field to sort by (id, date)
  -q, --query string             Free-text search
      --start int                Offset to start from (pagination)
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

