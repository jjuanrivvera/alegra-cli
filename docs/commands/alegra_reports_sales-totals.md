---
title: alegra reports sales-totals
---

## alegra reports sales-totals

Sales totals grouped by period (day, month, or year)

```
alegra reports sales-totals [flags]
```

### Options

```
      --document-status string   Document status: open, closed, applied
      --from string              Range start date (YYYY-MM-DD)
      --group-by string          Temporal grouping: day, month, year (default "month")
  -h, --help                     help for sales-totals
      --limit int                Number of rows per page (default 10)
      --start int                Pagination offset
      --to string                Range end date (YYYY-MM-DD)
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

* [alegra reports](alegra_reports.md)	 - Read-only Alegra sales reports

