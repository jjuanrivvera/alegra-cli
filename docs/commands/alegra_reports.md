---
title: alegra reports
---

## alegra reports

Read-only Alegra sales reports

### Synopsis

reports fetches read-only Alegra sales reports.

Each subcommand GETs a report subpath under /reports and aggregates sales
documents over a date range. Use --from / --to to bound the range and
--start / --limit to paginate.

### Options

```
  -h, --help   help for reports
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

* [alegra](alegra.md)	 - Alegra accounting system CLI
* [alegra reports sales-by-client](alegra_reports_sales-by-client.md)	 - Sales grouped by client over a date range
* [alegra reports sales-by-seller](alegra_reports_sales-by-seller.md)	 - Sales grouped by seller over a date range
* [alegra reports sales-documents](alegra_reports_sales-documents.md)	 - List individual sales documents (invoices, credit/debit notes) in a range
* [alegra reports sales-totals](alegra_reports_sales-totals.md)	 - Sales totals grouped by period (day, month, or year)

