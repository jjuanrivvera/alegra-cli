---
title: alegra recurring-payments
---

## alegra recurring-payments

View recurring payments (pagos recurrentes)

### Synopsis

View Alegra recurring payments, which are automatically registered against a bank account on a recurring schedule. This resource is read-only.

### Options

```
  -h, --help   help for recurring-payments
```

### Options inherited from parent commands

```
      --base-url string             Override the API base URL (env: ALEGRA_BASE_URL)
      --columns strings             Comma-separated columns for table/csv output
      --dry-run                     Print the equivalent curl request without sending it
      --no-color                    Disable colored output (also respects the NO_COLOR env var)
  -o, --output string               Output format: table, json, yaml, csv (env: ALEGRA_OUTPUT)
      --profile string              Configuration profile to use (env: ALEGRA_PROFILE)
      --requests-per-second float   Client-side rate limit (default from config)
      --show-token                  In --dry-run, do not redact the Authorization header
  -v, --verbose                     Enable verbose (debug) logging to stderr
```

### SEE ALSO

* [alegra](alegra.md)	 - Alegra accounting system CLI
* [alegra recurring-payments export](alegra_recurring-payments_export.md)	 - Export all recurring-payments to CSV or JSON
* [alegra recurring-payments get](alegra_recurring-payments_get.md)	 - Get a single recurring-payment by ID
* [alegra recurring-payments list](alegra_recurring-payments_list.md)	 - List recurring-payments

