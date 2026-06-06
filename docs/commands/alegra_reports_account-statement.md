---
title: alegra reports account-statement
---

## alegra reports account-statement

Account statement for a contact (estado de cuenta) — plan-dependent

```
alegra reports account-statement [flags]
```

### Examples

```
  alegra reports account-statement --client-id 12 --from 2026-01-01 --to 2026-12-31
```

### Options

```
      --client-id string   Contact (client) ID
      --from string        Range start date (YYYY-MM-DD)
  -h, --help               help for account-statement
      --limit int          Number of rows per page (default 10)
      --start int          Pagination offset
      --to string          Range end date (YYYY-MM-DD)
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

