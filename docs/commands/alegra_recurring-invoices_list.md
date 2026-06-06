---
title: alegra recurring-invoices list
---

## alegra recurring-invoices list

List recurring-invoices

```
alegra recurring-invoices list [flags]
```

### Options

```
      --all                      Fetch all pages
      --client-id string         Filter by client ID
      --end-date string          Filter by recurring invoice end date (YYYY-MM-DD)
  -h, --help                     help for list
      --limit int                Max records per page (max 30)
      --name string              Filter by client name
      --order-direction string   Sort direction: ASC or DESC
      --order-field string       Field to sort by (name, startDate, endDate, repeatEvery, term)
  -q, --query string             Free-text search
      --repeat-every string      Filter by recurrence interval
      --start int                Offset to start from (pagination)
      --start-date string        Filter by recurring invoice start date (YYYY-MM-DD)
      --term string              Filter by payment term
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

* [alegra recurring-invoices](alegra_recurring-invoices.md)	 - Manage recurring invoices (facturas recurrentes)

