---
title: alegra recurring-payments list
---

## alegra recurring-payments list

List recurring-payments

```
alegra recurring-payments list [flags]
```

### Options

```
      --all                      Fetch all pages
      --client-id string         Filter by client ID
  -h, --help                     help for list
      --limit int                Max records per page (max 30)
      --order-direction string   Sort direction: ASC or DESC
      --order-field string       Field to sort by (id, number, date, type)
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

* [alegra recurring-payments](alegra_recurring-payments.md)	 - View recurring payments (pagos recurrentes)

