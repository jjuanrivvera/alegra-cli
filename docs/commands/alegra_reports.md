---
title: alegra reports
---

## alegra reports

Read-only Alegra reports

### Synopsis

reports fetches read-only Alegra reports.

Each subcommand GETs a report subpath under /reports. Use --from / --to to bound
the date range and --start / --limit to paginate the paginated reports.

Note: some reports (income-statement, account-statement) are only available on
certain Alegra plans and may return HTTP 403 otherwise.

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
* [alegra reports account-statement](alegra_reports_account-statement.md)	 - Account statement for a contact (estado de cuenta) — plan-dependent
* [alegra reports income-statement](alegra_reports_income-statement.md)	 - Income statement / profit & loss (estado de resultados) — plan-dependent
* [alegra reports sales-by-client](alegra_reports_sales-by-client.md)	 - Sales grouped by client over a date range
* [alegra reports sales-by-client-totals](alegra_reports_sales-by-client-totals.md)	 - Grand totals of the sales-by-client report
* [alegra reports sales-by-seller](alegra_reports_sales-by-seller.md)	 - Sales grouped by seller over a date range

