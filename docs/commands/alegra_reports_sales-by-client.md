---
title: alegra reports sales-by-client
---

## alegra reports sales-by-client

Sales grouped by client over a date range

```
alegra reports sales-by-client [flags]
```

### Examples

```
  alegra reports sales-by-client --from 2026-01-01 --to 2026-03-31
```

### Options

```
      --client-name string       Filter by client name
      --from string              Range start date (YYYY-MM-DD)
  -h, --help                     help for sales-by-client
      --limit int                Number of rows per page (default 10)
      --order-direction string   Order direction: ASC or DESC
      --order-field string       Order by: totalDocuments, subTotal, total
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

* [alegra reports](alegra_reports.md)	 - Read-only Alegra reports

