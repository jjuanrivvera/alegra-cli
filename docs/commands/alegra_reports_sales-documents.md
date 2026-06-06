---
title: alegra reports sales-documents
---

## alegra reports sales-documents

List individual sales documents (invoices, credit/debit notes) in a range

```
alegra reports sales-documents [flags]
```

### Options

```
      --document-number string   Filter by document number
      --document-status string   Document status: open, closed, applied
      --document-types string    Comma-separated types: invoice, creditNote, incomeDebitNote
      --from string              Range start date (YYYY-MM-DD)
  -h, --help                     help for sales-documents
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

