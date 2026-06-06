---
title: alegra payments
---

## alegra payments

Manage payments (incomes and expenses)

### Options

```
  -h, --help   help for payments
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
* [alegra payments create](alegra_payments_create.md)	 - Create a payment
* [alegra payments delete](alegra_payments_delete.md)	 - Delete a payment by ID
* [alegra payments export](alegra_payments_export.md)	 - Export all payments to CSV or JSON
* [alegra payments get](alegra_payments_get.md)	 - Get a single payment by ID
* [alegra payments import](alegra_payments_import.md)	 - Bulk-create payments from a CSV file
* [alegra payments list](alegra_payments_list.md)	 - List payments
* [alegra payments open](alegra_payments_open.md)	 - Revert a voided payment
* [alegra payments stamp](alegra_payments_stamp.md)	 - Emit electronic payment receipt (REP)
* [alegra payments update](alegra_payments_update.md)	 - Update a payment by ID
* [alegra payments void](alegra_payments_void.md)	 - Void a payment

